<script>
import Header from '../components/Header.vue';
import TableCard from '../components/TableCard.vue';
import tableService from '../services/tableService';
// import BottomNav from '../components/BottomNav.vue';
import BaseLayout from '../components/layout/BaseLayout.vue';
import NotificationPopup from '../components/NotificationPopup.vue';
import notificationService from '../services/notificationService';
import { store } from '../store'

export default {
  components: {
    Header,
    TableCard,
    BaseLayout,
    NotificationPopup,
    //BottomNav,
  },
  data() {
    return {
      timerId: null,
      tables: [],
      notifications: false,
      intervalTime: 1000,
      dataUpdateInterval: null,
    }
  },
  created() {
    this.fetchData();
    this.dataUpdateInterval = setInterval(this.fetchData, this.intervalTime);
  },
  beforeDestroy() {
    if (this.dataUpdateInterval != null) clearInterval(this.dataUpdateInterval)
  },
  beforeRouteLeave (to, from, next) {
    if (this.dataUpdateInterval != null) clearInterval(this.dataUpdateInterval)
    next();
  },
  methods: {
    fetchData() {
      try {
        tableService.all()
          .then((response) => {
			if (store.user != null && store.user.hasOwnProperty('role_id') && store.user.role_id == 3) {
					return this.$router.push('/shop')
			}
            if (!response.data.success) {
              return;
            }
            this.tables = response.data.tables;

            let found = false
            for (let i = 0; i < this.tables.length; i++) {
              if (this.tables[i].game != null) {
                found = true
                break
              }
            }

            if (!found) {
              if (this.dataUpdateInterval != null) clearInterval(this.dataUpdateInterval) 

              this.intervalTime = 10000

              this.dataUpdateInterval = setInterval(this.fetchData, 10000)
            } else {
              if (this.intervalTime == 10000) {
                if (this.dataUpdateInterval != null) clearInterval(this.dataUpdateInterval) 
                this.intervalTime = 1000
                this.dataUpdateInterval = setInterval(this.fetchData, 1000)
              }
            }
          });
      } catch(error) {
			if (store.user != null && store.user.hasOwnProperty('role_id') && store.user.role_id == 3) {
					return this.$router.push('/shop')
			}
        console.error(error);
      }
    }
  },
}
</script>

<template>
  <BaseLayout>

  <Header />
  <main>
    <div class="container">
      <div class="tables">
        <router-link :to="'/table/' + table.id" v-for="table in tables" style="text-decoration: none;">
          <TableCard :id="table.id" :title="table.name" :game="table.game" :gameUsers="table.game_users" v-if="table.status === 1" />
        </router-link>
      </div>
    </div>
    <div style="height: 200px;"></div>
  </main>
  </BaseLayout>
</template>

<style scoped>
.tables {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin: 10px 0;
}

@media only screen and (max-width: 1200px) {
  .tables {
    grid-template-columns: 1fr;
  }
}
</style>
