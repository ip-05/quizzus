<template>
  <div class="wrapper">
    <header class="header">
      <dynamic-island />
    </header>
    <main class="main">
      <NuxtPage />
    </main>
    <footer class="footer">
      <div>
        <a class="footer__link" href="https://github.com/ip-05/quizzus">
          <img class="footer__icon" src="~/assets/svg/icon-github.svg" alt="GitHub Icon" />
        </a>
      </div>
      <span class="footer__copyright">© 2022 Quizzus. All rights reserved</span>
    </footer>
  </div>
</template>

<style scoped>
.wrapper {
  width: 100%;
  min-height: 100vh;
  display: grid;
  grid-template-rows: 70px 1fr auto;
}

.main {
  /* margin-top: 70px; */
  width: 100%;
  height: 100%;
}

.main__content {
  width: 100%;
  height: 100%;
}

.footer {
  width: 100%;
  height: auto;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 10px;
  align-items: center;
  font-family: 'Inter-Medium', sans-serif;
  font-size: var(--font-quaternary-size);
  color: var(--font-secondary-color);
  padding-bottom: 30px;
}

.footer__link {
  display: block;
  width: 24px;
  height: 24px;
}

.footer__icon {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>
