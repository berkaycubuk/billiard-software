<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import moment from 'moment';
import { useRoute } from 'vue-router';
import Card from '../../../components/Card.vue';
import DeletePopup from '../../../components/DeletePopup.vue';
import DataTable from '../../../components/DataTable.vue';
import DataRow from '../../../components/DataRow.vue';
import DataCol from '../../../components/DataCol.vue';
import ProductEditPopup from '../../../components/popups/ProductEditPopup.vue';
import { UploadsURL } from '../../../utils';
import pricingService from '../../../services/pricingService';

export default {
  setup() {
    const route = useRoute();
    return { route, moment };
  },
  components: {
    AdminLayout,
    Card,
    DeletePopup,
    DataTable,
    DataRow,
    DataCol,
    ProductEditPopup
},
  data() {
    return {
      pricings: [],
      deletePopup: false,
      editPopup: false,
      UploadsURL: UploadsURL(),
    }
  },
  mounted() {
    this.fetchProduct();
  },
  methods: {
    fetchProduct() {
      pricingService.get(parseInt(this.route.params.id))
        .then((res) => {
          this.pricings = res;
        });
    },
    onEdit() {
      this.fetchProduct();
      this.editPopup = false;
    },
    onDelete() {
    }
  }
}
</script>

<template>
  <AdminLayout>

    <ProductEditPopup v-if="editPopup === true" :product="product" @closed="() => editPopup = false" @saved="onEdit" />
    <DeletePopup v-if="deletePopup === true" :title="$t('delete_product')" @closed="() => deletePopup = false" @deleted="onDelete" />

    <div class="cards">

      <div class="grid">

        <Card :title="$t('role_pricing')">

          <template v-slot:actions>
            <router-link to="/admin/pricing" class="button button--small button--secondary">{{ $t('go_back') }}</router-link>
          </template>

          <div v-if="pricings != null && pricings.length > 0" class="grid">

            <Card :title="`${pricing.subscription_name != '' ? pricing.subscription_name : 'No subscription'}`" v-for="pricing of pricings">

              <Card v-for="block of pricing.blocks">
                <div><b>{{ $t('player_count') }}:</b> {{ block.player_count }}</div>
                <div><b>{{ $t('per_minute') }}:</b> {{ block.per_minute }}</div>
              </Card>
              
            </Card>

          </div>

        </Card>

      </div>

      <Card v-if="product != null">
        <div class="grid grid-cols-1">
          <div style="display: flex; align-items: center; gap: 6px;">
            <b>{{ $t('created_at') }}:</b>
            {{ moment(product.created_at).format('DD.MM.YYYY - HH:mm') }}
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
.product__image {
  max-width: 600px;
  width: 100%;
}
</style>
