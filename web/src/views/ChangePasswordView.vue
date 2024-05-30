<script>
import { toast } from 'vue3-toastify';
import HeaderSimple from '../components/HeaderSimple.vue';
import authService from '../services/authService';
import BaseLayout from '../components/layout/BaseLayout.vue';

export default {
  components: {
    HeaderSimple,
    BaseLayout,
  },
  data() {
    return {
      current: null,
      newPassword: null,
      newPasswordConfirm: null,
    };
  },
  methods: {
    submit() {
      if (this.current === null || this.newPassword === null || this.newPasswordConfirm === null) {
        return;
      }

      if (this.newPassword != this.newPasswordConfirm) {
        toast.error("errors.passwords_not_match");
        return;
      }

      authService.updatePassword(this.current, this.newPassword)
        .then(() => {
          this.current = null;
          this.newPassword = null;
          this.newPasswordConfirm = null;
          toast.success("success.password_updated");
        });
    },
  },
}
</script>

<template>
  <BaseLayout>
  <HeaderSimple :title="$t('change_password')" backUrl="/profile" />
  <main class="container" style="flex: 1;">
    <div class="change-pass-form">
      <form class="form" @submit.prevent="submit">
        <div class="form__row">
          <div class="form__col">
            <label class="form__label">{{ $t('current_password') }}</label>
            <input type="password" class="form__input" v-model="current" />
          </div>
        </div>
        <div class="form__row">
          <div class="form__col">
            <label class="form__label">{{ $t('new_password') }}</label>
            <input type="password" class="form__input" v-model="newPassword" />
          </div>
          <div class="form__col">
            <label class="form__label">{{ $t('confirm_new_password') }}</label>
            <input type="password" class="form__input" v-model="newPasswordConfirm" />
          </div>
        </div>
        <div class="form__row">
          <button class="button button--primary" type="submit">{{ $t('save') }}</button>
        </div>
      </form>
    </div>
  </main>
  </BaseLayout>
</template>

<style scoped>
.change-pass-form {
  margin: 20px 0;
  padding: 20px;
  background-color: #191919;
  border-radius: 8px;
}
</style>
