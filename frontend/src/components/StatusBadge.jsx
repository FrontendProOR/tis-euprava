export default function StatusBadge({ status }) {
  if (!status) return null;

  let className = "badge";

  switch (status) {
    case "SUBMITTED":
      className += " badge-submitted";
      break;
    case "IN_PROCESS":
      className += " badge-process";
      break;
    case "APPROVED":
      className += " badge-approved";
      break;
    case "REJECTED":
      className += " badge-rejected";
      break;
    case "SCHEDULED":
      className += " badge-process";
      break;
    case "CANCELLED":
      className += " badge-submitted";
      break;
    case "PAID":
      className += " badge-approved";
      break;
    default:
      className += " badge-submitted";
  }

  return <span className={className}>{status}</span>;
}
