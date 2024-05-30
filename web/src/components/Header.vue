<script>
import permissionService from '../services/permissionService';
import notificationService from '../services/notificationService';

export default {
  data() {
    return {
      permissions: {
        accessAdmin: false,
        showNotifications: false,
        showProfile: false,
      },  
      notifications: [],
    }
  },
  methods: {
    fetchNotifications() {
      notificationService.active()
        .then((res) => {
          this.notifications = res;
        });
    },
  },
  mounted() {
    permissionService.permission('access_admin')
      .then((res) => {
        this.permissions.accessAdmin = res; 
      }); 
    permissionService.permission('header.show_notifications')
      .then((res) => {
        this.permissions.showNotifications = res; 
      }); 
    permissionService.permission('header.show_profile')
      .then((res) => {
        this.permissions.showProfile = res; 
      }); 
    this.fetchNotifications();
  },
}
</script>

<template>
  <header class="header">
   <div class="container">
    <div class="header__inner">
      <div class="header__left">
        <router-link class="header__logo" to="/">
          <img src="/logo.svg" style="height: 34px;" />
        </router-link>
        <nav class="header__nav">
          <router-link class="header__item" to="/">
            {{ $t('home') }}
          </router-link>
          <router-link class="header__item" to="/shop">
            {{ $t('shop') }}
          </router-link>
          <router-link v-if="permissions.accessAdmin == true" class="header__item" to="/admin">
            {{ $t('admin') }}
          </router-link>
        </nav>
      </div>
      <div class="header__right">
        <router-link v-if="permissions.showNotifications && notifications.length > 0" class="header__button header__button--animate" to="/notifications">
		<v-badge :content="notifications.length" color="red">
          <svg stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#000000"><path d="M18 8.4c0-1.697-.632-3.325-1.757-4.525C15.117 2.675 13.59 2 12 2c-1.591 0-3.117.674-4.243 1.875C6.632 5.075 6 6.703 6 8.4 6 15.867 3 18 3 18h18s-3-2.133-3-9.6zM13.73 21a1.999 1.999 0 01-3.46 0" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
		</v-badge>
        </router-link> 
        <router-link v-if="permissions.showNotifications && notifications.length == 0" class="header__button" to="/notifications">
          <svg stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#000000"><path d="M18 8.4c0-1.697-.632-3.325-1.757-4.525C15.117 2.675 13.59 2 12 2c-1.591 0-3.117.674-4.243 1.875C6.632 5.075 6 6.703 6 8.4 6 15.867 3 18 3 18h18s-3-2.133-3-9.6zM13.73 21a1.999 1.999 0 01-3.46 0" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
        </router-link> 
        <router-link v-if="permissions.showProfile" class="header__button" to="/profile">
          <svg stroke-width="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#000000"><path d="M5 20v-1a7 7 0 017-7v0a7 7 0 017 7v1M12 12a4 4 0 100-8 4 4 0 000 8z" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
        </router-link> 
      </div>
    </div>
   </div>
  </header>
</template>
