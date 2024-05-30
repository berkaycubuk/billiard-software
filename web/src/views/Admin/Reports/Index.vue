<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import Card from '../../../components/Card.vue';
import reportService from '../../../services/reportService';
import Filter from '../../../components/Filter.vue';
import FilterTags from '../../../components/FilterTags.vue';

export default {
  components: {
    AdminLayout,
    Filter,
    FilterTags,
    Card,
  },
  data() {
    return {
      total: null,
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  methods: {
    fetchData(start = null, end = null) {
      reportService.total(start, end)
        .then((res) => {
          this.total = res.data.total;
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
  mounted() {
    this.fetchData(this.$route.query.start, this.$route.query.end);
  },
  beforeRouteUpdate(to, from, next) {
    this.fetchData(to.query.start, to.query.end);

    next();
  }
}
</script>

<template>
  <AdminLayout>

    <Card :title="$t('total_sales')">
      <template v-slot:actions>
        <Filter />
      </template>

      <FilterTags />

      <div>{{ parseFloat(total).toFixed(2) }}{{ convertCurrency() }}</div>
    </Card>
    <br />

    <Card>
      <div class="grid grid-cols-2">
        <router-link class="button button--blue" to="/admin/reports/tables">{{ $t("tables") }}</router-link>
        <router-link class="button button--blue" to="/admin/reports/kiosk">{{ $t("kiosk") }}</router-link>
        <router-link class="button button--blue" to="/admin/reports/subscriptions">{{ $t("subscriptions") }}</router-link>
        <!--
        <router-link class="button button--blue" to="/admin/reports/games">{{ $t("games") }}</router-link>
        -->
        <router-link class="button button--blue" to="/admin/reports/users">{{ $t("users") }}</router-link>
      </div>
    </Card>

  </AdminLayout>
</template>
