import { NavLink, Outlet } from "react-router-dom";
import BackendStatus from "./BackendStatus";
import { useCitizenContext } from "../context/CitizenContext";

function Crest() {
  return (
    <div className="crest" aria-hidden="true">
      <div className="crest-inner">RS</div>
    </div>
  );
}

function TopTabs({ roleMode }) {
  const tabs = [
    { to: "/", label: "Početna" },
    { to: "/citizens", label: "Građani" },
    { to: "/requests", label: "Zahtevi" },
    { to: "/payments", label: "Plaćanja" },
    { to: "/appointments", label: "Termini" },
    { to: "/certificates", label: "Sertifikati" },
  ];

  if (roleMode === "officer") {
    tabs.splice(3, 0, { to: "/officer", label: "Službenik" });
  }

  return (
    <nav className="top-tabs">
      {tabs.map((t) => (
        <NavLink
          key={t.to}
          to={t.to}
          end={t.to === "/"}
          className={({ isActive }) => (isActive ? "tab active" : "tab")}
        >
          {t.label}
        </NavLink>
      ))}
    </nav>
  );
}

export default function Layout() {
  const {
    citizens,
    selectedCitizenId,
    setSelectedCitizenId,
    roleMode,
    setRoleMode,
  } = useCitizenContext();

  return (
    <div className="mup-shell">
      <header className="mup-topbar">
        <div className="top-left">
          <Crest />
          <div className="brand-title">
            <div className="brand-line-1">E-Servisi MUP</div>
            <div className="brand-line-2">
              Građanski servis za upravne postupke
            </div>
          </div>
        </div>

        <div className="top-right">
          <div className="role-toggle" title="Simulacija bez logina">
            <button
              className={roleMode === "citizen" ? "pill active" : "pill"}
              onClick={() => setRoleMode("citizen")}
              type="button"
            >
              Građanin
            </button>
            <button
              className={roleMode === "officer" ? "pill active" : "pill"}
              onClick={() => setRoleMode("officer")}
              type="button"
            >
              Službenik
            </button>
          </div>

          <div className="active-citizen">
            <div className="small">Aktivni građanin</div>
            <select
              className="select select-compact"
              value={selectedCitizenId}
              onChange={(e) => setSelectedCitizenId(e.target.value)}
              title="Bez logina: simulacija prijavljenog korisnika"
            >
              <option value="">(nije izabran)</option>
              {citizens.map((c) => {
                const name = `${c.first_name ?? c.firstName ?? ""} ${
                  c.last_name ?? c.lastName ?? ""
                }`.trim();

                return (
                  <option key={c.id} value={c.id}>
                    {name || "(bez imena)"}
                    {c.jmbg ? ` · ${c.jmbg}` : ""}
                  </option>
                );
              })}
            </select>
          </div>

          <BackendStatus />
        </div>
      </header>

      <div className="mup-subbar">
        <TopTabs roleMode={roleMode} />
      </div>

      <main className="mup-container">
        <Outlet />
      </main>

      <footer className="mup-footer">
        Portal MUP · Informacioni sistem za upravne postupke
      </footer>
    </div>
  );
}