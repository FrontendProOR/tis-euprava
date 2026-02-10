import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import StatusBadge from "../components/StatusBadge";
import { getRequests, createRequest } from "../api/requests";

// Pomocne funkcije (lokalne, da ne zavisis od drugih fajlova)
function shorten(id = "") {
  if (!id) return "";
  return id.length > 10 ? `${id.slice(0, 6)}…${id.slice(-4)}` : id;
}

function fmtDate(s) {
  if (!s) return "";
  const d = new Date(s);
  if (Number.isNaN(d.getTime())) return s;
  return d.toLocaleString();
}

export default function Requests() {
  // Filters
  const [filters, setFilters] = useState({
    citizenId: "",
    status: "",
    type: "",
  });

  // Data
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(false);

  // Modal (Novi zahtev)
  const [open, setOpen] = useState(false);
  const [createForm, setCreateForm] = useState({
    citizenId: "",
    type: "ID_CARD",
  });

  const statusOptions = useMemo(
    () => ["SUBMITTED", "IN_PROCESS", "APPROVED", "REJECTED"],
    []
  );
  const typeOptions = useMemo(
    () => ["ID_CARD", "PASSPORT", "DRIVER_LICENSE"],
    []
  );

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function load(overrideFilters) {
    setLoading(true);
    try {
      const used = overrideFilters ?? filters;
      const data = await getRequests({
        citizenId: used.citizenId || undefined,
        status: used.status || undefined,
        type: used.type || undefined,
      });

      // backend može vratiti {items: []} ili direktno []
      const list = Array.isArray(data) ? data : data?.items ?? [];
      setRows(list);
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri učitavanju zahteva");
    } finally {
      setLoading(false);
    }
  }

  async function applyFilters() {
    await load(filters);
  }

  async function submitNewRequest() {
    if (!createForm.citizenId.trim()) {
      alert("Unesi citizenId");
      return;
    }
    try {
      await createRequest({
        citizenId: createForm.citizenId.trim(),
        type: createForm.type,
      });
      setOpen(false);
      setCreateForm((p) => ({ ...p, citizenId: "" }));
      await load(filters);
      alert("Zahtev uspešno kreiran");
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri kreiranju zahteva");
    }
  }

  async function copy(text) {
    try {
      await navigator.clipboard.writeText(text);
      alert("Kopirano!");
    } catch {
      alert("Ne mogu da kopiram (browser dozvole)");
    }
  }

  return (
    <div className="grid" style={{ gap: 16 }}>
      {/* Header */}
      <div className="row row-between">
        <div>
          <h1 className="h1">Zahtevi</h1>
          <p className="p">
            Filteri + tabela zahteva. Otvori detalje da uradiš status → uplatu → sertifikat.
          </p>
        </div>

        <button className="btn btn-primary" onClick={() => setOpen(true)}>
          + Novi zahtev
        </button>
      </div>

      {/* Filter bar */}
      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Filteri</h3>
        </div>
        <div className="card-c">
          <div className="grid grid-3" style={{ alignItems: "end" }}>
            <div>
              <div className="small" style={{ marginBottom: 6 }}>Citizen ID (opciono)</div>
              <input
                className="input"
                placeholder="npr. 8b3d..."
                value={filters.citizenId}
                onChange={(e) =>
                  setFilters((p) => ({ ...p, citizenId: e.target.value }))
                }
              />
            </div>

            <div>
              <div className="small" style={{ marginBottom: 6 }}>Status</div>
              <select
                className="select"
                value={filters.status}
                onChange={(e) =>
                  setFilters((p) => ({ ...p, status: e.target.value }))
                }
              >
                <option value="">Svi</option>
                {statusOptions.map((s) => (
                  <option key={s} value={s}>{s}</option>
                ))}
              </select>
            </div>

            <div>
              <div className="small" style={{ marginBottom: 6 }}>Tip</div>
              <select
                className="select"
                value={filters.type}
                onChange={(e) =>
                  setFilters((p) => ({ ...p, type: e.target.value }))
                }
              >
                <option value="">Svi</option>
                {typeOptions.map((t) => (
                  <option key={t} value={t}>{t}</option>
                ))}
              </select>
            </div>
          </div>

          <div className="row" style={{ marginTop: 12 }}>
            <button className="btn btn-secondary" onClick={applyFilters}>
              Primeni
            </button>
            <button
              className="btn btn-outline"
              onClick={() => {
                const reset = { citizenId: "", status: "", type: "" };
                setFilters(reset);
                load(reset);
              }}
            >
              Reset
            </button>

            <div className="small" style={{ marginLeft: "auto" }}>
              {loading ? "Učitavanje..." : `Ukupno: ${rows.length}`}
            </div>
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Lista zahteva</h3>
        </div>

        <div className="card-c" style={{ paddingTop: 0 }}>
          <table className="table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Citizen</th>
                <th>Type</th>
                <th>Status</th>
                <th>Submitted</th>
                <th>Processed</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {rows.map((r) => (
                <tr key={r.id}>
                  <td>
                    <div className="row">
                      <span style={{ fontWeight: 800 }}>{shorten(r.id)}</span>
                      <button
                        className="btn btn-outline"
                        style={{ padding: "6px 10px" }}
                        onClick={() => copy(r.id)}
                      >
                        Copy
                      </button>
                    </div>
                  </td>
                  <td style={{ color: "var(--muted)" }}>{r.citizenId}</td>
                  <td>{r.type}</td>
                  <td>
                    <StatusBadge status={r.status} />
                  </td>
                  <td style={{ color: "var(--muted)" }}>{fmtDate(r.submittedAt)}</td>
                  <td style={{ color: "var(--muted)" }}>{fmtDate(r.processedAt)}</td>
                  <td style={{ whiteSpace: "nowrap" }}>
                    <Link to={`/requests/${r.id}`} className="btn btn-primary" style={{ padding: "8px 12px" }}>
                      Detalji
                    </Link>
                  </td>
                </tr>
              ))}

              {!loading && rows.length === 0 && (
                <tr>
                  <td colSpan={7} style={{ color: "var(--muted)" }}>
                    Nema rezultata. Probaj da kreiraš novi zahtev ili promeni filtere.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Modal: Novi zahtev */}
      {open && (
        <div
          role="dialog"
          aria-modal="true"
          onClick={() => setOpen(false)}
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
            style={{ width: "min(520px, 100%)" }}
          >
            <div className="card-h">
              <h3 className="card-t">Novi zahtev</h3>
            </div>
            <div className="card-c">
              <div className="grid" style={{ gap: 10 }}>
                <div>
                  <div className="small" style={{ marginBottom: 6 }}>Citizen ID</div>
                  <input
                    className="input"
                    placeholder="Unesi citizenId"
                    value={createForm.citizenId}
                    onChange={(e) =>
                      setCreateForm((p) => ({ ...p, citizenId: e.target.value }))
                    }
                  />
                </div>

                <div>
                  <div className="small" style={{ marginBottom: 6 }}>Tip</div>
                  <select
                    className="select"
                    value={createForm.type}
                    onChange={(e) =>
                      setCreateForm((p) => ({ ...p, type: e.target.value }))
                    }
                  >
                    {typeOptions.map((t) => (
                      <option key={t} value={t}>{t}</option>
                    ))}
                  </select>
                </div>

                <div className="row" style={{ marginTop: 6 }}>
                  <button className="btn btn-primary" onClick={submitNewRequest}>
                    Pošalji zahtev
                  </button>
                  <button className="btn btn-secondary" onClick={() => setOpen(false)}>
                    Otkaži
                  </button>
                </div>

                <div className="small">
                  Tok demo: napravi građanina → kreiraj zahtev → Detalji → status → uplata → sertifikat.
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
