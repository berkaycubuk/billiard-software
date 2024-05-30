<script>
  export default {
    methods: {
      getFilters() {
        const filters = Object.entries(this.$route.query)
        return filters;
      }
    },
  }
</script>

<template>
  <div class="filters-title" v-if="getFilters() != null && getFilters().length > 0">Filters</div>
  <div class="filters" v-if="getFilters() != null && getFilters().length > 0">
    <div class="filter" v-for="filter of getFilters()"><b>{{ filter[0] }}:</b> {{ filter[1] }}</div>
  </div>
</template>

<style scoped>
.filters {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 10px;
}

.filters-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 6px;
}

.filter {
  padding: 4px 6px;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}
</style>
