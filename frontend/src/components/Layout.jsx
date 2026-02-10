import { NavLink, Outlet } from "react-router-dom";
import BackendStatus from "./BackendStatus";

export default function Layout() {
  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <div className="brand-badge">e</div>
          e-Uprava
        </div>

        <nav className="nav">
          <NavLink to="/" end className={({ isActive }) => (isActive ? "active" : "")}>Dashboard</NavLink>
          <NavLink to="/citizens" className={({ isActive }) => (isActive ? "active" : "")}>Građani</NavLink>
          <NavLink to="/requests" className={({ isActive }) => (isActive ? "active" : "")}>Zahtevi</NavLink>
          <NavLink to="/payments" className={({ isActive }) => (isActive ? "active" : "")}>Uplate</NavLink>
          <NavLink to="/appointments" className={({ isActive }) => (isActive ? "active" : "")}>Termini</NavLink>
          <NavLink to="/certificates" className={({ isActive }) => (isActive ? "active" : "")}>Sertifikati</NavLink>
        </nav>

        <div style={{ marginTop: 14 }} className="small">
          Portal MUP · Demo UI
        </div>
      </aside>

      <div className="main">
        <header className="header">
          <div className="row">
            <div style={{ fontWeight: 900 }}>MUP</div>
            <span className="small">Uprava zahteva</span>
          </div>
          <BackendStatus />
        </header>

        <main className="container">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
