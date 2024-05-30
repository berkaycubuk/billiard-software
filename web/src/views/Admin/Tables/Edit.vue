<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import tableService from '../../../services/tableService';
import moment from 'moment';
import { useRoute } from 'vue-router';
import Card from '../../../components/Card.vue';
import DeletePopup from '../../../components/DeletePopup.vue';
import TableEditPopup from '../../../components/popups/TableEditPopup.vue';

export default {
  setup() {
    const route = useRoute();
    return { route, moment };
  },
  components: {
    AdminLayout,
    Card,
    DeletePopup,
    TableEditPopup
},
  data() {
    return {
      table: null,
      deletePopup: false,
      editPopup: false,
    }
  },
  mounted() {
    this.fetchTable();
  },
  methods: {
    fetchTable() {
      tableService.get(this.route.params.id)
        .then((res) => {
          this.table = res.data.table;
        });
    },
    onEdit() {
      this.fetchTable();
      this.editPopup = false;
    },
    onDelete() {
      tableService.del(parseInt(this.route.params.id))
        .then(() => {
          this.$router.push('/admin/tables');
        });
    }
  }
}
</script>

<template>
  <AdminLayout>

    <TableEditPopup v-if="editPopup === true" :table="table" @closed="() => editPopup = false" @saved="onEdit" />
    <DeletePopup v-if="deletePopup === true" :title="$t('delete_table')" @closed="() => deletePopup = false" @deleted="onDelete" />

    <div class="cards">

      <div class="grid">

        <Card :title="table.name" v-if="table != null">

          <template v-slot:actions>
            <router-link to="/admin/tables" class="button button--small button--secondary">{{ $t('go_back') }}</router-link>
          </template>

          <div class="grid">
          </div>

        </Card>

      </div>

      <Card v-if="table != null">
        <div class="grid grid-cols-1">
          <div style="display: flex; align-items: center; gap: 6px;">
            <b>{{ $t('created_at') }}:</b>
            {{ moment(table.created_at).format('DD.MM.YYYY - HH:mm') }}
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
