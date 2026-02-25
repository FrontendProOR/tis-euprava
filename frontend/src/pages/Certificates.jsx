import { useEffect, useMemo, useState } from "react";
import { downloadCertificate } from "../api/certificates";
import CitizenSelect from "../components/CitizenSelect";
import { useCitizenContext } from "../context/CitizenContext";
import { getRequests } from "../api/requests";
import StatusBadge from "../components/StatusBadge";
import { fmtRSD } from "../utils/pricing";
import { labelRequestType, labelRequestStatus } from "../utils/labels";

export default function Certificates() {
  const { selectedCitizenId } = useCitizenContext();
  const [requestId, setRequestId] = useState("");
  const [requests, setRequests] = useState([]);

  useEffect(() => {
    setRequestId("");
    setRequests([]);
    if (!selectedCitizenId) return;

    (async () => {
      const data = await getRequests({ citizenId: selectedCitizenId });
      const list = Array.isArray(data) ? data : data?.items ?? [];
      setRequests(list);
    })().catch(console.error);
  }, [selectedCitizenId]);

  const eligible = useMemo(() => {
    return requests.filter((r) => (r.status === "APPROVED" || r.status === "COMPLETED") && r.paid);
  }, [requests]);

  const requestOptions = useMemo(() => {
    return eligible.map((r) => ({
      id: r.id,
      label: `${labelRequestType(r.type)} · ${fmtRSD(r.price)} · ${labelRequestStatus(r.status)} · Plaćeno`,
    }));
  }, [eligible]);

  async function download() {
    if (!requestId) {
      alert("Izaberi zahtev");
      return;
    }
    try {
      await downloadCertificate(requestId);
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Ne mogu da preuzmem sertifikat");
    }
  }

  return (
    <div className="grid" style={{ gap: 16 }}>
      <div className="row row-between">
        <div>
          <h1 className="h1">Sertifikati</h1>
          <p className="p">Sertifikat je dostupan samo za zahteve koji su odobreni i plaćeni.</p>
        </div>
        <StatusBadge status="APPROVED" />
      </div>

      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Preuzmi sertifikat</h3>
        </div>
        <div className="card-c">
          <div className="grid" style={{ gap: 10 }}>
            <CitizenSelect label="Građanin" required={true} />

            <div>
              <div className="small" style={{ marginBottom: 6 }}>Zahtev</div>
              <select className="select" value={requestId} onChange={(e) => setRequestId(e.target.value)}>
                <option value="">Izaberi zahtev</option>
                {requestOptions.map((o) => (
                  <option key={o.id} value={o.id}>{o.label}</option>
                ))}
              </select>
            </div>

            <button className="btn btn-primary" onClick={download} disabled={!requestId}>
              Preuzmi PDF
            </button>

            {eligible.length === 0 && (
              <div className="small" style={{ color: "var(--muted)" }}>
                Trenutno nema dostupnih sertifikata. Potrebno je: uplata + status APPROVED.
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
