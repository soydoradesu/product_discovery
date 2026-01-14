import { toast } from "sonner";
import { 
    useLogout, 
    useMe } from "@/features/auth/hooks";
import { Button } from "@/components/ui/button";
import { 
    Card, 
    CardContent, 
    CardHeader, 
    CardTitle } from "@/components/ui/card";

export function HomePage() {
  const me = useMe();
  const logout = useLogout();

  async function onLogout() {
    try {
      await logout.mutateAsync();
      toast.success("Logged out");
      window.location.href = "/login";
    } catch {
      toast.error("Logout failed");
    }
  }

  return (
    <div className="max-w-5xl mx-auto p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold">Home</h1>
        <Button variant="outline" onClick={onLogout} disabled={logout.isPending}>
          {logout.isPending ? "Logging out…" : "Logout"}
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Session</CardTitle>
        </CardHeader>
        <CardContent className="text-sm">
          <div className="text-muted-foreground">Logged in as</div>
          <div className="font-medium mt-1">{me.data?.email}</div>
        </CardContent>
      </Card>

      <div className="text-sm text-muted-foreground">
        TODO: Search page (debounce + URL state + filters + pagination) + Products.
      </div>
    </div>
  );
}