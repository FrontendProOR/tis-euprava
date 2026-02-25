export default function Modal({ open, title, children, onClose, width = 760 }) {
  if (!open) return null;
  return (
    <div
      role="dialog"
      aria-modal="true"
      onClick={onClose}
      style={{
        position: "fixed",
        inset: 0,
        background: "rgba(15,23,42,.35)",
        display: "grid",
        placeItems: "center",
        zIndex: 50,
        padding: 18,
      }}
    >
      <div
        className="card"
        onClick={(e) => e.stopPropagation()}
        style={{ width: `min(${width}px, 100%)` }}
      >
        <div className="card-h" style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
          <h3 className="card-t">{title}</h3>
          <button className="btn btn-outline" onClick={onClose}>
            Zatvori
          </button>
        </div>
        <div className="card-c">{children}</div>
      </div>
    </div>
  );
}
