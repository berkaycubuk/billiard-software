<script>
import { toast } from 'vue3-toastify';
import authService from '../services/authService';

export default {
  mounted() {
    authService.verify(this.$route.params.token)
      .then((res) => {
        if (res.data.success == false) {
          toast.error(res.data.message);
          return;
        }

        this.$router.push("/login?msg=success.account_verified");
      });
  }
}
</script>

<template>
  <main class="container" style="flex: 1;">
    <div class="ne__container">
      <p class="ne__description">loading...</p>
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
