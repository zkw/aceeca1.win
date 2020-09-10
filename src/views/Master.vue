<template lang="pug">
.master
  b-form.m-3
    b-form-group(label="请输入网站的主密码: ")
      b-form-input(v-model="master" type="password")
  b-card.m-3(no-body)
    b-tabs(card)
      b-tab(title="增加权限")
        b-form
          b-form-group(label="请输入用户的唯一标识符: ")
            b-form-input(v-model="from" type="text" required)
          b-form-group(label="请输入目标权限: ")
            b-form-input(v-model="to" type="text" required)
          b-form-group(label="请输入权限类别: ")
            b-form-radio-group(v-model="role" required)
              b-form-radio(value="ADMINISTRATOR") 管理员
              b-form-radio(value="MEMBER") 成员
          b-button(type="submit" variant="danger") 提交
</template>
