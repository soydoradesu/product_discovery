import { useMemo, useState } from "react";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import { toast } from "sonner";
import { Star, PackageCheck, PackageX, ChevronLeft } from "lucide-react";

import { ApiError } from "@/lib/http";
import { useProductDetail } from "@/features/catalog/hooks";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";

function safeDecodeFrom(from: string | null): string {
  if (!from) return "/";
  try {
    const decoded = decodeURIComponent(from);
    if (!decoded.startsWith("/")) return "/";
    return decoded;
  } catch {
    return "/";
  }
}

function Skeleton() {
  return (
    <div className="grid gap-4 lg:grid-cols-[1.2fr_0.8fr]">
      <div className="rounded-xl border bg-card overflow-hidden">
        <div className="aspect-[4/3] bg-muted animate-pulse" />
        <div className="p-4 space-y-2">
          <div className="h-5 w-2/3 bg-muted animate-pulse rounded" />
          <div className="h-4 w-1/3 bg-muted animate-pulse rounded" />
          <div className="h-20 w-full bg-muted animate-pulse rounded mt-4" />
        </div>
      </div>
      <div className="rounded-xl border bg-card p-4 space-y-3">
        <div className="h-6 w-1/2 bg-muted animate-pulse rounded" />
        <div className="h-10 w-full bg-muted animate-pulse rounded" />
        <div className="h-10 w-full bg-muted animate-pulse rounded" />
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

  const [activeIdx, setActiveIdx] = useState(0);

  function goBack() {
    nav(from);
  }

  if (q.isLoading) {
    return (
      <div className="min-h-screen bg-muted/20">
        <div className="w-full px-6 py-4">
          <Button variant="outline" onClick={goBack}>
            <ChevronLeft className="mr-1 h-4 w-4" />
            Back
          </Button>
        </div>
        <div className="w-full px-6 pb-8">
          <Skeleton />
        </div>
      </div>
    );
  }

  if (q.isError) {
    const err = q.error;
    const msg = err instanceof ApiError ? err.message : "Failed to load product";

    if (err instanceof ApiError && err.status === 401) {
      const next = encodeURIComponent(`/products/${id}?from=${encodeURIComponent(from)}`);
      toast.error("Session expired. Please login again.");
      nav(`/login?next=${next}`, { replace: true });
      return null;
    }

    return (
      <div className="min-h-screen bg-muted/20">
        <div className="w-full px-6 py-4">
          <Button variant="outline" onClick={goBack}>
            <ChevronLeft className="mr-1 h-4 w-4" />
            Back
          </Button>
        </div>

        <div className="w-full px-6 pb-8">
          <Card className="rounded-xl">
            <CardContent className="p-6 space-y-3">
              <div className="text-base font-semibold">Error</div>
              <div className="text-sm text-muted-foreground">{msg}</div>
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
      </div>
    );
  }

  const p = q.data!;
  const images = (p.images ?? []).slice().sort((a, b) => a.position - b.position);
  const hasImages = images.length > 0;

  const active = hasImages ? images[Math.min(activeIdx, images.length - 1)] : null;

  return (
    <div className="min-h-screen bg-muted/20">
      {/* Top bar */}
      <div className="sticky top-0 z-20 border-b bg-background/95 backdrop-blur">
        <div className="w-full px-6 py-3 flex items-center justify-between gap-3">
          <Button variant="outline" onClick={goBack}>
            <ChevronLeft className="mr-1 h-4 w-4" />
            Back
          </Button>
          <div className="text-sm text-muted-foreground truncate">{p.name}</div>
        </div>
      </div>

      <div className="w-full px-6 py-6">
        <div className="grid gap-4 lg:grid-cols-[1.2fr_0.8fr]">
          {/* Left: gallery + description */}
          <div className="space-y-4">
            <Card className="rounded-xl overflow-hidden">
              <CardContent className="p-0">
                <div className="bg-muted">
                  {active ? (
                    <img
                      src={active.url}
                      alt={p.name}
                      className="w-full aspect-[4/3] object-cover"
                      loading="lazy"
                    />
                  ) : (
                    <div className="w-full aspect-[4/3] bg-gradient-to-br from-muted to-background" />
                  )}
                </div>

                {hasImages ? (
                  <div className="p-3 border-t bg-background">
                    <div className="flex gap-2 overflow-auto">
                      {images.slice(0, 10).map((img, idx) => (
                        <button
                          key={img.position}
                          className={[
                            "h-14 w-20 rounded-lg overflow-hidden border shrink-0",
                            idx === activeIdx ? "ring-2 ring-emerald-600 border-emerald-600" : "",
                          ].join(" ")}
                          onClick={() => setActiveIdx(idx)}
                          aria-label={`Image ${idx + 1}`}
                        >
                          <img
                            src={img.url}
                            alt={`${p.name} ${idx + 1}`}
                            className="h-full w-full object-cover"
                            loading="lazy"
                          />
                        </button>
                      ))}
                    </div>
                  </div>
                ) : null}
              </CardContent>
            </Card>

            <Card className="rounded-xl">
              <CardContent className="p-5 space-y-3">
                <div className="text-base font-semibold">Description</div>
                <div className="text-sm text-foreground/90 whitespace-pre-wrap">{p.description}</div>
              </CardContent>
            </Card>
          </div>

          {/* Right: buy/info card */}
          <div className="lg:sticky lg:top-20 h-fit">
            <Card className="rounded-xl">
              <CardContent className="p-5 space-y-4">
                <div className="space-y-1">
                  <div className="text-lg font-semibold leading-tight">{p.name}</div>

                  <div className="flex flex-wrap items-center gap-2">
                    <div className="inline-flex items-center gap-1 text-sm">
                      <Star className="h-4 w-4 fill-yellow-400 stroke-yellow-400" />
                      <span className="font-semibold">{p.rating.toFixed(1)}</span>
                      <span className="text-muted-foreground">/ 5</span>
                    </div>

                    {p.inStock ? (
                      <Badge className="bg-emerald-600 hover:bg-emerald-600 text-white">
                        <PackageCheck className="mr-1 h-3.5 w-3.5" />
                        Stock available
                      </Badge>
                    ) : (
                      <Badge variant="secondary">
                        <PackageX className="mr-1 h-3.5 w-3.5" />
                        Out of stock
                      </Badge>
                    )}
                  </div>

                  {p.categories?.length ? (
                    <div className="flex flex-wrap gap-1 pt-1">
                      {p.categories.slice(0, 6).map((c) => (
                        <span
                          key={c.id}
                          className="rounded-full border bg-background px-2 py-0.5 text-[11px] text-muted-foreground"
                        >
                          {c.name}
                        </span>
                      ))}
                    </div>
                  ) : null}
                </div>

                <div className="rounded-xl border bg-muted/30 p-4">
                  <div className="text-xs text-muted-foreground">Price</div>
                  <div className="text-2xl font-bold">Rp{(p.price)}</div>
                </div>

                <div className="text-xs text-muted-foreground">
                  ID: <span className="font-medium text-foreground">{p.id}</span> - Date Posted:{" "}
                  <span className="font-medium text-foreground">
                    {new Date(p.createdAt).toLocaleDateString("id-ID")}
                  </span>
                </div>

                <Button variant="outline" className="w-full" onClick={goBack}>
                  Back to catalog
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}