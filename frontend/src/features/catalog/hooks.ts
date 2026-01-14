import { keepPreviousData, useQuery } from "@tanstack/react-query";
import * as api from "./api";

export function useCategories() {
  return useQuery({
    queryKey: ["categories"],
    queryFn: api.listCategories,
    staleTime: 5 * 60_000,
    refetchOnWindowFocus: false
  });
}

export function useProductSearch(params: api.SearchParams, enabled: boolean) {
  const key = [
    "productSearch",
    params.q,
    params.categories.slice().sort((a, b) => a - b).join(","),
    params.minPrice,
    params.maxPrice,
    params.inStock,
    params.sort,
    params.method,
    params.page,
    params.pageSize
  ];

  return useQuery({
    queryKey: key,
    queryFn: () => api.searchProducts(params),
    enabled,
    staleTime: 10_000,
    refetchOnWindowFocus: false,
    placeholderData: keepPreviousData
  });
}
