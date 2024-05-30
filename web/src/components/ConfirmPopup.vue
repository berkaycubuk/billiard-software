<script>
  export default {
  }
</script>

<template>
  <v-card title="Are you sure?" style="width: fit-content; margin: 0 auto;">
    <v-card-actions>
      <v-btn @click="() => $emit('confirm')" variant="flat" color="indigo-darken-2">Confirm</v-btn>
      <v-btn @click="() => $emit('cancel')">Cancel</v-btn>
    </v-card-actions>
  </v-card>
</template>
