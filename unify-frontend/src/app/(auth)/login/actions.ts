import { loginRequest } from "@/services/auth.service"


export async function loginAction(payload: { username: string; password: string }) {
    if (!payload.username || !payload.password) throw new Error("Invalid credentials");

    const res = await loginRequest(payload);
    return res.data;
}
