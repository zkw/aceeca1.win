import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

routes = [
  {
    path: '/'
    component: -> import('../views/Home.vue')
  }
  {
    path: '/master'
    component: -> import('../views/Master.vue')
  }
]

router = new VueRouter(
  mode: 'history'
  base: process.env.BASE_URL
  routes: routes
)

export default router
