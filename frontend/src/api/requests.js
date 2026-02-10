import { apiGet, apiPost, apiPatch } from "./http";

function qs(params = {}) {
  const sp = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== null && `${v}`.trim() !== "") sp.set(k, v);
  });
  const s = sp.toString();
  return s ? `?${s}` : "";
}

export async function getRequests(params) {
  return apiGet(`/api/requests${qs(params)}`);
}

export async function createRequest(payload) {
  return apiPost(`/api/requests`, payload);
}

export async function getRequestById(id) {
  return apiGet(`/api/requests/${id}`);
}

export async function updateRequestStatus(id, status) {
  return apiPatch(`/api/requests/${id}/status`, { status });
}
