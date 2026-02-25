import { apiPost } from "./http";

// POST /api/payments
// payload: { requestId, amount, reference }
export async function createPayment(payload) {
  return apiPost(`/api/payments`, payload);
}
