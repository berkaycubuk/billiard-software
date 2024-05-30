<script>
import HeaderSimple from '../../components/HeaderSimple.vue';
import userService from '../../services/userService';
import subscriptionService from '../../services/subscriptionService';
import Popup from '../../components/Popup.vue';
import moment from 'moment';
import Card from '../../components/Card.vue';
import UpdateAccountPopup from '../../components/popups/UpdateAccountPopup.vue';
import BaseLayout from '../../components/layout/BaseLayout.vue';

export default {
  components: {
    HeaderSimple,
    Popup,
    BaseLayout,
    Card,
    UpdateAccountPopup
},
  setup() {
    return { moment };
  },
  data() {
    return {
      user: null,
      editPopup: false,
    }
  },
  methods: {
    fetchPage() {
      userService.getProfile()
        .then((response) => {
          this.user = response.data.user;
        });
    },
    onSaved() {
      this.editPopup = false;
      this.fetchPage();
    }
  },
  mounted() {
    this.fetchPage();
  },
}
</script>

<template>
  <BaseLayout>
  <HeaderSimple :title="$t('account')" />

  <UpdateAccountPopup :user="user" v-if="editPopup === true" @closed="() => editPopup = false" @saved="onSaved" />

  <main class="container" style="flex: 1;">

    <div class="cards">

      <Card :title="$t('my_information')" v-if="user != null">

        <div class="grid">

          <div>
            <b>{{ $t('full_name') }}:</b> {{ user.name }} {{  user.surname }}
          </div>

          <div>
            <b>{{ $t('email') }}:</b> {{ user.email }}
          </div>

          <div>
            <b>{{ $t('phone') }}:</b> {{ user.phone }}
          </div>

          <div>
            <b>{{ $t('account_created_at') }}:</b> {{ moment(user.created_at).format('DD.MM.YYYY - HH:mm') }}
          </div>

        </div>

      </Card>

      <Card>
        <div class="grid">
          <button class="button button--primary" @click="() => editPopup = true">{{ $t('update_account') }}</button>
        </div>
      </Card>

    </div>
  
  </main>
  </BaseLayout>
</template>

<style scoped>
.subscriptions {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
}

.complete-popup__buttons {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

@media only screen and (max-width: 650px) {
  .complete-popup__buttons {
    grid-template-columns: 1fr;
  }
}
</style>
