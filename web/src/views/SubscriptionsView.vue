<script>
import HeaderSimple from '../components/HeaderSimple.vue';
import userService from '../services/userService';
import subscriptionService from '../services/subscriptionService';
import Popup from '../components/Popup.vue';
import moment from 'moment';
import BaseLayout from '../components/layout/BaseLayout.vue';
import { store } from '../store'

export default {
  components: {
    HeaderSimple,
    BaseLayout,
    Popup,
  },
  setup() {
    return { moment, store };
  },
  data() {
    return {
      activeSubscription: null,
      currency: import.meta.env.VITE_CURRENCY,
      subscriptions: [],
      selectedSub: null,
      buySubPopup: false,
    }
  },
  methods: {
    openSubPopup(id) {
      this.selectedSub = id;
      this.buySubPopup = true;
    },
    closeSubPopup() {
      this.buySubPopup = false;
    },
      convertCurrency() {
        if (!this.currency) return '';

        if (this.currency == 'kr') return 'kr';

        if (this.currency == 'eur') return '€';

        if (this.currency == 'usd') return '$';

        if (this.currency == 'try') return '₺';

        return this.currency;
      },
    pause() {
      if (this.activeSubscription === null) {
        return;
      }
      subscriptionService.pause(this.activeSubscription.id)
        .then(() => {
          this.fetchPage();
        });
    },
    unpause() {
      if (this.activeSubscription === null) {
        return;
      }
      subscriptionService.unpause(this.activeSubscription.id)
        .then(() => {
          this.fetchPage();
        });
    },
    buySub() {
      subscriptionService.buy(this.selectedSub)
        .then((res) => {
          if (res.data.success === false) {
            toast.error(this.$t(res.data.message));
            console.error(res.data.message);
            return;
          }

          this.$router.push('/payment/overview/' + res.data.order_id);
        });
    },
    fetchPage() {
      userService.myActiveSubscription()
        .then((response) => {
          if (response.data.success === false) {
            return;
          }

          this.activeSubscription = response.data.subscription;
        });
    }
  },
  mounted() {
    this.fetchPage();
    subscriptionService.all()
      .then((subscriptions) => {
        this.subscriptions = subscriptions;
      });
  },
}
</script>

<template>
  <BaseLayout>
  <HeaderSimple :title="$t('subscriptions')" backUrl="/profile" />

  <Popup v-if="buySubPopup" @closed="closeSubPopup" :title="$t('buy_subscription')">
    <p>{{ $t('buy_subscription_popup_text') }}</p>
    <br />
    <div class="complete-popup__buttons">
      <button class="button button--primary" @click="buySub">{{ $t('buy') }}</button>
      <button class="button button--secondary" @click="closeSubPopup">{{ $t('cancel') }}</button>
    </div>
  </Popup>

  <main class="container" style="flex: 1;">

    <div class="cards">

      <div class="card">
        <div class="card__header">
          <div class="card__title">{{ $t('available_subscriptions') }}</div>
        </div>
        <div class="card__body">

          <div class="subscriptions" v-if="subscriptions.length">

	  		<template v-for="sub of subscriptions">
            <div class="subscription-card" v-if="sub.hidden == false && ((store.user.role_id == 1 && sub.role == 2) || store.user.role_id == sub.role)">
              <div class="subscription-card__title">{{ sub.name }}</div>
              <div class="subscription-card__grid">
                <div class="subscription-card__time">
                  {{ $t('duration') }}: {{ sub.hours / 24 }} {{ $t('days') }}
                </div>
                    <button class="button button--primary" @click="() => openSubPopup(sub.id)" :disabled="activeSubscription">{{ $t('buy') }}: {{ parseFloat(sub.price).toFixed(2) }}{{ convertCurrency() }}</button>
              </div>
            </div>
			</template>
          
          </div>
      
        </div>
      </div>

      <div class="card">
        <div class="card__header">
          <div class="card__title">{{ $t('active_subscription') }}</div>
        </div>
        <div class="card__body">

          <div class="subscription-card" v-if="activeSubscription">
            <div class="subscription-card__title">{{ activeSubscription.name }}</div>

            <div class="subscription-card__time" v-if="activeSubscription.status === 1">
              {{ $t('ending_at') }}: {{ moment(activeSubscription.ending_at).format('DD.MM.YYYY - HH:mm') }}
            </div>

            <div class="subscription-card__time" v-if="activeSubscription.status === 2">
              {{ $t('paused') }}
            </div>

            <!--
            <br />

            <div class="grid">
              <button @click="pause" v-if="activeSubscription.status !== 2" class="button button--blue">{{ $t('pause_subscription') }}</button>
              <button @click="unpause" v-if="activeSubscription.status === 2" class="button button--blue">{{ $t('unpause_subscription') }}</button>
            </div>
            -->

          </div>
      
        </div>
      </div>

    </div>
  
  </main>
  </BaseLayout>
</template>

<style scoped>
.subscriptions {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
}

.complete-popup__buttons {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

@media only screen and (max-width: 650px) {
  .complete-popup__buttons {
    grid-template-columns: 1fr;
  }
}
</style>
