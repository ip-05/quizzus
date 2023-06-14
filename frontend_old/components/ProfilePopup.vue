<template>
  <div class="profile__wrapper">
    <Transition name="backdrop">
      <div v-if="isOpen" class="backdrop" @click="closePopup"></div>
    </Transition>
    <left-action-button
      mode="avatar"
      :img-src="authStore.user ? authStore.user.profilePicture : null"
      @click="openPopup"
    />
    <Transition name="backdrop">
      <div v-if="isOpen" class="popup">
        <medium-button v-if="!authStore.isAuthed" src="svg/icon-google.svg" @click="authStore.signInGoogle"
          >Sign In</medium-button
        >
        <NuxtLink v-if="authStore.isAuthed" to="/" class="popup__link">
          <medium-button src="svg/icon-settings.svg">Settings</medium-button>
        </NuxtLink>
        <NuxtLink v-if="authStore.isAuthed" to="/" class="popup__link">
          <medium-button src="svg/icon-logout.svg" @click="handleLogout">Logout</medium-button>
        </NuxtLink>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';

const authStore = useAuthStore();
const isOpen = ref(false);

const closePopup = () => (isOpen.value = !isOpen.value);
const openPopup = () => (isOpen.value = true);

const handleLogout = () => {
  isOpen.value = false;
  authStore.logout();
};
</script>

<style scoped>
.backdrop {
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  position: fixed;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.2);
  transition: all 250ms ease;
  z-index: -1;
}

.backdrop-enter-active,
.backdrop-leave-active {
  transition: opacity 0.5s ease;
}

.backdrop-enter-from,
.backdrop-leave-to {
  opacity: 0;
}

.popup {
  margin-top: 10px;
  width: 200px;
  background: var(--background-main-color);
  border-radius: 10px;
  outline: solid 3px var(--border-color);
  outline-offset: -3px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  box-shadow: 0px 8px 32px rgba(0, 0, 0, 0.08);
  /* gap: 10px; */
}

.popup__link {
  text-decoration: none;
}
</style>
