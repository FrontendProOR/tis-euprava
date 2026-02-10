import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import StatusBadge from "../components/StatusBadge";
import { getRequestById, updateRequestStatus } from "../api/requests";
import { createPayment } from "../api/payments";
import { downloadCertificate } from "../api/certificates";

export default function RequestDetails() {
  const { id } = useParams();
  const [request, setRequest] = useState(null);
  const [nextStatus, setNextStatus] = useState("");

  useEffect(() => {
    load();
  }, [id]);

  async function load() {
    const data = await getRequestById(id);
    setRequest(data);
  }

  async function changeStatus() {
    await updateRequestStatus(id, nextStatus);
    await load();
  }

  async function pay() {
    await createPayment({
      requestId: id,
      amount: 3500,
      reference: "WEB-001",
    });
    await load();
  }

  async function download() {
    await downloadCertificate(id);
  }

  if (!request) return <div>Učitavanje...</div>;

  return (
    <div className="grid grid-3">
      {/* PODACI O ZAHTEVU */}
      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Podaci o zahtevu</h3>
        </div>
        <div className="card-c">
          <p><b>ID:</b> {request.id}</p>
          <p><b>Građanin:</b> {request.citizenId}</p>
          <p><b>Tip:</b> {request.type}</p>
          <p>
            <b>Status:</b>{" "}
            <StatusBadge status={request.status} />
          </p>
        </div>
      </div>

      {/* STATUS */}
      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Promena statusa</h3>
        </div>
        <div className="card-c">
          <select
            className="select"
            value={nextStatus}
            onChange={(e) => setNextStatus(e.target.value)}
          >
            <option value="">Izaberi sledeći status</option>
            {request.status === "SUBMITTED" && (
              <option value="IN_PROCESS">IN_PROCESS</option>
            )}
            {request.status === "IN_PROCESS" && (
              <>
                <option value="APPROVED">APPROVED</option>
                <option value="REJECTED">REJECTED</option>
              </>
            )}
          </select>

          <button className="btn btn-primary" style={{ marginTop: 10 }} onClick={changeStatus}>
            Sačuvaj status
          </button>
        </div>
      </div>

      {/* UPLATA + SERTIFIKAT */}
      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Uplata & Sertifikat</h3>
        </div>
        <div className="card-c">
          {!request.payment && (
            <button className="btn btn-secondary" onClick={pay}>
              Plati 3.500 RSD
            </button>
          )}

          {request.payment && (
            <>
              <StatusBadge status="PAID" />
              <br /><br />
            </>
          )}

          <button
            className="btn btn-primary"
            disabled={request.status !== "APPROVED" || !request.payment}
            onClick={download}
          >
            Preuzmi sertifikat (PDF)
          </button>
        </div>
      </div>
    </div>
  );
}
