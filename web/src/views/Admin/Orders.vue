<script>
import AdminLayout from '../../components/layout/AdminLayout.vue';
import moment from 'moment';
import Card from '../../components/Card.vue';
import orderService from '../../services/orderService';
import Pagination from '../../components/Pagination.vue';
import Filter from '../../components/Filter.vue';
import FilterTags from '../../components/FilterTags.vue';

export default {
  setup() {
    return { moment, orderService };
  },
  components: {
    AdminLayout,
    Card,
    Pagination,
    Filter,
    FilterTags,
  },
  data() {
    return {
      orders: [],
      pagination: null,
      currentPage: 1,
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  methods: {
    fetchOrders(customPage = null, start = null, end = null, method = null, zeroOrders = false) {
      orderService.all(customPage == null ? this.$route.query.page : customPage, start, end, method, zeroOrders)
        .then((res) => {
          this.orders = res.orders;
          this.pagination = res.pagination;
        });
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
      this.$router.push({name: 'admin.orders', query: {page: this.currentPage, end: this.$route.query.end, start: this.$route.query.start, method: this.$route.query.method, zeroOrders: this.$route.query.zeroOrders}});
    }
  },
  mounted() {
    this.fetchOrders(this.$route.query.page, this.$route.query.start, this.$route.query.end, this.$route.query.method, this.$route.query.zeroOrders);
  },
  beforeRouteUpdate(to, from, next) {
    this.fetchOrders(to.query.page, to.query.start, to.query.end, to.query.method, to.query.zeroOrders);

    next();
  }
}
</script>

<template>
  <AdminLayout>
    <div class="">

      <Card :title="$t('orders')">
        <template v-slot:actions>
          <Filter :method="true" :zeroOrders="true" />
        </template>

        <FilterTags />

        <v-table>
          <tbody>
            <tr v-for="order of orders" cols="6" @click="() => $router.push('/admin/orders/' + order.id)" class="hoverable-row">
              <td>{{ order.reference }}</td>
              <td>
                <div style="display: flex; align-items: center; gap: 6px;">
                  <div :class="'order-status order-status--' + order.status"></div>
                  {{ $t(orderService.statusString(order.status)) }}
                </div>
              </td>
              <td>{{ order.detail.user_name }}</td>
              <td>{{ parseFloat(order.price).toFixed(2) }}{{ convertCurrency() }}</td>
              <td>{{ moment(order.created_at).format('DD.MM.YYYY - HH:mm') }}</td>
              <!--
              <td>
                <div class="user__right">
                  <router-link :to="'/admin/orders/' + order.id" class="button button--small button--primary">Details</router-link>
                </div>
              </td>
              -->
            </tr>
          </tbody>
        </v-table>

        <v-pagination v-if="pagination != null" v-model="currentPage" :length="pagination.size" @update:model-value="paginationHit"></v-pagination>

        <!--
        <Pagination v-if="pagination != null && pagination.size != 1" routeName="admin.orders" :pagination="pagination" />
        -->

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

.hoverable-row:hover {
  cursor: pointer;
  background-color: #111;
}
</style>
