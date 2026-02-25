import { labelAppointmentStatus, labelPaymentStatus, labelRequestStatus } from "../utils/labels";

export default function StatusBadge({ status, kind = "request" }) {
  if (!status) return null;

  let className = "badge";
  const s = String(status).toUpperCase();

  switch (s) {
    case "PENDING":
      className += " badge-pending";
      break;
    case "IN_REVIEW":
      className += " badge-review";
      break;
    case "APPROVED":
      className += " badge-approved";
      break;
    case "COMPLETED":
      className += " badge-completed";
      break;
    case "REJECTED":
      className += " badge-rejected";
      break;

    case "SCHEDULED":
      className += " badge-review";
      break;
    case "CANCELLED":
      className += " badge-pending";
      break;
    case "PAID":
      className += " badge-paid";
      break;
    case "UNPAID":
      className += " badge-unpaid";
      break;

    default:
      className += " badge-pending";
  }

  let label = s;
  if (kind === "payment") label = labelPaymentStatus(s);
  else if (kind === "appointment") label = labelAppointmentStatus(s);
  else label = labelRequestStatus(s);

  return <span className={className} title={s}>{label}</span>;
}
