<script>
import AdminLayout from '../../components/layout/AdminLayout.vue';
import userService from '../../services/userService';
import orderService from '../../services/orderService';
import moment from 'moment';
import { useRoute } from 'vue-router';
import Card from '../../components/Card.vue';
import UserEditPopup from '../../components/popups/UserEditPopup.vue';
import DeletePopup from '../../components/DeletePopup.vue';
import roleService from '../../services/roleService';
import subscriptionService from '../../services/subscriptionService'

export default {
  setup() {
    const route = useRoute();
    return { route, moment };
  },
  components: {
    AdminLayout,
    Card,
    UserEditPopup,
    DeletePopup,
  },
  data() {
    return {
      user: null,
      orders: [],
      roles: [],
      subscriptions: [],
      deletePopup: false,
      editPopup: false,
      pausePopup: false,
      pausePopupValue: null,
      subDeletePopup: false,
      subDeleteId: null,
      addSubPopup: false,
      subPopupValue: null,
    }
  },
  mounted() {
    this.fetchUser()
    this.fetchRoles()
    this.fetchOrders()
    this.fetchSubscriptions()
  },
  methods: {
    fetchUser() {
      userService.get(this.route.params.id)
        .then((user) => {
          this.user = user;
        });
    },
    fetchSubscriptions() {
      subscriptionService.all()
        .then((res) => {
          this.subscriptions = res
        })
    },
    fetchOrders() {
      orderService.userOrders(this.route.params.id)
        .then((data) => {
          this.orders = data.orders;
        });
    },
    fetchRoles() {
      roleService.all()
        .then((roles) => {
          this.roles = roles;
        });
    },
    onEdit() {
      this.fetchUser();
      this.editPopup = false;
    },
    onDelete() {
      userService.del(this.route.params.id)
        .then(() => {
          this.$router.push('/admin/users');
        });
    },
    pauseSub() {
      if (this.pausePopupValue == null) return

      subscriptionService.pauseUserSub(this.user.id, this.pausePopupValue)
        .then(() => {
          this.pausePopup = false
          this.fetchUser()
        })
    },
    addSub() {
      if (this.subPopupValue == null) return

      subscriptionService.addUserSub(this.user.id, this.subPopupValue)
        .then(() => {
          this.addSubPopup = false
          this.fetchUser()
        })
    },
    deleteSub() {
      subscriptionService.delUserSub(this.subDeleteId)
        .then(() => {
          this.fetchUser()
          this.subDeletePopup = false
        })
    },
  }
}
</script>

<template>
  <AdminLayout>

    <UserEditPopup v-if="editPopup === true" :roles="roles" :user="user" @closed="() => editPopup = false" @saved="onEdit" />
    <DeletePopup v-if="deletePopup === true" :title="$t('delete_user')" @closed="() => deletePopup = false" @deleted="onDelete" />

    <v-dialog v-model="addSubPopup">
      <v-card style="max-width: 600px; width: 100%; margin: 0 auto;" title="Select subscription">
        <v-card-text>
          <v-select
            :label="$t('subscription')"
            variant="outlined"
            :items="subscriptions"
            item-title="name"
            item-value="id"
            hide-details
            v-model="subPopupValue"
          ></v-select>
        </v-card-text>
        <v-card-actions>
          <v-btn variant="flat" color="blue" @click="addSub">{{$t('add_subscription') }}</v-btn>
          <v-btn @click="addSubPopup = false">{{ $t('cancel') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="pausePopup">
      <v-card style="max-width: 600px; width: 100%; margin: 0 auto;">
        <v-card-text>
          {{ $t('subscription_pause_popup_text') }} 
        </v-card-text>
        <v-card-actions>
          <v-btn variant="flat" color="blue" @click="pauseSub">{{ $t('pause_unpause_subscription') }}</v-btn>
          <v-btn @click="pausePopup = false">{{ $t('cancel') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="subDeletePopup">
      <v-card style="max-width: 600px; width: 100%; margin: 0 auto;">
        <v-card-text>
          {{ $t('delete_subscription_popup_text') }}
        </v-card-text>
        <v-card-actions>
          <v-btn variant="flat" color="red" @click="deleteSub">{{ $t('delete_subscription') }}</v-btn>
          <v-btn @click="subDeletePopup = false">{{ $t('cancel') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <div class="cards">

      <div class="grid">

        <Card :title="$t('user')" v-if="user != null">

          <template v-slot:actions>
            <router-link to="/admin/users" class="button button--small button--secondary">{{ $t('go_back') }}</router-link>
          </template>

          <div class="grid">
            <div><b>{{ $t('name') }}:</b> {{ user.name }}</div>
            <div><b>{{ $t('surname') }}:</b> {{ user.surname }}</div>
            <div><b>{{ $t('email') }}:</b> {{ user.email }}</div>
            <div><b>{{ $t('phone') }}:</b> {{ user.phone }}</div>
            <div><b>{{ $t('role') }}:</b> {{ user.role.role.name }}</div>
          </div>

        </Card>

        <Card :title="$t('subscription')" v-if="user != null">
        <p v-if="user.subscription === null">{{ $t('user_not_have_subscription') }}</p><br v-if="user.subscription === null" />
        <button v-if="user.subscription === null" class="button button--small button--blue" @click="addSubPopup = true">{{ $t('add_subscription') }}</button>

          <div class="grid" v-if="user.subscription != null">

            <div class="grid">
              <div><b>{{ $t('name') }}:</b> {{ user.subscription.name }}</div>
              <div><b>{{ $t('started_at') }}:</b> {{ moment(user.subscription.created_at).format('DD.MM.YYYY - HH:mm') }}</div>
              <div><b>{{ $t('ending_at') }}:</b> {{ moment(user.subscription.ending_at).format('DD.MM.YYYY - HH:mm') }}</div>
              <div><b>{{ $t('status') }}:</b> {{ user.subscription.status == 1 ? 'Active' : (user.subscription.status == 2 ? $t('paused') : null) }}</div>
            </div>

            <!-- options -->
            <div>
              <hr />
              <div style="display: flex; align-items: center; gap: 10px; margin-top: 10px;">
                <button class="button button--small button--red" @click="() => {subDeleteId = user.subscription.id; subDeletePopup = true;}">Delete Subscription</button>
                <button class="button button--small button--blue" @click="() => {pausePopupValue = user.subscription.id; pausePopup = true;}">{{ user.subscription.status == 1 ? $t('pause') : $t('unpause') }} {{ $t('subscription') }}</button>
              </div>
            </div>

          </div>

        </Card>

        <Card :title="$t('orders')" v-if="user != null" maxHeight="300">

          <p v-if="orders && orders.length <= 0">{{ $t('no_orders_found') }}</p>

          <v-table>
            <tbody>
              <tr v-for="order of orders" @click="() => $router.push('/admin/orders/' + order.id)" style="cursor: pointer;">
                <td>{{ order.reference }}</td>
                <td>{{ moment(order.created_at).format('DD.MM.YYYY - HH:mm') }}</td>
              </tr>
            </tbody>
          </v-table>

        </Card>

      </div>

      <Card v-if="user != null">
        <div class="grid grid-cols-1">
          <div style="display: flex; align-items: center; gap: 6px;">
            <b>{{ $t('created_at') }}:</b>
            {{ moment(user.created_at).format('DD.MM.YYYY - HH:mm') }}
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
