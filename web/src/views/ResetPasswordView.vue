<template>
  <main class="container">
    <div class="auth-box-outside">
      <div class="auth-box">
        <h1 class="auth-box__title">{{ $t('reset_password') }}</h1>

        <p class="auth-box__description">
          {{ $t('reset_password_text') }}
        </p>

        <form @submit.prevent="submit" class="auth-box__form form">
          <div class="form__row">
            <div class="form__col">
              <label class="form__label">{{ $t('email') }}</label>
              <input class="form__input" type="text" v-model="email" />
            </div>
          </div>
          <button class="button button--primary" type="submit">{{ $t('send') }}</button>
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
        email: null,
      }
    },
    methods: {
      submit() {
        if (this.email === "" || this.email === null)
          return;
        authService.forgotPassword(this.email)
          .then((res) => {
            if (res.data.success === true)
              toast.success("success.mail_sent");
          });
      }
    }
  }
</script>
