<template lang="pug">
.master
  b-card.m-3(no-body)
    b-tabs(card)
      b-tab(title="主密码增加权限")
        b-form
          b-form-group(label="请输入网站的主密码: ")
            b-form-input(v-model="master" type="password")
          b-form-group(label="请选择用户: " description="帮你")
            b-input-group
              b-form-input(v-model="from" type="text")
              b-dropdown(v-slot:append text="搜索并选择" right)
          b-form-group(label="请输入目标权限: ")
            b-form-input(v-model="to" type="text")
          b-form-group(label="请输入权限类别: ")
            b-form-radio-group(v-model="role")
              b-form-radio(value="ADMINISTRATOR") 管理员
              b-form-radio(value="MEMBER") 成员
          b-button(variant="danger" @click="submit") 提交
</template>

<script lang="coffee">
export default
  data: ->
    master: null
    from: null
    to: null
    role: null
  methods:
    submit: -> console.log(from, to, role)
</script>
