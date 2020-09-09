<template lang="pug">
#app
  b-navbar(toggleable="lg" type="dark" variant="dark")
    b-navbar-brand(to="/") 张昆玮的博客
    b-navbar-toggle(target="nav-collapse")
    b-collapse#nav-collapse(is-nav)
      b-navbar-nav
        b-nav-item(to="/about") 关于我
        b-nav-item-dropdown(text="特色内容")
      b-navbar-nav.ml-auto
        b-nav-item-dropdown(:text="getUser()" right)
          b-dropdown-item(v-if="!user" v-b-toggle.login) 微信登录
          b-dropdown-item(v-if="user") 修改昵称
          b-dropdown-item(v-if="user") 登出
  b-sidebar#login(lazy bg-variant="dark" text-variant="light" right)
    b-card.mt-3(bg-variant="info")
      | 请微信关注“张昆玮”公众号，
      br
      | 并发送验证码：
      h1.text-center {{ getToken() }}
      b-button.m-2(variant="primary") 我已发送，继续登录
      b-button.m-2(variant="danger") 取消
  router-view
</template>

<script lang="coffee">
export default
  data: ->
    token: null
    user: null
  methods:
    getToken: -> if @token then @token else @requireToken()
    getUser: -> if @user then @user else '登录'
    requireToken: ->
      @token = await @axios.get('https://wx.aceeca1.win/ajax/user-login-1')
</script>
