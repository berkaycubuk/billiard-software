<script>
import userService from '../services/userService';
export default {
  data() {
    return {
      user: null,
    }
  },
  mounted() {
    userService.getProfile()
      .then((response) => {
        this.user = response.data.user; 

        if (this.user.is_verified === true) {
          this.$router.push('/');
          return;
        }
      }).catch((err) => {
        // logout if user not found
        if (err.response.data.message === "error.user_not_found" || err.response.data.message === "error.auth_required") {
          this.cookies.remove('token');
          this.$router.push('/login');
          return;
        }
        toast.error(err.response.data.message); 
      });
  }
}
</script>

<template>
  <main class="container" style="flex: 1;">
    <div class="ne__container">
      <h1 class="ne__title">Verify your account</h1>
      <p class="ne__description" v-if="user != null">You have to verify your account. An e-mail sent to <b class="green">{{ user.email }}</b> address.</p>
      <button class="button button--primary">Re-send confirmation</button>
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

.green {
  color: var(--green);
}
</style>
