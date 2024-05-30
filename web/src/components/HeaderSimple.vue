<script>
export default {
  props: {
    title: String,
    backUrl: String,
  }
}
</script>

<template>
  <header class="profile-header">
    <div class="container">
      <div class="profile-header__inner">
        <div class="profile-header__left">
          <router-link class="profile-header__go-back" :to="backUrl" v-if="backUrl">
            <svg width="30px" height="30px" viewBox="0 0 24 24" stroke-width="1.5" fill="none" xmlns="http://www.w3.org/2000/svg" color="#f8faed"><path d="M21 12H3m0 0l8.5-8.5M3 12l8.5 8.5" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
          </router-link>
          <div class="profile-header__go-back" @click="() => $router.back()" v-if="!backUrl">
            <svg width="30px" height="30px" viewBox="0 0 24 24" stroke-width="1.5" fill="none" xmlns="http://www.w3.org/2000/svg" color="#f8faed"><path d="M21 12H3m0 0l8.5-8.5M3 12l8.5 8.5" stroke="#f8faed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>
          </div>
          <div class="profile-header__name">{{ title }}</div>
        </div>
        <div class="profile-header__right">

        </div>
      </div>
    </div>
  </header>
</template>

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
  cursor: pointer;
}

.profile-header__name {
  font-family: 'Bricolage Grotesque', sans-serif;
  font-size: 1.8rem;
  font-weight: 700;
}

@media only screen and (max-width: 600px) {
  .profile-header__inner {
    height: 50px;
  }

  .profile-header__go-back {
    width: 24px;
    height: 24px;
  }

  .profile-header__go-back svg {
    width: 24px;
    height: 24px;
  }

  .profile-header__left {
    flex-direction: row;
    align-items: center;
  }

  .profile-header__name {
    font-size: 1.4rem;
  }
}
</style>
