<script>
import Popup from './Popup.vue';
import tableService from '../services/tableService';

export default {
  components: {
    Popup,
  },
  props: {
    isGuest: Boolean,
    currentTable: Number,
    tables: [],
    guests: [],
  },
  data() {
    return {
      selectedTable: null,
      guestPage: 0,
      selectedGuest: null,
    }
  },
  methods: {
    acceptTransfer() {
      if (this.selectedTable === null) { return; }

      if (this.isGuest === true) {
        if (this.selectedGuest === null) { return; }

        tableService.transferGuest(this.currentTable, this.selectedTable, this.selectedGuest)
          .then((res) => {
            //this.$router.push("/table/" + this.selectedTable);
            this.$router.push('/');
            this.$emit('closed');
          });
        return;
      }

      tableService.transfer(this.currentTable, this.selectedTable)
        .then((res) => {
          //this.$router.push("/table/" + this.selectedTable);
          this.$router.push('/');
          this.$emit('closed');
        });
    },
  },
}
</script>

<template>
  <Popup :title="$t('transfer_to_other_table')" @closed="$emit('closed')">
    <div v-if="isGuest">

      <div v-if="guestPage === 0">

        <div class="user-selector">

          <div class="user-selector__item" :class="{ 'user-selector__item--selected': user.name === selectedGuest }" v-for="user of guests" @click="() => selectedGuest = (selectedGuest === user.name ? null : user.name)">
            {{ user.name }}
          </div>
          
        </div>
        <br />
        <div class="complete-popup__buttons">
          <button class="button button--primary" @click="() => guestPage = 1" :disabled="selectedGuest === null">{{ $t('next') }}</button>
          <button class="button button--secondary" @click="$emit('closed')">{{ $t('cancel') }}</button>
        </div>

      </div>

      <div v-if="guestPage === 1">

        <div class="user-selector">

          <div class="user-selector__item" :class="{ 'user-selector__item--selected': table.id === selectedTable }" v-for="table of tables" @click="() => selectedTable = (selectedTable === table.id ? null : table.id)">
            {{ table.name }}
          </div>
          
        </div>
        <br />
        <div class="complete-popup__buttons">
          <button class="button button--primary" @click="acceptTransfer" :disabled="selectedTable === null">{{ $t('transfer') }}</button>
          <button class="button button--secondary" @click="() => guestPage = 0">{{ $t('back') }}</button>
        </div>

      </div>

    </div>

    <div v-if="!isGuest">
      <div class="user-selector">

        <div class="user-selector__item" :class="{ 'user-selector__item--selected': table.id === selectedTable }" v-for="table of tables" @click="() => selectedTable = (selectedTable === table.id ? null : table.id)">
          {{ table.name }}
        </div>
        
      </div>
      <br />
      <div class="complete-popup__buttons">
        <button class="button button--primary" @click="acceptTransfer" :disabled="selectedTable === null">{{ $t('transfer') }}</button>
        <button class="button button--secondary" @click="$emit('closed')">{{ $t('cancel') }}</button>
      </div>
    </div>
  </Popup>
</template>

<style scoped>

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
