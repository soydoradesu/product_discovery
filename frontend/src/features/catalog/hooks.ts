import { 
    keepPreviousData,
    useQuery 
} from "@tanstack/react-query";
import * as api from "./api";

export function useCategories() {
    return useQuery({
        queryKey: ["categories"],
        queryFn: api.listCategories,
        staleTime: 5 * 60_000,
        refetchOnWindowFocus: false
    });
}

export function useProductSearch(params: api.SearchParams) {
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
        staleTime: 10_000,
        refetchOnWindowFocus: false,
        placeholderData: keepPreviousData
    });
}

export function useProductDetail(id: number) {
    return useQuery({
        queryKey: ["productDetail", id],
        queryFn: () => api.getProductDetail(id),
        staleTime: 30_000,
        refetchOnWindowFocus: false,
        enabled: Number.isFinite(id) && id > 0
    });
}
