import { 
    useEffect, 
    useMemo 
} from "react";
import { 
    useNavigate, 
    useSearchParams,
    useLocation,
} from "react-router-dom";
import { toast } from "sonner";

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
import { 
    Card,
    CardContent, 
    CardHeader, 
    CardTitle 
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";

function clampInt(n: number, min: number, max: number) {
    return Math.max(min, Math.min(max, n));
}

function parseNumberList(values: string[]): number[] {
    return values
        .map((v) => Number(v))
        .filter((n) => Number.isFinite(n) && n > 0);
}

function getInt(sp: URLSearchParams, key: string, def: number) {
    const v = sp.get(key);
    if (!v) return def;
    const n = Number(v);
    return Number.isFinite(n) ? Math.trunc(n) : def;
}

function getEnum<T extends string>(sp: URLSearchParams, key: string, allowed: readonly T[], def: T): T {
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

function SkeletonGrid() {
    return (
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="rounded-lg border p-4">
            <div className="h-4 w-2/3 bg-muted animate-pulse rounded" />
            <div className="mt-3 h-3 w-1/3 bg-muted animate-pulse rounded" />
            <div className="mt-6 h-24 w-full bg-muted animate-pulse rounded" />
            </div>
        ))}
        </div>
    );
}

export function HomePage() {
    const me = useMe();
    const logout = useLogout();
    const loc = useLocation();
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

    const method = getEnum(
        sp,
        "method",
        ["asc", "desc"] as const,
        sort === "price" ? "asc" : "desc"
    );

    const page = clampInt(getInt(sp, "page", 1), 1, 1_000_000);
    const pageSize = clampInt(getInt(sp, "pageSize", 20), 1, 100);

    const normalized: SearchParams = useMemo(
        () => ({
        q: debouncedQ,
        categories,
        minPrice,
        maxPrice,
        inStock,
        sort,
        method,
        page,
        pageSize
        }),
        [debouncedQ, categories, minPrice, maxPrice, inStock, sort, method, page, pageSize]
    );

    const catsQ = useCategories();

    const searchQ = useProductSearch(normalized);

    useEffect(() => {
        if (!searchQ.isError) return;
        const err = searchQ.error;
        const msg = err instanceof ApiError ? err.message : "Search failed";
        toast.error(msg);
    }, [searchQ.isError]);

    function updateSearchParams(mut: (p: URLSearchParams) => void) {
        setSp(
        (prev) => {
            const next = new URLSearchParams(prev);
            mut(next);
            return next;
        },
        { replace: true }
        );
    }

    function resetToFirstPage(p: URLSearchParams) {
        p.set("page", "1");
    }

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

    return (
        <div className="max-w-6xl mx-auto p-6 space-y-6">
        <div className="flex items-center justify-between gap-3">
            <div>
            <div className="text-xl font-semibold">Product Search</div>
            <div className="text-sm text-muted-foreground">
                Logged in as <span className="font-medium">{me.data?.email}</span>
            </div>
            </div>

            <Button variant="outline" onClick={onLogout} disabled={logout.isPending}>
            {logout.isPending ? "Logging out…" : "Logout"}
            </Button>
        </div>

        {/* Controls */}
        <Card>
            <CardHeader>
            <CardTitle className="text-base">Search</CardTitle>
            </CardHeader>

            <CardContent className="space-y-4">
            <div className="flex flex-col gap-3 md:flex-row md:items-center">
                <div className="flex-1">
                <Input
                    value={qRaw}
                    onChange={(e) => {
                    const v = e.target.value;
                    updateSearchParams((p) => {
                        setParam(p, "q", v);
                        resetToFirstPage(p);
                    });
                    }}
                    placeholder="Search name or description…"
                />
                <div className="mt-1 text-xs text-muted-foreground">
                    Debounced typing (350ms). URL updates immediately; query uses debounced value.
                </div>
                </div>

                <div className="flex gap-2 items-center">
                <select
                    className="h-10 rounded-md border bg-background px-3 text-sm"
                    value={sort}
                    onChange={(e) => {
                    const v = e.target.value;
                    updateSearchParams((p) => {
                        p.set("sort", v);
                        resetToFirstPage(p);
                    });
                    }}
                >
                    <option value="relevance">relevance</option>
                    <option value="created_at">created_at</option>
                    <option value="price">price</option>
                    <option value="rating">rating</option>
                </select>

                <select
                    className="h-10 rounded-md border bg-background px-3 text-sm"
                    value={method}
                    onChange={(e) => {
                    const v = e.target.value;
                    updateSearchParams((p) => {
                        p.set("method", v);
                        resetToFirstPage(p);
                    });
                    }}
                >
                    <option value="asc">asc</option>
                    <option value="desc">desc</option>
                </select>

                <select
                    className="h-10 rounded-md border bg-background px-3 text-sm"
                    value={String(pageSize)}
                    onChange={(e) => {
                    const v = e.target.value;
                    updateSearchParams((p) => {
                        p.set("pageSize", v);
                        resetToFirstPage(p);
                    });
                    }}
                >
                    <option value="20">20</option>
                    <option value="50">50</option>
                    <option value="100">100</option>
                </select>
                </div>
            </div>

            {/* Filters row */}
            <div className="grid gap-3 md:grid-cols-3">
                <div className="space-y-2">
                <div className="text-sm font-medium">Price</div>
                <div className="flex gap-2">
                    <Input
                    inputMode="decimal"
                    placeholder="min"
                    value={minPrice}
                    onChange={(e) => {
                        const v = e.target.value;
                        updateSearchParams((p) => {
                        setParam(p, "minPrice", v);
                        resetToFirstPage(p);
                        });
                    }}
                    />
                    <Input
                    inputMode="decimal"
                    placeholder="max"
                    value={maxPrice}
                    onChange={(e) => {
                        const v = e.target.value;
                        updateSearchParams((p) => {
                        setParam(p, "maxPrice", v);
                        resetToFirstPage(p);
                        });
                    }}
                    />
                </div>
                </div>

                <div className="space-y-2">
                <div className="text-sm font-medium">In stock</div>
                <select
                    className="h-10 w-full rounded-md border bg-background px-3 text-sm"
                    value={inStock}
                    onChange={(e) => {
                    const v = e.target.value as "any" | "true" | "false";
                    updateSearchParams((p) => {
                        if (v === "any") p.delete("inStock");
                        else p.set("inStock", v);
                        resetToFirstPage(p);
                    });
                    }}
                >
                    <option value="any">Any</option>
                    <option value="true">In stock only</option>
                    <option value="false">Out of stock only</option>
                </select>
                </div>

                <div className="flex items-end">
                <Button
                    variant="outline"
                    onClick={() => {
                    setSp(new URLSearchParams(), { replace: true });
                    }}
                >
                    Clear all filters
                </Button>
                </div>
            </div>

            {/* Categories */}
            <div className="space-y-2">
                <div className="text-sm font-medium">Categories</div>

                {catsQ.isLoading ? (
                <div className="text-sm text-muted-foreground">Loading categories…</div>
                ) : catsQ.isError ? (
                <div className="text-sm text-destructive">Failed to load categories</div>
                ) : (
                <div className="flex flex-wrap gap-2">
                    {catsQ.data?.items.map((c) => {
                    const active = categories.includes(c.id);
                    return (
                        <Button
                        key={c.id}
                        variant={active ? "default" : "outline"}
                        size="sm"
                        onClick={() => {
                            updateSearchParams((p) => {
                            const next = new Set(categories);
                            if (next.has(c.id)) next.delete(c.id);
                            else next.add(c.id);

                            setMulti(p, "category", Array.from(next));
                            resetToFirstPage(p);
                            });
                        }}
                        >
                        {c.name}
                        </Button>
                    );
                    })}
                </div>
                )}
            </div>

            {/* Fetch indicator */}
            {searchQ.isFetching && !searchQ.isLoading && (
                <div className="text-xs text-muted-foreground">Updating results…</div>
            )}
            </CardContent>
        </Card>

        {/* Results */}
        <Card>
            <CardHeader className="flex-row items-center justify-between space-y-0">
            <CardTitle className="text-base">Results</CardTitle>
            <div className="text-sm text-muted-foreground">{searchQ.data ? `${searchQ.data.total} total` : ""}</div>
            </CardHeader>

            <CardContent className="space-y-4">
            {searchQ.isLoading ? (
                <SkeletonGrid />
            ) : searchQ.isError ? (
                <div className="rounded-md border p-4 text-sm">
                <div className="font-medium">Something went wrong</div>
                <div className="text-muted-foreground mt-1">
                    {searchQ.error instanceof ApiError ? searchQ.error.message : "Search failed"}
                </div>
                </div>
            ) : (searchQ.data?.items?.length ?? 0) === 0 ? (
                <div className="rounded-md border p-6 text-sm space-y-3">
                <div className="font-medium">No results</div>
                <div className="text-muted-foreground">Try removing filters or using a different keyword.</div>
                <div className="flex gap-2">
                    <Button
                    variant="outline"
                    onClick={() => {
                        updateSearchParams((p) => {
                        p.delete("q");
                        p.delete("minPrice");
                        p.delete("maxPrice");
                        p.delete("inStock");
                        p.delete("category");
                        p.set("page", "1");
                        });
                    }}
                    >
                    Clear filters
                    </Button>
                </div>
                </div>
            ) : (
                <>
                <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                    {(searchQ.data?.items ?? []).map((p) => (
                    <div key={p.id} className="rounded-lg border p-4 space-y-2">
                        <div className="font-medium line-clamp-1">{p.name}</div>
                        <div className="text-sm text-muted-foreground">
                        ${p.price.toFixed(2)} • ⭐ {p.rating.toFixed(1)} • {p.inStock ? "In stock" : "Out of stock"}
                        </div>

                        {p.categories?.length ? (
                        <div className="text-xs text-muted-foreground line-clamp-1">
                            {p.categories.map((c) => c.name).join(", ")}
                        </div>
                        ) : null}

                        <div className="pt-2">
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={() => {
                                    const from = encodeURIComponent(loc.pathname + loc.search);
                                    nav(`/products/${p.id}?from=${from}`);
                                }}
                            >
                                View details
                            </Button>
                        </div>
                    </div>
                    ))}
                </div>

                {/* Pagination */}
                <div className="flex items-center justify-between gap-3">
                    <div className="text-sm text-muted-foreground">
                    Page <span className="font-medium text-foreground">{page}</span>
                    {totalPages ? (
                        <>
                        {" "}
                        of <span className="font-medium text-foreground">{totalPages}</span>
                        </>
                    ) : null}
                    </div>

                    <div className="flex gap-2">
                    <Button
                        variant="outline"
                        disabled={page <= 1}
                        onClick={() => {
                        updateSearchParams((p) => {
                            p.set("page", String(Math.max(1, page - 1)));
                        });
                        }}
                    >
                        Prev
                    </Button>
                    <Button
                        variant="outline"
                        disabled={totalPages !== 0 && page >= totalPages}
                        onClick={() => {
                        updateSearchParams((p) => {
                            p.set("page", String(page + 1));
                        });
                        }}
                    >
                        Next
                    </Button>
                    </div>
                </div>
                </>
            )}
            </CardContent>
        </Card>
        </div>
    );
}