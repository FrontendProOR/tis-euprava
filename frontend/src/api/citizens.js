import { apiGet, apiPost } from "./http";

export async function createCitizen(payload) {
  return apiPost(`/api/citizens`, payload);
}

export async function getCitizens() {
  return apiGet(`/api/citizens`);
}

export async function getCitizenById(id) {
  return apiGet(`/api/citizens/${id}`);
}
