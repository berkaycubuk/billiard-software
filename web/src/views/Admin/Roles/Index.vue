<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import moment from 'moment';
import Card from '../../../components/Card.vue';
import roleService from '../../../services/roleService';
import RoleNewPopup from '../../../components/popups/RoleNewPopup.vue';

export default {
  setup() {
    return { moment };
  },
  components: {
    AdminLayout,
    Card,
    RoleNewPopup
},
  data() {
    return {
      roles: [],
      newPopup: false,
    }
  },
  mounted() {
    this.fetchRoles();
  },
  methods: {
    fetchRoles() {
      roleService.all()
        .then((roles) => {
          this.roles = roles;
        });
    },
    onCreate() {
      this.fetchRoles();
      this.newPopup = false;
    }
  }
}
</script>

<template>
  <AdminLayout>

    <RoleNewPopup v-if="newPopup === true" @closed="() => newPopup = false" @saved="onCreate" />

    <div class="users-page">

      <Card :title="$t('roles')">

        <template v-slot:actions>
          <button @click="() => newPopup = true" class="button button--small button--primary">{{ $t('new_role') }}</button>
        </template>

        <v-table>
          <thead>
            <tr cols="3">
              <th><b>{{ $t('name') }}</b></th>
              <th><b>{{ $t('created_at') }}</b></th>
              <th>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="role of roles" cols="3">
              <td>{{ role.name }}</td>
              <td>{{ moment(role.created_at).format('DD.MM.YYYY') }}</td>
              <td>
                <div class="user__right">
                  <router-link :to="'/admin/roles/' + role.id" class="button button--small button--primary">{{ $t('details') }}</router-link>
                </div>
              </td>
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
