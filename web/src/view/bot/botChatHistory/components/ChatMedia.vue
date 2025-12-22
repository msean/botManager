<template>
  <div class="media">
    <!-- 图片 -->
    <img
      v-if="isPhoto"
      :src="msg.fileUrl"
      class="photo"
    />

    <!-- 表情包占位 -->
    <div
      v-else-if="isSticker"
      class="sticker-placeholder"
    >
      &lt;表情包&gt;
    </div>

    <video
        v-else-if="isVideo"
        :src="msg.fileUrl"
        controls
        playsinline
        preload="metadata"
        class="video"
    ></video>
    <!-- 以后要用再放开 -->
    <!--
    <video
      v-else-if="isVideo"
      :src="msg.fileUrl"
      controls
      class="video"
    />

    <audio
      v-else-if="isVoice || isAudio"
      :src="msg.fileUrl"
      controls
    />
    -->
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  msg: Object
})

const isPhoto = computed(() => props.msg.fileType === 'photo')
const isVideo = computed(() => props.msg.fileType === 'video')
const isVoice = computed(() => props.msg.fileType === 'voice')
const isAudio = computed(() => props.msg.fileType === 'audio')
const isSticker = computed(() => props.msg.fileType === 'sticker')
</script>

<style scoped>
.photo {
  max-width: 240px;
  max-height: 3200px;
  border-radius: 6px;
}

.video {
  max-width: 320px;
  max-height: 240px;
  border-radius: 6px;
  background: #000;
}

.sticker-placeholder {
  padding: 6px 10px;
  font-size: 13px;
  color: #666;
  background: #f0f2f5;
  border-radius: 6px;
  display: inline-block;
}
</style>
