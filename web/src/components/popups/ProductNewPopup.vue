<script>
import Popup from '../Popup.vue';
import productService from '../../services/productService';
import FileUpload from '../FileUpload.vue';

export default {
    components: {
      Popup,
      FileUpload
    },
    data() {
      return {
        form: {
          name: "",
          price: "",
        }
      };
    },
    methods: {
      onSave() {
        productService.create(this.form.name, parseFloat(this.form.price), this.getUploadedImageID())
          .then(() => {
            this.$emit('saved');
          });
      },
      getUploadedImageID() {
        let found = document.querySelector('input[name="file"]');
        if (found && found.value) {
          return parseInt(found.value);
        }

        return null;
      }
    }
}
</script>

<template>
  <Popup :title="$t('new_product')" @closed="$emit('closed')">
    <form class="form">

      <div class="form__row">
        <div class="form__col">
          <label class="form__label">{{ $t('name') }}</label>
          <input type="text" class="form__input" v-model="form.name" />
        </div>
      </div>

      <div class="form__row">
        <div class="form__col">
          <label class="form__label">{{ $t('price') }}</label>
          <input type="number" class="form__input" v-model="form.price" />
        </div>
      </div>

      <div class="form__row">
        <div class="form__col">
          <label class="form__label">{{ $t('image') }}</label>
          <FileUpload />
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