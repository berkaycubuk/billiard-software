<script>
import AdminLayout from '../../components/layout/AdminLayout.vue';
import moment from 'moment';
import Card from '../../components/Card.vue';
import subscriptionService from '../../services/subscriptionService';
import DeletePopup from '../../components/DeletePopup.vue';
import OrderApprovePopup from '../../components/OrderApprovePopup.vue';
import SubscriptionEditPopup from '../../components/SubscriptionEditPopup.vue';

export default {
  setup() {
    return { moment, subscriptionService };
  },
  components: {
    AdminLayout,
    Card,
    DeletePopup,
    OrderApprovePopup,
    SubscriptionEditPopup
},
  data() {
    return {
      subscription: null,
      deletePopup: false,
      editPopup: false,
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  mounted() {
    this.fetchSubscription();
  },
  methods: {
    fetchSubscription() {
      subscriptionService.get(parseInt(this.$route.params.id))
        .then((res) => {
          this.subscription = res;
        });
    },
    onSaved() {
      this.fetchSubscription();
      this.editPopup = false;
    },
    onDelete() {
      if (this.order === null) {
        return;
      }

      subscriptionService.del(this.subscription.id)
        .then(() => {
          this.$router.push('/admin/subscriptions');
        });
    },
  },
}
</script>

<template>
  <AdminLayout>

    <DeletePopup v-if="deletePopup === true" :title="$t('delete_subscription')" @closed="() => deletePopup = false" @deleted="onDelete" />
    <SubscriptionEditPopup v-if="editPopup === true" @closed="() => editPopup = false" :subscription="subscription" @saved="onSaved" />

    <div class="cards">

      <Card :title="$t('subscription')" v-if="subscription != null">
        <template v-slot:actions>
          <router-link to="/admin/subscriptions" class="button button--small button--secondary">{{ $t('go_back') }}</router-link>
        </template>

        <div class="grid">
          <div><b>{{ $t('name') }}:</b> {{ subscription.name }}</div>
          <div><b>{{ $t('price') }}:</b> {{ parseFloat(subscription.price).toFixed(2) }}{{ currency }}</div>
          <div><b>{{ $t('time') }}:</b> {{ subscriptionService.hoursToHumanReadable(subscription.hours) }}</div>
        </div>
      </Card>

      <Card v-if="subscription != null">
        <div class="grid grid-cols-1">
          <div style="display: flex; align-items: center; gap: 6px;">
            <b>{{ $t('created_at') }}:</b>
            {{ moment(subscription.created_at).format('DD.MM.YYYY - HH:mm') }}
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
.users {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.user {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  align-items: center;
  padding: 10px;
  background-color: #232323;
  border-radius: 8px;
}

.user__right {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
}
</style>
