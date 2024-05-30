<script>
import Popup from '../Popup.vue';
import userService from '../../services/userService';

export default {
    props: ['dates', 'method', 'user', 'zeroOrders'],
    components: {
        Popup,
    },
    data() {
      return {
        form: {
          startDate: this.$route.query.start,
          endDate: this.$route.query.end,
          method: this.$route.query.method,
          user: this.$route.query.user,
          zeroOrders: this.$route.query.zeroOrders == "true" ? true : false
        },
        users: [],
        showDates: true,
      };
    },
    mounted() {
      /*
      if (this.user) {
        this.fetchUsers();
      }
      */
      
      if (this.dates != null && this.dates == false) {
        this.showDates = false;
      }
    },
    methods: {
      fetchUsers() {
        userService.allWithoutPagination()
          .then((res) => {
            this.users = res.users
          })
      },
        onApply() {
            const currentQuery = { ...this.$route.query };

            if (this.form.startDate != "") {
                currentQuery.start = this.form.startDate;
            } else {
                delete currentQuery.start;
            }

            if (this.form.endDate != "") {
                currentQuery.end = this.form.endDate;
            } else {
                delete currentQuery.end;
            }

            if (this.form.user != "") {
                currentQuery.user = this.form.user;
            } else {
                delete currentQuery.user;
            }

            if (this.form.method != "") {
                currentQuery.method = this.form.method;
            } else {
                delete currentQuery.method;
            }

            if (this.form.zeroOrders != false) {
              currentQuery.zeroOrders = true
            } else {
              delete currentQuery.zeroOrders
            }

            currentQuery.page = 1

            this.$router.push({ path: this.$route.path, query: currentQuery });

            this.$emit('saved');
        },
        onClear() {
            this.$router.push({ path: this.$route.path, query: {} });
            this.$emit('saved');
        }
    },
}
</script>

<template>
  <Popup :title="$t('filter')" @closed="$emit('closed')">
    <form class="form">

      <div class="form__row" v-if="showDates">
        <div class="form__col">
          <label class="form__label">{{ $t('start_date') }}</label>
          <input type="date" class="form__input" v-model="form.startDate" />
        </div>
        <div class="form__col">
          <label class="form__label">{{ $t('end_date') }}</label>
          <input type="date" class="form__input" v-model="form.endDate" />
        </div>
      </div>

      <div class="form__row" v-if="user">
        <div class="form__col">
          <label class="form__label">{{ $t('user') }}</label>
          <input type="text" class="form__input" v-model="form.user" />
          <!--
          <select class="form__input" v-model="form.user">
            <option v-for="user of users" :value="user.id">{{ user.name }} {{ user.surname }}</option>
          </select>
          -->
        </div>
      </div>

      <div class="form__row" v-if="method">
        <div class="form__col">
          <label class="form__label">{{ $t('method') }}</label>
          <select class="form__input" v-model="form.method">
            <option value="">{{ $t('all') }}</option>
            <option value="1">{{ $t('physical') }}</option>
            <option value="2">{{ $t('vipps') }}</option>
            <option value="3">{{ $t('system') }}</option>
          </select>
        </div>
      </div>

      <div style="display: flex; align-items: center; gap: 10px;" v-if="zeroOrders">
        <v-checkbox
          label="Show zero orders"
          v-model="form.zeroOrders"
          hide-details></v-checkbox>
      </div>

    </form>

    <br />

    <div class="grid grid-cols-2">
      <button class="button button--primary" @click="onApply">{{ $t('apply') }}</button>
      <button class="button button--red" @click="onClear">{{ $t('clear') }}</button>
    </div>
  </Popup>
</template>
