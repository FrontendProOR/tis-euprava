import { apiGet, apiPost, apiDelete } from "./http";

export async function createAppointment(payload) {
  return apiPost(`/api/appointments`, payload);
}

export async function getAppointmentsByCitizen(citizenId) {
  return apiGet(`/api/appointments?citizenId=${encodeURIComponent(citizenId)}`);
}

export async function cancelAppointment(id) {
  return apiDelete(`/api/appointments/${id}`);
}
