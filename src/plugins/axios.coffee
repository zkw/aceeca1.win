import Vue from 'vue'
import axios from 'axios'

_axios = axios.create()

Plugin.install = (Vue, options) ->
  Vue.axios = _axios
  window.axios = _axios
  Object.defineProperties Vue.prototype,
    axios: get: -> _axios
    $axios: get: -> _axios

Vue.use Plugin
export default Plugin;
