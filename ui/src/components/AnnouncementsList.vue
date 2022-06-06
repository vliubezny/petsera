<script setup>
import { formatIsoTimestamp } from "@/util/time";

defineProps({
  items: {
    type: Array,
    default() {
      return [];
    },
  },
});

const emit = defineEmits(["item-selected"]);

const onItemClick = (item) => {
  emit("item-selected", item.id);
};
</script>
<template>
  <div class="h-full p-2">
    <div class="h-full sm:mx-auto sm:w-full sm:max-w-md">
      <ul
        class="flex max-h-full flex-col divide-y overflow-y-auto rounded-lg bg-white py-4 shadow"
      >
        <li
          v-for="i in items"
          :key="i.id"
          @click="onItemClick(i)"
          class="h-21 flex cursor-pointer p-1 hover:bg-violet-100"
        >
          <div class="h-20 w-20 flex-none">
            <img
              :src="i.imageURL"
              class="h-full w-full rounded-sm object-cover"
            />
          </div>

          <div class="flex grow flex-col px-2">
            <div class="text-right text-xs text-slate-500">
              {{ formatIsoTimestamp(i.createdAt) }}
            </div>
            <div class="text-xs line-clamp-4">
              {{ i.text }}
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>
