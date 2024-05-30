<script>
  import notificationService from '../services/notificationService';
  import moment from 'moment';

  export default {
    data() {
      return {
        notifications: [],
      }
    },
    setup() {
      return { moment }
    },
    methods: {
      fetchNotifications() {
        notificationService.game()
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
    }
  }
</script>

<template>
  <v-card style="max-width: 600px; width: 100%; margin: 0 auto;">
    <v-card-text v-if="notifications.length > 0">

      <div style="font-size: 20px; font-weight: bold;">{{ notifications[0].message }}</div>
      <br />
      <p>This notification will continue to show up while you're playing.</p>

    </v-card-text>

    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn @click="() => $emit('close')">Close</v-btn>
    </v-card-actions>
  </v-card>
</template>
