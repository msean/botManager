<template>
  <div class="text">
    {{ msg.text }}
  </div>
</template>

<script setup>
defineProps({ msg: Object })
</script>

<style scoped>
.text {
  background: #fff;
  padding: 8px 12px;
  border-radius: 6px;
  word-break: break-word;
}
</style>
