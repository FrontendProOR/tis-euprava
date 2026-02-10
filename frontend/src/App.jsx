import { Routes, Route } from "react-router-dom";
import Layout from "./components/Layout";

import Dashboard from "./pages/Dashboard";
import Citizens from "./pages/Citizens";
import Requests from "./pages/Requests";
import RequestDetails from "./pages/RequestDetails";
import Payments from "./pages/Payments";
import Appointments from "./pages/Appointments";
import Certificates from "./pages/Certificates";

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<Dashboard />} />
        <Route path="/citizens" element={<Citizens />} />
        <Route path="/requests" element={<Requests />} />
        <Route path="/requests/:id" element={<RequestDetails />} />
        <Route path="/payments" element={<Payments />} />
        <Route path="/appointments" element={<Appointments />} />
        <Route path="/certificates" element={<Certificates />} />
      </Route>
    </Routes>
  );
}
