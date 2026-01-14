import { 
    type FormEvent, 
    useMemo, 
    useState 
} from "react";
import { 
    useNavigate, 
    useSearchParams 
} from "react-router-dom";
import { toast } from "sonner";
import { ApiError } from "@/lib/http";
import { 
    useLogin, 
    useMe 
} from "@/features/auth/hooks";
import { Button } from "@/components/ui/button";
import { 
    Card, 
    CardContent, 
    CardHeader, 
    CardTitle, 
    CardDescription 
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";

export function LoginPage() {
  const me = useMe();
  const login = useLogin();
  const nav = useNavigate();
  const [sp] = useSearchParams();

  const [email, setEmail] = useState("demo@example.com");
  const [password, setPassword] = useState("Password123!");

  const next = useMemo(() => sp.get("next") || "/", [sp]);
  const oauth = sp.get("oauth");

  if (me.isSuccess) {
    nav(next, { replace: true });
  }

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    try {
      await login.mutateAsync({ email, password });
      nav(next, { replace: true });
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Login failed";
      toast.error(msg);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-6">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Login</CardTitle>
          <CardDescription>
            Sign in with Google or email/password. Cookies are HttpOnly (no token in localStorage).
          </CardDescription>
        </CardHeader>

        <CardContent className="space-y-4">
          {oauth === "success" && (
            <div className="rounded-md border p-3 text-sm">
              Google login success. If you’re not redirected, go to Home.
            </div>
          )}

          <Button asChild variant="outline" className="w-full">
            <a href="/api/auth/google/start">Continue with Google</a>
          </Button>

          <div className="text-center text-sm text-muted-foreground">or</div>

          <form onSubmit={onSubmit} className="space-y-3">
            <div className="space-y-1">
              <div className="text-sm font-medium">Email</div>
              <Input
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                autoComplete="email"
                placeholder="you@example.com"
              />
            </div>

            <div className="space-y-1">
              <div className="text-sm font-medium">Password</div>
              <Input
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                type="password"
                autoComplete="current-password"
              />
            </div>

            <Button type="submit" className="w-full" disabled={login.isPending}>
              {login.isPending ? "Signing in…" : "Sign in"}
            </Button>

            <div className="text-xs text-muted-foreground">
              Demo user: demo@example.com / Password123!
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}