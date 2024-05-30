<script>
import Header from '../../components/Header.vue';
import BottomNav from '../../components/BottomNav.vue';

export default {
  components: {
    Header,
    BottomNav,
  }
}
</script>

<template>
  <Header />
  <div class="admin-layout">
    <div class="container">

      <nav class="admin__nav">
        <router-link to="/admin">{{ $t('dashboard') }}</router-link>
        <router-link to="/admin/reports">{{ $t('reports') }}</router-link>
        <router-link to="/admin/users">{{ $t('users') }}</router-link>
        <router-link to="/admin/roles">{{ $t('roles') }}</router-link>
        <router-link to="/admin/tables">{{ $t('tables') }}</router-link>
        <router-link to="/admin/products">{{ $t('products') }}</router-link>
        <router-link to="/admin/pricing">{{ $t('pricing') }}</router-link>
        <router-link to="/admin/orders">{{ $t('orders') }}</router-link>
        <router-link to="/admin/subscriptions">{{ $t('subscriptions') }}</router-link>
      </nav>

      <main class="admin-layout__main">
        <slot></slot> 
        <div style="height: 80px;"></div>
      </main>
    </div>
  </div>
  <BottomNav />
</template>

<style scoped>
.admin__nav {
  display: flex;
  gap: 20px;
  padding: 16px 20px;
  border-radius: 8px;
  background-color: #191919;
  margin: 10px 0;
  overflow-x: auto;
}

.admin__nav a {
  color: var(--white);
  text-decoration: none;
  width: fit-content;
  font-weight: 600;
}
</style>
