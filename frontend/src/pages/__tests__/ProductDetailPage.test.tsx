// Test:
// Test “Back to results” restores exact prior URL state.

import { describe, it, expect, vi } from "vitest";
import userEvent from "@testing-library/user-event";
import { screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router-dom";

import { ProductDetailPage } from "@/pages/ProductDetailPage";
import { LocationDisplay } from "@/test/render";

const useProductDetailMock = vi.fn();

vi.mock("@/features/catalog/hooks", () => ({
  useProductDetail: (id: number) => useProductDetailMock(id)
}));

describe("ProductDetailPage", () => {
  it("Back to results navigates to the from= URL", async () => {
    const user = userEvent.setup();

    useProductDetailMock.mockReturnValue({
      isLoading: false,
      isError: false,
      data: {
        id: 1,
        name: "Nova Phone 0001",
        price: 199.99,
        description: "desc",
        rating: 4.5,
        inStock: true,
        createdAt: new Date().toISOString(),
        images: [],
        categories: []
      }
    });

    const from = encodeURIComponent("/?q=Nova&page=2&sort=relevance&method=desc");
    const route = `/products/1?from=${from}`;

    renderWithRoute();

    async function renderWithRoute() {
      return (
        <MemoryRouter initialEntries={[route]}>
          <LocationDisplay />
          <Routes>
            <Route path="/products/:id" element={<ProductDetailPage />} />
            {/* target route */}
            <Route path="/" element={<div>Home</div>} />
          </Routes>
        </MemoryRouter>
      );
    }

    // We need actual render call:
    const { render } = await import("@testing-library/react");
    render(await renderWithRoute());

    await user.click(screen.getByRole("button", { name: "← Back to results" }));
    expect(screen.getByTestId("location").textContent).toBe("/?q=Nova&page=2&sort=relevance&method=desc");
  });
});
