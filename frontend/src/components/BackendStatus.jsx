import { useEffect, useState } from "react";

export default function BackendStatus() {
  const [online, setOnline] = useState(null);

  useEffect(() => {
    let alive = true;

    async function ping() {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL}/health`);
        if (!alive) return;
        setOnline(res.ok);
      } catch {
        if (!alive) return;
        setOnline(false);
      }
    }

    ping();
    const t = setInterval(ping, 4000);
    return () => {
      alive = false;
      clearInterval(t);
    };
  }, []);

  return (
    <div className="row">
      <span className={`dot ${online ? "dot-green" : "dot-red"}`} />
      <span className="small">{online ? "Backend online" : "Backend offline"}</span>
    </div>
  );
}
