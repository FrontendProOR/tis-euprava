import { useState } from "react";
import { downloadCertificate } from "../api/certificates";

export default function Certificates() {
  const [requestId, setRequestId] = useState("");

  async function download() {
    if (!requestId.trim()) return alert("Unesi requestId");
    try {
      await downloadCertificate(requestId.trim());
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri preuzimanju sertifikata");
    }
  }

  return (
    <div className="grid">
      <div>
        <h1 className="h1">Sertifikati</h1>
        <p className="p">Prečica za preuzimanje PDF sertifikata po requestId.</p>
      </div>

      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Preuzmi sertifikat</h3>
        </div>
        <div className="card-c">
          <div className="row">
            <input
              className="input"
              placeholder="Request ID"
              value={requestId}
              onChange={(e) => setRequestId(e.target.value)}
            />
            <button className="btn btn-primary" onClick={download}>
              Preuzmi PDF
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
