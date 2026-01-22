import os
import json
import time
import psycopg2
import logging
import logging.handlers
import traceback
from datetime import datetime
from telethon import TelegramClient, events
from telethon.tl.types import Message

# ================== 基础配置 ==================
api_id = 21561671
api_hash = "d007b96854cc2ba7b4e0776e90705fea"

PG_DSN = "dbname=bot_manager user=admin password=123456 host=16.162.105.189 port=5432"

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

# ================== PostgreSQL ==================
try:
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
    """关键词热加载"""
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
        logger.info("关键词加载成功 count=%d keywords=%s", len(KEYWORDS), list(KEYWORDS))

    except Exception as e:
        logger.error("关键词加载失败: %s", e)
        KEYWORDS = set()


# ================== 表相关 ==================
def message_table(group_id: int) -> str:
    return f'bot_messages_{abs(group_id)}'


def ensure_chat_map_table(cur):
    """群/频道 映射表（只启动时一次）"""
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
    """索引（完全对齐你 Go 里的那套）"""
    index_suffix = abs(hash(table)) % 100000

    cur.execute(f"""
        CREATE INDEX IF NOT EXISTS idx_{index_suffix}_time
        ON "{table}"(timestamp)
    """)

    cur.execute(f"""
        CREATE INDEX IF NOT EXISTS idx_{index_suffix}_user_id
        ON "{table}"(user_id)
    """)

    cur.execute("""
        CREATE EXTENSION IF NOT EXISTS pg_trgm
    """)

    cur.execute(f"""
        CREATE INDEX IF NOT EXISTS idx_{index_suffix}_username_trgm
        ON "{table}" USING gin (username gin_trgm_ops)
    """)

    cur.execute(f"""
        CREATE INDEX IF NOT EXISTS idx_{index_suffix}_text_caption_trgm
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
        caption = getattr(msg, "message", None) or ""
        content = f"{text} {caption}"

        load_keywords()

        hit = None
        for kw in KEYWORDS:
            if kw in content:
                hit = kw
                break

        if not hit:
            return

        chat = await event.get_chat()
        sender = await event.get_sender()

        group_id = chat.id
        group_type = 2 if event.is_channel else 1
        group_name = getattr(chat, "title", "")

        table = message_table(group_id)
        cur = conn.cursor()

        # ===== 群映射表 =====
        cur.execute("""
            INSERT INTO bot_chat_map (group_id, group_type, group_name)
            VALUES (%s,%s,%s)
            ON CONFLICT (group_id) DO UPDATE
            SET group_name = EXCLUDED.group_name,
                updated_at = NOW()
        """, (group_id, group_type, group_name))

        # ===== 消息表 + 索引 =====
        if table not in MESSAGE_TABLE_READY:
            ensure_message_table(cur, table)
            ensure_indexes(cur, table)
            MESSAGE_TABLE_READY.add(table)
            logger.info("消息表初始化完成 table=%s", table)

        raw_json = json.dumps(
            msg.to_dict(),
            ensure_ascii=False,
            default=str  # 解决 datetime 序列化问题
        )

        cur.execute(
            f"""
            INSERT INTO "{table}" (
                id,
                user_id, username, first_name, last_name, nick_name, is_bot,
                reply_to_message_id,
                message_type, text, caption,
                file_id, file_unique_id, file_type,
                timestamp, raw
            )
            VALUES (
                %s,
                %s,%s,%s,%s,%s,%s,
                %s,
                %s,%s,%s,
                %s,%s,%s,
                %s,%s
            )
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

                msg.date.replace(tzinfo=None),  # ✅ UTC 原样入库
                raw_json
            )
        )

        logger.info(
            "命中关键词=%s group_id=%s msg_id=%s",
            hit, group_id, msg.id
        )

    except Exception:
        logger.error("消息处理异常\n%s", traceback.format_exc())


# ================== 启动 ==================
if __name__ == "__main__":
    logger.info("Telegram Listener 启动")

    cur = conn.cursor()
    ensure_chat_map_table(cur)

    load_keywords(force=True)

    if not KEYWORDS:
        logger.warning("关键词为空，不会入库任何消息")

    client.start()
    logger.info("Telegram Listener 已启动，开始监听")
    client.run_until_disconnected()
