<script>
export default {
  props: ['title', 'maxHeight'],
  methods: {
    maxHeightStyle() {
      if (this.maxHeight && this.maxHeight != null && this.maxHeight != "") {
        return `max-height: ${this.maxHeight}px; overflow-y: auto;`;
      }

      return "";
    },
  },
}
</script>

<template>
  <div class="card">
    <div class="card__header">
      <div class="card__title">{{ title }}</div>
      <div class="card__actions">
        <slot name="actions"></slot>
      </div>
    </div>
    <div class="card__body" :style="maxHeightStyle()">
      <slot></slot>
    </div>
  </div>
</template>

<style scoped>
.card {
  padding: 20px;
  background-color: #191919;
  border-radius: 8px;
}

.card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card__title {
  font-size: 1.2rem;
  font-weight: 600;
}

.card__action {
}

.card__body {
  margin-top: 10px;
}
</style>
