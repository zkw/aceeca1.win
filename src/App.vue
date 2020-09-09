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
          b-dropdown-item(v-if="!user" v-b-toggle.login @click="requireToken") 微信登录
          b-dropdown-item(v-if="user" v-b-toggle.set-nick) 修改昵称
          b-dropdown-item(v-if="user" @click="logout") 登出
  b-alert(:show="successCountDown" @dismiss-count-down="successCountDownChanged") 成功
  b-alert(variant="danger" :show="failureCountDown" @dismiss-count-down="failureCountDownChanged") 失败
  b-sidebar#login(lazy bg-variant="dark" text-variant="light" right)
    b-card.mt-3(bg-variant="info")
      | 请微信关注“张昆玮”公众号，
      br
      | 并发送验证码：
      h1.text-center {{ token }}
      b-button.m-2(variant="primary" v-b-toggle.login @click="login") 我已发送，继续登录
      b-button.m-2(variant="danger" v-b-toggle.login) 取消
  b-sidebar#set-nick(lazy bg-variant="dark" text-variant="light" right)
    b-card.mt-3(bg-variant="info")
      .m-2 当前昵称: {{ user }}
      b-form-input.m-2(v-model="nick" placeholder="新的昵称")
      b-button.m-2(variant="primary" v-b-toggle.set-nick @click="setNick") 修改
      b-button.m-2(variant="danger" v-b-toggle.set-nick) 取消
  router-view
</template>

<script lang="coffee">
export default
  data: ->
    successCountDown: 0
    failureCountDown: 0
    nick: null
    token: null
    user: null
  methods:
    getToken: -> if @token then @token else @requireToken()
    getUser: ->
      if @user then @user else '登录'
    login: ->
      ajax = await @axios.get('/ajax/user-login-2')
      @user = ajax.data
      if @user then @successCountDown = 2 else @failureCountDown = 2
    logout: ->
      ajax = await @axios.get('/ajax/user-logout')
      @user = null
    requireToken: ->
      console.log('require token')
      ajax = await @axios.get('/ajax/user-login-1')
      @token = ajax.data
    requireUser: ->
      ajax = await @axios.get('/ajax/user-status')
      @user = ajax.data
    successCountDownChanged: (dismissCountDown) ->
      @successCountDown = dismissCountDown
    failureCountDownChanged: (dismissCountDown) ->
      @failureCountDown = dismissCountDown
    setNick: ->
      ajax = await @axios.get('/ajax/user-set-nick', nick: @nick)
      console.log(ajax.data)
      if ajax.data then @successCountDown = 2 else @failureCountDown = 2
  mounted: -> @requireUser()
</script>
