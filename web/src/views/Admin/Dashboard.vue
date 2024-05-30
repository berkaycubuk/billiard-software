<script>
import AdminLayout from '../../components/layout/AdminLayout.vue';
import dashboardService from '../../services/dashboardService';
import Card from '../../components/Card.vue';
import ConfirmPopup from '../../components/ConfirmPopup.vue';
import { toast } from 'vue3-toastify';
import moment from 'moment';
import notificationService from '../../services/notificationService';

export default {
  components: {
    AdminLayout,
    Card,
    ConfirmPopup,
  },
  setup() {
    return { moment };
  },
  data() {
    return {
      orders: [],
      currency: import.meta.env.VITE_CURRENCY,
      games: [],
      gameDetailsPopup: false,
      selectedGame: null,
      kickUserPopup: false,
      pauseUserPopup: false,
      selectedUser: null,
      endUserPopup: false,
      notiConfirmPopup: false,
      notiMessage: null,
      messageLinkInput: null,
      messageLinkTextInput: null,
    }
  },
  methods: {
      convertCurrency() {
        if (!this.currency) return '';

        if (this.currency == 'kr') return 'kr';

        if (this.currency == 'eur') return '€';

        if (this.currency == 'usd') return '$';

        if (this.currency == 'try') return '₺';

        return this.currency;
      },
    gameClick(game) {
      this.selectedGame = game;
      this.gameDetailsPopup = true;
    },
    kickUserClick(id) {
      this.selectedUser = id;
      this.kickUserPopup = true;
    },
    kickUserConfirm() {
      dashboardService.kickUser(this.selectedUser)
        .then(() => {
          this.fetchGames()
          this.kickUserPopup = false
          this.gameDetailsPopup = false
        })
    },
    pauseUserClick(id) {
      this.selectedUser = id;
      this.pauseUserPopup = true;
    },
    pauseUserConfirm() {
      this.pauseUserPopup = false;
      dashboardService.pauseUnpause(this.selectedUser)
        .then(() => {
          this.gameDetailsPopup = false;
          this.fetchGames();
        });
    },
    endUserClick(id) {
      this.selectedUser = id;
      this.endUserPopup = true;
    },
    endUserConfirm() {
      this.endUserPopup = false;
      dashboardService.leaveTable(this.selectedUser)
        .then(() => {
          toast.success("Order created for the user.");
          this.fetchGames();
          this.endUserPopup = false;
          this.gameDetailsPopup = false;
        });
    },
    publishNoti() {
      notificationService.create(this.notiMessage)
        .then(() => {
          this.notiMessage = null
          this.notiConfirmPopup = false
        })
    },
    fetchGames() {
      dashboardService.activeGames()
        .then((res) => {
          this.games = res.data.games;
        });
    },
    fetchOrders() {
      dashboardService.activeOrders()
        .then((res) => {
          this.orders = res.data.data.orders;
        });
    },
    appendLinkToMessage() {
      if (this.messageLinkInput == null || this.messageLinkInput == "") return

      const anchor_text = this.messageLinkTextInput != null && this.messageLinkTextInput != "" ? this.messageLinkTextInput : this.messageLinkInput

      const anchor_tag = `<a href="${this.messageLinkInput}" target="_blank">${anchor_text}</a>`

      if (this.notiMessage == null) {
        this.notiMessage = anchor_tag
      } else {
        this.notiMessage = this.notiMessage + anchor_tag
      }

      this.messageLinkInput = null
      this.messageLinkTextInput = null
    },
  },
  mounted() {
    this.fetchGames();
    this.fetchOrders();
  }
}
</script>

<template>
  <AdminLayout>

    <v-dialog v-model="kickUserPopup" persistent>
      <ConfirmPopup @cancel="kickUserPopup = false" @confirm="kickUserConfirm" />
    </v-dialog>

    <v-dialog v-model="pauseUserPopup" persistent>
      <ConfirmPopup @cancel="pauseUserPopup = false" @confirm="pauseUserConfirm" />
    </v-dialog>

    <v-dialog v-model="endUserPopup" persistent>
      <ConfirmPopup @cancel="endUserPopup = false" @confirm="endUserConfirm" />
    </v-dialog>

    <v-dialog v-model="gameDetailsPopup" persistent>
      <v-card :title="$t('game') + ' #' + selectedGame.id" style="width: 100%; margin: 0 auto; max-width: 1200px;" class="max-w-xl w-full mx-auto">
        <v-container>
          <h3>{{ $t('players') }}</h3>

          <template v-for="user in selectedGame.game_users">
            <v-card :title="user.name" v-if="user && user.ended_at == null">
              <v-list>
                <v-list-item
                  :title="$t('started_at') + ': ' + user.started_at"
                ></v-list-item>
                <v-list-item
                  :title="$t('status') + ': ' + (user.status == 1 ? $t('playing') : $t('paused'))"
                ></v-list-item>
              </v-list>

              <v-divider></v-divider>

              <v-card-actions>
                <v-btn variant="flat" color="green-darken-2" @click="() => endUserClick(user.id)">{{ $t('end_players_game') }}</v-btn>
                <v-btn variant="flat" color="indigo-darken-3" @click="() => pauseUserClick(user.id)">{{ user.status == 1 ? $t('pause') : $t('unpause') }}</v-btn>
                <v-btn variant="flat" color="red-darken-2" @click="() => kickUserClick(user.id)">{{ $t('kick') }}</v-btn>
              </v-card-actions>
            </v-card>
          </template>
        </v-container>

        <v-card-actions>
          <v-btn variant="flat" @click="gameDetailsPopup = false">{{ $t('close') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="notiConfirmPopup">
      <v-card style="max-width: 600px; width: 100%; margin: 0 auto;" :title="$t('are_you_sure')">
        <v-card-text>
          <p>{{ $t('notification_publish_text') }}</p>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="teal" variant="flat" @click="publishNoti">{{ $t('publish') }}</v-btn>
          <v-btn @click="notiConfirmPopup = false">{{ $t('cancel') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <br/>
    <Card :title="$t('publish_new_notification')">
      <div class="form">

        <div class="form__row">
          <div class="form__col">
            <label class="form__label">{{ $t('message') }}</label>
            <textarea v-model="notiMessage" class="form__input"></textarea>
          </div>
        </div>

        <div>
          <label>{{ $t('link') }}</label>
          <input type="text" class="form__input" v-model="messageLinkInput" placeholder="https://google.com" />
          <label>{{ $t('text') }}</label>
          <input type="text" class="form__input" v-model="messageLinkTextInput" placeholder="Google" />
          <button class="button button--secondary" @click="appendLinkToMessage">{{ $t('add_link') }}</button>
        </div>

        <div>
          <button @click="() => notiMessage != null ? notiConfirmPopup = true : null" style="padding: 8px 14px;" class="button button--primary">{{ $t('publish') }}</button>
        </div>

      </div> 
    </Card>
    <br/>

    <div style="display: flex; flex-direction: column; gap: 10px;">
      <Card :title="$t('active_games')">

        <div class="game">

          <v-list>
            <v-list-item
              v-for="game in games"
              :title="game.table_id"
              :subtitle="'Game #' + game.id"
              @click="() => gameClick(game)"
            >
            </v-list-item>
          </v-list>
          
        </div>

      </Card>

      <Card :title="$t('active_orders')">

          <v-list lines="two">
            <v-list-item
              v-for="order in orders"
              :title="$t('order') + ' #' + order.id + ' - ' + order.price + convertCurrency()"
              :subtitle="order.detail.user_name + ' ' + order.detail.user_surname + ' - ' + moment(order.created_at).format('DD.MM.YYYY - HH:mm')"
              @click="() => $router.push('/admin/orders/' + order.id)"
            >
            </v-list-item>
          </v-list>

      </Card>
    </div>

  </AdminLayout>
</template>

<style scoped>

</style>
