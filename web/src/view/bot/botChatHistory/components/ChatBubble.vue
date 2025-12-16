<template>
  <div class="bubble" :class="isSelf ? 'self' : 'other'">
    <div class="avatar">
      {{ avatar }}
    </div>

    <div class="content">
      <div class="username">
        {{ msg.firstName }} {{ msg.lastName }}
        <span v-if="msg.username">@{{ msg.username }}</span>
      </div>

      <!-- 文本 -->
      <ChatText v-if="msg.messageType === 'text'" :msg="msg" />

      <!-- 图片 / 视频 -->
      <ChatMedia
        v-else-if="['photo','video','document'].includes(msg.messageType)"
        :msg="msg"
      />

      <div class="time">
        {{ formatTime(msg.timestamp) }}
      </div>
    </div>
  </div>
</template>

<script setup>
import ChatText from './ChatText.vue'
import ChatMedia from './ChatMedia.vue'
import dayjs from 'dayjs'

const props = defineProps({
  msg: Object,
  selfId: Number
})

const isSelf = props.msg.userId === props.selfId
const avatar = props.msg.firstName?.[0] || '👤'

const formatTime = t => dayjs(t).format('HH:mm')
</script>

<style scoped>
.bubble {
  display: flex;
  margin-bottom: 12px;
}
.other {
  flex-direction: row;
}
.self {
  flex-direction: row-reverse;
}
.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #409eff;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}
.content {
  max-width: 70%;
  margin: 0 10px;
}
.username {
  font-size: 12px;
  color: #666;
}
.time {
  font-size: 11px;
  color: #999;
  text-align: right;
}
</style>
