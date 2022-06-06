<script setup>
import { formatIsoTimestamp } from "@/util/time";

defineProps({
  item: {
    type: Object,
  },
});

const emit = defineEmits(["close"]);
</script>
<template>
  <div class="h-full p-2">
    <div class="h-full sm:mx-auto sm:w-full sm:max-w-md">
      <div class="flex h-full flex-col rounded-lg bg-white shadow">
        <div class="h-2/5 w-full flex-none">
          <img
            :src="item.imageURL"
            class="h-full w-full rounded-t-lg object-cover"
          />
        </div>
        <div class="py-2 px-4 text-right text-xs text-slate-500">
          {{ formatIsoTimestamp(item.createdAt) }}
        </div>
        <div class="grow overflow-y-auto px-4 py-2 shadow">
          {{ item.text }}
        </div>
        <div class="p-4">
          <button
            type="button"
            @click="emit('close')"
            class="focus:ring-3 flex w-full justify-center rounded-md border border-transparent bg-white py-2 px-4 text-sm font-medium text-indigo-600 shadow-sm ring-2 ring-indigo-500 hover:bg-indigo-700 hover:text-white focus:outline-none focus:ring-indigo-600 focus:ring-offset-2"
          >
            Back
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
