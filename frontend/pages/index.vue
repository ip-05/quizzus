<template>
  <div class="main">
    <div class="title">Quizzus</div>
    <div class="desciption">
      <div class="slogan">
        {{ slogan }}
      </div>
    </div>
    <div class="actions">
      <NuxtLink to="/new" class="link">
        <home-create-button />
      </NuxtLink>
      <enter-code-input has-icon />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const slogans = ref([
  'Fun, free, and educational!',
  'Learn, compete, and have fun!',
  'The ultimate quiz challenge platform!',
  'Knowledge is power, and Quizzus is your source!',
  'Join the quiz revolution and test your skills!',
  'Where brains and brawn come together in friendly competition!',
  'Sharpen your mind and take on the competition!',
]);

const slogan = ref('');

onMounted(() => {
  slogan.value = slogans.value[Math.floor(Math.random() * slogans.value.length)];
});
</script>

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

.slogan {
  min-height: 27px;
  height: auto;
}
</style>
