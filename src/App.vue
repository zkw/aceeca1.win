<template lang="pug">
#app
  b-navbar(toggleable="lg" type="dark" variant="dark")
    b-navbar-brand(to="/") {{ siteName }}
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
          b-dropdown-divider
          b-dropdown-item(to="/setting") 进入管理后台
  b-sidebar#login(bg-variant="dark" text-variant="light" right)
    b-card.mt-3(bg-variant="info")
      | 请微信关注“张昆玮”公众号，
      br
      | 并发送验证码：
      h1.text-center {{ token }}
      b-button.m-2(variant="primary" v-b-toggle.login @click="login") 我已发送，继续登录
      b-button.m-2(variant="danger" v-b-toggle.login) 取消
  b-sidebar#set-nick(bg-variant="dark" text-variant="light" right)
    b-card.mt-3(bg-variant="info")
      .m-2 <b>当前昵称:</b>  {{ user }}
      .m-2 <b>昵称要求:</b> 六个字符或两个汉字以上
      b-form-input.m-2(v-model="nick" placeholder="新的昵称")
      b-button.m-2(variant="primary" v-b-toggle.set-nick @click="setNick") 修改
      b-button.m-2(variant="danger" v-b-toggle.set-nick) 取消
  router-view
</template>

<script lang="coffee">
export default
  data: ->
    siteName: null
    user: null
    nick: null
    token: null
  methods:
    getToken: -> if @token then @token else @requireToken()
    getUser: ->
      if @user then @user else '登录'
    login: ->
      try
        ajax = await @axios.get('/ajax/user-login-2')
        @$root.user = @user = ajax.data
        @$bvModal.msgBoxOk('登录成功')
      catch error
        @$bvModal.msgBoxOk('登录失败', okVariant: 'danger')
    logout: ->
      ajax = await @axios.get('/ajax/user-logout')
      @$root.user = @user = null
    requireSiteInformation: ->
      ajax = await @axios.get('/ajax/site-information')
      document.title = @$root.siteName = @siteName = ajax.data.Name
    requireToken: ->
      ajax = await @axios.get('/ajax/user-login-1')
      @token = ajax.data
    requireUser: ->
      ajax = await @axios.get('/ajax/user-status')
      @$root.user = @user = ajax.data
    setNick: ->
      try
        ajax = await @axios.get('/ajax/user-set-nick', params: nick: @nick)
        @$root.user = @user = @nick
        @$bvModal.msgBoxOk('修改昵称成功')
      catch error
        @$bvModal.msgBoxOk('修改昵称失败', okVariant: 'danger')
  mounted: -> 
    @requireSiteInformation()
    @requireUser()
</script>
