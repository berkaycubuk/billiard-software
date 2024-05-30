import {createI18n} from "vue-i18n";
import messages from '@intlify/unplugin-vue-i18n/messages'

const i18n = createI18n({
  locale: 'no',
  fallbackLocale: 'en',
  legacy: false,
  globalInjection: true,
  messages
})

const $t = i18n.global.t

export { i18n, $t };
