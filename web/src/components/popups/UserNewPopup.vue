<script>
import userService from '../../services/userService';
import Popup from '../Popup.vue';

export default {
    components: {
        Popup,
    },
    props: {
      roles: [],
    },
    data() {
      return {
        errors: {},
        form: {
          name: "",
          surname: "",
          email: "",
          phone: "",
          emailVerified: false,
          role: "",
        }
      };
    },
    methods: {
      onSave() {
        try {
          userService.create(this.form.name, this.form.surname, this.form.email, this.form.phone, this.form.emailVerified, this.form.role)
            .then(() => {
              this.$emit('saved');
            });
        } catch(error) {
            console.log(error);
        }
      },
    }
}
</script>

<template>
  <Popup :title="$t('new_user')" @closed="$emit('closed')">
    <form class="form">

      <div class="form__row">
        <div class="form__col">
          <label class="form__label">{{ $t('name') }}</label>
          <input type="text" class="form__input" v-model="form.name" />
        </div>
        <div class="form__col">
          <label>{{ $t('surname') }}</label>
          <input type="text" class="form__input" v-model="form.surname" />
        </div>
      </div>

      <div class="form__row">
        <div class="form__col">
          <label>{{ $t('email') }}</label>
          <input type="email" class="form__input" v-model="form.email" />
        </div>
      </div>

      <div class="form__row">
        <div class="form__col">
          <label>{{ $t('phone') }}</label>
          <input type="tel" class="form__input" v-model="form.phone" />
        </div>
      </div>

      <div class="form__row">
        <div class="form__col">
          <label>{{ $t('role') }}</label>
          <select class="form__input" v-model="form.role">
            <option v-for="role of roles" :value="role.id">{{ role.name }}</option>
          </select>
        </div>
      </div>

      <!--
      <div class="form__row">
        <div class="form__col" style="width: fit-content;">
          <label>{{ $t('email_verified') }} <input type="checkbox" v-model="form.emailVerified" /></label>
        </div>
      </div>
    -->

    </form>

    <br />

    <div class="grid grid-cols-2">
      <button class="button button--primary" @click="onSave">{{ $t('save') }}</button>
      <button class="button button--secondary" @click="$emit('closed')">{{ $t('cancel') }}</button>
    </div>
  </Popup>
</template>
