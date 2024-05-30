<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import moment from 'moment';
import { useRoute } from 'vue-router';
import Card from '../../../components/Card.vue';
import DeletePopup from '../../../components/DeletePopup.vue';
import productService from '../../../services/productService';
import ProductEditPopup from '../../../components/popups/ProductEditPopup.vue';
import { UploadsURL } from '../../../utils';

export default {
  setup() {
    const route = useRoute();
    return { route, moment };
  },
  components: {
    AdminLayout,
    Card,
    DeletePopup,
    ProductEditPopup
},
  data() {
    return {
      product: null,
      deletePopup: false,
      editPopup: false,
      UploadsURL: UploadsURL(),
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  mounted() {
    this.fetchProduct();
  },
  methods: {
    fetchProduct() {
      productService.get(parseInt(this.route.params.id))
        .then((res) => {
          this.product = res;
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
    onEdit() {
      this.fetchProduct();
      this.editPopup = false;
    },
    onDelete() {
      productService.del(parseInt(this.route.params.id))
        .then(() => {
          this.$router.push('/admin/products');
        });
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

        <Card :title="product.name" v-if="product != null">

          <template v-slot:actions>
            <router-link to="/admin/products" class="button button--small button--secondary">{{ $t('go_back') }}</router-link>
          </template>

          <div class="grid">
            <div><b>{{ $t('price') }}:</b> {{ product.price }}{{ convertCurrency() }}</div>

            <img v-if="product.image.upload_filename != ''" class="product__image" :src="UploadsURL + '/' + product.image.upload_filename" />
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
