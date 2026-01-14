// Tests:
// URL updates immediately while typing
// debounce affects query hook params (q)
// changing filter resets page to 1
// empty state renders
// error state renders

import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { screen, act, fireEvent } from "@testing-library/react";

import { renderWithRoute } from "@/test/render";
import { HomePage } from "@/pages/HomePage";
import { ApiError } from "@/lib/http";

const useCategoriesMock = vi.fn();
const useProductSearchMock = vi.fn();

vi.mock("@/features/auth/hooks", () => ({
  useMe: () => ({ data: { email: "demo@example.com" }, isSuccess: true }),
  useLogout: () => ({ mutateAsync: vi.fn().mockResolvedValue({}), isPending: false }),
}));

vi.mock("@/features/catalog/hooks", () => ({
  useCategories: () => useCategoriesMock(),
  useProductSearch: (params: any) => useProductSearchMock(params),
}));

describe("HomePage", () => {
  beforeEach(() => {
    useCategoriesMock.mockReturnValue({
      isLoading: false,
      isError: false,
      data: { items: [{ id: 1, name: "Laptop" }, { id: 2, name: "Phone" }] },
    });

    useProductSearchMock.mockImplementation((params: any) => ({
      isLoading: false,
      isError: false,
      isFetching: false,
      data: {
        items: [],
        page: params.page,
        pageSize: params.pageSize,
        total: 0,
        totalPages: 0,
      },
    }));
  });

  afterEach(() => {
    vi.clearAllMocks();
    vi.useRealTimers();
  });

  it("updates URL immediately when typing", () => {
    renderWithRoute(<HomePage />, "/");

    const input = screen.getByPlaceholderText("Search name or description…");
    fireEvent.change(input, { target: { value: "N" } });

    expect(screen.getByTestId("location").textContent).toContain("q=N");
    expect(screen.getByTestId("location").textContent).toContain("page=1");
  });

  it("debounces q passed to search hook (350ms)", async () => {
    vi.useFakeTimers();

    renderWithRoute(<HomePage />, "/");

    const input = screen.getByPlaceholderText("Search name or description…");
    fireEvent.change(input, { target: { value: "Nova" } });

    // Immediately after change: debounced q should still be old value ("")
    const lastBefore = useProductSearchMock.mock.calls.at(-1)?.[0];
    expect(lastBefore.q).toBe("");

    await act(async () => {
      vi.advanceTimersByTime(400);
    });

    const lastAfter = useProductSearchMock.mock.calls.at(-1)?.[0];
    expect(lastAfter.q).toBe("Nova");
  });

  it("changing a filter resets page=1", () => {
    renderWithRoute(<HomePage />, "/?q=Nova&page=3&pageSize=20");

    const minInput = screen.getByPlaceholderText("min");
    fireEvent.change(minInput, { target: { value: "100" } });

    expect(screen.getByTestId("location").textContent).toContain("minPrice=100");
    expect(screen.getByTestId("location").textContent).toContain("page=1");
  });

  it("renders empty state when items=[]", () => {
    renderWithRoute(<HomePage />, "/?q=Nova&page=1");
    expect(screen.getByText("No results")).toBeInTheDocument();
  });

  it("renders error state when search hook errors", () => {
    useProductSearchMock.mockImplementationOnce(() => ({
      isLoading: false,
      isError: true,
      isFetching: false,
      error: new ApiError("boom", "INTERNAL", 500),
      data: undefined,
    }));

    renderWithRoute(<HomePage />, "/?q=Nova&page=1");
    expect(screen.getByText("Something went wrong")).toBeInTheDocument();
    expect(screen.getByText("boom")).toBeInTheDocument();
  });
});