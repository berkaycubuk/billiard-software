<script>
export default {
  props: {
    cols: Number = 1,
  },
  methods: {
    gridColString() {
      let string = "";

      for (let i = 0; i < this.cols; i++) {
        string += "1fr ";
      }

      return string;
    }
  }
}
</script>

<template>
  <div class="data-row" :style="'grid-template-columns: ' + gridColString()">
    <slot></slot>
  </div>
</template>

<style scoped>
.data-row {
  display: grid;
  align-items: center;
  padding: 10px;
  background-color: #232323;
  border-radius: 8px;
}

@media only screen and (max-width: 600px) {
  .data-row {
    grid-template-columns: 1fr !important;
    gap: 4px;
  }
}
</style>
