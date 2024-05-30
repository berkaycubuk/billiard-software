<template>
  <BaseLayout>
  <HeaderSimple :title="`${user.name} ${user.surname}`" backUrl="/" />
  <main class="container" style="flex: 1;">

    <section class="profile-section">

      <h2 class="profile-section__title">{{ $t('account') }}</h2>

      <nav>
        <!--
        <router-link to="/profile/account">
          <svg width="24px" height="24px" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#f8faed"><path d="M5 20v-1a7 7 0 017-7v0a7 7 0 017 7v1M12 12a4 4 0 100-8 4 4 0 000 8z" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
          {{ $t('my_information') }}
        </router-link>
        -->
        <router-link to="/profile/subscriptions">
          <svg width="24px" height="24px" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#f8faed"><path d="M3 20.4V3.6a.6.6 0 01.6-.6h16.8a.6.6 0 01.6.6v16.8a.6.6 0 01-.6.6H3.6a.6.6 0 01-.6-.6z" stroke="#f8faed" stroke-width="1.5"></path><path d="M15 8.5c-.685-.685-1.891-1.161-3-1.191M9 15c.644.86 1.843 1.35 3 1.391m0-9.082c-1.32-.036-2.5.561-2.5 2.191 0 3 5.5 1.5 5.5 4.5 0 1.711-1.464 2.446-3 2.391m0-9.082V5.5m0 10.891V18.5" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
          {{ $t('subscriptions') }}
        </router-link>
        <router-link to="/profile/waiting-orders">
          <svg width="24px" height="24px" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#f8faed"><path d="M19.5 22a1.5 1.5 0 100-3 1.5 1.5 0 000 3zM9.5 22a1.5 1.5 0 100-3 1.5 1.5 0 000 3z" fill="#f8faed" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path><path d="M5 4h17l-2 11H7L5 4zm0 0c-.167-.667-1-2-3-2M20 15H5.23c-1.784 0-2.73.781-2.73 2 0 1.219.946 2 2.73 2H19.5" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
          {{ $t('waiting_orders') }}
        </router-link>
        <router-link to="/profile/order-history">
          <v-icon>mdi-history</v-icon>
          {{ $t('order_history') }}
        </router-link>
      </nav>

    </section>

    <section class="profile-section">

      <h2 class="profile-section__title">{{ $t('security') }}</h2>

      <nav>
        <router-link to="/profile/change-password">
          <svg width="24px" height="24px" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#f8faed"><path d="M16 12h1.4a.6.6 0 01.6.6v6.8a.6.6 0 01-.6.6H6.6a.6.6 0 01-.6-.6v-6.8a.6.6 0 01.6-.6H8m8 0V8c0-1.333-.8-4-4-4S8 6.667 8 8v4m8 0H8" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
          {{ $t('change_password') }}
        </router-link>
      </nav>

    </section>

    <section class="profile-section">

      <h2 class="profile-section__title">Other</h2>

      <nav>
        <button @click.prevent="logout">
          <svg width="24px" height="24px" stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#f8faed"><path d="M12 12h7m0 0l-3 3m3-3l-3-3M19 6V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2h10a2 2 0 002-2v-1" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
          {{ $t('logout') }}
        </button>
      </nav>

    </section>

  </main>
  </BaseLayout>
  <!--
  <Footer />
  -->
</template>

<script>
import Footer from '../components/Footer.vue';
import { useCookies } from 'vue3-cookies';
import userService from '../services/userService';
import { toast } from 'vue3-toastify';
import HeaderSimple from '../components/HeaderSimple.vue';
import BaseLayout from '../components/layout/BaseLayout.vue';

export default {
  setup() {
    const { cookies } = useCookies();
    return { cookies };
  },
  components: {
    HeaderSimple,
    BaseLayout,
  },
  data() {
    return {
      user: {
        id: null,
        name: null,
        surname: null,
        email: null
      },
    }
  },
  mounted() {
    userService.getProfile()
      .then((response) => {
        if (!response.data.success) {
          return;
        }
        this.user = response.data.user; 
      }).catch((err) => {
        // logout if user not found
        if (err.response.data.message === "error.user_not_found" || err.response.data.message === "error.auth_required") {
          this.cookies.remove('token');
          this.$route.push('/login');
          return;
        }
        toast.error(err.response.data.message); 
      }); 
  },
  methods: {
    logout() {
      this.cookies.remove('token');
      this.$router.push('/login');
    }
  }
}
</script>

<style scoped>
.profile-header {
  background-color: #191919;
}

.profile-header__inner {
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.profile-header__left {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 10px;
}

.profile-header__go-back {
  display: inline-block;
  width: fit-content;
}

.profile-header__name {
  font-family: 'Bricolage Grotesque', sans-serif;
  font-size: 1.8rem;
  font-weight: 700;
}

.profile-section {
  margin: 20px 0;
}

.profile-section__title {
  font-size: 1.2rem;
  font-weight: 700;
  margin-bottom: 10px;
  color: #8A8B85;
}

.profile-section nav {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.profile-section nav a {
  padding: 20px 16px;
  display: flex;
  align-items: center;
  gap: 20px;
  color: #f8faed;
  text-decoration: none;
  font-size: 1.2rem;
  font-weight: 600;
  border-radius: 12px;
  background-color: #191919;
  transition: all .2s;
}

.profile-section nav a:hover {
  opacity: .8;
}

.profile-section nav button {
  cursor: pointer;
  padding: 20px 16px;
  display: flex;
  align-items: center;
  gap: 20px;
  color: #f8faed;
  text-decoration: none;
  font-size: 1.2rem;
  font-weight: 600;
  border-radius: 12px;
  background-color: #191919;
  transition: all .2s;
  outline: none;
  border: none;
}

.profile-section nav button:hover {
  opacity: .8;
}
</style>
