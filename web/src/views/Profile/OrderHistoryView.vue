<script>
  import HeaderSimple from '../../components/HeaderSimple.vue';
  import Card from '../../components/Card.vue';
  import BaseLayout from '../../components/layout/BaseLayout.vue';
  import orderService from '../../services/orderService';
  import moment from 'moment';
  import Filter from '../../components/Filter.vue';

  export default {
    components: {
      HeaderSimple,
      Card,
      BaseLayout,
      Filter,
    },
    setup() {
      return { moment };
    },
    data() {
      return {
        orders: [],
        detailsPopup: false,
        orderDetails: null,
        currentPage: 1,
        paginationLength: 1,
      currency: import.meta.env.VITE_CURRENCY,
      }
    },
    methods: {
      fetchOrders(customPage = null, start = null, end = null) {
        orderService.myOldOrders(customPage == null ? this.$route.query.page : null, start, end)
          .then((res) => {
            this.orders = res.orders
            this.paginationLength = res.pagination.size
          })
      },
      convertCurrency() {
        if (!this.currency) return '';

        if (this.currency == 'kr') return 'kr';

        if (this.currency == 'eur') return '€';

        if (this.currency == 'usd') return '$';

        if (this.currency == 'try') return '₺';

        return this.currency;
      },
      paginationHit() {
        this.$router.push({name: 'profile.orderHistory', query: {page: this.currentPage, end: this.$route.query.end, start: this.$route.query.start}});
      },
      seeDetails(id) {
        orderService.get(id)
          .then((res) => {
            this.orderDetails = res
          })

        this.detailsPopup = true;
      },
    },
    mounted() {
      this.fetchOrders(this.$route.query.page, this.$route.query.start, this.$route.query.end)
    },
    beforeRouteUpdate(to, from, next) {
      this.fetchOrders(to.query.page, to.query.start, to.query.end);

      next();
    }
  }
</script>

<template>
  <BaseLayout>
  <HeaderSimple :title="$t('order_history')" backUrl="/profile"/>

  <v-dialog v-model="detailsPopup">
    <v-card style="max-width: 600px; width: 100%; margin: 0 auto;">
      <v-card-text v-if="orderDetails != null">
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
            <tr cols="3" v-for="item of orderDetails.items">
              <td>{{ item.product_name }}</td>
              <td align="center">{{ item.product_amount }}</td>
                <td align="right">{{ item.product_price }}{{ convertCurrency() }}</td>
            </tr>
            <tr cols="3">
              <td></td>
              <td align="center"><b>{{ $t('total') }}</b></td>
              <td align="right">{{ orderDetails.price }}{{ convertCurrency() }}</td>
            </tr>
          </tbody>
        </v-table>
        <br/>

        <h3 style="margin-bottom: 8px;">{{ $t('discounts') }}</h3>
        <v-table>
          <tbody>
            <tr cols="2" v-for="discount of orderDetails.discounts.slice().reverse()">
              <td>{{ moment(discount.created_at).format('dd.mm.yyyy - hh:mm') }}</td>
              <td align="right">{{ discount.discount_price }}{{ convertCurrency() }}</td>
            </tr>
          </tbody>
        </v-table>
      </v-card-text> 

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn @click="detailsPopup = false">{{ $t('close') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <main class="container" style="flex: 1;">

    <br />

    <Card>
      <template v-slot:actions>
        <Filter />
      </template>
      <v-table v-if="orders != null && orders.length > 0">
        <thead cols="4">
          <tr>
            <th>{{ $t('date') }}</th>
            <th>{{ $t('reference') }}</th>
            <th>{{ $t('price') }}</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order of orders" cols="4">
            <td>{{ moment(order.created_at).format('DD.MM.YYYY - HH:mm') }}</td>
            <td>{{ order.reference }}</td>
            <td>{{ order.price }}{{ convertCurrency() }}</td>
            <td align="right">
              <button @click="() => seeDetails(order.id)" class="button button--smal button--primary">{{ $t('details') }}</button>
            </td>
          </tr>
        </tbody>
      </v-table>

      <div v-if="orders == null || orders.length == 0">
        <p>{{ $t('you_dont_have_orders') }}</p>
      </div>

      <v-pagination v-model="currentPage" :length="paginationLength"></v-pagination>
    </Card>

  </main>
  </BaseLayout>
</template>
