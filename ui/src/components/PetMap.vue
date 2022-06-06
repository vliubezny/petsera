<script setup>
// @ts-check

import { ref, onMounted, defineProps, watch, defineEmits } from "vue";
import { Loader } from "@googlemaps/js-api-loader";
import { MarkerClusterer } from "@googlemaps/markerclusterer";

// @ts-ignore
const { apiKey } = window.appConfig;

const props = defineProps({
  markers: {
    type: Array,
    default() {
      return [];
    },
  },
  selectable: {
    type: Boolean,
  },
});

const emit = defineEmits([
  "bounds-changed",
  "marker-selected",
  "position-selected",
]);

const map = ref();

/** @type {google.maps.Map} */
let gmap;

/** @type {MarkerClusterer} */
let clusterer;

/** @type {google.maps.Marker} */
let currentPosition;

const updateMarkers = () => {
  clusterer.clearMarkers(true);
  const markers = props.markers.map((m) => {
    const marker = new window.google.maps.Marker({
      // @ts-ignore
      position: m.position,
      // @ts-ignore
      label: m.title,
      // @ts-ignore
      icon: m.selected
        ? "https://maps.gstatic.com/mapfiles/ms2/micons/purple-dot.png"
        : "https://maps.gstatic.com/mapfiles/ms2/micons/red-dot.png",
    });
    marker.addListener("click", () => {
      // @ts-ignore
      emit("marker-selected", m.id);
    });
    return marker;
  });
  clusterer.addMarkers(markers);
};

watch(props.markers, updateMarkers);

watch(
  () => props.selectable,
  (selectable) => {
    if (!selectable && currentPosition) {
      clusterer.removeMarker(currentPosition);
      currentPosition = null;
    }
  }
);

onMounted(() => {
  /**
   * Init Google Map
   * @param {google} google
   */
  const init = (google) => {
    gmap = new google.maps.Map(map.value, {
      center: { lat: 53.906, lng: 27.555 },
      zoom: 12,
      streetViewControl: false,
    });

    clusterer = new MarkerClusterer({ map: gmap });
    updateMarkers();

    gmap.addListener("click", (e) => {
      if (!props.selectable) {
        return;
      }

      if (currentPosition) {
        clusterer.removeMarker(currentPosition);
      }

      const position = { lat: e.latLng.lat(), lng: e.latLng.lng() };
      const marker = new google.maps.Marker({
        position,
        title: "New position",
        icon: "https://maps.gstatic.com/mapfiles/ms2/micons/purple-dot.png",
      });
      clusterer.addMarker(marker);

      currentPosition = marker;

      emit("position-selected", position);
    });

    gmap.addListener("bounds_changed", () => {
      const bounds = gmap.getBounds();
      emit("bounds-changed", bounds.toJSON());
    });
  };

  if (window.google) {
    init(window.google);
  } else {
    const loader = new Loader({
      apiKey,
      version: "weekly",
      language: "en",
    });

    loader.load().then(init);
  }
});
</script>

<template>
  <div ref="map" data-testId="map"></div>
</template>
