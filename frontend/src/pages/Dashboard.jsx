import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import StatusBadge from "../components/StatusBadge";
import SearchInput from "../components/SearchInput";
import { getRequests, updateRequestStatus } from "../api/requests";
import { getCitizens } from "../api/citizens";
import { fmtRSD, priceFor } from "../utils/pricing";

function fmtDate(s) {
  if (!s) return "";
  const d = new Date(s);
  if (Number.isNaN(d.getTime())) return s;
  return d.toLocaleString();
}

export default function Dashboard() {
  const [citizens, setCitizens] = useState([]);
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(false);
  const [q, setQ] = useState("");

  useEffect(() => {
    refresh();
  }, []);

  async function refresh() {
    setLoading(true);
    try {
      const [cits, reqs] = await Promise.all([getCitizens(), getRequests()]);
      const list = Array.isArray(reqs) ? reqs : reqs?.items ?? [];
      setCitizens(Array.isArray(cits) ? cits : cits?.items ?? []);
      setRequests(list);
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri učitavanju");
    } finally {
      setLoading(false);
    }
  }

  const citizenById = useMemo(() => {
    const m = new Map();
    for (const c of citizens) m.set(c.id, c);
    return m;
  }, [citizens]);

  const stats = useMemo(() => {
    const out = {
      total: requests.length,
      PENDING: 0,
      IN_REVIEW: 0,
      APPROVED: 0,
      REJECTED: 0,
      COMPLETED: 0,
    };
    for (const r of requests) {
      const s = String(r.status ?? "").toUpperCase();
      if (out[s] !== undefined) out[s] += 1;
    }
    return out;
  }, [requests]);

  const recent = useMemo(() => {
    const list = [...requests].sort((a, b) => {
      const ta = new Date(a.submittedAt ?? a.submitted_at ?? 0).getTime();
      const tb = new Date(b.submittedAt ?? b.submitted_at ?? 0).getTime();
      return tb - ta;
    });
    return list.slice(0, 10);
  }, [requests]);

  const filtered = useMemo(() => {
    const term = q.trim().toLowerCase();
    if (!term) return recent;
    return recent.filter((r) => {
      const c = citizenById.get(r.citizenId ?? r.citizen_id);
      const name = c ? `${c.first_name ?? c.firstName ?? ""} ${c.last_name ?? c.lastName ?? ""}` : "";
      const jmbg = c?.jmbg ?? "";
      return (
        `${r.id}`.toLowerCase().includes(term) ||
        `${r.type}`.toLowerCase().includes(term) ||
        `${r.status}`.toLowerCase().includes(term) ||
        name.toLowerCase().includes(term) ||
        jmbg.toLowerCase().includes(term)
      );
    });
  }, [q, recent, citizenById]);

  async function moveToReview(r) {
    try {
      await updateRequestStatus(r.id, "IN_REVIEW");
      await refresh();
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Ne mogu da promenim status");
    }
  }

  return (
    <div className="grid" style={{ gap: 16 }}>
      <div className="row row-between" style={{ alignItems: "flex-end" }}>
        <div>
          <h1 className="h1">Dashboard</h1>
          <p className="p">Pregled MUP modula: građani, zahtevi, statusi i poslednje aktivnosti.</p>
        </div>
        <div className="row" style={{ gap: 10 }}>
          <SearchInput value={q} onChange={setQ} placeholder="Pretraga poslednjih zahteva..." />
          <button className="btn btn-secondary" onClick={refresh} disabled={loading}>
            Osveži
          </button>
        </div>
      </div>

      <div className="grid grid-3">
        <div className="card kpi">
          <div className="kpi-num">{stats.total}</div>
          <div className="kpi-lbl">Ukupno zahteva</div>
        </div>
        <div className="card kpi">
          <div className="kpi-num">{stats.PENDING + stats.IN_REVIEW}</div>
          <div className="kpi-lbl">Otvoreni (PENDING + IN_REVIEW)</div>
        </div>
        <div className="card kpi">
          <div className="kpi-num">{stats.APPROVED + stats.COMPLETED}</div>
          <div className="kpi-lbl">Pozitivni (APPROVED + COMPLETED)</div>
        </div>
      </div>

      <div className="grid grid-3">
        <div className="card kpi">
          <div className="kpi-num">{stats.PENDING}</div>
          <div className="kpi-lbl">PENDING</div>
        </div>
        <div className="card kpi">
          <div className="kpi-num">{stats.IN_REVIEW}</div>
          <div className="kpi-lbl">IN_REVIEW</div>
        </div>
        <div className="card kpi">
          <div className="kpi-num">{stats.REJECTED}</div>
          <div className="kpi-lbl">REJECTED</div>
        </div>
      </div>

      <div className="card">
        <div className="card-h" style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
          <h3 className="card-t">Poslednji zahtevi</h3>
          <div className="small">Prikaz: {filtered.length} / {recent.length}</div>
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
                <th></th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((r) => {
                const cid = r.citizenId ?? r.citizen_id;
                const c = citizenById.get(cid);
                const name = c
                  ? `${c.first_name ?? c.firstName ?? ""} ${c.last_name ?? c.lastName ?? ""}`.trim()
                  : "(nepoznat)";
                const submitted = r.submittedAt ?? r.submitted_at;
                const price = r.price ?? priceFor(r.type);
                return (
                  <tr key={r.id}>
                    <td>{name}</td>
                    <td>{r.type}</td>
                    <td>{fmtRSD(price)}</td>
                    <td><StatusBadge status={r.status} /></td>
                    <td><StatusBadge status={r.paid ? "PAID" : "UNPAID"} /></td>
                    <td style={{ color: "var(--muted)" }}>{fmtDate(submitted)}</td>
                    <td style={{ whiteSpace: "nowrap", display: "flex", justifyContent: "flex-end", gap: 8 }}>
                      {String(r.status).toUpperCase() === "PENDING" && (
                        <button className="btn btn-secondary" onClick={() => moveToReview(r)}>
                          U obradi
                        </button>
                      )}
                      <Link to={`/requests/${r.id}`} className="btn btn-primary" style={{ padding: "8px 12px" }}>
                        Detalji
                      </Link>
                    </td>
                  </tr>
                );
              })}

              {!loading && filtered.length === 0 && (
                <tr>
                  <td colSpan={7} style={{ color: "var(--muted)" }}>
                    Nema zahteva za prikaz.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
