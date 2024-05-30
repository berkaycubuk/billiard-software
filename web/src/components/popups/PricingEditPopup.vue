<script>
import pricingService from '../../services/pricingService';
import Popup from '../Popup.vue';

export default {
    components: {
      Popup,
    },
    props: {
      pricing: Object,
      roles: [],
      subscriptions: [],
    },
    data() {
      return {
        id: this.pricing ? this.pricing.id : null,
        form: {
          role: this.pricing ? this.pricing.role_id : null,
          subscription: this.pricing ? this.pricing.subscription_id : null,
          playerCount: this.pricing ? this.pricing.player_count : null,
          perMinute: this.pricing ? this.pricing.per_minute : null,
        }
      };
    },
    methods: {
      onSave() {
        pricingService.update(this.id, this.form.role, this.form.subscription, this.form.playerCount, this.form.perMinute)
          .then(() => {
            this.$emit('saved');
          });
      },
    }
}
</script>

<template>
  <Popup :title="$t('edit_pricing')" @closed="$emit('closed')">
    <form class="form">

      <div class="form__row">
        <div class="form__col">
          <label class="form__label">{{ $t('role') }}</label>
          <select class="form__input" v-model="form.role">
            <option value="null">No role</option>
            <option v-for="role of roles" :value="role.id">{{ role.name }}</option>
          </select>
        </div>
      </div>

      <div class="form__row">
        <div class="form__col">
          <label class="form__label">{{ $t('subscription') }}</label>
          <select class="form__input" v-model="form.subscription">
            <option value="null">No subscription</option>
            <option v-for="sub of subscriptions" :value="sub.id">{{ sub.name }}</option>
          </select>
        </div>
      </div>

      <div class="form__row">
        <div class="form__col">
          <label class="form__label">{{ $t('player_count') }}</label>
          <input type="number" class="form__input" v-model="form.playerCount" />
        </div>
        <div class="form__col">
          <label class="form__label">{{ $t('per_minute') }}</label>
          <input type="number" class="form__input" v-model="form.perMinute" />
        </div>
      </div>

    </form>

    <br />

    <div class="grid grid-cols-2">
      <button class="button button--primary" @click="onSave">{{ $t('save') }}</button>
      <button class="button button--secondary" @click="$emit('closed')">{{ $t('cancel') }}</button>
    </div>
  </Popup>
</template>