<template>
  <div class="chat-page">
    <!-- 查询条件 -->
    <el-card class="mb-2">
      <el-form :inline="true" size="small">
        <el-form-item label="机器人" required>
          <el-select
            v-model="search.botID"
            placeholder="请选择机器人"
            style="width: 220px"
            @change="onBotChange"
          >
            <el-option
              v-for="bot in botList"
              :key="bot.botID"
              :label="bot.name"
              :value="bot.botID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="群聊" required>
          <el-select
            v-model="search.chatGroupID"
            placeholder="请选择群聊"
            style="width: 260px"
          >
            <el-option
              v-for="g in chatGroupList"
              :key="g.chatGroupID"
              :label="g.chatGroupName"
              :value="g.chatGroupID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="用户名">
          <el-input v-model="search.username" clearable />
        </el-form-item>

        <el-form-item label="用户ID">
          <el-input v-model="search.userId" clearable />
        </el-form-item>

        <el-form-item label="消息">
          <el-input v-model="search.text" clearable />
        </el-form-item>

        <el-form-item label="时间">
          <el-date-picker
            v-model="search.timeRange"
            type="datetimerange"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="onInitLoad">
            加载聊天
          </el-button>
          <el-button @click="onReset">
            重置
        </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 聊天区域 -->
    <div class="chat-container" ref="chatContainer">
      <div v-if="loadingTop" class="loading">加载更早消息中...</div>

      <ChatBubble
        v-for="item in messageList"
        :key="item.id"
        :msg="item"
        :self-id="selfId"
      />

      <div v-if="loadingBottom" class="loading">加载更新消息中...</div>
      <el-empty v-if="!messageList.length && !loadingBottom" />
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import { getBotChoiceWithChatGroup } from '@/api/bot/bot'
import { getChatMessageList } from '@/api/bot/botChatGroup'
import ChatBubble from './components/ChatBubble.vue'

/* ===== 常量 ===== */
const LIMIT = 10
const SCROLL_OFFSET = 20
const selfId = 0 // 你的机器人 userId

/* ===== DOM ===== */
const chatContainer = ref(null)
let scrollEl = null

/* ===== 状态 ===== */
const botList = ref([])
const chatGroupList = ref([])
const messageList = ref([])

const beforeId = ref(null)
const afterId = ref(null)

const hasMoreUp = ref(true)
const hasMoreDown = ref(true)

const loadingTop = ref(false)
const loadingBottom = ref(false)
const isInitLoading = ref(false)

/* ===== 搜索条件 ===== */
const search = ref({
  botID: undefined,
  chatGroupID: undefined,
  username: '',
  userId: '',
  text: '',
  timeRange: null,
  startTime: null,
  endTime: null
})

/* ===== 初始化机器人 ===== */
const loadBots = async () => {
  const res = await getBotChoiceWithChatGroup()
  if (res.code === 0) botList.value = res.data
}

/* ===== 切换机器人 ===== */
const onBotChange = botID => {
  const bot = botList.value.find(b => b.botID === botID)
  chatGroupList.value = bot?.botChatGroups || []
  search.value.chatGroupID = undefined
}

/* ===== 初始化加载 ===== */
const onInitLoad = async () => {
  if (!search.value.botID || !search.value.chatGroupID) {
    return ElMessage.warning('请选择机器人和群聊')
  }

  isInitLoading.value = true
  messageList.value = []
  beforeId.value = null
  afterId.value = null
  hasMoreUp.value = true
  hasMoreDown.value = true

  if (search.value.timeRange?.length === 2) {
    search.value.startTime = search.value.timeRange[0]
    search.value.endTime = search.value.timeRange[1]
  } else {
    search.value.startTime = null
    search.value.endTime = null
  }

  await loadInitialMessages()
  await nextTick()
  scrollToBottom()
  isInitLoading.value = false
}

/* ===== 初次加载消息 ===== */
// 初次加载消息
const loadInitialMessages = async () => {
  loadingBottom.value = true
  const res = await getChatMessageList({
    ...search.value,
    limit: LIMIT
  })
  if (res.code === 0) {
    const list = res.data.list || []

    // 补充转发与回复字段
    list.forEach(msg => {
      if (msg.raw?.forward_from) {
        msg.ForwardFromUsername =
          msg.raw.forward_from.first_name +
          (msg.raw.forward_from.last_name ? ' ' + msg.raw.forward_from.last_name : '')
      }
      if (msg.raw?.reply_to_message) {
        msg.replyText = msg.raw.reply_to_message.text || ''
        msg.replyUsername =
          msg.raw.reply_to_message.from?.first_name ||
          msg.raw.reply_to_message.username ||
          '未知用户'
      }
    })

    messageList.value = list
    if (list.length) {
      beforeId.value = list[0].id
      afterId.value = list[list.length - 1].id
    }
    hasMoreUp.value = res.data.hasMore
    hasMoreDown.value = true
  }
  loadingBottom.value = false
}

const initialSearchState = () => ({
  botID: undefined,
  chatGroupID: undefined,
  username: '',
  userId: '',
  text: '',
  timeRange: null,
  startTime: null,
  endTime: null
})


const onReset = () => {
  // 重置搜索条件
  search.value = initialSearchState()

  // 清空下拉依赖
  chatGroupList.value = []

  // 清空消息数据
  messageList.value = []

  // 重置游标
  beforeId.value = null
  afterId.value = null

  // 重置分页状态
  hasMoreUp.value = true
  hasMoreDown.value = true

  // 重置 loading
  loadingTop.value = false
  loadingBottom.value = false
  isInitLoading.value = false

  // 滚动回顶部
  nextTick(() => {
    if (scrollEl) scrollEl.scrollTop = 0
  })
}


/* ===== 向上加载历史 ===== */
const loadMoreUp = async () => {
  if (!beforeId.value || !hasMoreUp.value || loadingTop.value) return
  loadingTop.value = true
  const prevHeight = scrollEl.scrollHeight

  const res = await getChatMessageList({
    ...search.value,
    beforeId: beforeId.value,
    limit: LIMIT
  })
  if (res.code === 0) {
    const list = res.data.list || []
    if (list.length) {
      messageList.value = [...list, ...messageList.value]
      beforeId.value = list[0].id
      await nextTick()
      scrollEl.scrollTop += scrollEl.scrollHeight - prevHeight
    } else {
      hasMoreUp.value = false
    }
  }
  loadingTop.value = false
}

/* ===== 向下加载更新消息 ===== */
const loadMoreDown = async () => {
  if (!afterId.value || !hasMoreDown.value || loadingBottom.value) return
  loadingBottom.value = true
  const res = await getChatMessageList({
    ...search.value,
    afterId: afterId.value,
    limit: LIMIT
  })
  if (res.code === 0) {
    const list = res.data.list || []
    if (list.length) {
      messageList.value = [...messageList.value, ...list]
      afterId.value = list[list.length - 1].id
    } else {
      hasMoreDown.value = false
    }
  }
  loadingBottom.value = false
}

/* ===== 滚动监听 ===== */
const onScroll = async e => {
  if (isInitLoading.value) return
  const el = e.target
  if (el.scrollTop <= SCROLL_OFFSET) await loadMoreUp()
  if (el.scrollHeight - el.scrollTop - el.clientHeight <= SCROLL_OFFSET) await loadMoreDown()
}

/* ===== 滚动到底 ===== */
const scrollToBottom = () => {
  if (scrollEl) scrollEl.scrollTop = scrollEl.scrollHeight
}

/* ===== 绑定滚动事件 ===== */
onMounted(() => {
  nextTick(() => {
    scrollEl = chatContainer.value
    scrollEl?.addEventListener('scroll', onScroll)
  })
})

onBeforeUnmount(() => {
  scrollEl?.removeEventListener('scroll', onScroll)
})

loadBots()
</script>

<style scoped>
.chat-container {
  height: 600px;
  overflow-y: auto;
  padding: 14px;
  background: #f5f7fa;
  border-radius: 6px;
  box-shadow: inset 0 0 5px rgba(0,0,0,0.05);
}

.loading {
  text-align: center;
  font-size: 12px;
  color: #999;
  padding: 8px 0;
}
</style>
