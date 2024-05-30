<script>
export default {
   props: {
    title: String
   }
}
</script>

<template>
  <div class="popup-outside">
    <div class="popup">
      <div class="popup__header">
        <h2 class="popup__title">{{ title }}</h2>
        <button @click="$emit('closed')" class="button button--secondary button--small button--flex-center">
          <svg width="24px" height="24px" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#ffffff"><path d="M6.758 17.243L12.001 12m5.243-5.243L12 12m0 0L6.758 6.757M12.001 12l5.243 5.243" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
        </button>
      </div>
      <div class="popup__content">
        <slot></slot>
      </div>
    </div> 
  </div>
</template>

<style scoped>
.popup-outside {
  position: fixed;
  z-index: 99;
  top: 0;
  left: 0;
  width: 100%;
  padding: 0 40px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
}

.popup {
  width: 100%;
  max-width: 560px;
  padding: 20px;
  background-color: #191919;
  border-radius: 8px;
}

.popup__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.popup__title {
  margin-bottom: 20px;
}

.popup__content {
}
</style>
