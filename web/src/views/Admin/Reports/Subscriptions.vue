<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import moment from 'moment';
import Card from '../../../components/Card.vue';
import reportService from '../../../services/reportService';
import Filter from '../../../components/Filter.vue';
import FilterTags from '../../../components/FilterTags.vue';

export default {
  setup() {
    return { moment };
  },
  components: {
    AdminLayout,
    Card,
    Filter,
    FilterTags,
  },
  data() {
    return {
      items: [],
      users: [],
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  mounted() {
    this.fetchData(this.$route.query.start, this.$route.query.end);
  },
  methods: {
    fetchData(start = null, end = null) {
      reportService.subscriptions(start, end)
        .then((res) => {
          this.items = res.data.products;
          this.users = res.data.users;
        });
    },
      convertCurrency() {
        if (!this.currency) return '';

        if (this.currency == 'kr') return 'kr';

        if (this.currency == 'eur') return '€';

        if (this.currency == 'usd') return '$';

        if (this.currency == 'try') return '₺';

        return this.currency;
      }
  },
  beforeRouteUpdate(to, from, next) {
    this.fetchData(to.query.start, to.query.end);

    next();
  }
}
</script>

<template>
  <AdminLayout>

    <div class="users-page">

      <Card :title="$t('subscriptions_report')">

        <template v-slot:actions>
          <Filter />
        </template>

        <FilterTags />

        <v-table>
          <thead>
            <tr cols="3">
              <th><b>{{ $t('subscription') }}</b></th>
              <th><b>{{ $t('count') }}</b></th>
              <th><b>{{ $t('total') }}</b></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item of items" cols="3">
              <td>{{ item.name }}</td>
              <td>{{ item.count }}</td>
              <td>{{ item.total }}{{ convertCurrency() }}</td>
            </tr>
          </tbody>
        </v-table>
        <br />

        <v-table>
          <thead>
            <tr cols="3">
              <th><b>{{ $t('user') }}</b></th>
              <th><b>{{ $t('subscription') }}</b></th>
              <th><b>{{ $t('started_at') }}</b></th>
              <th><b>{{ $t('ending_at') }}</b></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user of users" cols="3">
              <td>{{ user.user.name }} {{ user.user.surname }}</td>
              <td>{{ user.name }}</td>
              <td>{{ moment(user.created_at).format('DD.MM.YYYY - HH:mm') }}</td>
              <td>{{ moment(user.ending_at).format('DD.MM.YYYY - HH:mm') }}</td>
            </tr>
          </tbody>
        </v-table>

      </Card>

    </div>
  </AdminLayout>
</template>

<style scoped>
.user__right {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
}
</style>
