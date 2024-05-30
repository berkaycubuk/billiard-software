<script>
import AdminLayout from '../../../components/layout/AdminLayout.vue';
import moment from 'moment';
import Card from '../../../components/Card.vue';
import productService from '../../../services/productService';
import ProductNewPopup from '../../../components/popups/ProductNewPopup.vue';
import { UploadsURL } from '../../../utils';

export default {
  setup() {
    return { moment };
  },
  components: {
    AdminLayout,
    Card,
    ProductNewPopup
},
  data() {
    return {
      products: [],
      newPopup: false,
      UploadsURL: UploadsURL(),
      currency: import.meta.env.VITE_CURRENCY,
    }
  },
  mounted() {
    this.fetchProducts();
  },
  methods: {
    fetchProducts() {
      productService.all()
        .then((res) => {
          this.products = res.data.products;
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
      this.fetchProducts();
      this.newPopup = false;
    },
    moveUp(id) {
      productService.updateOrder(id, 1)
        .then(() => this.fetchProducts())
    },
    moveDown(id) {
      productService.updateOrder(id, 2)
        .then(() => this.fetchProducts())
    },
  }
}
</script>

<template>
  <AdminLayout>

    <ProductNewPopup v-if="newPopup === true" @closed="() => newPopup = false" @saved="onCreate" />

    <div class="users-page">

      <Card :title="$t('products')">

        <template v-slot:actions>
          <button @click="() => newPopup = true" class="button button--small button--primary">{{ $t('new_product') }}</button>
        </template>

        <v-table>
          <thead>
            <tr>
              <th><b>{{ $t('image') }}</b></th>
              <th><b>{{ $t('name') }}</b></th>
              <th><b>{{ $t('price') }}</b></th>
              <th><b>{{ $t('change_order') }}</b></th>
              <th>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="product of products">
              <td><img v-if="product.image.upload_filename != ''" class="product__image" :src="`${UploadsURL}/${product.image.upload_filename}`" /></td>
              <td>{{ product.name }}</td>
              <td>{{ product.price }}{{ convertCurrency() }}</td>
              <td>
                <div style="display: flex; gap: 5px">
                  <button class="order-btn" @click="() => moveUp(product.id)"><v-icon>mdi-chevron-up</v-icon></button>
                  <button class="order-btn" @click="() => moveDown(product.id)"><v-icon>mdi-chevron-down</v-icon></button>
                </div>
              </td>
              <td>
                <div class="user__right">
                  <router-link :to="'/admin/products/' + product.id" class="button button--small button--primary">{{ $t('details') }}</router-link>
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

.order-btn {
  width: fit-content;
  background-color: var(--green);
  border-radius: 4px;
}

.order-btn:hover {
  opacity: .8;
}
</style>
