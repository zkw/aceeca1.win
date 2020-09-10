<template lang="pug">
.master
  b-card.m-3(no-body)
    b-tabs(card)
      b-tab(title="主密码增加权限")
        b-form
          b-form-group(label="请输入网站的主密码: ")
            b-form-input(v-model="master" type="password")
          select-user(v-model="from")
          b-form-group(label="请输入目标权限: ")
            b-form-input(v-model="to" type="text")
          b-form-group(label="请选择权限类别: ")
            b-form-radio-group(v-model="role")
              b-form-radio(value="ADMINISTRATOR") 管理员
              b-form-radio(value="MEMBER") 成员
          b-button(variant="primary" @click="submit") 提交
</template>

<script lang="coffee">
import SelectUser from '../common/SelectUser.vue'

export default
  components: 
    'select-user': SelectUser
  data: ->
    master: null
    from: null
    to: null
    role: null
  methods:
    submit: -> 
      try
        ajax = await @axios.get('/ajax/user-grant-permission-by-root', params:
          'master-password': @master
          from: @from
          to: @to
          role: @role
        )
        @$bvModal.msgBoxOk('操作成功完成')
      catch error
        @$bvModal.msgBoxOk('操作失败: ' + error.response.statusText, okVariant: 'danger')
</script>
