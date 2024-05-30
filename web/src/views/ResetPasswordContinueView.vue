<template>
  <main class="container">
    <div class="auth-box-outside">
      <div class="auth-box">
        <h1 class="auth-box__title">{{ $t('reset_password') }}</h1>

        <form @submit.prevent="submit" class="auth-box__form form">
          <div class="form__row">
            <div class="form__col">
              <label class="form__label">{{ $t('new_password') }}</label>
              <input class="form__input" type="password" v-model="password" />
            </div>
          </div>
          <div class="form__row">
            <div class="form__col">
              <label class="form__label">{{ $t('new_password_repeat') }}</label>
              <input class="form__input" type="password" v-model="passwordRepeat" />
            </div>
          </div>
          <button class="button button--primary" type="submit">{{ $t('update') }}</button>
        </form>
      </div>
    </div>
  </main>
</template>

<script>
import { toast } from 'vue3-toastify';
import authService from '../services/authService';

  export default {
    data() {
      return {
        password: null,
        passwordRepeat: null,
      }
    },
    methods: {
      submit() {
        if (this.password === null || this.passwordRepeat === null)
          return;
        if (this.password != this.passwordRepeat) {
          toast.error("errors.passwords_not_match");
          return;
        }
        authService.passwordResetComplete(this.$route.params.token, this.password)
          .then((res) => {
            if (res.data.success === true) {
              this.$router.push('/login?msg=success.password_updated');
              return;
            }

            toast.error(res.data.message);
          });
      }
    }
  }
</script>
