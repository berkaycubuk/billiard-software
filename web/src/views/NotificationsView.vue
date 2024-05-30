<script>
import HeaderSimple from '../components/HeaderSimple.vue';
import notificationService from '../services/notificationService';
import moment from 'moment';
import BaseLayout from '../components/layout/BaseLayout.vue';

export default {
  setup() {
    return { moment };
  },
  components: {
    HeaderSimple,
    BaseLayout,
  },
  data() {
    return {
      notifications: [],
    }
  },
  methods: {
    genNot() {
      notificationService.genNot()
        .then(() => {
          this.fetchNotifications();
        });
    },
    fetchNotifications() {
      notificationService.active()
        .then((res) => {
          this.notifications = res;
        });
    },
    markRead(id) {
      notificationService.mark(id)
        .then(() => {
          this.fetchNotifications();
        });
    },
  },
  mounted() {
    this.fetchNotifications();
  },
}
</script>

<template>
  <BaseLayout>
  <HeaderSimple :title="$t('notifications')" backUrl="/" />
  <main class="container" style="flex: 1;">

    <br />

    <v-list lines="three" v-if="notifications.length != 0">
      <v-list-item
        v-for="notification in notifications"
        :key="notification.id"
        :title="moment(notification.created_at).format('DD.MM.YYYY - HH:mm')"
      >
        <template v-slot:subtitle>
          <v-dialog>
            <template v-slot:activator="{ props }">
              <div v-bind="props" class="notification-content" v-html="notification.message"></div>
            </template>

            <template v-slot:default="{ isActive }">
              <v-card>
                <v-card-text>
                  <div class="notification-content" v-html="notification.message"></div>
                </v-card-text>

                <v-card-actions>
                  <v-btn @click="isActive.value = false">Close</v-btn>
                </v-card-actions>
              </v-card>
            </template>
          </v-dialog>
        </template>
        <template v-slot:append>
          <v-btn @click="() => markRead(notification.id)" variant="text" icon="mdi-check"></v-btn>
        </template>
      </v-list-item>
    </v-list>

    <v-alert v-if="notifications.length === 0" :text="$t('no_notifications')"></v-alert>
  </main>
  </BaseLayout>
</template>

<style>
.notification-content a {
  color: #d6ffe2;
}
</style>
