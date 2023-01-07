<template>
  <div class="main">
    <div class="content">
      <create-game />
    </div>
  </div>
</template>

<script setup></script>

<style scoped>
.main {
  margin-top: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.content {
  max-width: 1080px;
  width: calc(100% - 20px);
  height: auto;
}
</style>
