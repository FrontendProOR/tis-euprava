import http from "./http";

// POST /api/payments
export async function createPayment(payload) {
  const { data } = await http.post("/api/payments", payload);
  return data;
}
