<script>
import orderService from '../services/orderService';
import paymentService from '../services/paymentService';
import { useRoute } from 'vue-router';
import { UploadsURL } from '../utils';
import { store } from '../store';

export default {
setup() {
    const route = useRoute();

    return { route };
},
data() {
  return {
    order: null,
    currency: import.meta.env.VITE_CURRENCY,
    UploadsURL: UploadsURL(),
  }
},
methods: {
      sendToCardReader() {
        orderService.get(this.$route.params.id)
          .then((res) => {
            const json_message = JSON.stringify(res)
            Payment.postMessage(json_message)
          })
      },
  proceedToPayment() {
    paymentService.init(parseInt(this.route.params.id), 2)
      .then((res) => {
        if (res.data.success == false) {
          console.error(res.data.message);
          return;
        }

        this.$router.push(`/payment/vipps/${this.route.params.id}?frontendUrl=${res.data.checkout_frontend_url}&token=${res.data.token}`);
      });
  },
  proceedToPaymentPhysical() {
    this.$router.push(`/payment/physical/${this.route.params.id}`);
  },
  cancelOrder() {
    orderService.cancel(parseInt(this.route.params.id))
      .then(() => {
        this.$router.push('/');
      });
  },
  canSeePhysicalPayment() {
    if (store.user) {
      return [3].includes(store.user.role_id)
    }

    return false
  },
},
mounted() {
 paymentService.overview(this.route.params.id)
  .then((res) => {
    this.order = res.data.order;

    if (this.order.price === "0") {
      this.$router.push(`/payment/result/${this.route.params.id}`);
    } else if (this.order.status === 3 || this.order.status === 7) {
      this.$router.push(`/payment/result/${this.route.params.id}`);
    }
  })
  .catch((err) => {
    this.$router.push("/404");
  });

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

    <div class="cards">

      <div class="card">
        <div class="card__header">
          <div class="card__title">{{ $t('products') }}</div>
        </div>
        <div class="card__body">

          <div class="products" v-if="order != null && order.items.length">

            <div class="product" v-for="item of order.items">

              <img v-if="false" class="product__image" src="" />

              <div class="product__info">
                <div class="product__title">{{ item.product_name }}</div>
                <div class="product__price">{{ item.product_price }}{{ currency }}</div>
                <div class="product__count">x{{ item.product_amount }}</div>
              </div>

            </div>
          
          </div>

        </div>
      </div>

      <div class="card">
        <div class="card__header">
          <div class="card__title">{{ $t('overview') }}</div>
        </div>
        <div class="card__body">
          <div class="payment-total" v-if="order != null">Total: {{ order.price }}{{ currency }}</div>

          <p>{{ $t('total_value_can_change') }}</p>

          <div class="payment-buttons">
            <button class="button button--primary" @click="proceedToPayment">{{ $t('pay_with_vipps') }}</button>
            <button class="button button--primary" @click="sendToCardReader" v-if="canSeePhysicalPayment()">{{ $t('pay_with_card') }}</button>
            <button class="button button--red" @click="() => $router.push('/')">{{ $t('go_back') }}</button>
          </div>
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
