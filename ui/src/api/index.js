import axios from "axios";

export function createAnnouncement(data, photo) {
  const formData = new FormData();
  formData.append("data", JSON.stringify(data));
  formData.append("file", photo);
  return axios.post("/api/announcements", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
}

export function getAnnouncements(filter) {
  return axios
    .get("/api/announcements", { params: { ...filter } })
    .then((resp) => resp.data);
}
