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
      tables: [],
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  mounted() {
    this.fetchData(this.$route.query.start, this.$route.query.end);
  },
  methods: {
    fetchData(start = null, end = null) {
      reportService.tables(start, end)
        .then((res) => {
          this.tables = res.data.sales;
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

      <Card :title="$t('tables_report')">

        <template v-slot:actions>
          <Filter />
        </template>

        <FilterTags />

        <v-table>
          <thead>
            <tr cols="3">
              <th><b>{{ $t('table') }}</b></th>
              <th><b>{{ $t('games') }}</b></th>
              <th><b>{{ $t('sales') }}</b></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="table of tables" cols="3">
              <td>{{ table.name }}</td>
              <td>{{ table.games }}</td>
              <td>{{ parseFloat(table.sales).toFixed(2) }}{{ convertCurrency() }}</td>
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
