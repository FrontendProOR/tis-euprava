export default function SearchInput({ value, onChange, placeholder = "Pretraga..." }) {
  return (
    <div style={{ position: "relative", width: "min(420px, 100%)" }}>
      <span
        aria-hidden="true"
        style={{
          position: "absolute",
          left: 12,
          top: "50%",
          transform: "translateY(-50%)",
          color: "var(--muted)",
          fontSize: 14,
        }}
      >
        ⌕
      </span>
      <input
        className="input"
        value={value}
        onChange={(e) => onChange?.(e.target.value)}
        placeholder={placeholder}
        style={{ paddingLeft: 34 }}
      />
    </div>
  );
}
