<template>
  <div class="main">
    <div class="title">Quizzus</div>
    <div class="desciption">
      <span>
        Completely free and user-friendly platform for creating<br />
        quizzes, sharing them with your friends and coworkers.
      </span>
    </div>
    <div class="actions">
      <NuxtLink to="/new" class="link">
        <home-create-button />
      </NuxtLink>
      <enter-code-input has-icon />
    </div>
  </div>
</template>

<script setup></script>

<style scoped>
.main {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.title {
  font-size: 64px;
  font-family: 'Inter-Bold', sans-serif;
  margin-bottom: 20px;
  color: var(--font-primary-color);
}

.desciption {
  font-size: var(--font-secondary-size);
  font-family: 'Inter-Medium', sans-serif;
  color: var(--font-secondary-color);
  line-height: 1.5;
  margin-bottom: 40px;
}

.actions {
  display: flex;
  gap: 10px;
}

.link {
  text-decoration: none;
}
</style>
