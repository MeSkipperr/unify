import axios from "axios";
import { toast } from "sonner";

const isClient = typeof window !== "undefined";

const api = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
    headers: { "Content-Type": "application/json" },
    withCredentials: true,
});



export default api;
