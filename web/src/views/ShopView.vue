<script>
import HeaderSimple from '../components/HeaderSimple.vue';
import ShopItem from '../components/ShopItem.vue';
import productService from '../services/productService';
import { cart } from '../cart';
import { store } from '../store';
import Popup from '../components/Popup.vue';
import shopService from '../services/shopService';
import { toast } from 'vue3-toastify';
import userService from '../services/userService';
import BaseLayout from '../components/layout/BaseLayout.vue';

export default {
  components: {
    HeaderSimple,
    ShopItem,
    BaseLayout,
    Popup,
  },
  setup() {
    return { store };
  },
  data() {
    return {
      cart,
      completePopupOpen: false,
      addToAccountPopup: false,
      selectedUser: null,
      products: [],
      guests: [],
      items: [],
      currency: import.meta.env.VITE_CURRENCY
    }
  },
  mounted() {
    try {
      productService.all()
        .then((response) => {
          if (!response.data.success) {
            return;
          }
          this.products = response.data.products;
          this.fetchGuests();
        }); 
    } catch(error) {
      console.error(error); 
    }
  },
  methods: {
    fetchGuests() {
      userService.getGuests()
        .then((res) => {
          this.guests = res.data.guests != null ? res.data.guests : [];
        });
    },
    clickComplete() {
      this.completePopupOpen = true;
    },
    completePopupClosed() {
      this.completePopupOpen = false;
    },
      convertCurrency() {
        if (!this.currency) return '';

        if (this.currency == 'kr') return 'kr';

        if (this.currency == 'eur') return '€';

        if (this.currency == 'usd') return '$';

        if (this.currency == 'try') return '₺';

        return this.currency;
      },
    emptyCart() {
      this.cart.empty()
    },
    payNow() {
      shopService.buy(this.cart.items.map((item) => {
        return {
          id: item.id,
          count: item.count,
        };
      }))
        .then((res) => {
          if (res.data.success === false) {
            console.error(res.data.message);
            toast.error(this.$t(res.data.message));
            return;
          }

          this.cart.empty();
          this.$router.push('/payment/overview/' + res.data.order_id);
        })
        .catch((err) => {
          console.error(err);
        });
    },
    guestAddToAccount() {
      shopService.addToGuestAccount(this.selectedUser,this.cart.items.map((item) => {
        return {
          id: item.id,
          count: item.count,
        };
      }))
        .then(() => {
          this.completePopupOpen = false;
          this.addToAccountPopup = false;
          toast.success(this.$t('order_added_to_account'));
          this.cart.empty()
          // redirect back to home
          // client requested this
          setTimeout(() => {
            this.$router.push('/')
          }, 2000)
        })
        .catch((err) => console.error(err))
    },
    addToAccount() {
      shopService.addToAccount(this.cart.items.map((item) => {
        return {
          id: item.id,
          count: item.count,
        };
      }))
        .then((res) => {
          if (res.data.success === false) {
            console.error(res.data.message);
            toast.error(this.$t(res.data.message));
            return;
          }

          toast.success(this.$t('order_added_to_account'));
          this.completePopupOpen = false;
          this.cart.empty();
          // redirect back to home
          // client requested this
          setTimeout(() => {
            this.$router.push('/')
          }, 2000)
        })
        .catch((err) => {
          console.error(err);
        });
    }
  },
}
</script>

<template>
  <BaseLayout>
  <HeaderSimple :title="$t('shop')" backUrl="/" />

  <Popup v-if="completePopupOpen === true" :title="$t('complete_cart')" @closed="completePopupClosed">
    <div class="complete-popup__buttons">
      <button class="button button--primary" @click="payNow">{{ $t('pay_now') }}</button>
      <button v-if="!(store.user && store.user.role_id === 4) && store.user.role_id != 3" class="button button--secondary" @click="addToAccount">{{ $t('add_to_account') }}</button>
      <button v-if="store.user && store.user.role_id === 4" class="button button--secondary" @click="() => addToAccountPopup = true">{{ $t('add_to_guest_account') }}</button>
    </div>
  </Popup>

  <Popup v-if="addToAccountPopup === true" :title="$t('add_to_guest_account')" @closed="() => addToAccountPopup = false">
    <div class="user-selector">

      <div class="user-selector__item" :class="{ 'user-selector__item--selected': user.name === selectedUser }" v-for="user of guests" @click="() => selectedUser = (selectedUser === user.name ? null : user.name)">
        {{ user.name }}
      </div>
      
    </div>
    <br />
    <div class="complete-popup__buttons">
      <button class="button button--primary" @click="guestAddToAccount" :disabled="selectedUser === null">{{ $t('add_to_account') }}</button>
      <button class="button button--secondary" @click="() => addToAccountPopup = false">{{ $t('cancel') }}</button>
    </div>
  </Popup>

  <main class="container" style="flex: 1;">

    <div class="cards">

      <div class="shop-items" v-if="products.length > 0">
        <ShopItem v-for="product of products" :image="product.image.upload_filename" :id="product.id" :name="product.name" :price="product.price" :currency="currency" />
      </div>

      <div class="small-cart">
        <v-dialog>
          <template v-slot:activator="{props}">
            <v-badge v-if="cart.items.length > 0" :content="cart.items.length" color="red">
              <v-btn v-bind="props" icon="mdi-cart" size="large" color="green"></v-btn>
            </v-badge>
            <v-btn v-if="cart.items.length == 0" v-bind="props" icon="mdi-cart" size="large" color="green"></v-btn>
          </template>

          <template v-slot:default="{ isActive }">
            <v-card :title="$t('cart')">
              <template v-slot:append>
                <v-btn icon="mdi-close" @click="isActive.value = false"></v-btn>
              </template>
              <v-card-text>
              <div class="cart__items">
                <div class="cart__item" v-for="item of cart.items">
                  <div class="cart-item__name">{{ item.name }}</div>
                  <div class="cart-item__value">
                    <div class="cart-item__price">{{ item.price }}{{ convertCurrency() }}</div>
                    <div class="cart-item__count">{{ item.count }}</div>
                  </div>
                </div>
              </div>

              <div class="cart__total">{{ $t('total') }}: {{ cart.total }}{{ convertCurrency() }}</div>
              </v-card-text>
              <v-card-actions v-if="cart.items.length > 0">
                <div style="display: flex; flex-direction: column; gap: 6px; width: 100%;">
                <v-btn block color="green" variant="flat" v-if="store.user && store.user.role_id == 3" @click="() => {payNow(); isActive.value = false;}">{{ $t('complete_cart') }}</v-btn>
                <v-btn block color="green" variant="flat" v-if="store.user && store.user.role_id != 3" @click="() => {clickComplete(); isActive.value = false;}">{{ $t('complete_cart') }}</v-btn>
                <v-btn block color="red" variant="flat" style="margin-left: 0;" @click="() => {emptyCart(); isActive.value = false;}">{{ $t('empty_cart') }}</v-btn>
                </div>
              </v-card-actions>
            </v-card>
          </template>
        </v-dialog>
      </div>

      <div class="cart-total-container">
        <div class="card">
          <div class="card__header">
            <div class="card__title">{{ $t('cart') }}</div>
          </div>
          <div class="card__body">
            <div class="cart">

              <div class="cart__items">
                <div class="cart__item" v-for="item of cart.items">
                  <div class="cart-item__name">{{ item.name }}</div>
                  <div class="cart-item__value">
                    <div class="cart-item__price">{{ item.price }}{{ convertCurrency() }}</div>
                    <div class="cart-item__count">{{ item.count }}</div>
                  </div>
                </div>
              </div>

              <div class="cart__total">{{ $t('total') }}: {{ cart.total }}{{ convertCurrency() }}</div>

              <div class="cart__buttons" v-if="store.user && store.user.role_id == 3">
                <button v-if="cart.items.length > 0" class="button button--primary" @click="payNow">{{ $t('complete_cart') }}</button>
                <button v-if="cart.items.length > 0" class="button button--red" @click="cart.empty">{{ $t('empty_cart') }}</button>
              </div>

              <div class="cart__buttons" v-if="store.user && store.user.role_id != 3">
                <button v-if="cart.items.length > 0" class="button button--primary" @click="clickComplete">{{ $t('complete_cart') }}</button>
                <button v-if="cart.items.length > 0" class="button button--red" @click="cart.empty">{{ $t('empty_cart') }}</button>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>

    <div style="height: 200px"></div>

  </main>
  </BaseLayout>
</template>

<style scoped>
.small-cart {
  display: none;
  position: fixed;
  bottom: 70px;
  right: 20px;
}

.shop-items {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
  gap: 12px;
  margin: 10px 0;
}

.cart {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cart__items {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.cart__item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  border-radius: 8px;
  background-color: #232323;
}

.cart-item__value {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

.cart-item__count {
  text-align: right;
}

.cart__total {
  font-size: 1.4rem;
}

.cart__buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.user-selector {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
}

.user-selector__item {
  cursor: pointer;
  padding: 14px;
  background-color: #232323;
  border-radius: 8px;
}

.user-selector__item--selected {
  outline: 2px solid var(--green);
}

.complete-popup__buttons {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.cart-total-container {
}

@media only screen and (max-width: 1200px) {
  .shop-items {
    grid-template-columns: 1fr 1fr 1fr;
	order: 1;
  }
}

@media only screen and (max-width: 650px) {
  .shop-items {
    grid-template-columns: 1fr 1fr;
	order: 1;
  }

  .complete-popup__buttons {
    grid-template-columns: 1fr;
  }
}

@media only screen and (max-width: 600px) {
  .cart-total-container {
    display: none;
  }

  .small-cart {
    display: block;
  }
}

@media only screen and (max-width: 450px) {
  .shop-items {
    grid-template-columns: 1fr 1fr;
	order: 1;
  }
}

@media only screen and (max-width: 330px) {
  .shop-items {
    grid-template-columns: 1fr;
	order: 1;
  }
}
</style>
