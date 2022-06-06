<script setup>
import _ from "lodash";
import { reactive, ref, watch, computed } from "vue";
import { getAnnouncements } from "@/api";
import PetMap from "@/components/PetMap.vue";
import AnnouncementForm from "@/components/AnnouncementForm.vue";
import AnnouncementPanel from "@/components/AnnouncementPanel.vue";
import AnnouncementsList from "@/components/AnnouncementsList.vue";

const submitMode = ref(false);

const lat = ref();
const lng = ref();

const after = ref("2018-05-30T00:00:00Z");
const bounds = ref();

const loading = ref(false);

const announcements = reactive([]);

const selectedId = ref(null);

const selectedItem = computed(() => {
  return announcements.find((a) => a.id === selectedId.value);
});

const onPositionSelected = (pos) => {
  lat.value = pos.lat;
  lng.value = pos.lng;
};

const onSubmitted = () => {
  switchSelectMode(false);
  fetchAnnouncements();
};

const onBoundsChanged = (newBounds) => {
  bounds.value = newBounds;
};

const fetchAnnouncements = _.debounce(async () => {
  loading.value = true;
  try {
    const items = await getAnnouncements({
      ...bounds.value,
      after: after.value,
    });
    announcements.length = 0;
    announcements.push(...items);
  } catch (err) {
    console.error("fail to fetch announcements", err);
  } finally {
    loading.value = false;
  }
}, 1000);

watch([after, bounds], fetchAnnouncements);

const deselectItem = () => {
  if (selectedId.value) {
    const item = announcements.find((a) => a.id === selectedId.value);
    if (item) {
      item.selected = false;
    }
    selectedId.value = null;
  }
};

const selectItem = (id) => {
  deselectItem();
  const item = announcements.find((a) => a.id === id);
  if (item) {
    item.selected = true;
    selectedId.value = id;
  }
};

const switchSelectMode = (mode) => {
  submitMode.value = mode ?? !submitMode.value;
  if (!submitMode.value) {
    lat.value = undefined;
    lng.value = undefined;
  }
};
</script>

<template>
  <header class="flex h-14 w-full items-center justify-between bg-white p-4">
    <h1 class="font-joan text-4xl font-bold text-indigo-700">Petsera</h1>
    <button
      type="button"
      class="flex w-32 justify-center rounded-md border border-transparent bg-indigo-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
      @click="switchSelectMode()"
    >
      {{ submitMode ? "Cancel" : "Report" }}
    </button>
  </header>
  <main class="container m-auto flex">
    <PetMap
      :selectable="submitMode"
      :markers="announcements"
      @position-selected="onPositionSelected"
      @bounds-changed="onBoundsChanged"
      @marker-selected="selectItem"
      :class="{ 'animate-pulse': loading }"
      class="m-2 h-[calc(100vh-theme(height.14)-1rem)] w-[calc(100%-theme(width.72))] rounded-lg shadow"
    ></PetMap>
    <div class="h-[calc(100vh-theme(height.14))] w-72">
      <AnnouncementForm
        v-if="submitMode"
        :lat="lat"
        :lng="lng"
        @submitted="onSubmitted"
      ></AnnouncementForm>
      <AnnouncementPanel
        v-if="selectedItem"
        :item="selectedItem"
        @close="deselectItem"
      ></AnnouncementPanel>
      <AnnouncementsList
        v-if="announcements.length"
        :items="announcements"
        @item-selected="selectItem"
      ></AnnouncementsList>
      <p
        v-else
        class="m-2 rounded-lg bg-white p-2 text-center text-xs text-slate-500"
      >
        There is no announcements in selected area.
      </p>
    </div>
  </main>
</template>
