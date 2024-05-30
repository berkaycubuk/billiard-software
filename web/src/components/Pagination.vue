<script>
export default {
    props: {
        pagination: Object,
        routeName: String,
    },
    methods: {
        getList() {
            let list = [];

            for (let i = 0; i < this.pagination.size; i++) {
                list.push({
                    number: i + 1,
                });
            }
            
            return list;
        }
    }
}
</script>

<template>
    <div class="pagination">

        <div v-for="item of getList()">
            <router-link class="button button--small button--blue" :to="{ name: routeName, query: {page: item.number}}">{{ item.number }}</router-link>
        </div>

    </div>
</template>

<style scoped>
.pagination {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-top: 10px;
}
</style>