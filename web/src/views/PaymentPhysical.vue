<script>
import paymentService from '../services/paymentService';
import orderService from  '../services/orderService'
import { store } from '../store'

export default {
    data() {
      return {
        paymentResult: null,
      }
    },
    methods: {
      canSeeAcceptPayment() {
        if (store.user && store.user.role_id && store.user.role_id == 1) {
          return true
        }

        return false
      },
        acceptPayment() {
            paymentService.accept(parseInt(this.$route.params.id))
                .then(() => {
                    this.$router.push(`/payment/result/${this.$route.params.id}`);
                });
        },
        sendToCardReader() {
          orderService.get(this.$route.params.id)
            .then((res) => {
              const json_message = JSON.stringify(res)
              Payment.postMessage(json_message)
            })
        },
    },
    mounted() {
      window.acceptSumupPayment = function() {
        const url_items = window.location.pathname.split("/")
        const order_id = url_items[url_items.length - 1]
        paymentService.accept(parseInt(order_id))
          .then(() => {
            window.location.href = `/payment/result/${order_id}`
          });
      }
    },
}
</script>

<template>
  <main class="container" style="flex: 1;">

      <div class="card" style="margin-top: 20px;">
        <div class="card__header">
          <div class="card__title">{{ $t('physical_payment') }}</div>
        </div>
        <div class="card__body">
          <div class="payment-buttons">
            <button v-if="canSeeAcceptPayment()" class="button button--primary" @click="acceptPayment">{{ $t('accept_payment') }}</button>
            <button class="button button--primary" @click="sendToCardReader">Send to card reader</button>
            <router-link class="button button--secondary" :to="`/payment/overview/${$route.params.id}`">{{ $t('back') }}</router-link>
          </div>
        </div>
      </div>

  </main>
</template>

<style scoped>
.overview-page {
  margin: 10px 0;
}

.products {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
}

.product {
  display: flex;
  gap: 10px;
}

.product__image {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 8px;
}

.product__info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.product__title {
  font-size: 1.2rem;
  font-weight: 600;
}

.product__price {
  font-size: 1.2rem;
}

.payment-total {
  font-size: 1.4rem;
}

.payment-buttons {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>
