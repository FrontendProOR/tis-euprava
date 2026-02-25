import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { useCitizenContext } from "../context/CitizenContext";
import { labelRequestType } from "../utils/labels";
import StatusBadge from "../components/StatusBadge";
import { getRequestById, updateRequestStatus } from "../api/requests";
import { createPayment } from "../api/payments";
import { downloadCertificate } from "../api/certificates";
import { getCitizenById } from "../api/citizens";
import { fmtRSD } from "../utils/pricing";

export default function RequestDetails() {
  const { roleMode } = useCitizenContext();
  const { id } = useParams();
  const [request, setRequest] = useState(null);
  const [citizen, setCitizen] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  async function load() {
    setLoading(true);
    try {
      const data = await getRequestById(id);
      setRequest(data);

      try {
        if (data?.citizenId) {
          const c = await getCitizenById(data.citizenId);
          setCitizen(c);
        } else {
          setCitizen(null);
        }
      } catch {
        setCitizen(null);
      }
    } finally {
      setLoading(false);
    }
  }

  async function pay() {
    try {
      await createPayment({
        requestId: id,
        amount: Number(request?.price ?? 0),
        reference: `WEB-${String(Math.floor(Math.random()*1e6)).padStart(6,"0")}`,
      });
      await load();
      alert("Uplata evidentirana. Status je automatski ažuriran.");
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri uplati");
    }
  }

  async function startReview() {
    try {
      await updateRequestStatus(id, "IN_REVIEW");
      await load();
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Ne mogu da promenim status");
    }
  }

  async function download() {
    try {
      await downloadCertificate(id);
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Ne mogu da preuzmem sertifikat");
    }
  }

  if (loading && !request) return <div>Učitavanje...</div>;
  if (!request) return <div>Ne postoji zahtev.</div>;

  const citizenName = citizen
    ? `${(citizen.first_name ?? citizen.firstName ?? "")} ${(citizen.last_name ?? citizen.lastName ?? "")}`.trim()
    : "(nepoznat)";

  const canDownload = (request.status === "APPROVED" || request.status === "COMPLETED") && request.paid;

  return (
    <div className="grid grid-3">
      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Podaci o zahtevu</h3>
        </div>
        <div className="card-c">
          <p><b>Građanin:</b> {citizenName}</p>
          <p><b>Tip:</b> {labelRequestType(request.type)}</p>
          <p><b>Cena:</b> {fmtRSD(request.price)}</p>
          <p><b>Status:</b> <StatusBadge status={request.status} /></p>
          <p><b>Plaćeno:</b> {request.paid ? <StatusBadge status="PAID" /> : <StatusBadge status="UNPAID" />}</p>
        </div>
      </div>

      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Obrada & Uplata</h3>
        </div>
        <div className="card-c">
          {roleMode === "officer" && request.status === "PENDING" && (
            <>
              <div className="small" style={{ marginBottom: 10, color: "var(--muted)" }}>
                Prvi korak je da službenik prebaci zahtev iz <b>PENDING</b> u <b>IN_REVIEW</b>.
              </div>
              <button className="btn btn-outline" onClick={startReview}>
                Počni obradu (IN_REVIEW)
              </button>
              <div className="divider" style={{ margin: "14px 0" }} />
              <div className="row" style={{ gap: 10, flexWrap: "wrap" }}>
                <button className="btn btn-outline" type="button" onClick={async () => { await updateRequestStatus(id, "APPROVED"); await load(); }}>
                  Odobri (APPROVED)
                </button>
                <button className="btn btn-danger" type="button" onClick={async () => { await updateRequestStatus(id, "REJECTED"); await load(); }}>
                  Odbij (REJECTED)
                </button>
              </div>

            </>
          )}

          {!request.paid && request.status !== "REJECTED" && (
            <>
              <div className="small" style={{ marginBottom: 10, color: "var(--muted)" }}>
                Nakon uplate, sistem automatski odlučuje <b>APPROVED</b> ili <b>REJECTED</b> (u zavisnosti od iznosa).
              </div>
              <button className="btn btn-primary" onClick={pay}>
                Plati {fmtRSD(request.price)}
              </button>
            </>
          )}

          {request.paid && (
            <>
              <StatusBadge status="PAID" />
              <div className="small" style={{ marginTop: 10, color: "var(--muted)" }}>
                Reference: {request.payment?.reference ?? "WEB-001"} · Uplaćeno: {fmtRSD(request.payment?.amount)}
              </div>
            </>
          )}
        </div>
      </div>

      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Sertifikat</h3>
        </div>
        <div className="card-c">
          <button className="btn btn-primary" disabled={!canDownload} onClick={download}>
            Preuzmi sertifikat (PDF)
          </button>
          <div className="small" style={{ marginTop: 10, color: "var(--muted)" }}>
            Sertifikat je dostupan kada je zahtev <b>APPROVED</b> (ili <b>COMPLETED</b>) i kada je uplata evidentirana.
          </div>
        </div>
      </div>
    </div>
  );
}
