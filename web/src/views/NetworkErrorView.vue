<script>
  import userService from '../services/userService';
  export default {
  mounted() {
    userService.getProfile()
      .then((res) => {
        if (res.data.success === true) {
          this.$router.push('/');
        }
      }); 
  },
  }
</script>

<template>
  <main class="container" style="flex: 1;">
    <div class="ne__container">
      <h1 class="ne__title">Network error</h1>
      <p class="ne__description">Currently unable to access server. If still not working after refresh, server could be down.</p>
      <router-link to="/" class="button button--primary">
        <svg width="24px" height="24px" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#f8faed"><path d="M21.888 13.5C21.164 18.311 17.013 22 12 22 6.477 22 2 17.523 2 12S6.477 2 12 2c4.1 0 7.625 2.468 9.168 6" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path><path d="M17 8h4.4a.6.6 0 00.6-.6V3" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
        Retry
      </router-link>
    </div>
  </main>
</template>

<style scoped>
.ne__container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
}

.ne__title {
  margin-bottom: 20px;
}

.ne__description {
  font-size: 1.2rem;
  margin-bottom: 20px;
}
</style>
