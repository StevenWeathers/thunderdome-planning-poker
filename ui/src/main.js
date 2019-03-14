import Vue from 'vue'
import VueRouter from 'vue-router'
import Buefy from 'buefy'
import 'buefy/dist/buefy.css'

import App from './App.vue'
import Landing from './components/Landing.vue'
import Battle from './components/Battle.vue'

Vue.config.productionTip = false

Vue.use(Buefy)
Vue.use(VueRouter)

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', component: Landing },
    { path: '/battle/:id', component: Battle }
  ]
})

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
