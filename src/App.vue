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
        b-nav-item-dropdown(:text="user" right)
          b-dropdown-item 登出
  router-view
</template>
