import { SearchX } from "lucide-react";
import { Button } from "@/components/ui/button";

export function EmptyState({
  title,
  description,
  onClear,
}: {
  title: string;
  description?: string;
  onClear?: () => void;
}) {
  return (
    <div className="rounded-xl border bg-card p-8 text-center">
      <div className="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-muted">
        <SearchX className="h-6 w-6 text-muted-foreground" />
      </div>
      <div className="text-base font-semibold">{title}</div>
      {description ? <div className="mt-1 text-sm text-muted-foreground">{description}</div> : null}

      {onClear ? (
        <div className="mt-4">
          <Button variant="outline" onClick={onClear}>
            Reset Filter
          </Button>
        </div>
      ) : null}
    </div>
  );
}
