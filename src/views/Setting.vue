<template lang="pug">
.master
  b-card.m-3(no-body)
    b-tabs(card)
      b-tab(title="主密码修改权限")
        set-permission-as-master
      b-tab(title="管理员修改权限")
</template>

<script lang="coffee">
import SetPermissionAsMaster from './SetPermissionAsMaster.vue'

export default
  components: 
    'set-permission-as-master': SetPermissionAsMaster
</script>
