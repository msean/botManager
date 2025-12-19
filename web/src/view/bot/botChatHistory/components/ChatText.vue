<template>
  <div class="text-bubble">
    {{ text }}
  </div>
</template>

<script setup>
defineProps({
  text: String
})
</script>

<style scoped>
.text-bubble {
  background: #fff;
  padding: 8px 12px;
  border-radius: 12px;
  word-break: break-word;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  line-height: 1.5;
}
</style>
