import React from "react";
import ReactDOM from "react-dom/client";
import { 
  BrowserRouter 
} from "react-router-dom";
import { 
  QueryClient, 
  QueryClientProvider 
} from "@tanstack/react-query";
import { 
  Toaster 
} from "@/components/ui/sonner";
import App from "./App";
import "./index.css";

document.documentElement.classList.add("light"); 
document.documentElement.style.colorScheme = "light"; 

const qc = new QueryClient();

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={qc}>
      <BrowserRouter>
        <App />
        <Toaster richColors />
      </BrowserRouter>
    </QueryClientProvider>
  </React.StrictMode>
);
