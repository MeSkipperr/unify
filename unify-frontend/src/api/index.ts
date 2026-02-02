import axios from "axios"
import { toast } from "sonner"

const isClient = typeof window !== "undefined"

const api = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
    headers: {
        "Content-Type": "application/json",
    },
    withCredentials: true,
})

api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (!error.response) {
            if (isClient) {
                toast.error("Network error. Please check your connection.")
            }
            return Promise.reject(error)
        }
        return Promise.reject(error)
    }
)

export default api
