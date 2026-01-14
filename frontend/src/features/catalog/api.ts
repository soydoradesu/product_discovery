import { http } from "@/lib/http";

export type Category = { 
    id: number; 
    name: string 
};

export type ProductCategory = { 
    id: number; 
    name: string 
};

export type ProductSummary = {
    id: number;
    name: string;
    price: number;
    rating: number;
    inStock: boolean;
    createdAt: string;
    thumbnail?: string | null;
    categories: ProductCategory[];
};

export type ProductImage = { 
    url: string; 
    position: number 
};

export type ProductDetail = {
    id: number;
    name: string;
    price: number;
    description: string;
    rating: number;
    inStock: boolean;
    createdAt: string;
    images: ProductImage[];
    categories: ProductCategory[];
};

export type SearchResponse = {
    items: ProductSummary[];
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
};

export type SearchParams = {
    q: string;
    categories: number[];
    minPrice: string; 
    maxPrice: string;
    inStock: "any" | "true" | "false";
    sort: "relevance" | "price" | "created_at" | "rating";
    method: "asc" | "desc";
    page: number; 
    pageSize: number;
};

export async function listCategories(): Promise<{ items: Category[] }> {
    return http<{ items: Category[] }>(
        "/api/categories", 
        { method: "GET" }
    );
}

export async function searchProducts(params: SearchParams): Promise<SearchResponse> {
    const sp = new URLSearchParams();

    const q = params.q.trim();
    if (q) sp.set("q", q);

    for (const c of params.categories) {
        sp.append("category", String(c));
    }

    if (params.minPrice.trim() !== "") sp.set("minPrice", params.minPrice.trim());
    if (params.maxPrice.trim() !== "") sp.set("maxPrice", params.maxPrice.trim());

    if (params.inStock !== "any") sp.set("inStock", params.inStock);

    sp.set("sort", params.sort);
    sp.set("method", params.method);
    sp.set("page", String(params.page));
    sp.set("pageSize", String(params.pageSize));

    return http<SearchResponse>(
        `/api/products/search?${sp.toString()}`, 
        { method: "GET" }
    );
}

export async function getProductDetail(id: number): Promise<ProductDetail> {
  return http<ProductDetail>(
    `/api/products/${id}`, 
    { method: "GET" }
);
}