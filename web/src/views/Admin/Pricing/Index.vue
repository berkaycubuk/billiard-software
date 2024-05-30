<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import moment from 'moment';
import Card from '../../../components/Card.vue';
import pricingService from '../../../services/pricingService';
import roleService from '../../../services/roleService';
import subscriptionService from '../../../services/subscriptionService';
import PricingEditPopup from '../../../components/popups/PricingEditPopup.vue';
import PricingNewPopup from '../../../components/popups/PricingNewPopup.vue';
import DeletePopup from '../../../components/DeletePopup.vue';

export default {
  setup() {
    return { moment };
  },
  components: {
    AdminLayout,
    Card,
    PricingEditPopup,
    PricingNewPopup,
    DeletePopup
},
  data() {
    return {
      pricings: [],
      roles: [],
      subscriptions: [],
      editing: null,
      deleting: null,
      newPopup: false,
      editPopup: false,
      deletePopup: false,
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  mounted() {
    this.fetchPricings();
    this.fetchRoles();
    this.fetchSubscriptions();
  },
  methods: {
    fetchPricings() {
      pricingService.all()
        .then((res) => {
          this.pricings = res;
        });
    },
    fetchRoles() {
      roleService.all()
        .then((res) => {
          this.roles = res;
        });
    },
    fetchSubscriptions() {
      subscriptionService.all()
        .then((res) => {
          this.subscriptions = res;
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
    onCreate() {
      this.fetchPricings();
      this.newPopup = false;
    },
    onCreateComplete() {
      this.fetchPricings();
      this.newPopup = false;
    },
    onEditClick(pricing) {
      this.editing = pricing;
      this.editPopup = true;
    },
    onDeleteClick(id) {
      this.deleting = id;
      this.deletePopup = true;
    },
    onEditComplete() {
      this.fetchPricings();
      this.editPopup = false;
      this.editing = null;
    },
    onEditClose() {
      this.editPopup = false;
      this.editing = null;
    },
    onDeleteClose() {
      this.deletePopup = false;
      this.deleting = null;
    },
    onDeleteComplete() {
      pricingService.del(this.deleting)
        .then(() => {
          this.fetchPricings();
          this.deletePopup = false;
          this.deleting = null;
        });
    }
  }
}
</script>

<template>
  <AdminLayout>

    <PricingNewPopup v-if="newPopup === true" :roles="roles" :subscriptions="subscriptions" @closed="() => newPopup = false" @saved="onCreateComplete" />
    <PricingEditPopup v-if="editPopup === true" :roles="roles" :subscriptions="subscriptions" :pricing="editing" @closed="onEditClose" @saved="onEditComplete" />
    <DeletePopup v-if="deletePopup === true" @closed="onDeleteClose" @deleted="onDeleteComplete" />

    <div class="users-page">

      <Card :title="$t('pricing')">

        <template v-slot:actions>
          <button @click="() => newPopup = true" class="button button--small button--primary">{{ $t('new_pricing') }}</button>
        </template>

        <v-table>
          <thead>
            <tr cols="5">
              <th><b>{{ $t('role') }}</b></th>
              <th><b>{{ $t('subscription') }}</b></th>
              <th><b>{{ $t('player_count') }}</b></th>
              <th><b>{{ $t('per_minute') }}</b></th>
              <th>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="pricing of pricings" cols="5">
              <td>{{ pricing.role !== null ? pricing.role.name : "No role" }}</td>
              <td>{{ pricing.subscription !== null ? pricing.subscription.name : "No subscription" }}</td>
              <td>{{ pricing.player_count }}</td>
              <td>{{ pricing.per_minute }}{{ convertCurrency() }}</td>
              <td>
                <div class="user__right">
                  <button @click="() => onEditClick(pricing)" class="button button--small button--primary">{{ $t('edit') }}</button>
                  <button @click="() => onDeleteClick(pricing.id)" class="button button--small button--red">{{ $t('delete') }}</button>
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

.product__image {
  width: 60px;
  height: 60px;
  object-fit: cover;
  border-radius: 8px;
}
</style>
