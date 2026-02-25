import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import StatusBadge from "../components/StatusBadge";
import SearchInput from "../components/SearchInput";
import { getRequests, updateRequestStatus } from "../api/requests";
import { useCitizenContext } from "../context/CitizenContext";
import { labelRequestType } from "../utils/labels";

function fmtDate(s) {
  if (!s) return "";
  const d = new Date(s);
  if (Number.isNaN(d.getTime())) return s;
  return d.toLocaleString();
}

export default function OfficerPanel() {
  const { citizens } = useCitizenContext();
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(false);
  const [q, setQ] = useState("");

  function citizenName(id) {
    const c = citizens.find((x) => x.id === id);
    if (!c) return "(nepoznat)";
    return `${c.first_name ?? c.firstName ?? ""} ${c.last_name ?? c.lastName ?? ""}`.trim() || "(bez imena)";
  }

  async function load() {
    setLoading(true);
    try {
      const data = await getRequests();
      const list = Array.isArray(data) ? data : data?.items ?? [];
      // službenik: vidi samo otvorene / za obradu
      setRows(list.filter((r) => ["PENDING", "IN_REVIEW"].includes(String(r.status).toUpperCase())));
    } catch (e) {
      console.error(e);
      alert(e?.message ?? "Greška pri učitavanju");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  const visible = useMemo(() => {
    const term = q.trim().toLowerCase();
    if (!term) return rows;
    return rows.filter((r) => {
      const cn = citizenName(r.citizenId ?? r.citizen_id).toLowerCase();
      return (
        cn.includes(term) ||
        String(r.type || "").toLowerCase().includes(term) ||
        String(r.status || "").toLowerCase().includes(term)
      );
    });
  }, [rows, q, citizens]);

  async function setStatus(id, status) {
    try {
      await updateRequestStatus(id, status);
      await load();
    } catch (e) {
      console.error(e);
      alert(e?.message ?? "Greška pri promeni statusa");
    }
  }

  return (
    <div className="page">
      <div className="page-head">
        <div>
          <h1 className="h1">Službenik MUP-a</h1>
          <p className="p">Panel za obradu i validaciju zahteva (bez logina – demo režim).</p>
        </div>
        <button className="btn" onClick={load} disabled={loading} type="button">
          {loading ? "Učitavam..." : "Osveži"}
        </button>
      </div>

      <div className="card pad">
        <SearchInput value={q} onChange={setQ} placeholder="Pretraga (građanin / tip / status)" />
      </div>

      <div className="card" style={{ marginTop: 14 }}>
        <div className="table-wrap">
          <table className="table">
            <thead>
              <tr>
                <th>Građanin</th>
                <th>Tip</th>
                <th>Status</th>
                <th>Podnet</th>
                <th>Akcije</th>
              </tr>
            </thead>
            <tbody>
              {visible.map((r) => (
                <tr key={r.id}>
                  <td>{citizenName(r.citizenId ?? r.citizen_id)}</td>
                  <td>{labelRequestType(r.type)}</td>
                  <td><StatusBadge status={r.status} /></td>
                  <td className="muted">{fmtDate(r.submittedAt)}</td>
                  <td>
                    <div className="row" style={{ gap: 8, flexWrap: "wrap" }}>
                      <Link className="btn btn-outline" to={`/requests/${r.id}`}>Detalji</Link>
                      <button className="btn btn-outline" type="button" onClick={() => setStatus(r.id, "IN_REVIEW")}>
                        U obradi
                      </button>
                      <button className="btn btn-outline" type="button" onClick={() => setStatus(r.id, "APPROVED")}>
                        Odobri
                      </button>
                      <button className="btn btn-danger" type="button" onClick={() => setStatus(r.id, "REJECTED")}>
                        Odbij
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {!visible.length && (
                <tr>
                  <td colSpan={5} className="muted" style={{ padding: 16 }}>
                    Nema zahteva za obradu.
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
