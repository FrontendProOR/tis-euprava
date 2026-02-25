import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import StatusBadge from "../components/StatusBadge";
import SearchInput from "../components/SearchInput";
import { getRequests } from "../api/requests";
import { createPayment } from "../api/payments";
import { getCitizens } from "../api/citizens";
import { fmtRSD, priceFor } from "../utils/pricing";

function fmtDate(s) {
  if (!s) return "";
  const d = new Date(s);
  if (Number.isNaN(d.getTime())) return s;
  return d.toLocaleString();
}

export default function Payments() {
  const [q, setQ] = useState("");
  const [rows, setRows] = useState([]);
  const [citizens, setCitizens] = useState([]);
  const [loading, setLoading] = useState(false);

  async function load() {
    setLoading(true);
    try {
      const [reqs, cits] = await Promise.all([getRequests({}), getCitizens()]);
      const list = Array.isArray(reqs) ? reqs : reqs?.items ?? [];
      setRows(list);
      setCitizens(Array.isArray(cits) ? cits : cits?.items ?? []);
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri učitavanju uplata");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  const citizenName = (id) => {
    const c = citizens.find((x) => x.id === id);
    if (!c) return "(nepoznat)";
    const fn = c.first_name ?? c.firstName ?? "";
    const ln = c.last_name ?? c.lastName ?? "";
    return `${fn} ${ln}`.trim() || "(bez imena)";
  };

  const visible = useMemo(() => {
    const term = q.trim().toLowerCase();
    if (!term) return rows;
    return rows.filter((r) => {
      const cid = r.citizenId ?? r.citizen_id;
      return (
        `${r.id}`.toLowerCase().includes(term) ||
        `${r.type}`.toLowerCase().includes(term) ||
        `${r.status}`.toLowerCase().includes(term) ||
        citizenName(cid).toLowerCase().includes(term)
      );
    });
  }, [rows, q, citizens]);

  async function pay(r) {
    try {
      const amount = Number(r.price ?? priceFor(r.type));
      await createPayment({ requestId: r.id, amount, reference: "WEB-001" });
      await load();
      alert("Uplata evidentirana. Status zahteva je automatski ažuriran.");
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri uplati");
    }
  }

  return (
    <div className="grid" style={{ gap: 16 }}>
      <div className="row row-between" style={{ alignItems: "flex-end" }}>
        <div>
          <h1 className="h1">Uplate</h1>
          <p className="p">
            Ovde vidiš sve zahteve i da li je uplata evidentirana. Kada je uplata dovoljna, status postaje APPROVED; u suprotnom REJECTED.
          </p>
        </div>
        <div className="row" style={{ gap: 10 }}>
          <SearchInput value={q} onChange={setQ} placeholder="Pretraga (građanin, tip, status, id)..." />
          <button className="btn btn-secondary" onClick={load}>Osveži</button>
        </div>
      </div>

      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Lista uplata</h3>
        </div>
        <div className="card-c" style={{ paddingTop: 0 }}>
          <table className="table">
            <thead>
              <tr>
                <th>Zahtev</th>
                <th>Građanin</th>
                <th>Tip</th>
                <th>Iznos</th>
                <th>Uplata</th>
                <th>Status</th>
                <th>Podnet</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {visible.map((r) => {
                const cid = r.citizenId ?? r.citizen_id;
                const price = Number(r.price ?? priceFor(r.type));
                return (
                  <tr key={r.id}>
                    <td style={{ color: "var(--muted)" }}>{r.id}</td>
                    <td>{citizenName(cid)}</td>
                    <td>{r.type}</td>
                    <td>{fmtRSD(price)}</td>
                    <td>{r.paid ? <StatusBadge status="PAID" /> : <StatusBadge status="UNPAID" />}</td>
                    <td><StatusBadge status={r.status} /></td>
                    <td style={{ color: "var(--muted)" }}>{fmtDate(r.submittedAt ?? r.submitted_at)}</td>
                    <td style={{ whiteSpace: "nowrap", display: "flex", gap: 8, justifyContent: "flex-end" }}>
                      {!r.paid && r.status !== "REJECTED" && (
                        <button className="btn btn-primary" onClick={() => pay(r)}>
                          Plati {fmtRSD(price)}
                        </button>
                      )}
                      <Link to={`/requests/${r.id}`} className="btn btn-outline" style={{ padding: "8px 12px" }}>
                        Detalji
                      </Link>
                    </td>
                  </tr>
                );
              })}

              {!loading && visible.length === 0 && (
                <tr>
                  <td colSpan={8} style={{ color: "var(--muted)" }}>
                    Nema rezultata.
                  </td>
                </tr>
              )}
            </tbody>
          </table>

          {loading && <div className="small" style={{ padding: 12 }}>Učitavanje...</div>}
        </div>
      </div>
    </div>
  );
}
