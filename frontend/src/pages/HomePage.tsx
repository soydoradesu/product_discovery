import { 
    useCallback, 
    useEffect, 
    useMemo 
} from "react";
import { 
    useNavigate, 
    useSearchParams 
} from "react-router-dom";
import { toast } from "sonner";
import { Search } from "lucide-react";

import { ApiError } from "@/lib/http";
import { 
    useLogout, 
    useMe 
} from "@/features/auth/hooks";
import { 
    useCategories, 
    useProductSearch 
} from "@/features/catalog/hooks";
import type { SearchParams } from "@/features/catalog/api";
import { useDebouncedValue } from "@/hooks/useDebouncedValue";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ProductCard } from "@/components/catalog/ProductCard";
import { ProductCardSkeleton } from "@/components/catalog/ProductCardSkeleton";
import { EmptyState } from "@/components/catalog/EmptyState";
import { FiltersPanel } from "@/components/catalog/FiltersPanel";
import { PaginationBar } from "@/components/catalog/PaginationBar";
import { SortBar } from "@/components/catalog/SortBar";

function clampInt(n: number, min: number, max: number) {
  return Math.max(min, Math.min(max, n));
}

function parseNumberList(values: string[]): number[] {
  return values.map((v) => Number(v)).filter((n) => Number.isFinite(n) && n > 0);
}

function getInt(sp: URLSearchParams, key: string, def: number) {
  const v = sp.get(key);
  if (!v) return def;
  const n = Number(v);
  return Number.isFinite(n) ? Math.trunc(n) : def;
}

function getEnum<T extends string>(
  sp: URLSearchParams,
  key: string,
  allowed: readonly T[],
  def: T
): T {
  const v = (sp.get(key) || "").toLowerCase() as T;
  return (allowed as readonly string[]).includes(v) ? v : def;
}

function setParam(sp: URLSearchParams, key: string, value: string) {
  if (value.trim() === "") sp.delete(key);
  else sp.set(key, value);
}

function setMulti(sp: URLSearchParams, key: string, values: number[]) {
  sp.delete(key);
  for (const v of values) sp.append(key, String(v));
}

export function HomePage() {
  const me = useMe();
  const logout = useLogout();
  const nav = useNavigate();
  const [sp, setSp] = useSearchParams();

  const qRaw = sp.get("q") || "";
  const qTrim = qRaw.trim();
  const debouncedQ = useDebouncedValue(qTrim, 350);

  const categories = parseNumberList(sp.getAll("category"));
  const minPrice = sp.get("minPrice") || "";
  const maxPrice = sp.get("maxPrice") || "";
  const inStock = getEnum(sp, "inStock", ["any", "true", "false"] as const, "any");

  const sort = getEnum(
    sp,
    "sort",
    ["relevance", "price", "created_at", "rating"] as const,
    qTrim.length > 0 ? "relevance" : "created_at"
  );

  const method = getEnum(sp, "method", ["asc", "desc"] as const, sort === "price" ? "asc" : "desc");
  const page = clampInt(getInt(sp, "page", 1), 1, 1_000_000);
  const pageSize = clampInt(getInt(sp, "pageSize", 20), 1, 100);

  const params: SearchParams = useMemo(
    () => ({
      q: debouncedQ,
      categories,
      minPrice,
      maxPrice,
      inStock,
      sort,
      method,
      page,
      pageSize,
    }),
    [debouncedQ, categories, minPrice, maxPrice, inStock, sort, method, page, pageSize]
  );

  const catsQ = useCategories();
  const searchQ = useProductSearch(params);

  useEffect(() => {
    if (!searchQ.isError) return;
    const err = searchQ.error;
    const msg = err instanceof ApiError ? err.message : "Search failed";
    toast.error(msg);
  }, [searchQ.isError]); // eslint-disable-line react-hooks/exhaustive-deps

  const updateSearchParams = useCallback(
    (mut: (p: URLSearchParams) => void) => {
      setSp(
        (prev) => {
          const next = new URLSearchParams(prev);
          mut(next);
          return next;
        },
        { replace: true }
      );
    },
    [setSp]
  );

  const resetToFirstPage = useCallback((p: URLSearchParams) => p.set("page", "1"), []);

  const onClearAll = useCallback(() => {
    setSp(new URLSearchParams(), { replace: true });
  }, [setSp]);

  const onToggleCategory = useCallback(
    (id: number) => {
      updateSearchParams((p) => {
        const next = new Set(parseNumberList(p.getAll("category")));
        if (next.has(id)) next.delete(id);
        else next.add(id);

        setMulti(p, "category", Array.from(next));
        resetToFirstPage(p);
      });
    },
    [updateSearchParams, resetToFirstPage]
  );

  const onMinPrice = useCallback(
    (v: string) => {
      updateSearchParams((p) => {
        setParam(p, "minPrice", v);
        resetToFirstPage(p);
      });
    },
    [updateSearchParams, resetToFirstPage]
  );

  const onMaxPrice = useCallback(
    (v: string) => {
      updateSearchParams((p) => {
        setParam(p, "maxPrice", v);
        resetToFirstPage(p);
      });
    },
    [updateSearchParams, resetToFirstPage]
  );

  const onInStock = useCallback(
    (v: "any" | "true" | "false") => {
      updateSearchParams((p) => {
        if (v === "any") p.delete("inStock");
        else p.set("inStock", v);
        resetToFirstPage(p);
      });
    },
    [updateSearchParams, resetToFirstPage]
  );

  const onSort = useCallback(
    (v: "relevance" | "price" | "created_at" | "rating") => {
      updateSearchParams((p) => {
        p.set("sort", v);
        resetToFirstPage(p);
      });
    },
    [updateSearchParams, resetToFirstPage]
  );

  const onMethod = useCallback(
    (v: "asc" | "desc") => {
      updateSearchParams((p) => {
        p.set("method", v);
        resetToFirstPage(p);
      });
    },
    [updateSearchParams, resetToFirstPage]
  );

  const onPageSize = useCallback(
    (n: number) => {
      updateSearchParams((p) => {
        p.set("pageSize", String(n));
        resetToFirstPage(p);
      });
    },
    [updateSearchParams, resetToFirstPage]
  );

  const onPrev = useCallback(() => {
    updateSearchParams((p) => {
      p.set("page", String(Math.max(1, page - 1)));
    });
  }, [updateSearchParams, page]);

  const onNext = useCallback(() => {
    updateSearchParams((p) => {
      p.set("page", String(page + 1));
    });
  }, [updateSearchParams, page]);

  async function onLogout() {
    try {
      await logout.mutateAsync();
      toast.success("Logged out");
      nav("/login", { replace: true });
    } catch {
      toast.error("Logout failed");
    }
  }

  const totalPages = searchQ.data?.totalPages ?? 0;
  const totalLabel = searchQ.data ? `${searchQ.data.total.toLocaleString("id-ID")} product` : undefined;

  return (
    <div className="min-h-screen bg-muted/20">
      <div className="sticky top-0 z-20 border-b bg-background/95 backdrop-blur">
        <div className="w-full px-6 py-3">
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2">
              <div className="h-8 w-8 rounded-lg bg-emerald-600" />
              <div className="leading-tight">
                <div className="text-sm font-semibold">Product Discovery</div>
              </div>
            </div>

            <div className="ml-auto flex items-center gap-3">
              <div className="hidden sm:block text-xs text-muted-foreground">
                {me.data?.email}
              </div>
              <Button variant="destructive" size="sm" onClick={onLogout} disabled={logout.isPending}>
                {logout.isPending ? "Logging out…" : "Logout"}
              </Button>
            </div>
          </div>

          {/* Big search */}
          <div className="mt-3 relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              value={qRaw}
              onChange={(e) => {
                const v = e.target.value;
                updateSearchParams((p) => {
                  setParam(p, "q", v);
                  resetToFirstPage(p);
                });
              }}
              placeholder="Find product… (name / description)"
              className="pl-9 h-11"
              data-testid="search-input"
            />
          </div>
        </div>
      </div>

      <div className="w-full px-6 py-6">
        <div className="grid gap-4 lg:grid-cols-[280px_1fr]">
          {/* Sidebar Filters */}
          <div className="lg:sticky lg:top-24 h-fit">
            <FiltersPanel
              categories={catsQ.data?.items ?? []}
              selectedCategoryIds={categories}
              minPrice={minPrice}
              maxPrice={maxPrice}
              inStock={inStock}
              onToggleCategory={onToggleCategory}
              onMinPrice={onMinPrice}
              onMaxPrice={onMaxPrice}
              onInStock={onInStock}
              onClear={onClearAll}
              loadingCategories={catsQ.isLoading}
              categoriesError={catsQ.isError}
            />
          </div>

          {/* Main */}
          <div className="space-y-3">
            <div className="rounded-xl border bg-card p-4">
              <SortBar
                sort={sort}
                method={method}
                pageSize={pageSize}
                totalLabel={totalLabel}
                onSort={onSort}
                onMethod={onMethod}
                onPageSize={onPageSize}
              />

              {searchQ.isFetching && !searchQ.isLoading ? (
                <div className="mt-2 text-xs text-muted-foreground">Updating results…</div>
              ) : null}
            </div>

            {searchQ.isLoading ? (
              <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
                {Array.from({ length: 9 }).map((_, i) => (
                  <ProductCardSkeleton key={i} />
                ))}
              </div>
            ) : searchQ.isError ? (
              <EmptyState
                title="Terjadi kesalahan"
                description={
                  searchQ.error instanceof ApiError ? searchQ.error.message : "Search failed"
                }
                onClear={onClearAll}
              />
            ) : (searchQ.data?.items?.length ?? 0) === 0 ? (
              <EmptyState
                title="Tidak ada hasil"
                description="Coba ganti keyword atau reset filter."
                onClear={onClearAll}
              />
            ) : (
              <>
                <div className="grid gap-3 grid-cols-2 sm:grid-cols-3 lg:grid-cols-5">
                  {(searchQ.data?.items ?? []).map((p) => (
                    <ProductCard key={p.id} p={p} />
                  ))}
                </div>

                <div className="rounded-xl border bg-card p-4">
                  <PaginationBar page={page} totalPages={totalPages} onPrev={onPrev} onNext={onNext} />
                </div>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
