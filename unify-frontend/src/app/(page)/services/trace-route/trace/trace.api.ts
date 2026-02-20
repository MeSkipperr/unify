import api from "@/api"

type ServicesQueryProps = {
    page?: number
    pageSize?: number
}
export const getTraceSession = async (query: ServicesQueryProps) => {
    const params: Record<string, any> = {}

    params.page = query?.page ?? 1;
    params.pageSize = query?.pageSize ?? 50;

    const res = await api.get("/api/services/mtr-sessions/active", { params })
    return res.data
}

export const changeStatus = async (status: boolean) => {
    const res = await api.put("/api/services/mtr-sessions/active", {  status })

    return res.data
}   