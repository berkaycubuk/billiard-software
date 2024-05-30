<template>
  <BaseLayout>
  <HeaderSimple :title="$t('waiting_orders')" backUrl="/profile"/>

  <v-dialog v-model="detailsPopup" persistent>
    <v-card style="max-width: 600px; width: 100%; margin: 0 auto;" :title="$t('order_details')">
      <v-card-text>
        <div v-if="selectedOrder != null">
        <div v-for="item of selectedOrder.items">
          <div style="margin-bottom: 4px;">{{ item.product_name }} ({{ item.product_price }} kr) x {{ item.product_amount }}</div>
          <div style="margin-bottom: 4px; font-weight: bold;">{{ item.product_price * item.product_amount }}kr</div>
          <hr />
        </div>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-btn variant="flat" color="blue" v-if="selectedOrder != null" @click="$router.push(`/payment/overview/${selectedOrder.id}`)">{{ $t('pay') }}</v-btn>
        <v-btn @click="() => closeDetails()">{{ $t('close') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <main class="container" style="flex: 1;">

    <br />

    <Card>
      <v-table v-if="orders != null && orders.length > 0">
        <thead cols="4">
          <tr>
            <th>{{ $t('created_at') }}</th>
            <th>{{ $t('price') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order of orders" cols="4" @click="() => openDetails(order)">
            <td>{{ moment(order.created_at).format('DD.MM.YYYY - hh:mm') }}</td>
            <td>{{ order.price }}kr</td>
            <!--
            <td align="right">
              <router-link :to="`/payment/overview/${order.id}`" class="button button--smal button--primary">{{ $t('pay') }}</router-link>
            </td>
            -->
          </tr>
        </tbody>
      </v-table>

      <div v-if="orders == null || orders.length == 0">
        <p>{{ $t('you_dont_have_waiting_orders') }}</p>
      </div>
    </Card>

  </main>
  </BaseLayout>
</template>

<script>
import { useCookies } from 'vue3-cookies';
import orderService from '../services/orderService';
import HeaderSimple from '../components/HeaderSimple.vue';
import Card from '../components/Card.vue';
import BaseLayout from '../components/layout/BaseLayout.vue';
import moment from 'moment';

export default {
  setup() {
    const { cookies } = useCookies();
    return { cookies, moment };
  },
  components: {
    HeaderSimple,
    Card,
    BaseLayout,
  },
  data() {
    return {
      detailsPopup: false,
      selectedOrder: null,
      orders: [],
    }
  },
  mounted() {
    this.fetchOrders();
  },
  methods: {
    fetchOrders() {
      orderService.myWaitingOrders()
        .then((res) => {
          this.orders = res;
        });
    },
    openDetails(order) {
      this.selectedOrder = order
      this.detailsPopup = true
    },
    closeDetails() {
      this.detailsPopup = false
      this.selectedOrder = null
    },
    getStatusString(status) {
      if (status === 1) {
        return "Waiting";
      }
      return "Added to account";
    }
  }
}
</script>

<style scoped>
.profile-header {
  background-color: #191919;
}

.profile-header__inner {
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.profile-header__left {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 10px;
}

.profile-header__go-back {
  display: inline-block;
  width: fit-content;
}

.profile-header__name {
  font-family: 'Bricolage Grotesque', sans-serif;
  font-size: 1.8rem;
  font-weight: 700;
}

.profile-section {
  margin: 20px 0;
}

.profile-section__title {
  font-size: 1.2rem;
  font-weight: 700;
  margin-bottom: 10px;
  color: #8A8B85;
}

.profile-section nav {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.profile-section nav a {
  padding: 20px 16px;
  display: flex;
  align-items: center;
  gap: 20px;
  color: #f8faed;
  text-decoration: none;
  font-size: 1.2rem;
  font-weight: 600;
  border-radius: 12px;
  background-color: #191919;
  transition: all .2s;
}

.profile-section nav a:hover {
  opacity: .8;
}

.profile-section nav button {
  cursor: pointer;
  padding: 20px 16px;
  display: flex;
  align-items: center;
  gap: 20px;
  color: #f8faed;
  text-decoration: none;
  font-size: 1.2rem;
  font-weight: 600;
  border-radius: 12px;
  background-color: #191919;
  transition: all .2s;
  outline: none;
  border: none;
}

.profile-section nav button:hover {
  opacity: .8;
}
</style>
