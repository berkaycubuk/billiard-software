<script>
import AdminLayout from '../../components/layout/AdminLayout.vue';
import subscriptionService from '../../services/subscriptionService';
import moment from 'moment';
import Card from '../../components/Card.vue';
import SubscriptionNewPopup from '../../components/SubscriptionNewPopup.vue';

export default {
  setup() {
    return { moment, subscriptionService };
  },
  components: {
    AdminLayout,
    Card,
    SubscriptionNewPopup
},
  data() {
    return {
      addPopup: false,
      subscriptions: [],  
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  mounted() {
    this.fetchSubscriptions();
  },
  methods: {
      convertCurrency() {
        if (!this.currency) return '';

        if (this.currency == 'kr') return 'kr';

        if (this.currency == 'eur') return '€';

        if (this.currency == 'usd') return '$';

        if (this.currency == 'try') return '₺';

        return this.currency;
      },
    fetchSubscriptions() {
      subscriptionService.all()
        .then((subscriptions) => {
          this.subscriptions = subscriptions;
        }); 
    },
    onAdd() {
      this.fetchSubscriptions();
      this.addPopup = false;
    }
  }
}
</script>

<template>
  <AdminLayout>

    <SubscriptionNewPopup v-if="addPopup === true" @closed="() => addPopup = false" @saved="onAdd" />

    <div class="users-page">
      <Card :title="$t('subscriptions')">

        <template v-slot:actions>
          <button @click="() => addPopup = true" class="button button--small button--primary">{{ $t('new_subscription') }}</button>
        </template>

        <v-table>
          <thead>
            <tr>
              <th><b>{{ $t('name') }}</b></th>
              <th><b>{{ $t('price') }}</b></th>
              <th><b>{{ $t('time') }}</b></th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr cols="4" v-for="sub of subscriptions">
              <td>{{ sub.name }}</td>
              <td>{{ parseFloat(sub.price).toFixed(2) }}kr</td>
              <td>{{ subscriptionService.hoursToHumanReadable(sub.hours) }}</td>
              <td align="right">
                <router-link :to="'/admin/subscriptions/' + sub.id" class="button button--small button--primary">{{ $t('details') }}</router-link>
              </td>
            </tr>
          </tbody>
        </v-table>

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
