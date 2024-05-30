<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import moment from 'moment';
import Card from '../../../components/Card.vue';
import reportService from '../../../services/reportService';
import Filter from '../../../components/Filter.vue';
import UserReportPopup from '../../../components/popups/UserReportPopup.vue';
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
    UserReportPopup,
},
  data() {
    return {
      users: [],
      currentPage: 1,
      pagination: null,
      popupOpen: false,
      selected: null,
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  mounted() {
    this.fetchData(this.$route.query.start, this.$route.query.end, this.$route.query.user);
  },
  methods: {
    fetchData(start = null, end = null, user = null, page = null) {
      reportService.users(start, end, user, page == null ? this.$route.query.page : page)
        .then((res) => {
          this.users = res.data.data.users;
          this.pagination = res.data.data.pagination;
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
    openPopup(selected) {
      this.popupOpen = true;
      this.selected = selected;
    },
    closePopup() {
      this.popupOpen = false;
      this.selected = null;
    },
    paginationHit() {
      this.$router.push({name: 'admin.reports.users', query: {page: this.currentPage, end: this.$route.query.end, start: this.$route.query.start, user: this.$route.query.user}});
    },
  },
  beforeRouteUpdate(to, from, next) {
    this.fetchData(to.query.start, to.query.end, to.query.user, to.query.page);

    next();
  }
}
</script>

<template>
  <AdminLayout>

    <UserReportPopup v-if="popupOpen === true" @closed="closePopup" :id="selected" />

    <div class="users-page">

      <Card :title="$t('users_report')">

        <template v-slot:actions>
          <Filter :user="true" />
        </template>

        <FilterTags />

        <v-table>
          <thead>
            <tr cols="3">
              <th><b>{{ $t('user') }}</b></th>
              <th><b>{{ $t('order_count') }}</b></th>
              <th><b>{{ $t('total') }}</b></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user of users" cols="3" @click="() => openPopup(user.user.id)" style="cursor: pointer;">
              <td>{{ user.user.name }} {{ user.user.surname }}</td>
              <td>{{ user.count }}</td>
              <td>{{ user.total }}{{ convertCurrency() }}</td>
            </tr>
          </tbody>
        </v-table>

        <v-pagination v-if="pagination != null" v-model="currentPage" :length="pagination.size" @update:model-value="paginationHit"></v-pagination>

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
