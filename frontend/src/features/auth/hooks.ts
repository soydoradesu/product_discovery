import { 
    useMutation, 
    useQuery, 
    useQueryClient 
} from "@tanstack/react-query";
import * as api from "./api";

export function useMe() {
  return useQuery({
    queryKey: ["me"],
    queryFn: api.getMe,
    retry: false,
    staleTime: 30_000,
    refetchOnWindowFocus: false
  });
}

export function useLogin() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ 
        email, 
        password 
    }: { email: string; password: string }) => api.login(email, password),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["me"] });
    }
  });
}

export function useLogout() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: api.logout,
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["me"] });
    }
  });
}
