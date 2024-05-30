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
      total: 0,
      items: [],
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  mounted() {
    this.fetchData(this.$route.query.start, this.$route.query.end);
  },
  methods: {
    fetchData(start = null, end = null) {
      reportService.kiosk(start, end)
        .then((res) => {
          this.items = res.data.products
          if (res.data.total != null) this.total = res.data.total
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

      <Card :title="$t('kiosk_report')">

        <template v-slot:actions>
          <Filter />
        </template>

        <FilterTags />

        <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 10px;">
          <b>{{ $t('total') }}:</b> <div>{{ total }}{{ convertCurrency() }}</div>
        </div>

        <v-table>
          <thead>
            <tr cols="3">
              <th><b>{{ $t('product') }}</b></th>
              <th><b>{{ $t('count') }}</b></th>
              <th><b>{{ $t('total') }}</b></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item of items" cols="3">
              <td>{{ item.name }}</td>
              <td>{{ item.count }}</td>
              <td>{{ parseFloat(item.total).toFixed(2) }}{{ convertCurrency() }}</td>
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
