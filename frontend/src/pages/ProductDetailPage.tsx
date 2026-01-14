import { useMemo } from "react";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import { toast } from "sonner";

import { ApiError } from "@/lib/http";
import { useProductDetail } from "@/features/catalog/hooks";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

function safeDecodeFrom(from: string | null): string {
  if (!from) return "/";
  try {
    const decoded = decodeURIComponent(from);
    // prevent open redirect
    if (!decoded.startsWith("/")) return "/";
    return decoded;
  } catch {
    return "/";
  }
}

function DetailSkeleton() {
  return (
    <div className="space-y-4">
      <div className="rounded-lg border p-6">
        <div className="h-6 w-2/3 bg-muted animate-pulse rounded" />
        <div className="mt-3 h-4 w-1/3 bg-muted animate-pulse rounded" />
        <div className="mt-6 h-24 w-full bg-muted animate-pulse rounded" />
      </div>
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="h-40 rounded-lg border bg-muted animate-pulse" />
        ))}
      </div>
    </div>
  );
}

export function ProductDetailPage() {
  const nav = useNavigate();
  const { id: idParam } = useParams();
  const [sp] = useSearchParams();

  const id = useMemo(() => Number(idParam), [idParam]);
  const from = useMemo(() => safeDecodeFrom(sp.get("from")), [sp]);

  const q = useProductDetail(id);

  function goBack() {
    nav(from);
  }

  if (q.isLoading) {
    return (
      <div className="max-w-5xl mx-auto p-6 space-y-6">
        <div className="flex items-center justify-between">
          <Button variant="outline" onClick={goBack}>
            ← Back to results
          </Button>
        </div>
        <DetailSkeleton />
      </div>
    );
  }

  if (q.isError) {
    const err = q.error;
    const msg = err instanceof ApiError ? err.message : "Failed to load product";

    // If session expired / not authorized, redirect to login preserving "next"
    if (err instanceof ApiError && err.status === 401) {
      const next = encodeURIComponent(`/products/${id}?from=${encodeURIComponent(from)}`);
      toast.error("Session expired. Please login again.");
      nav(`/login?next=${next}`, { replace: true });
      return null;
    }

    return (
      <div className="max-w-5xl mx-auto p-6 space-y-6">
        <div className="flex items-center justify-between">
          <Button variant="outline" onClick={goBack}>
            ← Back to results
          </Button>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Error</CardTitle>
          </CardHeader>
          <CardContent className="text-sm space-y-3">
            <div className="text-muted-foreground">{msg}</div>
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => q.refetch()}>
                Retry
              </Button>
              <Button variant="outline" onClick={goBack}>
                Back
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  const p = q.data;

  if (!p) {
    return null;
  }

  return (
    <div className="max-w-5xl mx-auto p-6 space-y-6">
      <div className="flex items-center justify-between gap-3">
        <Button variant="outline" onClick={goBack}>
          ← Back to results
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-xl">{p.name}</CardTitle>
          <div className="text-sm text-muted-foreground">
            ${p.price.toFixed(2)} • ⭐ {p.rating.toFixed(1)} • {p.inStock ? "In stock" : "Out of stock"}
          </div>
          {p.categories?.length ? (
            <div className="text-xs text-muted-foreground mt-2">
              {p.categories.map((c) => c.name).join(", ")}
            </div>
          ) : null}
        </CardHeader>

        <CardContent className="space-y-4">
          <div className="text-sm whitespace-pre-wrap">{p.description}</div>

          {p.images?.length ? (
            <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
              {p.images
                .slice()
                .sort((a, b) => a.position - b.position)
                .map((img) => (
                  <div key={img.position} className="rounded-lg border overflow-hidden">
                    <img
                      src={img.url}
                      alt={`${p.name} ${img.position}`}
                      className="w-full h-48 object-cover"
                      loading="lazy"
                    />
                  </div>
                ))}
            </div>
          ) : (
            <div className="text-sm text-muted-foreground">No images</div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
