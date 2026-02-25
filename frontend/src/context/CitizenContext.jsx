import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { getCitizens } from "../api/citizens";

const CitizenContext = createContext(null);

const LS_KEY = "euprava:selectedCitizenId";
const LS_ROLE = "euprava:roleMode"; // citizen | officer


export function CitizenProvider({ children }) {
  const [citizens, setCitizens] = useState([]);
  const [loading, setLoading] = useState(false);
  const [selectedCitizenId, setSelectedCitizenId] = useState(() => {
    try {
      return localStorage.getItem(LS_KEY) || "";
    } catch {
      return "";
    }
  });

  const [roleMode, setRoleMode] = useState(() => {
    try {
      return localStorage.getItem(LS_ROLE) || "citizen";
    } catch {
      return "citizen";
    }
  });

  useEffect(() => {
    try {
      if (selectedCitizenId) localStorage.setItem(LS_KEY, selectedCitizenId);
      else localStorage.removeItem(LS_KEY);
    } catch {
      // ignore
    }
  }, [selectedCitizenId]);

  useEffect(() => {
    try {
      if (roleMode) localStorage.setItem(LS_ROLE, roleMode);
    } catch {
      // ignore
    }
  }, [roleMode]);

  async function refreshCitizens() {
    setLoading(true);
    try {
      const data = await getCitizens();
      const list = Array.isArray(data) ? data : data?.items ?? [];
      setCitizens(list);
      return list;
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    refreshCitizens();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const selectedCitizen = useMemo(() => {
    return citizens.find((c) => c.id === selectedCitizenId) || null;
  }, [citizens, selectedCitizenId]);

  const value = useMemo(
    () => ({
      citizens,
      citizensLoading: loading,
      refreshCitizens,
      selectedCitizenId,
      setSelectedCitizenId,
      selectedCitizen,
      roleMode,
      setRoleMode,
    }),
    [citizens, loading, selectedCitizenId, selectedCitizen, roleMode]
  );

  return (
    <CitizenContext.Provider value={value}>{children}</CitizenContext.Provider>
  );
}

export function useCitizenContext() {
  const ctx = useContext(CitizenContext);
  if (!ctx) throw new Error("useCitizenContext must be used within CitizenProvider");
  return ctx;
}
