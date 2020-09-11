<template lang="pug">
b-form
  b-form-group(label="请输入网站的主密码: ")
    b-form-input(v-model="master" type="password")
  template(v-if="master")
    select-user(v-model="user" @input="userSelected")
    b-form-group(label="请修改该用户的权限: ")
      b-form-textarea(v-model="role" max-rows="8")
    b-button(variant="primary" @click="submit") 提交
</template>

<script lang="coffee">
import SelectUser from '../common/SelectUser.vue'

export default
  components: 
    'select-user': SelectUser
  data: ->
    master: null
    user: null
    role: null
  methods:
    userSelected: ->
      if !@user
        @role = null
        return
      ajax = await @axios.post('/ajax/user-view-permission-by-root',
        MasterPassword: @master
        ID: @user
      )
      @role = ajax.data
    submit: -> 
      try
        ajax = await @axios.post('/ajax/user-edit-permission-by-root',
          MasterPassword: @master
          ID: @user
          UserProto: @role
        )
        @$bvModal.msgBoxOk('操作成功完成')
      catch error
        @$bvModal.msgBoxOk('操作失败: ' + error.response.statusText + ' ' + error.response.data, okVariant: 'danger')
</script>
