import { useMemo, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { toast } from "sonner";
import { Mail, Lock } from "lucide-react";
import { FcGoogle } from "react-icons/fc";

import { http, ApiError } from "@/lib/http";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

type LoginReq = { email: string; password: string };

export function LoginPage() {
  const nav = useNavigate();
  const [sp] = useSearchParams();

  const next = useMemo(() => sp.get("next") || "/", [sp]);

  const [email, setEmail] = useState("demo@example.com");
  const [password, setPassword] = useState("password");
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    try {
      await http<void>("/api/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password } satisfies LoginReq),
      });
      toast.success("Login success");
      nav(next, { replace: true });
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Login failed";
      toast.error(msg);
    } finally {
      setLoading(false);
    }
  }

  function onGoogle() {
    window.location.assign("/api/auth/google/start");
  }

  return (
    <div className="min-h-screen w-full bg-primary/20">
      <div className="min-h-screen w-full grid place-items-center px-6 py-10">
        <div className="w-full max-w-md">
          <Card className="rounded-2xl shadow-sm">
            <CardContent className="p-6 space-y-5">
              <div className="space-y-1 text-center">
                <div className="text-xl font-semibold">Login</div>
              </div>

              <Button type="button" variant="outline" className="w-full" onClick={onGoogle}>
                <FcGoogle className="mr-2 h-4 w-4" />
                Continue with Google
              </Button>

              <div className="relative">
                <Separator />
                <div className="absolute left-1/2 -translate-x-1/2 -top-3 bg-card px-2 text-xs text-muted-foreground">
                  or
                </div>
              </div>

              <form className="space-y-3" onSubmit={onSubmit}>
                <div className="space-y-2">
                  <label className="text-xs font-medium text-muted-foreground">Email</label>
                  <div className="relative">
                    <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                      className="pl-9 h-11"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      placeholder="email@domain.com"
                      autoComplete="email"
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="text-xs font-medium text-muted-foreground">Password</label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                      className="pl-9 h-11"
                      type="password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      placeholder="••••••••"
                      autoComplete="current-password"
                    />
                  </div>
                </div>

                <Button className="w-full h-11" type="submit" disabled={loading}>
                  {loading ? "Logging in…" : "Login"}
                </Button>

                <div className="text-xs text-muted-foreground text-center">
                  Demo inquiry:{" "}
                  <span className="font-medium text-foreground">demo@example.com</span> /{" "}
                  <span className="font-medium text-foreground">Password123!</span>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}