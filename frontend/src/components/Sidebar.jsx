import { Link, useLocation } from "react-router-dom";

const items = [
  { to: "/", label: "Dashboard" },
  { to: "/citizens", label: "Građani" },
  { to: "/requests", label: "Zahtevi" },
  { to: "/appointments", label: "Termini" },
];

export default function Sidebar() {
  const loc = useLocation();
  return (
    <div style={{ width: 220, background: "#111827", color: "white", padding: 16 }}>
      <div style={{ fontWeight: 700, marginBottom: 16 }}>e-Uprava</div>
      <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
        {items.map((x) => (
          <Link
            key={x.to}
            to={x.to}
            style={{
              textDecoration: "none",
              color: "white",
              padding: "10px 12px",
              borderRadius: 10,
              background: loc.pathname === x.to ? "#1f2937" : "transparent",
            }}
          >
            {x.label}
          </Link>
        ))}
      </div>
    </div>
  );
}
