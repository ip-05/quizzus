<script setup lang="ts">
import { onMounted, ref } from 'vue';

interface Props {
  speed?: number;
}

const props = withDefaults(defineProps<Props>(), {
  speed: 275,
});

const canvas = ref();
const ctx = ref();

const width = ref();
const height = ref();

const grid = ref();
const newGrid = ref();
const cellSize = ref(20);
const cellsX = ref();
const cellsY = ref();

onMounted(() => {
  width.value = window.innerWidth;
  height.value = window.innerHeight;

  canvas.value = document.getElementById('canvas') as HTMLCanvasElement;
  canvas.value.width = width.value;
  canvas.value.height = height.value;

  ctx.value = canvas.value.getContext('2d');

  const createRadialGradient = (
    context: CanvasRenderingContext2D,
    x0: number,
    y0: number,
    r0: number,
    x1: number,
    y1: number,
    r1: number,
    colorStops: [number, string][],
  ) => {
    const gradient = context.createRadialGradient(x0, y0, r0, x1, y1, r1);
    colorStops.forEach((stop: [number, string]) => {
      gradient.addColorStop(stop[0], stop[1]);
    });
    return gradient;
  };

  const gradient = createRadialGradient(
    ctx.value,
    0.5 * width.value,
    0.5 * height.value,
    0,
    0.5 * width.value,
    0.5 * height.value,
    0.5 * width.value,
    [
      [0.1, '#FBFBFB'],
      [0.5208, `rgb(66, 133, 244)`],
      [0.6563, `rgb(231, 76, 60)`],
      [0.7969, `rgb(241, 196, 15)`],
      [1, `rgb(46, 204, 113)`],
    ],
  );

  ctx.value.fillStyle = gradient;

  cellsX.value = Math.floor(width.value / cellSize.value);
  cellsY.value = Math.floor(height.value / cellSize.value);

  grid.value = new Int8Array(cellsX.value * cellsY.value);
  newGrid.value = new Int8Array(cellsX.value * cellsY.value);

  randomizeGrid();
  setInterval(updateGameState, props.speed);
});

const pos = (x: number, y: number) => {
  return cellsY.value * x + y;
};

const updateGameState = () => {
  for (let i = 1; i < cellsX.value - 1; i++) {
    for (let j = 1; j < cellsY.value - 1; j++) {
      const c =
        grid.value[pos(i - 1, j - 1)] +
        grid.value[pos(i - 1, j)] +
        grid.value[pos(i - 1, j + 1)] +
        grid.value[pos(i, j - 1)] +
        grid.value[pos(i, j + 1)] +
        grid.value[pos(i + 1, j - 1)] +
        grid.value[pos(i + 1, j)] +
        grid.value[pos(i + 1, j + 1)];
      if ((grid.value[pos(i, j)] && (c == 2 || c == 3)) || (!grid.value[pos(i, j)] && c == 3))
        newGrid.value[pos(i, j)] = 1;
      else newGrid.value[pos(i, j)] = 0;
    }
  }
  [grid.value, newGrid.value] = [newGrid.value, grid.value];
  drawGrid();
};

function drawGrid() {
  ctx.value.clearRect(0, 0, width.value, height.value);
  for (let i = 1; i < cellsX.value - 1; i++) {
    for (let j = 1; j < cellsY.value - 1; j++) {
      if (grid.value[pos(i, j)]) {
        ctx.value.fillRect(
          cellSize.value * (i - 1),
          cellSize.value * (j - 1),
          cellSize.value,
          cellSize.value,
        );
      }
    }
  }
}

const randomizeGrid = () => {
  for (let i = 0; i < cellsX.value; i++) {
    for (let j = 0; j < cellsY.value; j++) {
      grid.value[pos(i, j)] = Math.random() < 0.5 ? 0 : 1;
    }
  }
};
</script>

<template>
  <canvas id="canvas" class="canvas"></canvas>
</template>

<style scoped>
.canvas {
  top: 0;
  right: 0;
  left: 0;
  bottom: 0;
  position: fixed;
  width: 100%;
  height: 100%;
  z-index: -1;
  overflow: hidden;
  opacity: 0;
  animation: canvas_appear 2s ease 0.5s 1 normal forwards;
}

@keyframes canvas_appear {
  to {
    opacity: 0.1;
  }
}
</style>
