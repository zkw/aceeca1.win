<template lang="pug">
b-form
  b-form-group(label="请输入网站的主密码: ")
    b-form-input(v-model="master" type="password")
  template(v-if="master")
    b-form-group(label="请输入站点标识: ")
      b-input-group
        b-form-input(v-model="key" type="text")
        template(v-slot:append)
          b-button(variant="primary" @click="keyInputed") 拉取当前权限
    b-form-group(label="请修改站点信息: ")
      b-form-textarea(v-model="value" max-rows="8")
    b-button(variant="primary" @click="submit") 提交
</template>

<script lang="coffee">
export default
  data: ->
    master: null
    key: null
    value: null
  methods:
    keyInputed: ->
      try
        ajax = await @axios.post('/ajax/site-view',
          MasterPassword: @master
          ID: @key
        )
        @value = ajax.data
        @$bvModal.msgBoxOk('操作成功完成')
      catch error
        @$bvModal.msgBoxOk('操作失败: ' + error.response.statusText + ' ' + error.response.data, okVariant: 'danger')
    submit: -> 
      try
        ajax = await @axios.post('/ajax/site-edit',
          MasterPassword: @master
          ID: @key
          SiteProto: @value
        )
        @$bvModal.msgBoxOk('操作成功完成')
      catch error
        @$bvModal.msgBoxOk('操作失败: ' + error.response.statusText + ' ' + error.response.data, okVariant: 'danger')
</script>
