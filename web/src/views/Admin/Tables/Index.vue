<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import Card from '../../../components/Card.vue';
import tableService from '../../../services/tableService';
import TableNewPopup from '../../../components/popups/TableNewPopup.vue';

export default {
  components: {
    AdminLayout,
    Card,
    TableNewPopup
  },
  data() {
    return {
      tables: [],
      newPopup: false,
    }
  },
  mounted() {
    this.fetchTables();
  },
  methods: {
    fetchTables() {
      tableService.all()
        .then((res) => {
          this.tables = res.data.tables;
        });
    },
    onCreate() {
      this.fetchTables();
      this.newPopup = false;
    }
  }
}
</script>

<template>
  <AdminLayout>

    <TableNewPopup v-if="newPopup === true" @closed="() => newPopup = false" @saved="onCreate" />

    <div class="users-page">

      <Card :title="$t('tables')">

        <template v-slot:actions>
          <button @click="() => newPopup = true" class="button button--small button--primary">{{ $t('new_table') }}</button>
        </template>

        <v-table>

          <thead>
            <tr>
              <th><b>{{ $t('name') }}</b></th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="table of tables">
              <td>{{ table.name }}</td>
              <td>
                <div class="user__right">
                  <router-link :to="'/admin/tables/' + table.id" class="button button--small button--primary">{{ $t('details') }}</router-link>
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
