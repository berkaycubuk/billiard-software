<script>
import tableService from '../../services/tableService';
import Popup from '../Popup.vue';

export default {
    components: {
        Popup,
    },
    data() {
      return {
        form: {
          name: "",
        }
      };
    },
    methods: {
      onSave() {
        tableService.create(this.form.name)
          .then(() => {
            this.$emit('saved');
          });
      },
    }
}
</script>

<template>
  <Popup :title="$t('new_table')" @closed="$emit('closed')">
    <form class="form">

      <div class="form__row">
        <div class="form__col">
          <label class="form__label">{{ $t('name') }}</label>
          <input type="text" class="form__input" v-model="form.name" />
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