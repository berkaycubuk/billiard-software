<script>
import FilterPopup from './popups/FilterPopup.vue';

export default {
    props: ['dates', 'method', 'user', 'zeroOrders'],
    data() {
        return {
            show: false,
        };
    },
    components: { FilterPopup }
}
</script>

<template>
    <button class="button button--primary button--small" @click="() => show = true">{{ $t('filter') }}</button>
    <FilterPopup v-if="show === true" @closed="() => show = false" @saved="() => show = false" :dates="dates" :method="method" :user="user" :zeroOrders="zeroOrders" />
</template>

<style scoped>
</style>
