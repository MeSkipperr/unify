import axios from "axios"
import { toast } from "sonner"

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
        // Network error (tidak ada response)
        if (!error.response) {
            toast.error("Network error. Please check your connection.")
            return Promise.reject(error)
        }

        const { status, data } = error.response

        switch (status) {
            case 500:
                toast.error("Internal server error. Please try again later.")
                break
            default:
                toast.error("Something went wrong.")
        }

        return Promise.reject(error)
    }
)

export default api
