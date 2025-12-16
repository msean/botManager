<template>
  <div class="chat-page">
    <!-- 查询条件 -->
    <el-card class="mb-2">
      <el-form :inline="true">
        <!-- 机器人 -->
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

        <!-- 群聊 -->
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

        <el-form-item>
          <el-button type="primary" @click="onInitLoad">加载聊天</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 聊天区域 -->
    <el-card class="chat-box" ref="chatBox" @scroll="onScroll">
      <div v-if="loadingTop" class="loading">加载更新消息中...</div>

      <div v-for="item in messageList" :key="item.id" class="chat-item">
        <div class="avatar">
          {{ item.username?.slice(0, 1) || 'U' }}
        </div>

        <div class="content">
          <div class="meta">
            <span class="username">
              {{ item.username || item.firstName || '未知用户' }}
              <small>({{ item.userId }})</small>
            </span>
            <span class="time">{{ formatDate(item.timestamp) }}</span>
          </div>

          <!-- 回复 -->
          <div v-if="item.replyId" class="reply-box">
            <strong>{{ item.replyUsername }}：</strong>
            {{ item.replyText }}
          </div>

          <!-- 正文 -->
          <div class="bubble">
            {{ item.text || item.caption || `[${item.messageType}]` }}
          </div>
        </div>
      </div>

      <div v-if="loadingBottom" class="loading">加载更早消息中...</div>
      <el-empty v-if="!messageList.length && !loadingBottom" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { getBotChoiceWithChatGroup } from '@/api/bot/bot'
import { getChatMessageList } from '@/api/bot/botChatGroup'
import { formatDate } from '@/utils/format'

const botList = ref([])
const chatGroupList = ref([])
const messageList = ref([])

const page = ref(1)
const pageSize = ref(100)

const hasMoreUp = ref(true)
const hasMoreDown = ref(true)

const loadingTop = ref(false)
const loadingBottom = ref(false)

const chatBox = ref(null)

const search = ref({
  botID: undefined,
  chatGroupID: undefined,
  username: '',
  userId: ''
})

/* 加载机器人 */
const loadBots = async () => {
  const res = await getBotChoiceWithChatGroup()
  if (res.code === 0) botList.value = res.data
}

/* 切换机器人 */
const onBotChange = botID => {
  const bot = botList.value.find(b => b.botID === botID)
  chatGroupList.value = bot?.botChatGroups || []
  search.value.chatGroupID = undefined
}

/* 初始化加载 */
const onInitLoad = async () => {
  if (!search.value.botID || !search.value.chatGroupID) {
    return ElMessage.warning('请选择机器人和群聊')
  }

  page.value = 1
  hasMoreUp.value = true
  hasMoreDown.value = true
  messageList.value = []

  await loadPage()
  scrollToBottom()
}

/* 加载当前页 */
const loadPage = async () => {
  const res = await getChatMessageList({
    ...search.value,
    page: page.value,
    pageSize: pageSize.value
  })

  if (res.code !== 0) return

  const list = res.data.list || []

  if (!list.length) {
    if (loadingTop.value) hasMoreUp.value = false
    if (loadingBottom.value) hasMoreDown.value = false
    return
  }

  if (loadingTop.value) {
    messageList.value = [...list, ...messageList.value]
  } else {
    messageList.value = [...messageList.value, ...list]
  }
}

/* 滚动监听 */
const onScroll = async e => {
  const el = e.target

  // 顶部
  if (el.scrollTop === 0 && hasMoreUp.value && !loadingTop.value) {
    loadingTop.value = true
    page.value--
    await loadPage()
    loadingTop.value = false
  }

  // 底部
  if (
    el.scrollHeight - el.scrollTop - el.clientHeight < 10 &&
    hasMoreDown.value &&
    !loadingBottom.value
  ) {
    loadingBottom.value = true
    page.value++
    await loadPage()
    loadingBottom.value = false
  }
}

/* 滚动到底 */
const scrollToBottom = async () => {
  await nextTick()
  const el = chatBox.value?.$el
  if (el) el.scrollTop = el.scrollHeight
}

loadBots()
</script>


<style scoped>
.chat-box {
  height: 600px;
  overflow-y: auto;
}
.chat-item {
  display: flex;
  margin-bottom: 14px;
}
.avatar {
  width: 40px;
  height: 40px;
  background: #409eff;
  color: #fff;
  border-radius: 50%;
  text-align: center;
  line-height: 40px;
  margin-right: 12px;
}
.meta {
  font-size: 12px;
  color: #999;
}
.bubble {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 6px;
}
.reply-box {
  background: #eef3ff;
  padding: 6px;
  font-size: 12px;
  margin-bottom: 6px;
  border-left: 3px solid #409eff;
}
.loading {
  text-align: center;
  color: #999;
  padding: 8px;
}
</style>
