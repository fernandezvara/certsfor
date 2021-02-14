import Vue from 'vue'
import App from './App.vue'

Vue.config.productionTip = false

import { MdApp, MdToolbar, MdDrawer, MdMenu, MdList, MdButton, MdContent, MdCard, MdIcon, MdBadge, MdField, MdDivider, MdEmptyState } from 'vue-material/dist/components'
import 'vue-material/dist/vue-material.min.css'
import 'vue-material/dist/theme/default.css'

Vue.use(MdApp)
Vue.use(MdToolbar)
Vue.use(MdDrawer)
Vue.use(MdMenu)
Vue.use(MdList)
Vue.use(MdButton)
Vue.use(MdContent)
Vue.use(MdCard)
Vue.use(MdIcon)
Vue.use(MdBadge)
Vue.use(MdField)
Vue.use(MdDivider)
Vue.use(MdEmptyState)

import Clipboard from 'v-clipboard'

Vue.use(Clipboard)

import store from './store'

new Vue({
  store,
  render: h => h(App),
}).$mount('#app')
