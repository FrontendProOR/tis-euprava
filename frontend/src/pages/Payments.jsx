import { useState } from "react";
import StatusBadge from "../components/StatusBadge";
import { createPayment } from "../api/payments";

export default function Payments() {
  const [form, setForm] = useState({
    requestId: "",
    amount: 3500,
    reference: "WEB-001",
  });

  async function pay() {
    if (!form.requestId.trim()) {
      alert("Unesi requestId");
      return;
    }
    try {
      const res = await createPayment({
        requestId: form.requestId.trim(),
        amount: Number(form.amount),
        reference: form.reference || "WEB-001",
      });
      alert("Uplata uspešna");
      console.log("Payment:", res);
      setForm((p) => ({ ...p, requestId: "" })); // opcionalno: očisti requestId
    } catch (e) {
      console.error(e);
      alert(e?.response?.data?.message ?? "Greška pri uplati");
    }
  }

  return (
    <div className="grid">
      <div className="row row-between">
        <div>
          <h1 className="h1">Uplate</h1>
          <p className="p">Forma za uplatu po requestId (prečica).</p>
        </div>
        <StatusBadge status="PAID" />
      </div>

      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Nova uplata</h3>
        </div>
        <div className="card-c">
          <div className="grid grid-3">
            <div>
              <div className="small" style={{ marginBottom: 6 }}>Request ID</div>
              <input
                className="input"
                placeholder="Unesi requestId"
                value={form.requestId}
                onChange={(e) =>
                  setForm((p) => ({ ...p, requestId: e.target.value }))
                }
              />
            </div>

            <div>
              <div className="small" style={{ marginBottom: 6 }}>Iznos</div>
              <input
                className="input"
                type="number"
                value={form.amount}
                onChange={(e) =>
                  setForm((p) => ({ ...p, amount: e.target.value }))
                }
              />
            </div>

            <div>
              <div className="small" style={{ marginBottom: 6 }}>Poziv na broj</div>
              <input
                className="input"
                value={form.reference}
                onChange={(e) =>
                  setForm((p) => ({ ...p, reference: e.target.value }))
                }
              />
            </div>
          </div>

          <div className="row" style={{ marginTop: 12 }}>
            <button className="btn btn-primary" onClick={pay}>Plati</button>
            <span className="small">
              Tipično: 3500 RSD · reference: WEB-001
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
