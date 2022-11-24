<template>
  <nav class="nav">
    <div class="header">
      <div class="header__icon">
        <img src="~/assets/svg/button-menu.svg" alt="Menu" />
      </div>
      <div class="header__status">Welcome to Quizzus</div>
      <div class="header__icon">
        <img src="~/assets/svg/button-login.svg" alt="Login" />
      </div>
    </div>
    <div class="body">
      <!-- <div>
        <NuxtLink to="/">Home page</NuxtLink>
        <NuxtLink to="/about">About page</NuxtLink>
        <NuxtLink to="/game">Game page</NuxtLink>
        <NuxtPage />
      </div> -->
    </div>
  </nav>
</template>

<script>
export default {};
</script>

<style scoped>
.nav {
  position: fixed !important;
  top: 10px;
  left: 50%;
  transform: translate(-50%);
  cursor: pointer;
  max-width: 1280px;
  width: calc(100% - 20px);
  height: auto;
  padding: 10px;
  outline: solid 3px var(--border-color);
  outline-offset: -3px;
  border-radius: 40px;
  background: var(--background-main-color);
}

.header {
  width: 100%;
  height: auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header__icon {
  width: 30px;
  height: 30px;
}

.header__icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.header__status {
  font-family: 'Inter-Medium', sans-serif;
  font-size: var(--font-secondary-size);
}
</style>
