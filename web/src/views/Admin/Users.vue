<script>
import AdminLayout from '../../components/layout/AdminLayout.vue';
import userService from '../../services/userService';
import moment from 'moment';
import Card from '../../components/Card.vue';
import UserNewPopup from '../../components/popups/UserNewPopup.vue';
import roleService from '../../services/roleService';
import Filter from '../../components/Filter.vue';
import FilterTags from '../../components/FilterTags.vue';

export default {
  setup() {
    return { moment };
  },
  components: {
    AdminLayout,
    Card,
    UserNewPopup,
    Filter,
    FilterTags,
  },
  data() {
    return {
      users: [],
      currentPage: 1,
      newPopup: false,
      pagination: null,
    }
  },
  mounted() {
    this.fetchUsers(this.$route.query.user);
    this.fetchRoles();
  },
  beforeRouteUpdate(to, from, next) {
    this.fetchUsers(to.query.user, to.query.page);

    next();
  },
  methods: {
    fetchUsers(user = null, customPage = null) {
      userService.all(user, customPage == null ? this.$route.query.page : customPage)
        .then((res) => {
          this.users = res.users;
          this.pagination = res.pagination;
        }); 
    },
    paginationHit() {
      this.$router.push({name: 'admin.users', query: {page: this.currentPage}});
    },
    fetchRoles() {
      roleService.all()
        .then((roles) => {
          this.roles = roles;
        });
    },
    onCreate() {
      this.fetchUsers();
      this.newPopup = false;
    }
  }
}
</script>

<template>
  <AdminLayout>

    <UserNewPopup :roles="roles" v-if="newPopup === true" @closed="() => newPopup = false" @saved="onCreate" />

    <div class="users-page">

      <Card :title="$t('users')">

        <template v-slot:actions>
          <Filter :user="true" :dates="false" /> &nbsp;
          <button @click="() => newPopup = true" class="button button--small button--primary">{{ $t('new_user') }}</button>
        </template>

        <FilterTags />

        <v-table>
          <tbody>
            <tr v-for="user of users" cols="3" @click="() => $router.push('/admin/users/' + user.id)" style="cursor: pointer;">
              <td>{{ user.name }} {{ user.surname }}</td>
              <td>{{ moment(user.created_at).format('DD.MM.YYYY') }}</td>
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
