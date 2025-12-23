<template>
  <div class="chat-bubble" :class="{ self: isSelf }">
    <div class="avatar">{{ fullName }}</div>

    <div class="content">
      <!-- 用户名 + 时间 -->
      <div class="meta">
        <span class="fullName">{{ fullName }}</span>
        <span class="userName">{{ msg.username }}</span>
        <span class="userID">{{ msg.userId }}</span>
        <span class="time">{{ formatDate(msg.timestamp) }}</span>
    </div>
      <!-- 转发 -->
      <div v-if="forwardFrom" class="forward-bubble">
        Forwarded from {{ forwardFrom }}
      </div>

      <!-- 回复 -->
      <div v-if="replyText" class="reply-bubble">
        <span class="reply-username">
          回复 {{ replyUsername || '未知用户' }}:
        </span>
        <span class="reply-text">{{ replyText }}</span>
      </div>

      <!-- 媒体消息 -->
      <ChatMedia v-if="!isText" :msg="msg" />

      <!-- 文本消息 -->
      <ChatText v-else :text="msg.text" />

      <!-- caption（媒体说明） -->
      <div v-if="msg.caption" class="caption">
        {{ msg.caption }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { formatDate } from '@/utils/format'
import ChatText from './ChatText.vue'
import ChatMedia from './ChatMedia.vue'

const props = defineProps({
  msg: Object,
  selfId: Number
})

const isSelf = computed(() => props.msg.userId === props.selfId)

// ✅ 关键修复点
const isText = computed(() => props.msg.messageType === 'text')

const fullName = computed(() => {
  const f = props.msg.firstName || ''
  const l = props.msg.lastName || ''
  return (f + ' ' + l).trim() || props.msg.username || '未知用户'
})

const forwardFrom = computed(() => {
  return props.msg.raw?.forward_from?.first_name || ''
})

const replyText = computed(() => props.msg.replyText)
const replyUsername = computed(() => props.msg.replyUsername)
</script>

<style scoped>
.caption {
  margin-top: 6px;
  font-size: 14px;
  color: #333;
}
</style>

<style scoped>
.chat-bubble {
  display: flex;
  margin-bottom: 12px;
  align-items: flex-start;
}

.chat-bubble.self {
  flex-direction: row-reverse;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #409eff;
  color: #fff;
  text-align: center;
  line-height: 40px;
  font-weight: bold;
  margin: 0 10px;
  flex-shrink: 0;
}

.meta {
  display: flex;
  align-items: center;
  gap: 8px; /* 想多大就多大 */
}

.content {
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

.meta {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
  display: flex;
  justify-content: space-between;
}

.chat-bubble.self .meta {
  text-align: right;
}

/* 回复消息气泡 */
.reply-bubble {
  background: #f0f0f0;
  border-left: 3px solid #409eff;
  padding: 6px 10px;
  margin-bottom: 4px;
  border-radius: 6px;
  font-size: 13px;
  color: #555;
}

.reply-username {
  font-weight: bold;
  margin-right: 4px;
}

.reply-text {
  color: #333;
}

/* 转发消息样式 */
.forward-bubble {
  background: #e0f7fa;
  color: #00796b;
  padding: 4px 8px;
  margin-bottom: 4px;
  border-radius: 4px;
  font-size: 12px;
}
</style>
