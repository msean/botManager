<template>
  <div class="media">
    <img
      v-if="msg.messageType === 'photo'"
      :src="fileUrl"
    />

    <video
      v-else-if="msg.messageType === 'video'"
      controls
      :src="fileUrl"
    />

    <a v-else :href="fileUrl" target="_blank">
      📎 下载文件
    </a>

    <div v-if="msg.caption" class="caption">
      {{ msg.caption }}
    </div>
  </div>
</template>

<script setup>
const props = defineProps({ msg: Object })

// 你后端需要提供一个 file proxy 接口
const fileUrl = `/api/chat/file/${props.msg.fileId}`
</script>

<style scoped>
.media img {
  max-width: 240px;
  border-radius: 6px;
}
.caption {
  margin-top: 4px;
  font-size: 12px;
}
</style>
