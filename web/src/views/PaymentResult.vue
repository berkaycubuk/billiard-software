<script>
import { useRoute } from 'vue-router';
import paymentService from '../services/paymentService';

function paymentStatusString(code) {
  let response = "";
  switch (code) {
    case 1:
      response = "waiting";
      break;
    case 2:
      response = "canceled";
      break;
    case 3:
      response = "paid";
      break;
    case 7:
      response = "paid";
      break;
    case 8:
      response = "paid";
      break;
    case 4:
      response = "deleted";
      break;
  }

  return response;
}

export default {
  data() {
    return {
      id: null,
      status: 1,
      timeoutId: null,
    }
  },
  mounted() {
    this.id = this.$route.params.id;

    paymentService.status(parseInt(this.id), 2)
      .then((res) => {
        if (res.data.success === false) {
          console.error(res.data.message);
          return;
        }

        this.status = res.data.status;
      });

    this.timeoutId = setTimeout(() => {
      this.$router.push("/")
    }, 10000)
  },
  beforeRouteLeave() {
    if (this.timeoutId != null) {
      clearTimeout(this.timeoutId)
      this.timeoutId = null
    }
  },
}
</script>

<template>
  <main class="container" style="flex: 1;">

    <div class="cards">

      <div class="card">
        <div class="card__body">

          <div class="payment-result" v-if="status === 1">
            <div class="payment-result__title">{{ $t('waiting_payment') }}</div>
            <div class="payment-result__icon">
              <svg stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#ffffff"><path d="M7 12.5l3 3 7-7" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path><path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
            </div>
          </div>

          <div class="payment-result" v-if="status === 2">
            <div class="payment-result__icon">
              <svg stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#ffffff"><path d="M9.172 14.828L12.001 12m2.828-2.828L12.001 12m0 0L9.172 9.172M12.001 12l2.828 2.828M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
            </div>
            <div class="payment-result__title">{{ $t('payment_canceled') }}</div>
          </div>

          <div class="payment-result" v-if="status === 3 || status === 7 || status === 8">
            <div class="payment-result__icon">
              <svg stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#ffffff"><path d="M7 12.5l3 3 7-7" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path><path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
            </div>
            <div class="payment-result__title">{{ $t('payment_successful') }}</div>
          </div>

          <div v-if="status === 4">
            {{ $t('payment_deleted') }}
          </div>

        </div>
      </div>

      <div class="card">
        <div class="card__body">
          <div class="buttons">
            <router-link to="/" class="button button--primary">{{ $t('ok') }}</router-link>
          </div>
        </div>
      </div>

    </div>

  </main>
</template>

<style scoped>
.buttons {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
}

.payment-result {
  display: flex;
  flex-direction: column;
  gap: 20px;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.payment-result__title {
  font-size: 1.8rem;
  font-weight: 600;
}

.payment-result__icon {
  width: 120px;
  height: 120px;
}

.payment-result__icon svg {
  width: 100%;
}
</style>
