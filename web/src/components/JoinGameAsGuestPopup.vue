<script>
import Popup from '../components/Popup.vue';
import tableService from '../services/tableService';
import { toast } from 'vue3-toastify';

export default {
  props: {
    tableID: Number
  },
  components: {
    Popup,
  },
  data() {
    return {
      name: "",
    }
  },
  methods: {
    submit() {
      if (this.name.length === 0) {
        return;
      }

      console.log(this.tableID, this.name);

      tableService.joinAsGuest(this.tableID, this.name)
        .then((res) => {
          if (res.data.success === false) {
            toast.error(res.data.message);
            console.error(res.data.message);
            return;
          }

          this.$emit('success');
        }).catch((err) => {
          if (err.response && err.response.data && err.response.data.message) {
            toast.error(err.response.data.message);
            return;
          }
          toast.error(err.code);
        });
    }
  },
}
</script>

<template>
  <div class="form__row">
    <div class="form__col">
      <label class="form__label">Name (should be different from current players)</label>
      <input class="form__input" type="text" v-model="name" />
    </div>
  </div>
  <br/>
  <p>By clicking "{{ $t('join_game') }}" timer will be started for the guest user.</p>
  <br/>
  <div class="complete-popup__buttons">
    <button class="button button--primary" @click="submit" :disabled="name.length === 0">{{ $t('join_game') }}</button>
    <button class="button button--secondary" @click="$emit('closed')">{{ $t('cancel') }}</button>
  </div>
</template>

<style scoped>
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
