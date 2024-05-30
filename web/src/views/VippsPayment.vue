<script>
  import { toast } from 'vue3-toastify';
  import paymentService from '../services/paymentService'

export default {
  data() {
    return {
      id: null,
      frontendUrl: null,
      token: null,
      pirateTimeout: null,
    }
  },
  methods: {
    timeoutForPirates() {
      const self = this
      this.pirateTimeout = setTimeout(function() {
        alert("Payment timeout, please try again.")
        toast.info("Payment timeout, please try again.")

        self.$router.push('/')
      }, 60000 * 5) // 5 minutes
    },
    getStatus() {
      paymentService.status(parseInt(this.$route.params.id), 2)
        .then((res) => {
          if (res.data.status == 3) {
            this.$router.push('/payment/result/' + this.$route.params.id)
            return
          }
        })
        .catch((err) => {
          console.error(err)
        })
    },
  },
    beforeRouteLeave() {
      if (this.pirateTimeout != null) {
        clearTimeout(this.pirateTimeout)
        this.pirateTimeout = null
      }
    },
  mounted() {
    //this.getStatus()

    this.id = this.$route.params.id;
    this.token = this.$route.query.token;
    this.frontendUrl = this.$route.query.frontendUrl;

    if (!this.id || !this.token || !this.frontendUrl) {
      console.error("not all parameters supplied.");
      return;
    }

    let vippsCheckout = VippsCheckout({
      checkoutFrontendUrl: this.frontendUrl,
      iFrameContainerId: "vipps-checkout-frame-container",
      language: "en",
      token: this.token,
    });

    // timeout for the pirates
    this.timeoutForPirates()
  }
}
</script>

<template>
  <main class="container" style="flex: 1;">

    <div class="card">
      <div class="card__body">
        <div id="vipps-checkout-frame-container"></div>
      </div>
    </div>

  </main>
</template>
