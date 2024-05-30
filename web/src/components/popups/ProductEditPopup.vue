<script>
import Popup from '../Popup.vue';
import productService from '../../services/productService';
import FileUpload from '../FileUpload.vue';

export default {
    components: {
      Popup,
      FileUpload
  },
    props: {
      product: Object,
    },
    data() {
      return {
        id: this.product ? this.product.id : null,
        form: {
          name: this.product ? this.product.name : "",
          price: this.product ? this.product.price : "",
        }
      };
    },
    methods: {
      onSave() {
        productService.update(this.id, this.form.name, this.form.order, this.form.price, this.getUploadedImageID())
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
  <Popup :title="$t('edit_product')" @closed="$emit('closed')">
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
          <label class="form__label">{{ $t('image') }} (Leave empty if you want to keep current image)</label>
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
