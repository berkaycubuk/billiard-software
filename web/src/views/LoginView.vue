<template>
  <main class="container">
    <div class="auth-box-outside">
      <div class="alert alert--info" v-if="message !== null">
        <p class="alert__text">{{ $t(message) }}</p>
      </div>
      <div class="auth-box">

        <h1 class="auth-box__title">{{ $t('title') }}</h1>

        <p class="auth-box__description">
          {{ $t('login_desc') }}
        </p>

        <form @submit.prevent="login" class="auth-box__form form">
          <div class="form__row">
            <div class="form__col">
              <label class="form__label">{{ $t('email') }}</label>
              <input class="form__input" type="text" v-model="email" :class="{ 'form__input--error': errors && errors.hasOwnProperty('email') }" autocapitalize="off" />
              <p v-if="errors.email" class="form__error-text">{{ $t('validation.' + errors.email, { field: $t('email') }) }}</p>
            </div>
          </div>
          <div class="form__row">
            <div class="form__col">
              <label class="form__label">{{ $t('password') }}</label>
              <input class="form__input" type="password" v-model="password" :class="{ 'form__input--error': errors.password }" autocapitalize="off" />
              <p v-if="errors.password" class="form__error-text">{{ $t('validation.' + errors.password, { field: $t('password') }) }}</p>
            </div>
          </div>

		  <div class="form__row">
				  <v-checkbox
					label="Remember me"
					hide-details
					v-model="remember"
				  ></v-checkbox>
		  </div>

          <button class="button button--primary" type="submit">{{ $t('login') }}</button>
          <p class="form__text">
            {{ $t('forgot_your_password') }} <router-link class="form__link" to="/reset-password">{{ $t('reset_password') }}</router-link>
          </p>
          <p class="form__text">
            {{ $t('dont_have_account') }} <router-link class="form__link" to="/register">{{ $t('create_an_account') }}</router-link>
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
  import { useCookies } from "vue3-cookies";
  import authService from "../services/authService"; 
  import { toast } from 'vue3-toastify';
  import Popup from '../components/Popup.vue';

  export default {
    components: {
      Popup
    },
    setup() {
      const { cookies } = useCookies();

      return { cookies };
    },
    data() {
      return {
        email: null,
        password: null,
        errors: {},
        message: null,
		remember: false,
      }
    },
    mounted() {
      if (this.$route.query.msg) {
        this.message = this.$route.query.msg;
      }
    },
    methods: {
      async login() {
        // clear errors
        this.errors = {}

        let validated = true
        if (this.email == null) {
          this.errors['email'] = 'required'
          validated = false
        }

        if (this.password == null) {
          this.errors['password'] = 'required'
          validated = false
        }

        if (!validated) return

        try {
          const response = await authService.login(this.email, this.password, this.remember);

          this.$router.push('/');
        } catch(error) {
          if (!error.response || !error.response.data) {
            toast.error(error.message);
            return;
          }

          if (error.response && error.response.data) {
            //toast.error(error.response.data.message);
            if (error.response.data.hasOwnProperty('errors')) {
              this.errors = error.response.data.errors;
            }
          }
        }
      }
    }
  }
</script>
