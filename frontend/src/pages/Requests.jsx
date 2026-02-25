import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import StatusBadge from "../components/StatusBadge";
import { labelRequestType } from "../utils/labels";
import SearchInput from "../components/SearchInput";
import { getRequests, createRequest } from "../api/requests";
import { createPayment } from "../api/payments";
import CitizenSelect from "../components/CitizenSelect";
import { useCitizenContext } from "../context/CitizenContext";
import { fmtRSD, priceFor } from "../utils/pricing";

function fmtDate(s) {
  if (!s) return "";
  const d = new Date(s);
  if (Number.isNaN(d.getTime())) return s;
  return d.toLocaleString();
}

export default function Requests() {
  const { citizens, selectedCitizenId } = useCitizenContext();

  function citizenLabelById(id) {
    const c = citizens.find((x) => x.id === id);
    if (!c) return "(nepoznat građanin)";
    const fn = c.first_name ?? c.firstName ?? "";
    const ln = c.last_name ?? c.lastName ?? "";
    return `${fn} ${ln}`.trim() || "(bez imena)";
  }

  const [filters, setFilters] = useState({ status: "", type: "" });
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(false);
  const [q, setQ] = useState("");

  const [open, setOpen] = useState(false);
  const [createForm, setCreateForm] = useState({ type: "ID_CARD" });

  const statusOptions = useMemo(
    () => ["PENDING", "IN_REVIEW", "APPROVED", "REJECTED", "COMPLETED"],
    []
  );
  const typeOptions = useMemo(
    () => ["ID_CARD", "PASSPORT", "DRIVER_LICENSE", "RESIDENCE_CHANGE", "CITIZENSHIP"],
    []
  );

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedCitizenId]);

  async function load(overrideFilters) {
    setLoading(true);
    try {
      const used = overrideFilters ?? filters;
      const data = await getRequests({
        citizenId: selectedCitizenId || undefined,
        status: used.status || undefined,
        type: used.type || undefined,
      });
      const list = Array.isArray(data) ? data : data?.items ?? [];
      setRows(list);
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri učitavanju zahteva");
    } finally {
      setLoading(false);
    }
  }

  const visibleRows = useMemo(() => {
    const term = q.trim().toLowerCase();
    if (!term) return rows;
    return rows.filter((r) => {
      const cid = r.citizenId ?? r.citizen_id;
      const label = citizenLabelById(cid);
      return (
        String(r.type).toLowerCase().includes(term) ||
        String(r.status).toLowerCase().includes(term) ||
        String(label).toLowerCase().includes(term)
      );
    });
  }, [rows, q, citizens]);

  async function applyFilters() {
    await load(filters);
  }

  async function submitNewRequest() {
    if (!selectedCitizenId) {
      alert("Izaberi građanina");
      return;
    }
    try {
      await createRequest({
        citizenId: selectedCitizenId,
        type: createForm.type,
      });
      setOpen(false);
      await load(filters);
      alert("Zahtev uspešno kreiran (status: PENDING)");
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri kreiranju zahteva");
    }
  }

  async function payNow(r) {
    try {
      const amount = Number(r.price ?? priceFor(r.type));
      await createPayment({
        requestId: r.id,
        amount,
        reference: "WEB-001",
      });
      await load(filters);
      alert("Uplata evidentirana. Status je automatski ažuriran.");
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri uplati");
    }
  }

  return (
    <div className="grid" style={{ gap: 16 }}>
      <div className="row row-between" style={{ alignItems: "flex-end" }}>
        <div>
          <h1 className="h1">Zahtevi</h1>
          <p className="p">
            Tok: građanin → zahtev (PENDING) → uplata → status (APPROVED/REJECTED) → sertifikat.
          </p>
        </div>

        <div className="row" style={{ gap: 10 }}>
          <SearchInput value={q} onChange={setQ} placeholder="Pretraga zahteva (id, građanin, tip, status)..." />
          <button className="btn btn-primary" onClick={() => setOpen(true)}>
            + Novi zahtev
          </button>
        </div>
      </div>

      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Filteri</h3>
        </div>
        <div className="card-c">
          <div className="grid grid-3" style={{ alignItems: "end" }}>
            <CitizenSelect label="Građanin (filter)" required={false} />

            <div>
              <div className="small" style={{ marginBottom: 6 }}>Status</div>
              <select
                className="select"
                value={filters.status}
                onChange={(e) => setFilters((p) => ({ ...p, status: e.target.value }))}
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
                onChange={(e) => setFilters((p) => ({ ...p, type: e.target.value }))}
              >
                <option value="">Svi</option>
                {typeOptions.map((t) => (
                  <option key={t} value={t}>{t}</option>
                ))}
              </select>
            </div>
          </div>

          <div className="row" style={{ marginTop: 12 }}>
            <button className="btn btn-secondary" onClick={applyFilters}>Primeni</button>
            <button
              className="btn btn-outline"
              onClick={() => {
                const reset = { status: "", type: "" };
                setFilters(reset);
                load(reset);
              }}
            >
              Reset
            </button>

            <div className="small" style={{ marginLeft: "auto" }}>
              {loading ? "Učitavanje..." : `Prikazano: ${visibleRows.length} / ${rows.length}`}
            </div>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Lista zahteva</h3>
        </div>

        <div className="card-c" style={{ paddingTop: 0 }}>
          <table className="table">
            <thead>
              <tr>
                <th>Građanin</th>
                <th>Tip</th>
                <th>Cena</th>
                <th>Status</th>
                <th>Plaćeno</th>
                <th>Podnet</th>
                <th>Obrađen</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {visibleRows.map((r) => (
                <tr key={r.id}>
                  <td style={{ color: "var(--muted)" }}>{citizenLabelById(r.citizenId ?? r.citizen_id)}</td>
                  <td>{labelRequestType(r.type)}</td>
                  <td>{fmtRSD(r.price ?? priceFor(r.type))}</td>
                  <td><StatusBadge status={r.status} /></td>
                  <td>{r.paid ? <StatusBadge status="PAID" /> : <StatusBadge status="UNPAID" />}</td>
                  <td style={{ color: "var(--muted)" }}>{fmtDate(r.submittedAt ?? r.submitted_at)}</td>
                  <td style={{ color: "var(--muted)" }}>{fmtDate(r.processedAt ?? r.processed_at)}</td>
                  <td style={{ whiteSpace: "nowrap", display: "flex", gap: 8, justifyContent: "flex-end" }}>
                    {!r.paid && r.status !== "REJECTED" && (
                      <button className="btn btn-secondary" onClick={() => payNow(r)}>
                        Plati
                      </button>
                    )}
                    <Link to={`/requests/${r.id}`} className="btn btn-primary" style={{ padding: "8px 12px" }}>
                      Detalji
                    </Link>
                  </td>
                </tr>
              ))}

              {!loading && visibleRows.length === 0 && (
                <tr>
                  <td colSpan={8} style={{ color: "var(--muted)" }}>
                    Nema rezultata. Izaberi građanina → kreiraj novi zahtev.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

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
                <CitizenSelect label="Građanin" required={true} />

                <div>
                  <div className="small" style={{ marginBottom: 6 }}>Tip</div>
                  <select
                    className="select"
                    value={createForm.type}
                    onChange={(e) => setCreateForm((p) => ({ ...p, type: e.target.value }))}
                  >
                    {typeOptions.map((t) => (
                      <option key={t} value={t}>
                        {t} · {fmtRSD(priceFor(t))}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="small" style={{ color: "var(--muted)" }}>
                  Cena: <b>{fmtRSD(priceFor(createForm.type))}</b> · početni status: <b>PENDING</b>
                </div>

                <div className="row" style={{ marginTop: 6 }}>
                  <button className="btn btn-primary" onClick={submitNewRequest}>
                    Kreiraj zahtev
                  </button>
                  <button className="btn btn-secondary" onClick={() => setOpen(false)}>
                    Otkaži
                  </button>
                </div>

                <div className="small">
                  Nakon uplate, status se automatski ažurira (APPROVED ako je uplaćeno dovoljno, inače REJECTED).
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
