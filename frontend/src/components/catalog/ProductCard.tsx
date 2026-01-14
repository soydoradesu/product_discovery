import { useLocation, useNavigate } from "react-router-dom";
import { Star, PackageCheck, PackageX } from "lucide-react";

import type { ProductSummary } from "@/features/catalog/api";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

export function ProductCard({ p }: { p: ProductSummary }) {
  const nav = useNavigate();
  const loc = useLocation();

  const from = encodeURIComponent(loc.pathname + loc.search);
  const goDetail = () => nav(`/products/${p.id}?from=${from}`);

  return (
    <div
      className={cn(
        "group rounded-xl border bg-card shadow-sm hover:shadow-md transition-shadow overflow-hidden",
        "focus-within:ring-2 focus-within:ring-ring"
      )}
    >
      <div className="relative aspect-[4/3] bg-muted">
        {p.thumbnail ? (
          <img
            src={p.thumbnail}
            alt={p.name}
            className="h-full w-full object-cover"
            loading="lazy"
          />
        ) : (
          <div className="h-full w-full bg-gradient-to-br from-muted to-background" />
        )}

        <div className="absolute left-2 top-2 flex gap-2">
          {p.inStock ? (
            <Badge className="bg-emerald-600 hover:bg-emerald-600 text-white">
              <PackageCheck className="mr-1 h-3.5 w-3.5" />
              Ready
            </Badge>
          ) : (
            <Badge variant="secondary">
              <PackageX className="mr-1 h-3.5 w-3.5" />
              Out
            </Badge>
          )}
        </div>
      </div>

      <div className="p-3 space-y-2">
        <div className="min-h-[36px]">
          <div className="line-clamp-2 text-sm font-medium leading-5">{p.name}</div>
        </div>

        <div className="text-base font-semibold">Rp{(p.price)}</div>

        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <div className="inline-flex items-center gap-1">
            <Star className="h-3.5 w-3.5 fill-yellow-400 stroke-yellow-400" />
            <span className="font-medium text-foreground">{p.rating.toFixed(1)}</span>
          </div>
          <div className="truncate max-w-[120px] text-right">
            {p.categories?.length ? p.categories[0].name : ""}
          </div>
        </div>

        {p.categories?.length ? (
          <div className="flex flex-wrap gap-1 pt-1">
            {p.categories.slice(0, 2).map((c) => (
              <span
                key={c.id}
                className="rounded-full border bg-background px-2 py-0.5 text-[11px] text-muted-foreground"
              >
                {c.name}
              </span>
            ))}
          </div>
        ) : null}

        <div className="pt-2">
          <Button
            size="sm"
            className="w-full"
            onClick={goDetail}
            aria-label={`View details ${p.name}`}
          >
            View Detail
          </Button>
        </div>
      </div>
    </div>
  );
}