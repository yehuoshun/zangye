<!--
  App.vue — 藏叶前端根组件
  作为 Vue Router 的根容器，仅包含 <router-view />，
  所有页面内容由路由决定。
-->
<template>
  <!-- 路由出口：根据当前 URL 渲染对应页面组件 -->
  <router-view />
</template>

<script setup lang="ts">
// 根组件无需额外逻辑，路由和布局由子组件处理
</script>