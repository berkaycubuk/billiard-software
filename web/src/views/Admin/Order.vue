<script>
import AdminLayout from '../../components/layout/AdminLayout.vue';
import moment from 'moment';
import Card from '../../components/Card.vue';
import orderService from '../../services/orderService';
import DeletePopup from '../../components/DeletePopup.vue';
import OrderApprovePopup from '../../components/OrderApprovePopup.vue';
import ApplyDiscountPopup from '../../components/popups/ApplyDiscountPopup.vue';

export default {
  setup() {
    return { moment, orderService };
  },
  components: {
    AdminLayout,
    Card,
    DeletePopup,
    ApplyDiscountPopup,
    OrderApprovePopup,
  },
  data() {
    return {
      order: null,
      deletePopup: false,
      currency: import.meta.env.VITE_CURRENCY,
      approvePopup: false,
      discountPopup: false,
    }
  },
  mounted() {
    this.fetchOrder();
  },
  methods: {
    captureOrder() {
      orderService.capture(this.order.reference)
        .then((res) => {
          console.log(res)
        })
        .catch((err) => {
          console.log(err)
        });
    },
    fetchOrder() {
      orderService.get(this.$route.params.id)
        .then((res) => {
          this.order = res;
        });
    },
    onDelete() {
      if (this.order === null) {
        return;
      }

      orderService.del(this.order.id)
        .then(() => {
          this.$router.push('/admin/orders');
        });
    },
    onDiscountApply() {
      this.fetchOrder();
      this.discountPopup = false;
    },
    onApprove() {
      if (this.order === null) {
        return;
      }

      orderService.approve(this.order.id)
        .then(() => {
          this.approvePopup = false;
          this.fetchOrder();
        });
    }
  },
}
</script>

<template>
  <AdminLayout>

  <!--
  <v-btn color="orange" @click="() => captureOrder()">Capture</v-btn>
  -->

    <DeletePopup v-if="deletePopup === true" :title="$t('delete_order')" @closed="() => deletePopup = false" @deleted="onDelete" />
    <OrderApprovePopup v-if="approvePopup === true" @closed="() => approvePopup = false" @approved="onApprove" />
    <ApplyDiscountPopup v-if="discountPopup === true && order != null" :id="order.id" @closed="() => discountPopup = false" @saved="onDiscountApply" />

    <div class="cards">

      <Card :title="order.reference" v-if="order != null">
        <template v-slot:actions>
          <router-link to="/admin/orders" class="button button--small button--secondary">{{ $t('go_back') }}</router-link>
        </template>

        <h3 style="margin-bottom: 8px;">{{ $t('user') }}</h3>
        <div class="grid">
          <div><b>{{ $t('full_name') }}:</b> {{ order.detail.user_name }} {{ order.detail.user_surname }}</div>
          <div><b>{{ $t('email') }}:</b> {{ order.detail.user_email }}</div>
          <div><b>{{ $t('phone') }}:</b> {{ order.detail.user_phone }}</div>
        </div>

        <br/>

        <h3 style="margin-bottom: 8px;">{{ $t('items') }}</h3>
        <v-table>
          <thead>
            <tr cols="3">
              <th><b>{{ $t('product') }}</b></th>
              <th align="center"><b>{{ $t('count') }}</b></th>
              <th align="right"><b>{{ $t('price') }}</b></th>
            </tr>
          </thead>
          <tbody>
            <tr cols="3" v-for="item of order.items">
              <td>{{ item.product_name }}</td>
              <td align="center">{{ item.product_amount }}</td>
              <td align="right">{{ item.product_price }}{{ currency }}</td>
            </tr>
            <tr cols="3">
              <td></td>
              <td align="center"><b>{{ $t('total') }}</b></td>
              <td align="right">{{ order.price }}{{ currency }}</td>
            </tr>
          </tbody>
        </v-table>
        <br />

        <h3 style="margin-bottom: 8px;">{{ $t('history') }}</h3>
        <v-table>
          <thead>
            <tr cols="3">
              <th><b>{{ $t('time') }}</b></th>
              <th align="center"><b>{{ $t('action') }}</b></th>
              <th align="right"><b>{{ $t('method') }}</b></th>
            </tr>
          </thead>
          <tbody>
            <tr cols="3" v-for="history of order.histories.slice().reverse()">
              <td>{{ moment(history.created_at).format('DD.MM.YYYY - HH:mm') }}</td>
              <td align="center">{{ $t(orderService.actionString(history.action)) }}</td>
              <td align="right">{{ $t(orderService.methodString(history.method)) }}</td>
            </tr>
          </tbody>
        </v-table>
        <br/>

        <h3 style="margin-bottom: 8px;">{{ $t('discounts') }}</h3>
        <v-table>
          <tbody>
            <tr cols="2" v-for="discount of order.discounts.slice().reverse()">
              <td>{{ moment(discount.created_at).format('dd.mm.yyyy - hh:mm') }}</td>
              <td align="right">{{ discount.discount_price }}kr</td>
            </tr>
          </tbody>
        </v-table>

      </Card>

      <Card v-if="order != null">
        <div class="grid grid-cols-1">
          <div style="display: flex; align-items: center; gap: 6px;">
            <b>{{ $t('status') }}:</b>
            <div :class="'order-status order-status--' + order.status"></div>
            {{ $t(orderService.statusString(order.status)) }}
          </div>

          <div style="display: flex; align-items: center; gap: 6px;">
            <b>{{ $t('created_at') }}:</b>
            {{ moment(order.created_at).format('DD.MM.YYYY - HH:mm') }}
          </div>
        </div>

        <br/>
        <div class="grid grid-cols-1">
          <button class="button button--primary" @click="() => approvePopup = true" v-if="order.status === 1">{{ $t('approve') }}</button>
          <button class="button button--primary" @click="() => $router.push('/payment/overview/' + order.id)" v-if="order.status === 1">{{ $t('pay_as_admin') }}</button>
          <button class="button button--blue" @click="() => discountPopup = true" v-if="order.status === 1">{{ $t('apply_discount') }}</button>
          <button class="button button--red" @click="() => deletePopup = true">{{ $t('delete') }}</button>
        </div>
      </Card>

    </div>
  </AdminLayout>
</template>

<style scoped>
.users {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.user {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  align-items: center;
  padding: 10px;
  background-color: #232323;
  border-radius: 8px;
}

.user__right {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
}
</style>
