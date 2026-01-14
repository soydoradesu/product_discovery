import { 
    http 
} from "@/lib/http";

export type Me = { 
    userId: number; 
    email: string 
};

export async function getMe(): Promise<Me> {
    return http<Me>(
        "/api/me", 
        { method: "GET" }
    );
}

export async function login(email: string, password: string): Promise<{ ok: boolean }> {
  return http<{ ok: boolean }>
  ("/api/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password })
  });
}

export async function logout(): Promise<{ ok: boolean }> {
  return http<{ ok: boolean }>(
    "/api/auth/logout", 
    { method: "POST" }
);
}