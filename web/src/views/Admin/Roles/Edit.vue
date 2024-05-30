<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import moment from 'moment';
import { useRoute } from 'vue-router';
import Card from '../../../components/Card.vue';
import DeletePopup from '../../../components/DeletePopup.vue';
import roleService from '../../../services/roleService';
import RoleEditPopup from '../../../components/popups/RoleEditPopup.vue';

export default {
  setup() {
    const route = useRoute();
    return { route, moment };
  },
  components: {
    AdminLayout,
    Card,
    DeletePopup,
    RoleEditPopup
},
  data() {
    return {
      role: null,
      deletePopup: false,
      editPopup: false,
    }
  },
  mounted() {
    this.fetchRole();
  },
  methods: {
    fetchRole() {
      roleService.get(this.route.params.id)
        .then((role) => {
          this.role = role;
        });
    },
    onEdit() {
      this.fetchRole();
      this.editPopup = false;
    },
    onDelete() {
      roleService.del(parseInt(this.route.params.id))
        .then(() => {
          this.$router.push('/admin/roles');
        });
    }
  }
}
</script>

<template>
  <AdminLayout>

    <RoleEditPopup v-if="editPopup === true" :role="role" @closed="() => editPopup = false" @saved="onEdit" />
    <DeletePopup v-if="deletePopup === true" :title="$t('delete_role')" @closed="() => deletePopup = false" @deleted="onDelete" />

    <div class="cards">

      <div class="grid">

        <Card :title="role.name" v-if="role != null">

          <template v-slot:actions>
            <router-link to="/admin/roles" class="button button--small button--secondary">{{ $t('go_back') }}</router-link>
          </template>

          <div class="grid">
          </div>

        </Card>

      </div>

      <Card v-if="role != null">
        <div class="grid grid-cols-1">
          <div style="display: flex; align-items: center; gap: 6px;">
            <b>{{ $t('created_at') }}:</b>
            {{ moment(role.created_at).format('DD.MM.YYYY - HH:mm') }}
          </div>
        </div>

        <br/>
        <div class="grid grid-cols-1">
          <button class="button button--primary" @click="() => editPopup = true">{{ $t('edit') }}</button>
          <button class="button button--red" @click="() => deletePopup = true">{{ $t('delete') }}</button>
        </div>
      </Card>

    </div>
  </AdminLayout>
</template>

<style scoped>
</style>
