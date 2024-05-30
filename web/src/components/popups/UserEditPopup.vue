<script>
import userService from '../../services/userService';
import Popup from '../Popup.vue';

export default {
    components: {
        Popup,
    },
    props: {
      user: Object,
      roles: [],
    },
    data() {
      return {
        id: this.user ? this.user.id : null,
        form: {
          name: this.user ? this.user.name : "",
          surname: this.user ? this.user.surname : "",
          email: this.user ? this.user.email : "",
          phone: this.user ? this.user.phone : "",
          //emailVerified: this.user && this.user.email_verified_at != null,
          role: this.user && this.user.role ? this.user.role.role_id : "",
        }
      };
    },
    methods: {
      onSave() {
        userService.update(this.id, this.form.name, this.form.surname, this.form.email, this.form.phone, this.form.role)
          .then((res) => {
            this.$emit('saved');
          });
      },
    }
}
</script>

<template>
  <Popup :title="$t('edit_user')" @closed="$emit('closed')">
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
