import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "@/index.css";
import App from "@/App.tsx";
import { AuthProvider } from "@/context/AuthContext.tsx";
import { ServerProvider } from "@/context/ServerContext";

import { Toaster } from "@/components/ui/sonner";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AuthProvider>
      <ServerProvider>
        <App />
        <Toaster />
      </ServerProvider>
    </AuthProvider>
  </StrictMode>,
);
