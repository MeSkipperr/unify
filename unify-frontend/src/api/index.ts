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

        const { status } = error.response

        if (isClient) {
            switch (status) {
                case 401:
                    toast.error("Unauthorized. Please login again.")
                    break
                case 403:
                    toast.error("Access denied.")
                    break
                case 500:
                    toast.error("Internal server error. Please try again later.")
                    break
                default:
                    toast.error("Something went wrong.")
            }
        }

        return Promise.reject(error)
    }
)

export default api
