<script>
import Popup from './Popup.vue';
export default {
    components: {
        Popup,
    }
}
</script>

<template>
  <Popup :title="$t('approve_order')" @closed="$emit('closed')">
    <div class="grid grid-cols-2">
      <button class="button button--primary" @click="$emit('approved')">{{ $t('approve') }}</button>
      <button class="button button--secondary" @click="$emit('closed')">{{ $t('cancel') }}</button>
    </div>
  </Popup>
</template>