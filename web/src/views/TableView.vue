<script>
import Header from '../components/Header.vue';
import TableCard from '../components/TableCard.vue';
import tableService from '../services/tableService';
import { useRoute } from 'vue-router';
import { store } from '../store';
import moment from 'moment';
import Popup from '../components/Popup.vue';
import JoinGameAsGuestPopup from '../components/JoinGameAsGuestPopup.vue';
import TransferPopup from '../components/TransferPopup.vue';
import BaseLayout from '../components/layout/BaseLayout.vue';

function checkCanLeave(table) {
  if (table === null || table.game === null || table.game.game_users === null || table.game.game_users.length === 0)
    return false;

  const users = table.game.game_users;
  const found = users.find((user) => user.user_id === store.user.id);
  return found != null;
}

function isPaused(table) {
  if (table === null || table.game === null || table.game.game_users === null || table.game.game_users.length === 0)
    return false;

  const users = table.game.game_users;
  const found = users.find((user) => user.user_id === store.user.id);
  if (found == null) {
    return false;
  }

  return found.status === 2;
}

export default {
  setup() {
    const route = useRoute();

    return { route, moment, store };
  },
  components: {
    Header,
    Popup,
    TableCard,
    JoinGameAsGuestPopup,
    TransferPopup,
    BaseLayout,
  },
  data() {
    return {
      table: null,
      tableId: null,
      tables: [],
      canLeave: false,
      isPaused: false,
      moveTablePopup: false,
      leaveGamePopup: false,
      joinGamePopup: false,
      joinGameGuestPopup: false,
      pauseGamePopup: false,
      unpauseGamePopup: false,
      transferPopup: false,
      leaveGameSelectedUser: null,
      pauseSelectedUser: null,
      unpauseSelectedUser: null,
    }
  },
  created() {
    this.fetchTable();
    this.dataUpdateInterval = setInterval(this.fetchTable, 5000);
    this.fetchTables();
  },
  beforeDestroy() {
    clearInterval(this.dataUpdateInterval);
  },
  beforeRouteLeave (to, from, next) {
    clearInterval(this.dataUpdateInterval);
    next();
  },
  methods: {
    isGuest() {
      return store.user && store.user.role_id === 4;
    },
    isEmpty() {
      if (this.table === null) {
        return true;
      }

      if (this.table.game_users === null || this.table.game_users === []) {
        return true;
      }

      return this.table.game_users.length === 0;
    },
    historyString(action, name) {
      let string;
      switch (action) {
        case 1:
          string = "joined the game."
          break;
        case 2:
          string = "left the game."
          break;
        case 3:
          string = "paused the game."
          break;
        case 4:
          string = "unpaused the game."
          break;
        case 5:
          string = "transferred to other table."
          break;
      }

      return name + " " + string;
    },
    fetchTables() {
      try {
        tableService.all()
          .then((response) => {
            if (!response.data.success) {
              return;
            }
            this.tables = response.data.tables;
          });
      } catch(error) {
        console.error(error);
      }
    },
    fetchTable() {
      try {
        tableService.get(this.route.params.id)
          .then((response) => {
            if (!response.data.success) {
              return;
            }
            this.table = response.data.table;

            this.canLeave = checkCanLeave(this.table);
            this.isPaused = isPaused(this.table);
          });
      } catch(error) {
        console.error(error);
      }
    },
    transferPopupOpen() {
      this.transferPopup = true;
    },
    transferPopupClose() {
      this.transferPopup = false;
    },
    pauseGame() {
      tableService.pause(this.table.id)
        .then((res) => {
          this.$router.push('/')
          /*
          this.fetchTable();
          this.pauseGamePopup = false;
          */
        });
    },
    pauseGameAsGuest() {
      tableService.pauseAsGuest(this.table.id, this.pauseSelectedUser)
        .then((res) => {
          this.$router.push('/')
          /*
          this.fetchTable();
          this.pauseSelectedUser = null;
          this.pauseGamePopup = false;
          */
        });
    },
    unpauseGame() {
      tableService.unpause(this.table.id)
        .then((res) => {
          this.$router.push('/')
          /*
          this.fetchTable();
          this.unpauseGamePopup = false;
          */
        });
    },
    unpauseGameAsGuest() {
      tableService.unpauseAsGuest(this.table.id, this.unpauseSelectedUser)
        .then((res) => {
          this.$router.push('/')
          /*
          this.fetchTable();
          this.unpauseSelectedUser = null;
          this.unpauseGamePopup = false;
          */
        });
    },
    joinGame() {
      tableService.join(this.table.id)
        .then((res) => {
          this.$router.push('/')
          /*
          this.fetchTable();
          this.joinGamePopup = false;
          */
        });
    },
    unpauseGamePopupOpen() {
      this.unpauseGamePopup = true;
    },
    unpauseGamePopupClose() {
      this.unpauseGamePopup = false;
    },
    pauseGamePopupOpen() {
      this.pauseGamePopup = true;
    },
    pauseGamePopupClose() {
      this.pauseGamePopup = false;
    },
    joinGamePopupOpen() {
      this.joinGamePopup = true;
    },
    joinGamePopupClose() {
      this.joinGamePopup = false;
    },
    joinGameGuestPopupOpen() {
      this.joinGameGuestPopup = true;
    },
    joinGameGuestPopupClose() {
      this.joinGameGuestPopup = false;
    },
    leaveGamePopupOpen() {
      this.leaveGamePopup = true;
    },
    leaveGamePopupClose() {
      this.leaveGameSelectedUser = null;
      this.leaveGamePopup = false;
    },
    leaveGame() {
      tableService.leave(this.table.id)
        .then((res) => {
          if (res.data.success === false) {
            this.fetchTable();
            toast.error(res.data.message);
            return;
          }

          this.$router.push(`/payment/overview/${res.data.order.id}`);
        });
    },
    joinAsGuestSuccess() {
      this.$router.push('/');
      /*
      this.fetchTable();
      this.joinGameGuestPopupClose();
      */
    },
    leaveGameAsGuest() {
      tableService.leaveAsGuest(this.table.id, this.leaveGameSelectedUser)
        .then((res) => {
          if (res.data.success === false) {
            this.fetchTable();
            toast.error(res.data.message);
            return;
          }

          this.$router.push(`/payment/overview/${res.data.order.id}`);
        });
    }
  },
}
</script>

<template>
  <BaseLayout>
  <Header />

  <Popup v-if="joinGamePopup" :title="$t('join_game')" @closed="joinGamePopupClose">
    <p></p>
    <div class="complete-popup__buttons">
      <button class="button button--primary" @click="joinGame">{{ $t('join_game') }}</button>
      <button class="button button--secondary" @click="joinGamePopupClose">{{ $t('cancel') }}</button>
    </div>
  </Popup>

  <Popup v-if="leaveGamePopup" :title="$t('leave_game')" @closed="leaveGamePopupClose">
  <p>{{ $t('leave_game_popup_text') }} <b>{{ $t('this_action_is_not_reversable') }}</b></p>
    <br />

    <div v-if="isGuest()">
      <div class="user-selector">

        <div class="user-selector__item" :class="{ 'user-selector__item--selected': user.name === leaveGameSelectedUser }" v-for="user of table.game_users.filter((i) => i.user_id === null)" @click="() => leaveGameSelectedUser = (leaveGameSelectedUser === user.name ? null : user.name)">
          {{ user.name }}
        </div>
        
      </div>
      <br />
      <div class="complete-popup__buttons">
        <button class="button button--primary" @click="leaveGameAsGuest" :disabled="leaveGameSelectedUser === null">{{ $t('leave_game') }}</button>
        <button class="button button--secondary" @click="leaveGamePopupClose">{{ $t('cancel') }}</button>
      </div>
    </div>

    <div v-if="!isGuest()">
      <div class="complete-popup__buttons">
        <button class="button button--primary" @click="leaveGame">{{ $t('leave_game') }}</button>
        <button class="button button--secondary" @click="leaveGamePopupClose">{{ $t('cancel') }}</button>
      </div>
    </div>
  </Popup>

  <Popup v-if="joinGameGuestPopup" :title="$t('join_game_as_guest')" @closed="joinGameGuestPopupClose">
    <JoinGameAsGuestPopup :tableID="parseInt(route.params.id)" @success="joinAsGuestSuccess" @closed="joinGameGuestPopupClose"/>
  </Popup>

  <Popup v-if="pauseGamePopup" :title="$t('pause_game')" @closed="pauseGamePopupClose">

    <div v-if="isGuest()">
      <div class="user-selector">

        <div class="user-selector__item" :class="{ 'user-selector__item--selected': user.name === pauseSelectedUser }" v-for="user of table.game_users.filter((i) => i.status === 1 && i.user_id === null)" @click="() => pauseSelectedUser = (pauseSelectedUser === user.name ? null : user.name)">
          {{ user.name }}
        </div>
        
      </div>
      <br />
      <div class="complete-popup__buttons">
        <button class="button button--primary" @click="pauseGameAsGuest" :disabled="pauseSelectedUser === null">{{ $t('pause_game') }}</button>
        <button class="button button--secondary" @click="pauseGamePopupClose">{{ $t('cancel') }}</button>
      </div>
    </div>

    <div v-if="!isGuest()" class="complete-popup__buttons">
      <button class="button button--primary" @click="pauseGame">{{ $t('pause_game') }}</button>
      <button class="button button--secondary" @click="pauseGamePopupClose">{{ $t('cancel') }}</button>
    </div>
  </Popup>

  <Popup v-if="unpauseGamePopup" :title="$t('unpause_game')" @closed="unpauseGamePopupClose">
    <div v-if="isGuest()">
      <div class="user-selector">

        <div class="user-selector__item" :class="{ 'user-selector__item--selected': user.name === unpauseSelectedUser }" v-for="user of table.game_users.filter((i) => i.status === 2 && i.user_id === null)" @click="() => unpauseSelectedUser = (unpauseSelectedUser === user.name ? null : user.name)">
          {{ user.name }}
        </div>
        
      </div>
      <br />
      <div class="complete-popup__buttons">
        <button class="button button--primary" @click="unpauseGameAsGuest" :disabled="unpauseSelectedUser === null">{{ $t('unpause_game') }}</button>
        <button class="button button--secondary" @click="unpauseGamePopupClose">{{ $t('cancel') }}</button>
      </div>
    </div>

    <div v-if="!isGuest()" class="complete-popup__buttons">
      <button class="button button--primary" @click="unpauseGame">{{ $t('unpause_game') }}</button>
      <button class="button button--secondary" @click="unpauseGamePopupClose">{{ $t('cancel') }}</button>
    </div>
  </Popup>

  <TransferPopup :isGuest="isGuest()" :currentTable="table.id" :guests="table.game_users.filter((i) => i.user_id === null)" :tables="tables" v-if="transferPopup" @closed="() => {transferPopup = false; fetchTable();}" />

  <main>
    <div class="container">
      <div class="table-page">
        <div class="table-page__table">
          <template v-if="table !== null">
            <TableCard :title="table.name" :game="table.game" :gameUsers="table.game_users" />
          </template>
        </div>
        <div class="table-page__actions">
          <h2>{{ $t('actions_title') }}</h2>
          <div class="actions">

            <button v-if="!isGuest() && !canLeave" @click="joinGamePopupOpen" class="button button--primary">{{ $t('join_game') }}</button>
            <button v-if="isGuest()" @click="joinGameGuestPopupOpen" class="button button--primary">{{ $t('join_game_as_guest') }}</button>

            <button v-if="!isGuest() && canLeave" @click="leaveGamePopupOpen" class="button button--red">{{ $t('leave_game') }}</button>
            <button v-if="isGuest() && !isEmpty()" @click="leaveGamePopupOpen" class="button button--red">{{ $t('leave_game') }}</button>

            <button v-if="isGuest() && !isEmpty()" @click="pauseGamePopupOpen" class="button button--blue">{{ $t('pause_game') }}</button>
            <button v-if="!isGuest() && canLeave && !isPaused" @click="pauseGamePopupOpen" class="button button--blue">{{ $t('pause_game') }}</button>

            <button v-if="isGuest() && !isEmpty()" @click="unpauseGamePopupOpen" class="button button--blue">{{ $t('unpause_game') }}</button>
            <button v-if="!isGuest() && canLeave && isPaused" @click="unpauseGamePopupOpen" class="button button--blue">{{ $t('unpause_game') }}</button>

            <button v-if="isGuest() && !isEmpty()" @click="transferPopupOpen" class="button button--orange">{{ $t('transfer_to_other_table') }}</button>
            <button v-if="!isGuest() && canLeave" @click="transferPopupOpen" class="button button--orange">{{ $t('transfer_to_other_table') }}</button>

            <router-link to="/" class="button button--secondary">{{ $t('go_back') }}</router-link>
          </div>
          <!--
          <h2>{{ $t('history') }}</h2>
          <div class="history">
            <ul v-if="table && table.game && table.game.game_histories && table.game.game_histories.length">
              <li v-for="history of table.game.game_histories.slice().reverse()"><b>{{ moment(history.created_at).format('HH:mm') }} :</b> {{ historyString(history.action, history.game_user.name) }}</li>
            </ul>
          </div>
          -->
        </div>
      </div>
    </div>
    <div style="height: 200px;"></div>
  </main>
  </BaseLayout>
</template>

<style scoped>
.table-page {
  margin: 10px 0;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

@media only screen and (max-width: 1200px) {
  .table-page {
    grid-template-columns: 1fr;
  }
}

.table-page__actions {
}

.table-page__actions h2 {
  font-size: 1.2rem;
  font-weight: 700;
  margin-bottom: 10px;
  color: #8A8B85;
}

.actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 20px;
}

.history {
  min-height: 83px;
  max-height: 300px;
  overflow-y: auto;
  background-color: #191919;
  border-radius: 8px;
}

.history ul {
  list-style: none;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.history ul li {
  font-size: 1rem;
}

.button--secondary {
  background-color: #191919;
}

.user-selector {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
}

.user-selector__item {
  cursor: pointer;
  padding: 14px;
  background-color: #232323;
  border-radius: 8px;
}

.user-selector__item--selected {
  outline: 2px solid var(--green);
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
