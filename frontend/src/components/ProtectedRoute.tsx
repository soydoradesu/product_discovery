import { 
    Navigate,
    useLocation 
} from "react-router-dom";
import { 
    useMe 
} from "@/features/auth/hooks";

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const me = useMe();
  const loc = useLocation();

  if (me.isLoading) {
    return <div className="p-6">Loading…</div>;
  }

  if (me.isError) {
    const next = encodeURIComponent(loc.pathname + loc.search);
    return <Navigate to={`/login?next=${next}`} replace />;
  }

  return <>{children}</>;
}
