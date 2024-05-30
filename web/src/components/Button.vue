<script>
export default {
  props: {
    text: String,
    type: String = "primary",
  }
}
</script>

<template>
  <button class="button" :class="{ 'button--primary': type === 'primary' }">{{ text }}</button>
</template>

<style scoped>
.button {
  border-radius:
}
</style>
