<script>
import Popup from './Popup.vue';
import subscriptionService from '../services/subscriptionService';
import { toast } from 'vue3-toastify';
import roleService from '../services/roleService';


export default {
    components: {
        Popup,
    },
    data() {
      return {
        form: {
          name: "",
          price: "",
          hours: "",
		  role: 2,
		  hidden: false,
        },
		roles: [],
      };
    },
		mounted() {
				this.fetchRoles()
		},
    methods: {
		fetchRoles() {
				roleService.all()
				.then((res) => {
						this.roles = res
				})
		},
      onSave() {
        subscriptionService.add(this.form.name, parseFloat(this.form.price), this.form.hours, this.form.role, this.form.hidden)
          .then((res) => {
            this.$emit('saved');
          });
      },
    }
}
</script>

<template>
  <Popup :title="$t('new_subscription')" @closed="$emit('closed')">
    <form class="form">

      <div class="form__row">
        <div class="form__col">
          <label>{{ $t('name') }}</label>
          <input type="text" class="form__input" v-model="form.name" />
        </div>
      </div>

      <div class="form__row">
        <div class="form__col">
          <label>{{ $t('price') }}</label>
          <input type="number" class="form__input" v-model="form.price" />
        </div>

        <div class="form__col">
          <label>{{ $t('hours') }}</label>
          <input type="number" class="form__input" v-model="form.hours" />
        </div>
      </div>

	  <div class="form__row">
			  <div class="form__col">
					  <v-select
						label="Available to this role"
						:items="roles"
						item-title="name"
						item-value="id"
						v-model="form.role"
						hide-details
					  ></v-select>
			  </div>
	  </div>

	  <div class="form__row">
			  <div class="form__col">
					  <v-checkbox
					  	label="Hidden"
						v-model="form.hidden"
					  ></v-checkbox>
			  </div>
	  </div>

    </form>

    <br />

    <div class="grid grid-cols-2">
      <button class="button button--primary" @click="onSave">{{ $t('save') }}</button>
      <button class="button button--secondary" @click="$emit('closed')">{{ $t('cancel') }}</button>
    </div>
  </Popup>
</template>
