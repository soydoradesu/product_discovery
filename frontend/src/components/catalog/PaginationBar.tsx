import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";

export function PaginationBar({
  page,
  totalPages,
  onPrev,
  onNext,
}: {
  page: number;
  totalPages: number;
  onPrev: () => void;
  onNext: () => void;
}) {
  return (
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

      <div className="flex items-center gap-2">
        <Button variant="outline" size="sm" disabled={page <= 1} onClick={onPrev}>
          <ChevronLeft className="mr-1 h-4 w-4" />
          Prev
        </Button>
        <Button
          variant="outline"
          size="sm"
          disabled={totalPages !== 0 && page >= totalPages}
          onClick={onNext}
        >
          Next
          <ChevronRight className="ml-1 h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}
