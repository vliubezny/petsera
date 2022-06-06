<script setup>
import { ref } from "vue";
import { createAnnouncement } from "@/api";
import SpinnerIcon from "@/components/SpinnerIcon.vue";

const props = defineProps({
  lat: {
    type: Number,
  },
  lng: {
    type: Number,
  },
});

const emit = defineEmits(["submitted"]);

const text = ref("");
const photo = ref();

const loading = ref(false);

const onSubmit = async () => {
  const image = photo.value.files[0];
  const data = {
    text: text.value,
    position: {
      lat: props.lat,
      lng: props.lng,
    },
  };

  try {
    loading.value = true;
    await createAnnouncement(data, image);
  } catch (err) {
    console.error("fail to create", err);
  } finally {
    loading.value = false;
  }

  text.value = "";
  photo.value.value = "";
  emit("submitted");
};
</script>
<template>
  <div class="h-full p-2">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <div class="rounded-lg bg-white py-8 px-4 shadow sm:px-6">
        <form class="mb-0 space-y-6" @submit.prevent="onSubmit">
          <div>
            <label for="details" class="block text-sm font-medium text-gray-700"
              >Details</label
            >
            <div class="mt-1">
              <textarea
                id="details"
                name="details"
                required
                rows="4"
                :disabled="loading"
                v-model="text"
              ></textarea>
            </div>
          </div>

          <div>
            <label for="lat" class="block text-sm font-medium text-gray-700"
              >Latitude</label
            >
            <div class="mt-1">
              <input
                id="lat"
                name="lat"
                type="text"
                disabled
                required
                :value="lat"
              />
            </div>
          </div>

          <div>
            <label for="lng" class="block text-sm font-medium text-gray-700"
              >Longitude</label
            >
            <div class="mt-1">
              <input
                id="lng"
                name="lng"
                type="text"
                disabled
                required
                :value="lng"
              />
            </div>
          </div>

          <div>
            <label for="photo" class="block text-sm font-medium text-gray-700"
              >Photo</label
            >
            <div class="mt-1">
              <input
                ref="photo"
                id="photo"
                name="photo"
                type="file"
                required
                :disabled="loading"
                class="block w-full text-sm text-slate-500 file:mr-4 file:rounded-full file:border-0 file:bg-violet-50 file:py-2 file:px-4 file:text-sm file:font-semibold file:text-indigo-600 hover:file:bg-violet-100"
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              :disabled="loading"
              class="flex w-full justify-center rounded-md border border-transparent bg-indigo-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
            >
              <SpinnerIcon v-if="loading"></SpinnerIcon>
              {{ loading ? "Loading" : "Submit" }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
