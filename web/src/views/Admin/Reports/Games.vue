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
      games: [],
      currentPage: 1,
      pagination: null,
    }
  },
  mounted() {
    this.fetchData(this.$route.query.start, this.$route.query.end);
  },
  methods: {
    fetchData(start = null, end = null, customPage = null) {
      reportService.games(start, end, customPage == null ? this.$route.query.page : customPage)
        .then((res) => {
          this.games = res.data.data.games;
          this.pagination = res.data.data.pagination;
        });
    },
    paginationHit() {
      this.$router.push({name: 'admin.reports.games', query: {page: this.currentPage, end: this.$route.query.end, start: this.$route.query.start}});
    },
  },
  beforeRouteUpdate(to, from, next) {
    this.fetchData(to.query.start, to.query.end, to.query.page);

    next();
  }
}
</script>

<template>
  <AdminLayout>

    <div class="users-page">

      <Card :title="$t('games_report')">

        <template v-slot:actions>
          <Filter />
        </template>

        <FilterTags />

        <v-table>
          <thead>
            <tr cols="3">
              <th><b>{{ $t('game') }}</b></th>
              <th><b>{{ $t('started_at') }}</b></th>
              <th><b>{{ $t('ended_at') }}</b></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="game of games" cols="3">
              <td>Game #{{ game.id }}</td>
              <td>{{ moment(game.started_at).format('DD.MM.YYYY - HH:mm') }}</td>
              <td>{{ game.ended_at != null ? moment(game.ended_at).format('DD.MM.YYYY - HH:mm') : "not ended" }}</td>
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
