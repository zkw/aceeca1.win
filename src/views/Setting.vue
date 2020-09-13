<template lang="pug">
.master
  b-card.m-3(no-body)
    b-tabs(card)
      b-tab(title="主密码修改权限")
        set-permission-as-master
      b-tab(title="管理员修改权限" v-if="$root.user")
        set-permission-as-administrator
</template>

<script lang="coffee">
import SetPermissionAsMaster from './SetPermissionAsMaster.vue'
import SetPermissionAsAdministrator from './SetPermissionAsAdministrator.vue'

export default
  components: 
    'set-permission-as-master': SetPermissionAsMaster
    'set-permission-as-administrator': SetPermissionAsAdministrator
</script>
