<template lang="pug">
b-form
  b-form-group(label="请输入网站的主密码: ")
    b-form-input(v-model="master" type="password")
  template(v-if="master")
    b-form-group(label="请输入权限名: ")
      b-input-group
        b-form-input(v-model="name" type="text")
        template(v-slot:append)
          b-button(variant="primary" @click="nameInputed") 拉取当前权限
    b-form-group(label="请修改该权限: ")
      b-form-textarea(v-model="role" max-rows="8")
    b-button(variant="primary" @click="submit") 提交
</template>

<script lang="coffee">
export default
  data: ->
    master: null
    name: null
    role: null
  methods:
    nameInputed: ->
      try
        ajax = await @axios.post('/ajax/permission-view',
          MasterPassword: @master
          Name: @name
        )
        @role = ajax.data
        @$bvModal.msgBoxOk('操作成功完成')
      catch error
        @$bvModal.msgBoxOk('操作失败: ' + error.response.statusText + ' ' + error.response.data, okVariant: 'danger')
    submit: -> 
      try
        ajax = await @axios.post('/ajax/permission-edit',
          MasterPassword: @master
          Name: @name
          PermissionProto: @role
        )
        @$bvModal.msgBoxOk('操作成功完成')
      catch error
        @$bvModal.msgBoxOk('操作失败: ' + error.response.statusText + ' ' + error.response.data, okVariant: 'danger')
</script>
