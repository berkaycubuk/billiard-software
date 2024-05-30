<script>
import TableCardPlayer from '../components/TableCardPlayer.vue';

export default {
  components: {
    TableCardPlayer,
  },
	props: ['id', 'title', 'game', 'gameUsers'],
		/*
  props: {
    id: Number,
    title: String,
    game: {
      id: Number,
      table_id: Number,
      game_users: [],
    },
    gameUsers: [],
  },
  */
}
</script>

<template>
  <div class="table-card" style="background-image: url(https://images.unsplash.com/photo-1589759118394-f5cfe6178fd3?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2070&q=80)">
    <h3 class="table-card__title">{{ title }}</h3>

    <div class="table-card__players" v-if="game && game.game_users.length && gameUsers && gameUsers.length">
      <TableCardPlayer v-for="user of gameUsers" :name="user.name" :time="user.time.toString()" :price="user.price" :status="user.status" />
    </div>
  </div>
</template>

<style scoped>
.table-card {
  color: #f8faed;
  text-decoration: none;
  height: fit-content;
  cursor: pointer;
  padding: 20px;
  background-color: #191919;
  background-size: cover; 
  border-radius: 8px;
  transition: opacity .2s;
}

.table-card:hover {
  opacity: .8;
}

.table-card__title {
  font-size: 1.4rem;
  font-weight: 600;
  font-family: 'Bricolage Grotesque', sans-serif;
  text-align: center;
}

.table-card__players {
  margin-top: 20px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

@media only screen and (max-width: 580px) {
  .table-card__players {
    grid-template-columns: 1fr;
  }
}
</style>
