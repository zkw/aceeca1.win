import Vue from 'vue'
import './plugins/bootstrap-vue.coffee'
import './plugins/axios.coffee'
import App from './App.vue'
import router from './router/index.coffee'

new Vue(
  router: router
  render: (h) -> h(App)
).$mount('#app')
