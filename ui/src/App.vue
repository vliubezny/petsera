<script setup>
import { ref, onMounted, defineProps, watch, defineEmits } from "vue";
import PetMap from "./components/PetMap.vue";

const submitMode = ref(false);
</script>

<template>
  <header class="h-14 w-full"></header>
  <main class="container m-auto">
    <PetMap
      :selectable="submitMode"
      class="my-3 mx-auto aspect-square h-[calc(100vh-theme(height.14)-1.5rem)]"
    ></PetMap>
  </main>
</template>
