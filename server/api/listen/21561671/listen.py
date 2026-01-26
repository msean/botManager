import os
import json
import time
import yaml
import psycopg2
import logging
import logging.handlers
import traceback
import asyncio
import signal
from telethon import TelegramClient, events
from telethon.tl.types import Message

# ================== 基础配置 ==================
api_id = 21561671
api_hash = "d007b96854cc2ba7b4e0776e90705fea"

CONFIG_FILE = "../../../config.yaml"

SESSION_NAME = "session"
LOG_FILE = "tg_listener.log"

KEYWORD_FILE = "match.txt"
KEYWORD_RELOAD_INTERVAL = 60  # 秒

# ================== 日志配置 ==================
logger = logging.getLogger("tg_listener")
logger.setLevel(logging.INFO)

file_handler = logging.handlers.RotatingFileHandler(
    LOG_FILE,
    maxBytes=50 * 1024 * 1024,
    backupCount=5,
    encoding="utf-8"
)

formatter = logging.Formatter(
    "%(asctime)s | %(levelname)s | %(message)s"
)

file_handler.setFormatter(formatter)
logger.addHandler(file_handler)

# ================== 读取 YAML ==================
def load_config(path: str) -> dict:
    if not os.path.exists(path):
        raise RuntimeError(f"配置文件不存在: {path}")
    with open(path, "r", encoding="utf-8") as f:
        return yaml.safe_load(f)


def build_pg_dsn(cfg: dict) -> str:
    pg = cfg.get("pgsql")
    if not pg:
        raise RuntimeError("配置文件中缺少 pgsql 节点")

    return (
        f"dbname={pg['db-name']} "
        f"user={pg['username']} "
        f"password={pg['password']} "
        f"host={pg['path']} "
        f"port={pg['port']} "
        f"sslmode=disable "
        f"options='-c timezone=Asia/Shanghai'"
    )

# ================== PostgreSQL ==================
try:
    config = load_config(CONFIG_FILE)
    PG_DSN = build_pg_dsn(config)

    conn = psycopg2.connect(PG_DSN)
    conn.autocommit = True
    logger.info("PostgreSQL 连接成功")
except Exception as e:
    logger.error("PostgreSQL 连接失败: %s", e)
    raise

# ================== Telegram ==================
client = TelegramClient(SESSION_NAME, api_id, api_hash)

# ================== 表缓存 ==================
MESSAGE_TABLE_READY = set()

# ================== 关键词缓存 ==================
KEYWORDS = set()
KEYWORD_FILE_MTIME = 0
LAST_KEYWORD_CHECK = 0


def load_keywords(force=False):
    global KEYWORDS, KEYWORD_FILE_MTIME, LAST_KEYWORD_CHECK

    now = time.time()
    if not force and now - LAST_KEYWORD_CHECK < KEYWORD_RELOAD_INTERVAL:
        return

    LAST_KEYWORD_CHECK = now

    if not os.path.exists(KEYWORD_FILE):
        logger.warning("关键词文件不存在: %s", KEYWORD_FILE)
        KEYWORDS = set()
        return

    mtime = os.path.getmtime(KEYWORD_FILE)
    if mtime == KEYWORD_FILE_MTIME and not force:
        return

    try:
        with open(KEYWORD_FILE, "r", encoding="utf-8") as f:
            KEYWORDS = {
                line.strip()
                for line in f
                if line.strip() and not line.startswith("#")
            }
        KEYWORD_FILE_MTIME = mtime
        logger.info("关键词加载成功 count=%d", len(KEYWORDS))
    except Exception as e:
        logger.error("关键词加载失败: %s", e)
        KEYWORDS = set()

# ================== 表相关 ==================
def message_table(group_id: int) -> str:
    return f'bot_messages_{abs(group_id)}'


def ensure_chat_map_table(cur):
    cur.execute("""
    CREATE TABLE IF NOT EXISTS bot_chat_map (
        group_id BIGINT PRIMARY KEY,
        group_type SMALLINT NOT NULL,
        group_name VARCHAR(255),
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    )
    """)


def ensure_message_table(cur, table: str):
    cur.execute(f"""
    CREATE TABLE IF NOT EXISTS "{table}" (
        id BIGINT PRIMARY KEY,
        user_id BIGINT,
        username VARCHAR(64),
        first_name VARCHAR(128),
        last_name VARCHAR(128),
        nick_name VARCHAR(64),
        is_bot BOOLEAN,
        reply_to_message_id BIGINT,
        message_type VARCHAR(50),
        text TEXT,
        caption TEXT,
        file_id TEXT,
        file_unique_id TEXT,
        file_type VARCHAR(32),
        timestamp TIMESTAMP,
        raw JSONB,
        created_at TIMESTAMP DEFAULT NOW()
    )
    """)


def ensure_indexes(cur, table: str):
    idx = abs(hash(table)) % 100000

    cur.execute(f"""CREATE INDEX IF NOT EXISTS idx_{idx}_time ON "{table}"(timestamp)""")
    cur.execute(f"""CREATE INDEX IF NOT EXISTS idx_{idx}_user_id ON "{table}"(user_id)""")
    cur.execute("CREATE EXTENSION IF NOT EXISTS pg_trgm")
    cur.execute(f"""
        CREATE INDEX IF NOT EXISTS idx_{idx}_text_trgm
        ON "{table}" USING gin (
            (coalesce(text,'') || ' ' || coalesce(caption,'')) gin_trgm_ops
        )
    """)

# ================== Telegram 监听 ==================
@client.on(events.NewMessage)
async def handler(event):
    try:
        if not (event.is_group or event.is_channel):
            return

        msg: Message = event.message
        text = msg.text or ""
        caption = getattr(msg, "message", "") or ""
        content = f"{text} {caption}"

        load_keywords()
        hit = next((kw for kw in KEYWORDS if kw in content), None)
        if not hit:
            return

        chat = await event.get_chat()
        sender = await event.get_sender()

        table = message_table(chat.id)
        cur = conn.cursor()

        cur.execute("""
            INSERT INTO bot_chat_map (group_id, group_type, group_name)
            VALUES (%s,%s,%s)
            ON CONFLICT (group_id) DO UPDATE
            SET group_name = EXCLUDED.group_name,
                updated_at = NOW()
        """, (chat.id, 2 if event.is_channel else 1, getattr(chat, "title", "")))

        if table not in MESSAGE_TABLE_READY:
            ensure_message_table(cur, table)
            ensure_indexes(cur, table)
            MESSAGE_TABLE_READY.add(table)

        cur.execute(
            f"""
            INSERT INTO "{table}" (
                id,user_id,username,first_name,last_name,nick_name,is_bot,
                reply_to_message_id,message_type,text,caption,
                file_id,file_unique_id,file_type,timestamp,raw
            )
            VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)
            ON CONFLICT (id) DO NOTHING
            """,
            (
                msg.id,
                sender.id if sender else None,
                getattr(sender, "username", None),
                getattr(sender, "first_name", None),
                getattr(sender, "last_name", None),
                f"{getattr(sender, 'first_name', '')}{getattr(sender, 'last_name', '')}",
                getattr(sender, "bot", False),
                msg.reply_to_msg_id,
                msg.__class__.__name__,
                text,
                caption,
                None, None, None,
                msg.date,
                json.dumps(msg.to_dict(), ensure_ascii=False, default=str)
            )
        )

        logger.info("命中关键词=%s group_id=%s msg_id=%s", hit, chat.id, msg.id)

    except Exception:
        logger.error("消息处理异常\n%s", traceback.format_exc())

# ================== 优雅退出 ==================
async def shutdown(sig):
    logger.info("收到退出信号 %s，关闭 Telegram 客户端", sig)
    await client.disconnect()

# ================== 启动 ==================
if __name__ == "__main__":
    logger.info("Telegram Listener 启动")

    cur = conn.cursor()
    ensure_chat_map_table(cur)
    load_keywords(force=True)

    loop = asyncio.get_event_loop()
    for s in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(s, lambda sig=s: asyncio.create_task(shutdown(sig)))

    client.start()
    logger.info("开始监听 Telegram")

    try:
        loop.run_until_complete(client.run_until_disconnected())
    finally:
        loop.run_until_complete(client.disconnect())
        loop.close()
        logger.info("监听器已安全退出")
