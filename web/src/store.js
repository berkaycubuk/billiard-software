import { reactive } from 'vue';

export const store = reactive({
  user: null,
  setUser(val) {
    this.user = val;
  }
});
