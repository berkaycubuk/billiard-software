<template>
  <main class="container">
    <div class="auth-box-outside">
      <div class="auth-box">
        <h1 class="auth-box__title">{{ $t('title') }}</h1>

        <p class="auth-box__description">
          {{ $t('register_desc') }}
        </p>

        <form @submit.prevent="register" class="auth-box__form form">
          <div class="form__row">
            <div class="form__col">
              <label class="form__label">{{ t('name') }}</label>
              <input class="form__input" type="text" v-model="name" :class="{ 'form__input--error': errors.name }" />
              <p v-if="errors.name" class="form__error-text">{{ t('validation.' + errors.name, { field: t('name') }) }}</p>
            </div>
            <div class="form__col">
              <label class="form__label">{{ t('surname') }}</label>
              <input class="form__input" type="text" v-model="surname" :class="{ 'form__input--error': errors.surname }" />
              <p v-if="errors.surname" class="form__error-text">{{ t('validation.' + errors.surname, { field: t('surname') }) }}</p>
            </div>
          </div>
          <div class="form__row">
            <div class="form__col">
              <label class="form__label">{{ t('email') }}</label>
              <input class="form__input" type="text" v-model="email" :class="{ 'form__input--error': errors.email }" autocapitalize="off" />
              <p v-if="errors.email" class="form__error-text">{{ t('validation.' + errors.email, { field: t('email') }) }}</p>
            </div>
          </div>
          <div class="form__row">
            <div class="form__col">
              <label class="form__label">{{ t('phone') }}</label>
              <input class="form__input" type="tel" v-model="phone" :class="{ 'form__input--error': errors.phone }" />
              <p v-if="errors.phone" class="form__error-text">{{ t('validation.' + errors.phone, { field: t('phone') }) }}</p>
            </div>
          </div>
          <div class="form__row">
            <div class="form__col">
              <label class="form__label">{{ t('password') }}</label>
              <input class="form__input" type="password" v-model="password" :class="{ 'form__input--error': errors.password }" autocapitalize="off" />
              <p v-if="errors.password" class="form__error-text">{{ t('validation.' + errors.password, { field: t('password') }) }}</p>
            </div>
          </div>
          <div class="form__row">
            <div class="form__col">
              <label class="form__label">{{ t('confirm_password') }}</label>
              <input class="form__input" type="password" v-model="confirmPassword" :class="{ 'form__input--error': errors.confirm_password }" autocapitalize="off" />
              <p v-if="errors.confirm_password" class="form__error-text">{{ t('validation.' + errors.confirm_password, { field: t('confirm_password') }) }}</p>
            </div>
          </div>
          <button class="button button--primary">{{ t('register') }}</button>
          <p class="form__text">
            {{ $t('accept_policy_text') }}
          </p>
          <p class="form__text">
            {{ $t('already_have_an_account') }} <router-link class="form__link" to="/login">{{ $t('login') }}</router-link>
          </p>
          <p class="form__text">
            <router-link class="form__link" to="/privacy-policy">Privacy Policy</router-link>
          </p>
        </form>
      </div>
    </div>
  </main>
</template>

<script>
  import axios from 'axios';
  import { useCookies } from "vue3-cookies";
  import authService from '../services/authService';
  import { toast } from 'vue3-toastify';
  import { useI18n } from 'vue-i18n';

  export default {
    setup() {
      const { cookies } = useCookies();
      const { t } = useI18n();
      return { cookies, t };
    },
    data() {
      return {
        name: null,
        surname: null,
        email: null,
        phone: null,
        password: null,
        confirmPassword: null,
        errors: {},
      }
    },
    methods: {
      async register() {
        try {
          this.errors = {};
          const response = await authService.register(
            this.name, this.surname, this.email,
            this.phone,
            this.password, this.confirmPassword
          ); 

          if (!response.data.success) {
            toast.error(response.data.message);
            return;
          }

          this.$router.push('/login?msg=success.account_created');
        } catch(error) {
          if (!error.response || !error.response.data || !error.response.data.errors) {
            toast.error(error.message);
            return;
          }

          this.errors = error.response.data.errors;
        }
      }
    }
  }
</script>
