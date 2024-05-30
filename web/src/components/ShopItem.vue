<script>
import { cart } from '../cart';
import { UploadsURL } from '../utils';

export default {
  props: ['id', 'image', 'name', 'price', 'currency'],
  data() {
    return {
      UploadsURL: UploadsURL(),
      currency: import.meta.env.VITE_CURRENCY
    };
  },
  setup() {
    return { cart };
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
    add() {
      this.cart.addItem(this.id, this.name, parseFloat(this.price));
    },
    bigAdd() {
      if (this.cart.getCount(this.id) > 0) {
        return;
      }
      this.cart.addItem(this.id, this.name, parseFloat(this.price));
    },
    remove() {
      this.cart.removeItem(this.id, parseFloat(this.price));
    }
  },
}
</script>

<template>
  <div class="shop-item">
    <div class="shop-item__inner" @click="bigAdd">
      <img v-if="image != ''" :src="`${UploadsURL}/${image}`" />
      <div class="shop-item__bottom">
        <div class="shop-item__info">
          <div class="shop-item__name">{{ name }}</div>
          <div class="shop-item__price">{{ price }}{{ convertCurrency() }}</div>
        </div>
      </div>
    </div>

    <div class="shop-item__actions" v-if="cart.getCount(this.id) > 0">
      <div class="shop-item__value">
        <button class="button button--primary button--flex-center" @click="add">
          <svg width="24px" height="24px" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#ffffff"><path d="M8 12h4m4 0h-4m0 0V8m0 4v4M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
        </button>
        <div class="shop-item__value-count">{{ cart.getCount(id) }}</div>
        <button class="button button--primary button--flex-center" @click="remove">
          <svg width="24px" height="24px" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#ffffff"><path d="M8 12h8M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.shop-item {
  position: relative;
  height: fit-content;
}

.shop-item__inner {
  height: fit-content;
  cursor: pointer;
  padding: 20px;
  border-radius: 8px;
  background-color: #191919;
  transition: opacity .2s;
}

.shop-item__inner:hover {
  opacity: .8;
}

.shop-item img {
  width: 100%;
  height: 150px;
  object-fit: cover;
  border-radius: 8px;
}

.shop-item__actions {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0,0,0,.5);
  border-radius: 8px;
}

.shop-item__value {
  display: flex;
  align-items: center;
  gap: 10px;
}

.shop-item__value-count {
  font-size: 1.2rem;
}

.shop-item__bottom {
  margin-top: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.shop-item__info {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.shop-item__name {
  font-size: 1.2rem;
  font-weight: 600;
}

.shop-item__price {
  font-size: .9rem;
  color: #8A8B85;
}

@media only screen and (max-width: 450px) {
  .shop-item__info {
    gap: 3px;
  }

  .shop-item__name {
    font-size: .9rem;
    font-weight: 500;
  }
}
</style>
