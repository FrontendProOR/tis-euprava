import { useState } from "react";
import StatusBadge from "../components/StatusBadge";
import CitizenSelect from "../components/CitizenSelect";
import { useCitizenContext } from "../context/CitizenContext";
import {
  createAppointment,
  getAppointmentsByCitizen,
  cancelAppointment,
} from "../api/appointments";

export default function Appointments() {
  const { selectedCitizenId } = useCitizenContext();

  const [form, setForm] = useState({
    policeStation: "",
    dateTime: "",
  });

  const [list, setList] = useState([]);

  async function create() {
    if (!selectedCitizenId) {
      alert("Izaberi građanina");
      return;
    }
    await createAppointment({ ...form, citizenId: selectedCitizenId });
    alert("Termin zakazan");
  }

  async function load() {
    if (!selectedCitizenId) {
      alert("Izaberi građanina");
      return;
    }
    const data = await getAppointmentsByCitizen(selectedCitizenId);
    setList(data);
  }

  async function cancel(id) {
    await cancelAppointment(id);
    await load();
  }

  return (
    <div className="grid grid-2">
      {/* ZAKAŽI */}
      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Zakaži termin</h3>
        </div>
        <div className="card-c">
          <CitizenSelect label="Građanin" required={true} />
          <input className="input" placeholder="Policijska stanica"
            onChange={(e) => setForm({ ...form, policeStation: e.target.value })} />
          <input className="input" type="datetime-local"
            onChange={(e) => setForm({ ...form, dateTime: e.target.value })} />

          <button className="btn btn-primary" onClick={create}>
            Zakaži
          </button>
        </div>
      </div>

      {/* LISTA */}
      <div className="card">
        <div className="card-h">
          <h3 className="card-t">Termini za građanina</h3>
        </div>
        <div className="card-c">
          <CitizenSelect label="Građanin" required={true} />
          <div className="row" style={{ marginTop: 10 }}>
            <button className="btn btn-secondary" onClick={load}>Prikaži</button>
          </div>

          <table className="table" style={{ marginTop: 10 }}>
            <thead>
              <tr>
                <th>Datum</th>
                <th>Stanica</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {list.map((a) => (
                <tr key={a.id}>
                  <td>{a.dateTime}</td>
                  <td>{a.policeStation}</td>
                  <td><StatusBadge status={a.status} /></td>
                  <td>
                    {a.status === "SCHEDULED" && (
                      <button className="btn btn-danger" onClick={() => cancel(a.id)}>
                        Otkaži
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
