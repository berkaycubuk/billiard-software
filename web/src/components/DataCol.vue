<script>
export default {
  props: {
    align: String,
  },
  methods: {
    convertAlignToFlex() {
      let string = "display: flex;";
      if (this.align === "right") {
        return string + " justify-content: flex-end;";
      } else if (this.align === "center") {
        return string + " justify-content: center;";
      }

      return "";
    }
  },
}
</script>

<template>
  <div class="data-col" :style="convertAlignToFlex()">
    <slot></slot>
  </div>
</template>

<style scoped>
@media only screen and (max-width: 600px) {
  .data-col {
    justify-content: flex-start !important;
  }
}
</style>
