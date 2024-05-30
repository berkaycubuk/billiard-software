<script>
import Popup from './Popup.vue';

export default {
  components: {
    Popup,
  },
  props: {
    title: String,
    description: String,
  },
}
</script>

<template>
  <Popup :title="title" @closed="$emit('closed')">
    <div v-if="description != null || description != ''">
      <p>{{ description }}</p>
      <br />
    </div>

    <div class="grid grid-cols-2">
      <button class="button button--red" @click="$emit('deleted')">{{ $t('delete') }}</button>
      <button class="button button--secondary" @click="$emit('closed')">{{ $t('cancel') }}</button>
    </div>
  </Popup>
</template>
