export function ProductCardSkeleton() {
  return (
    <div className="rounded-xl border bg-card shadow-sm overflow-hidden">
      <div className="aspect-[4/3] bg-muted animate-pulse" />
      <div className="p-3 space-y-2">
        <div className="h-4 w-5/6 bg-muted animate-pulse rounded" />
        <div className="h-4 w-3/4 bg-muted animate-pulse rounded" />
        <div className="h-5 w-1/2 bg-muted animate-pulse rounded mt-2" />
        <div className="h-4 w-2/3 bg-muted animate-pulse rounded" />
        <div className="h-9 w-full bg-muted animate-pulse rounded mt-3" />
      </div>
    </div>
  );
}
