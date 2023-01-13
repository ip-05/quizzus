<template>
  <tr class="row">
    <td class="cell place">{{ place }}</td>
    <td class="cell points">{{ points }}</td>
    <td class="cell name">{{ name }}</td>
  </tr>
</template>

<script setup>
import { defineProps } from 'vue';
const props = defineProps({
  points: {
    type: Number,
    default: 0,
  },
  place: {
    type: Number,
    default: 0,
  },
  name: {
    type: String,
    default: 'Name',
  },
});
</script>

<style scoped>
.row {
  background: var(--background-main-color);
  position: relative;
}

.row::after {
  position: absolute;
  content: '';
  height: 2px;
  border-radius: 2px;
  background: var(--background-secondary-color);
  left: 20px;
  right: 20px;
  bottom: 0;
}

.row:only-of-type::after,
.row:last-of-type::after {
  height: 0;
}

.row:first-child > .cell:first-of-type {
  border-radius: 15px 0 0 0;
}

.row:first-child > .cell:last-of-type {
  border-radius: 0 15px 0 0;
}

.row:last-child > .cell:first-of-type {
  border-radius: 0 0 0 15px;
}

.row:last-child > .cell:last-of-type {
  border-radius: 0 0 15px 0;
}

.row:only-child > .cell:first-of-type {
  border-radius: 15px 0 0 15px;
}

.row:only-child > .cell:last-of-type {
  border-radius: 0 15px 15px 0;
}

.cell {
  /* border: 1px solid black; */
  padding: 20px;
}

.cell:nth-of-type(2) {
  text-align: center;
}

.cell:last-of-type {
  width: 100%;
}
</style>
