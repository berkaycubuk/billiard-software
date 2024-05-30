<script>
import roleService from '../../services/roleService';
import Popup from '../Popup.vue';

export default {
    components: {
        Popup,
    },
    props: {
      role: Object,
    },
    data() {
      return {
        id: this.role ? this.role.id : null,
        form: {
          name: this.role ? this.role.name : "",
        }
      };
    },
    methods: {
      onSave() {
        roleService.update(this.id, this.form.name)
          .then(() => {
            this.$emit('saved');
          });
      },
    }
}
</script>

<template>
  <Popup :title="$t('edit_role')" @closed="$emit('closed')">
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