import api from "@/api"

export type LoginPayload = {
    username: string
    password: string
}

export const loginRequest = async (payload: LoginPayload) => {
    const res = await api.post("/auth/login", payload);
    return res.data;
};
