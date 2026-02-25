import { useMemo, useState } from "react";
import { useCitizenContext } from "../context/CitizenContext";

function labelFor(c) {
  if (!c) return "";
  const name = `${c.first_name ?? c.firstName ?? ""} ${c.last_name ?? c.lastName ?? ""}`.trim();
  const jmbg = c.jmbg ? ` · ${c.jmbg}` : "";
  return `${name || "(bez imena)"}${jmbg}`;
}

export default function CitizenSelect({ label = "Građanin", required = false }) {
  const { citizens, citizensLoading, selectedCitizenId, setSelectedCitizenId, refreshCitizens } =
    useCitizenContext();
  const [q, setQ] = useState("");

  const filtered = useMemo(() => {
    const s = q.trim().toLowerCase();
    if (!s) return citizens;
    return citizens.filter((c) => labelFor(c).toLowerCase().includes(s));
  }, [citizens, q]);

  return (
    <div>
      <div className="small" style={{ marginBottom: 6 }}>
        {label}{required ? " *" : ""}
      </div>

      <div className="row" style={{ gap: 10 }}>
        <input
          className="input"
          placeholder="Pretraga (ime / prezime / JMBG)"
          value={q}
          onChange={(e) => setQ(e.target.value)}
          style={{ flex: 1 }}
        />
        <button
          className="btn btn-outline"
          type="button"
          onClick={() => refreshCitizens()}
          disabled={citizensLoading}
          style={{ padding: "8px 12px" }}
        >
          {citizensLoading ? "..." : "Osveži"}
        </button>
      </div>

      <select
        className="select"
        style={{ marginTop: 8, width: "100%" }}
        value={selectedCitizenId}
        onChange={(e) => setSelectedCitizenId(e.target.value)}
      >
        <option value="">{required ? "Izaberi građanina" : "(Svi građani)"}</option>
        {filtered.map((c) => (
          <option key={c.id} value={c.id}>
            {labelFor(c)}
          </option>
        ))}
      </select>
    </div>
  );
}
