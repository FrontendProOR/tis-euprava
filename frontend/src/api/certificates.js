const API = import.meta.env.VITE_API_URL;

function todayStamp() {
  const d = new Date();
  const pad = (n) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
}

export async function downloadCertificate(requestId) {
  const res = await fetch(`${API}/api/requests/${requestId}/certificate`);
  if (!res.ok) throw new Error(await res.text());

  const blob = await res.blob();
  const url = URL.createObjectURL(blob);

  const a = document.createElement("a");
  a.href = url;
  a.download = `sertifikat-${todayStamp()}.pdf`; // ne prikazujemo interne ID-jeve
  document.body.appendChild(a);
  a.click();
  a.remove();

  URL.revokeObjectURL(url);
}
