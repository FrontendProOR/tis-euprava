import { useEffect, useMemo, useState } from "react";
import StatusBadge from "../components/StatusBadge";
import SearchInput from "../components/SearchInput";
import Modal from "../components/Modal";
import { createCitizen, getCitizens } from "../api/citizens";
import { getRequests } from "../api/requests";
import { useCitizenContext } from "../context/CitizenContext";

function fmtDate(s) {
  if (!s) return "";
  const d = new Date(s);
  if (Number.isNaN(d.getTime())) return s;
  return d.toLocaleString();
}

export default function Citizens() {
  const { selectedCitizenId, setSelectedCitizenId } = useCitizenContext();

  const [citizens, setCitizens] = useState([]);
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(false);
  const [q, setQ] = useState("");

  // Create citizen modal
  const [openCreate, setOpenCreate] = useState(false);
  const [form, setForm] = useState({
    firstName: "",
    lastName: "",
    jmbg: "",
    email: "",
    phone: "",
    address: "",
  });

  // Requests modal per citizen
  const [openReq, setOpenReq] = useState(false);
  const [reqCitizenId, setReqCitizenId] = useState(null);

  useEffect(() => {
    refresh();
  }, []);

  async function refresh() {
    setLoading(true);
    try {
      const [cits, reqs] = await Promise.all([getCitizens(), getRequests()]);
      setCitizens(Array.isArray(cits) ? cits : cits?.items ?? []);
      setRequests(Array.isArray(reqs) ? reqs : reqs?.items ?? []);
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri učitavanju građana");
    } finally {
      setLoading(false);
    }
  }

  const requestsByCitizen = useMemo(() => {
    const m = new Map();
    for (const r of requests) {
      const cid = r.citizenId ?? r.citizen_id;
      if (!cid) continue;
      if (!m.has(cid)) m.set(cid, []);
      m.get(cid).push(r);
    }
    return m;
  }, [requests]);

  const rows = useMemo(() => {
    const term = q.trim().toLowerCase();
    const list = citizens.map((c) => {
      const rs = requestsByCitizen.get(c.id) ?? [];
      const active = rs.filter((r) => !["COMPLETED", "REJECTED"].includes(String(r.status).toUpperCase())).length;
      return {
        ...c,
        __requests: rs,
        __active: active,
        __total: rs.length,
      };
    });
    if (!term) return list;
    return list.filter((c) => {
      const fn = c.first_name ?? c.firstName ?? "";
      const ln = c.last_name ?? c.lastName ?? "";
      const jmbg = c.jmbg ?? "";
      const email = c.email ?? "";
      return (
        `${c.id}`.toLowerCase().includes(term) ||
        fn.toLowerCase().includes(term) ||
        ln.toLowerCase().includes(term) ||
        jmbg.toLowerCase().includes(term) ||
        email.toLowerCase().includes(term)
      );
    });
  }, [citizens, requestsByCitizen, q]);

  function citizenName(c) {
    const fn = c.first_name ?? c.firstName ?? "";
    const ln = c.last_name ?? c.lastName ?? "";
    return `${fn} ${ln}`.trim();
  }

  async function submit() {
    if (!form.firstName.trim() || !form.lastName.trim() || !form.jmbg.trim()) {
      alert("Ime, prezime i JMBG su obavezni");
      return;
    }
    try {
      await createCitizen({
        firstName: form.firstName,
        lastName: form.lastName,
        jmbg: form.jmbg,
        email: form.email,
        phone: form.phone,
        address: form.address,
      });
      setOpenCreate(false);
      setForm({ firstName: "", lastName: "", jmbg: "", email: "", phone: "", address: "" });
      await refresh();
      alert("Građanin uspešno kreiran");
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri kreiranju građanina");
    }
  }

  const reqCitizen = useMemo(() => citizens.find((c) => c.id === reqCitizenId), [citizens, reqCitizenId]);
  const reqList = useMemo(() => {
    const list = reqCitizenId ? (requestsByCitizen.get(reqCitizenId) ?? []) : [];
    return [...list].sort((a, b) => {
      const ta = new Date(a.submittedAt ?? a.submitted_at ?? 0).getTime();
      const tb = new Date(b.submittedAt ?? b.submitted_at ?? 0).getTime();
      return tb - ta;
    });
  }, [reqCitizenId, requestsByCitizen]);

  return (
    <div className="grid" style={{ gap: 16 }}>
      <div className="row row-between" style={{ alignItems: "flex-end" }}>
        <div>
          <h1 className="h1">Građani</h1>
          <p className="p">Lista građana + njihovi zahtevi. Klikom postavljaš aktivnog građanina za brži rad.</p>
        </div>
        <div className="row" style={{ gap: 10 }}>
          <SearchInput value={q} onChange={setQ} placeholder="Pretraga (ime, prezime, JMBG, email, id)..." />
          <button className="btn btn-secondary" onClick={refresh} disabled={loading}>
            Osveži
          </button>
          <button className="btn btn-primary" onClick={() => setOpenCreate(true)}>
            + Novi građanin
          </button>
        </div>
      </div>

      <div className="card">
        <div className="card-h" style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
          <h3 className="card-t">Lista građana</h3>
          <div className="small">{loading ? "Učitavanje..." : `Ukupno: ${rows.length}`}</div>
        </div>
        <div className="card-c" style={{ paddingTop: 0 }}>
          <table className="table">
            <thead>
              <tr>
                <th>Ime i prezime</th>
                <th>JMBG</th>
                <th>Email</th>
                <th>Aktivni zahtevi</th>
                <th>Ukupno</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {rows.map((c) => (
                <tr key={c.id}>
                  <td>
                    <div style={{ fontWeight: 800 }}>{citizenName(c) || "(bez imena)"}</div>
                    <div className="small">ID: {c.id}</div>
                  </td>
                  <td>{c.jmbg}</td>
                  <td style={{ color: "var(--muted)" }}>{c.email ?? ""}</td>
                  <td>{c.__active}</td>
                  <td>{c.__total}</td>
                  <td style={{ whiteSpace: "nowrap", display: "flex", justifyContent: "flex-end", gap: 8 }}>
                    <button
                      className={selectedCitizenId === c.id ? "btn btn-primary" : "btn btn-secondary"}
                      onClick={() => setSelectedCitizenId(c.id)}
                    >
                      {selectedCitizenId === c.id ? "Aktivan" : "Postavi"}
                    </button>
                    <button
                      className="btn btn-outline"
                      onClick={() => {
                        setReqCitizenId(c.id);
                        setOpenReq(true);
                      }}
                    >
                      Zahtevi
                    </button>
                  </td>
                </tr>
              ))}

              {!loading && rows.length === 0 && (
                <tr>
                  <td colSpan={6} style={{ color: "var(--muted)" }}>
                    Nema građana. Kreiraj novog građanina.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Modal: Create citizen */}
      <Modal open={openCreate} title="Novi građanin" onClose={() => setOpenCreate(false)} width={700}>
        <div className="grid grid-2" style={{ gap: 12 }}>
          <div>
            <div className="small" style={{ marginBottom: 6 }}>Ime *</div>
            <input className="input" value={form.firstName} onChange={(e) => setForm((p) => ({ ...p, firstName: e.target.value }))} />
          </div>
          <div>
            <div className="small" style={{ marginBottom: 6 }}>Prezime *</div>
            <input className="input" value={form.lastName} onChange={(e) => setForm((p) => ({ ...p, lastName: e.target.value }))} />
          </div>
          <div>
            <div className="small" style={{ marginBottom: 6 }}>JMBG *</div>
            <input className="input" value={form.jmbg} onChange={(e) => setForm((p) => ({ ...p, jmbg: e.target.value }))} />
          </div>
          <div>
            <div className="small" style={{ marginBottom: 6 }}>Email</div>
            <input className="input" value={form.email} onChange={(e) => setForm((p) => ({ ...p, email: e.target.value }))} />
          </div>
          <div>
            <div className="small" style={{ marginBottom: 6 }}>Telefon</div>
            <input className="input" value={form.phone} onChange={(e) => setForm((p) => ({ ...p, phone: e.target.value }))} />
          </div>
          <div>
            <div className="small" style={{ marginBottom: 6 }}>Adresa</div>
            <input className="input" value={form.address} onChange={(e) => setForm((p) => ({ ...p, address: e.target.value }))} />
          </div>
        </div>

        <div className="row" style={{ marginTop: 12, justifyContent: "flex-end" }}>
          <button className="btn btn-secondary" onClick={() => setOpenCreate(false)}>Otkaži</button>
          <button className="btn btn-primary" onClick={submit}>Sačuvaj</button>
        </div>
      </Modal>

      {/* Modal: Citizen requests */}
      <Modal
        open={openReq}
        title={`Zahtevi – ${reqCitizen ? citizenName(reqCitizen) : ""}`}
        onClose={() => setOpenReq(false)}
        width={900}
      >
        <table className="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Tip</th>
              <th>Status</th>
              <th>Plaćeno</th>
              <th>Podnet</th>
              <th>Obrađen</th>
            </tr>
          </thead>
          <tbody>
            {reqList.map((r) => (
              <tr key={r.id}>
                <td>{r.id}</td>
                <td>{r.type}</td>
                <td><StatusBadge status={r.status} /></td>
                <td><StatusBadge status={r.paid ? "PAID" : "UNPAID"} /></td>
                <td style={{ color: "var(--muted)" }}>{fmtDate(r.submittedAt ?? r.submitted_at)}</td>
                <td style={{ color: "var(--muted)" }}>{fmtDate(r.processedAt ?? r.processed_at)}</td>
              </tr>
            ))}
            {reqList.length === 0 && (
              <tr>
                <td colSpan={6} style={{ color: "var(--muted)" }}>Građanin nema zahteva.</td>
              </tr>
            )}
          </tbody>
        </table>
      </Modal>
    </div>
  );
}
