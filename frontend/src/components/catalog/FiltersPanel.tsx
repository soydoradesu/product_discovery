import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";

type Category = { id: number; name: string };

export function FiltersPanel({
  categories,
  selectedCategoryIds,
  minPrice,
  maxPrice,
  inStock,
  onToggleCategory,
  onMinPrice,
  onMaxPrice,
  onInStock,
  onClear,
  loadingCategories,
  categoriesError,
}: {
  categories: Category[];
  selectedCategoryIds: number[];
  minPrice: string;
  maxPrice: string;
  inStock: "any" | "true" | "false";
  onToggleCategory: (id: number) => void;
  onMinPrice: (v: string) => void;
  onMaxPrice: (v: string) => void;
  onInStock: (v: "any" | "true" | "false") => void;
  onClear: () => void;
  loadingCategories: boolean;
  categoriesError: boolean;
}) {
  return (
    <div className="rounded-xl border bg-card p-4 space-y-5">
      <div className="flex items-center justify-between">
        <div className="text-sm font-semibold">Filter</div>
        <Button variant="ghost" size="sm" onClick={onClear}>
          Reset
        </Button>
      </div>

      {/* Price */}
      <div className="space-y-2">
        <div className="text-xs font-medium text-muted-foreground">Price</div>
        <div className="grid grid-cols-2 gap-2">
          <input
            className="h-9 w-full rounded-md border bg-background px-3 text-sm"
            placeholder="min"
            inputMode="decimal"
            value={minPrice}
            onChange={(e) => onMinPrice(e.target.value)}
            data-testid="min-price-input"
          />
          <input
            className="h-9 w-full rounded-md border bg-background px-3 text-sm"
            placeholder="max"
            inputMode="decimal"
            value={maxPrice}
            onChange={(e) => onMaxPrice(e.target.value)}
            data-testid="max-price-input"
          />
        </div>
      </div>

      {/* Stock */}
      <div className="space-y-2">
        <div className="text-xs font-medium text-muted-foreground">Stock</div>
        <select
          className="h-9 w-full rounded-md border bg-background px-3 text-sm"
          value={inStock}
          onChange={(e) => onInStock(e.target.value as any)}
        >
          <option value="any">All</option>
          <option value="true">Ready</option>
          <option value="false">Out</option>
        </select>
      </div>

      {/* Categories */}
      <div className="space-y-2">
        <div className="text-xs font-medium text-muted-foreground">Category</div>

        {loadingCategories ? (
          <div className="text-sm text-muted-foreground">Loading…</div>
        ) : categoriesError ? (
          <div className="text-sm text-destructive">Fail to load category</div>
        ) : (
          <div className="max-h-64 overflow-auto rounded-md border bg-background p-2 space-y-2">
            {categories.map((c) => {
              const checked = selectedCategoryIds.includes(c.id);
              return (
                <label
                  key={c.id}
                  className="flex items-center gap-2 rounded-md px-2 py-1 hover:bg-muted/50 cursor-pointer"
                >
                  <Checkbox checked={checked} onCheckedChange={() => onToggleCategory(c.id)} />
                  <span className="text-sm">{c.name}</span>
                </label>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
