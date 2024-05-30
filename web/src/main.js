import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'

import * as Sentry from "@sentry/vue";

// Vuetify
import 'vuetify/styles';
import '@mdi/font/css/materialdesignicons.css';
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import router from './router'
/*
import { createI18n } from 'vue-i18n'
import messages from '@intlify/unplugin-vue-i18n/messages'
*/
import { i18n } from './config'
import Vue3Toastify, { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

/*
const i18n = createI18n({
  locale: 'en',
  fallbackLocale: 'en',
  legacy: false,
  globalInjection: true,
  messages
});
*/

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'dark',
  },
})

const app = createApp(App)

app.use(vuetify);

app.use(i18n);

app.use(Vue3Toastify, {
  autoClose: 3000,
  position: toast.POSITION.TOP_CENTER,
  theme: 'dark',
});

Sentry.init({
  app,
  dsn: "https://194345c9273118aa3cf1ea0439dd0f99@o4506500501340160.ingest.us.sentry.io/4507011936878592",
  integrations: [
    Sentry.browserTracingIntegration(),
    Sentry.replayIntegration({
      maskAllText: false,
      blockAllMedia: false,
    }),
  ],
  // Performance Monitoring
  //tracesSampleRate: 1.0, //  Capture 100% of the transactions
  // Set 'tracePropagationTargets' to control for which URLs distributed tracing should be enabled
  //tracePropagationTargets: ["localhost", /^https:\/\/yourserver\.io\/api/],
  // Session Replay
  replaysSessionSampleRate: 0.1, // This sets the sample rate at 10%. You may want to change it to 100% while in development and then sample at a lower rate in production.
  replaysOnErrorSampleRate: 1.0, // If you're not already sampling the entire session, change the sample rate to 100% when sampling sessions where errors occur.
});

app.use(router);

app.mount('#app')
