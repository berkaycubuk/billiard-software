<script>
import reportService from '../../services/reportService';
import Popup from '../Popup.vue';

export default {
    props: ['id'],
    components: {
        Popup,
    },
    data() {
      return {
        user: null,
        products: null,
        total: null,
      currency: import.meta.env.VITE_CURRENCY,
      };
    },
    mounted() {
      this.fetchData();
    },
    methods: {
      fetchData() {
        reportService.user(this.id, this.$route.query.start, this.$route.query.end)
          .then((res) => {
            this.user = res.data.user;
            this.products = res.data.products;
            this.total = res.data.total;
          });
      },
      convertCurrency() {
        if (!this.currency) return '';

        if (this.currency == 'kr') return 'kr';

        if (this.currency == 'eur') return '€';

        if (this.currency == 'usd') return '$';

        if (this.currency == 'try') return '₺';

        return this.currency;
      }
    }
}
</script>

<template>
  <Popup :title="$t('user_report')" @closed="$emit('closed')">
    <div v-if="user != null">
      <h3>{{ this.user.name }} {{ this.user.surname }}</h3>
    </div>

    <div v-if="products != null" style="margin-top: 10px;">

      <div v-for="product of products">
        {{ product.count }} x {{ product.name }} | {{ product.total }}{{ convertCurrency() }}
      </div>

    </div>

    <div v-if="total != null" style="margin-top: 10px;">
      <b>{{ $t('total') }}:</b> {{ this.total }}{{ convertCurrency() }}
    </div>

    <div class="grid grid-cols-2" style="margin-top: 16px;">
      <button class="button button--secondary" @click="$emit('closed')">{{ $t('close') }}</button>
    </div>
  </Popup>
</template>
