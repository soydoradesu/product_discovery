import { SlidersHorizontal } from "lucide-react";

export type SortValue = "relevance" | "created_at" | "price" | "rating";
export type MethodValue = "asc" | "desc";

export function SortBar({
  sort,
  method,
  pageSize,
  totalLabel,
  onSort,
  onMethod,
  onPageSize,
}: {
  sort: SortValue;
  method: MethodValue;
  pageSize: number;
  totalLabel?: string;
  onSort: (v: SortValue) => void;
  onMethod: (v: MethodValue) => void;
  onPageSize: (n: number) => void;
}) {
  return (
    <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div className="flex items-center gap-2 text-sm">
        <div className="inline-flex items-center gap-2 text-muted-foreground">
          <SlidersHorizontal className="h-4 w-4" />
          <span>Sort</span>
        </div>

        <select
          className="h-9 rounded-md border bg-background px-3 text-sm"
          value={sort}
          onChange={(e) => onSort(e.target.value as SortValue)}
          aria-label="Sort"
        >
          <option value="relevance">Relevant</option>
          <option value="created_at">Date Posted</option>
          <option value="price">Price</option>
          <option value="rating">Rating</option>
        </select>

        <select
          className="h-9 rounded-md border bg-background px-3 text-sm"
          value={method}
          onChange={(e) => onMethod(e.target.value as MethodValue)}
          aria-label="Method"
        >
          <option value="asc">Ascending</option>
          <option value="desc">Descending</option>
        </select>

        <select
          className="h-9 rounded-md border bg-background px-3 text-sm"
          value={String(pageSize)}
          onChange={(e) => onPageSize(Number(e.target.value))}
          aria-label="Page size"
        >
          <option value="20">20</option>
          <option value="50">50</option>
          <option value="100">100</option>
        </select>
      </div>

      {totalLabel ? <div className="text-sm text-muted-foreground">{totalLabel}</div> : null}
    </div>
  );
}
