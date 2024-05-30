<script>
import permissionService from '../services/permissionService';
import { store } from '../store'

export default {
  data() {
    return {
      permissions: {
        accessAdmin: false,
      },  
    }
  },
  methods: {
		  canSeeHome() {
					if (store.user != null && store.user.hasOwnProperty('role_id') && store.user.role_id == 3) {
							return false
					}

				  return true
		  },
  },
  mounted() {
    permissionService.permission('access_admin')
      .then((res) => {
        this.permissions.accessAdmin = res; 
      }); 
  },
}
</script>

<template>
  <div class="bottom-nav" v-if="canSeeHome() == true">
    <v-bottom-navigation>
      <v-btn @click="() => $router.push('/')">
        <v-icon>mdi-home-variant</v-icon>
        <span>Home</span>
      </v-btn>
      <v-btn @click="() => $router.push('/shop')">
        <v-icon>mdi-cart</v-icon>
        <span>Shop</span>
      </v-btn>
      <v-btn v-if="permissions.accessAdmin" @click="() => $router.push('/admin')">
        <v-icon>mdi-security</v-icon>
        <span>Admin</span>
      </v-btn>
    </v-bottom-navigation>
    <!--
    <div class="nav__inner">
      <router-link to="/">{{ $t('home') }}</router-link>
      <router-link to="/shop">{{ $t('shop') }}</router-link>
      <router-link to="/admin" v-if="permissions.accessAdmin">{{ $t('admin') }}</router-link>
    </div>    
    -->
  </div>
</template>

<style scoped>
.bottom-nav {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  background-color: #191919;
}

.nav__inner {
  padding: 16px;
  display: flex;
  gap: 30px;
  align-items: center;
  justify-content: center;
}

.bottom-nav a {
  text-decoration: none;
  font-size: 1.2rem;
  font-weight: 600;
  font-family: 'Bricolage Grotesque', sans-serif;
  color: var(--white);
}

@media only screen and (max-width: 600px) {
  .bottom-nav {
    display: block;
  }
}
</style>
