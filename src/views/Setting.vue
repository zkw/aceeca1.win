<template lang="pug">
.master
  b-card.m-3(no-body)
    b-tabs(card)
      b-tab(title="权限指派")
        set-permission
      b-tab(title="用户权限")
        set-user-permission-as-master
      b-tab(title="站点信息")
        set-site-information
      b-tab(title="管理员修改用户权限" v-if="$root.user")
        set-user-permission-as-administrator
</template>

<script lang="coffee">
import SetPermission from './SetPermission.vue'
import SetSiteInformation from './SetSiteInformation.vue'
import SetUserPermissionAsAdministrator from './SetUserPermissionAsAdministrator.vue'
import SetUserPermissionAsMaster from './SetUserPermissionAsMaster.vue'

export default
  components: 
    'set-permission': SetPermission
    'set-site-information': SetSiteInformation
    'set-user-permission-as-administrator': SetUserPermissionAsAdministrator
    'set-user-permission-as-master': SetUserPermissionAsMaster
</script>
