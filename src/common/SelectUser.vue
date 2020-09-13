<template lang="pug">
b-form-group(label="请选择用户: ")
  b-input-group(v-if="!value")
    b-form-input(v-model="query" type="text")
    template(v-slot:append)
      b-dropdown(text="搜索并选择" variant="primary" right @show="search")
        b-dropdown-item(v-for="(nick, id) in userList" :key="id" @click="select(id, nick)") {{ nick }}
  b-button(v-if="value" variant="info" @click="select(null, null)") {{nick}}
</template>

<script lang="coffee">
export default 
  data: ->
    query: null
    nick: null
    userList: {}
  props:
    value: String
  methods:
    search: ->
      ajax = await @axios.get('/ajax/user-list', params: query: @query)
      @userList = ajax.data
    select: (id, nick) ->
      @nick = nick
      @$emit('input', id)
</script>
