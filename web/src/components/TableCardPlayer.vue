<script>
  export default {
    props: {
      name: String,
      time: String,
      price: String,
      status: Number,
    },
    data() {
      return {
        currency: import.meta.env.VITE_CURRENCY
      }
    },
    methods: {
      convertDuration(duration) {
        const hour = Math.floor(duration / 60)
        const minute = Math.floor(duration - (hour * 60))
        const second = Math.floor((duration - (hour * 60) - minute) * 60)

        return `${hour < 10 ? '0' + hour : hour}:${minute < 10 ? '0' + minute : minute}:${second < 10 ? '0' + second : second}`
      },
      convertCurrency() {
        if (!this.currency) return '';

        if (this.currency == 'kr') return 'kr';

        if (this.currency == 'eur') return '€';

        if (this.currency == 'usd') return '$';

        if (this.currency == 'try') return '₺';

        return this.currency;
      }
    },
  }
</script>

<template>
  <div class="table-player" :class="{ 'table-player--paused': status === 2 }">
    <div class="table-player__name">{{ name }}</div>
    <div class="table-player__info">
      <div class="table-player__time">{{ convertDuration(time) }}</div>
      <div class="table-player__price">{{ price }}{{ convertCurrency() }}</div>
    </div>
  </div>
</template>

<style scoped>
@keyframes blink {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: .4;
  }
}

.table-player {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 10px;
  background-color: rgba(25, 25, 25, .95);
  border-radius: 8px;
}

.table-player--paused {
  background-color: rgba(100, 23, 23, .95);
  animation: blink 3s linear infinite;
}

.table-player--paused .table-player__time {
  color: var(--white);
}

.table-player__name {
  text-align: center;
  font-weight: 600;
  font-size: 1rem;
}

.table-player__info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.table-player__time {
  font-size: .9rem;
  color: var(--green);
}

.table-player__price {
  font-size: .9rem;
}
</style>
