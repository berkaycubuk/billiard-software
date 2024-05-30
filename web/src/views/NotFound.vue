<template>
  <main class="container" style="flex: 1;">
    <div class="ne__container">
      <h1 class="ne__title">{{ $t('page_not_found') }}</h1>
      <p class="ne__description"></p>
      <router-link to="/" class="button button--primary">
        {{ $t('go_back') }}
      </router-link>
    </div>
  </main>
</template>

<style scoped>
.ne__container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
}

.ne__title {
  margin-bottom: 20px;
}

.ne__description {
  font-size: 1.2rem;
  margin-bottom: 20px;
}
</style>
