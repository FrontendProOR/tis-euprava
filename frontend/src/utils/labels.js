export const REQUEST_TYPE_LABEL = {
  ID_CARD: "Lična karta",
  PASSPORT: "Pasoš",
  DRIVER_LICENSE: "Vozačka dozvola",
  RESIDENCE_CHANGE: "Promena prebivališta",
  CITIZENSHIP: "Zahtev za državljanstvo",
};

export const REQUEST_STATUS_LABEL = {
  PENDING: "Podnet",
  IN_REVIEW: "U obradi",
  APPROVED: "Odobren",
  REJECTED: "Odbijen",
  COMPLETED: "Završen",
};

export const PAYMENT_STATUS_LABEL = {
  PAID: "Plaćeno",
  UNPAID: "Nije plaćeno",
};

export const APPOINTMENT_STATUS_LABEL = {
  SCHEDULED: "Zakazan",
  CANCELLED: "Otkazan",
};

export function labelRequestType(t) {
  return REQUEST_TYPE_LABEL[String(t || "").toUpperCase()] || String(t || "");
}

export function labelRequestStatus(s) {
  return REQUEST_STATUS_LABEL[String(s || "").toUpperCase()] || String(s || "");
}

export function labelPaymentStatus(s) {
  return PAYMENT_STATUS_LABEL[String(s || "").toUpperCase()] || String(s || "");
}

export function labelAppointmentStatus(s) {
  return APPOINTMENT_STATUS_LABEL[String(s || "").toUpperCase()] || String(s || "");
}
