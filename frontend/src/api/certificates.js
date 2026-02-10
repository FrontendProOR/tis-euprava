const API = import.meta.env.VITE_API_URL;

export async function downloadCertificate(requestId) {
  const res = await fetch(`${API}/api/requests/${requestId}/certificate`);
  if (!res.ok) throw new Error(await res.text());

  const blob = await res.blob();
  const url = URL.createObjectURL(blob);

  const a = document.createElement("a");
  a.href = url;
  a.download = `certificate-${requestId}.pdf`;
  document.body.appendChild(a);
  a.click();
  a.remove();

  URL.revokeObjectURL(url);
}
