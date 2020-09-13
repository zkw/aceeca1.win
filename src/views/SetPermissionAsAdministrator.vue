<template lang="pug">
b-form
  b-form-group(label="请选择要修改的权限: ")
    b-dropdown(v-if="!permission" text="请选择" variant="primary" @show="search")
      b-dropdown-item(v-for="permission in permissionList" :key="permission" @click="select(permission)") {{ permission }}
    b-button(v-if="permission" variant="info" @click="select(null)") {{ permission }}
  select-user(v-model="user")
  b-form-group(v-if="user && permission" label="请选择操作: ")
    b-button.mr-3(variant="success" @click="setPermission('ADMINISTRATOR')") 设为管理员
    b-button.mr-3(variant="warning" @click="setPermission('MEMBER')") 设为成员
    b-button.mr-3(variant="danger" @click="removePermission") 删除
</template>

<script lang="coffee">
import SelectUser from '../common/SelectUser.vue'

export default
  components: 
    'select-user': SelectUser
  data: ->
    permission: null
    permissionList: []
    user: null
  methods:
    search: ->
      ajax = await @axios.get('/ajax/user-role-as-administrator')
      @permissionList = ajax.data
    select: (permission) ->
      @permission = permission
    setPermission: (role) ->
      try
        ajax = await @axios.post('/ajax/user-edit-permission-by-administrator', 
          User: @user
          Permission: @permission
          Role: role
        )
        @$bvModal.msgBoxOk('操作成功完成')
        @user = null
      catch error
        @$bvModal.msgBoxOk('操作失败: ' + error.response.statusText + ' ' + error.response.data, okVariant: 'danger')
    removePermission: ->
      try
        ajax = await @axios.post('/ajax/user-remove-permission-by-administrator', 
          User: @user
          Permission: @permission
        )
        @$bvModal.msgBoxOk('操作成功完成')
        @user = null
      catch error
        @$bvModal.msgBoxOk('操作失败: ' + error.response.statusText + ' ' + error.response.data, okVariant: 'danger')
</script>
