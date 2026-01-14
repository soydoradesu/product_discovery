export type ApiErrorEnvelope = {
    error: { 
        code: string; 
        message: string 
    };
};

export class ApiError extends Error {
    code: string;
    status: number;

    constructor(message: string, code: string, status: number) {
        super(message);
        this.code = code;
        this.status = status;
    }
}

async function parseApiError(res: Response): Promise<ApiError> {
    let msg = `Request failed (${res.status})`;
    let code = "HTTP_ERROR";

    const data = (await res.json()) as ApiErrorEnvelope;
    if (data?.error?.message) {
        msg = data.error.message;
    }
    if (data?.error?.code) {
        code = data.error.code;
    }
    
    return new ApiError(msg, code, res.status);
}

export async function http<T>(input: RequestInfo, init?: RequestInit): Promise<T> {
    const res = await fetch(input, {
        ...init,
        credentials: "include",
        headers: {
        "Content-Type": "application/json",
        ...(init?.headers ?? {})
        }
    });

    if (!res.ok){ 
        throw await parseApiError(res);
    }
    const text = await res.text();

    if (!text) {
        return {} as T;
    }
    
    return JSON.parse(text) as T;
}