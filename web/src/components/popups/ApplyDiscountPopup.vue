<script>
import Popup from '../Popup.vue';
import orderService from '../../services/orderService';

export default {
    components: {
        Popup,
    },
    props: {
      id: Number,
    },
    data() {
      return {
        form: {
          amount: null,
        }
      };
    },
    methods: {
      onSave() {
        orderService.applyDiscount(this.id, this.form.amount)
          .then(() => {
            this.$emit('saved');
          });
      },
    }
}
</script>

<template>
  <Popup :title="$t('apply_discount')" @closed="$emit('closed')">
    <form class="form">

      <div class="form__row">
        <div class="form__col">
          <label class="form__label">{{ $t('discount_amount') }}</label>
          <input type="number" class="form__input" v-model="form.amount" />
        </div>
      </div>

    </form>

    <br />

    <div class="grid grid-cols-2">
      <button class="button button--primary" @click="onSave">{{ $t('apply') }}</button>
      <button class="button button--secondary" @click="$emit('closed')">{{ $t('cancel') }}</button>
    </div>
  </Popup>
</template>
