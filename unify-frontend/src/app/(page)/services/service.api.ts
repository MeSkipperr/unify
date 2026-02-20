import api from "@/api"

type ServicesQueryProps = {
    page?: number
    pageSize?: number
}

export const getServices = async (query: ServicesQueryProps) => {
    const params: Record<string, any> = {}

    params.page = query?.page ?? 1;
    params.pageSize = query?.pageSize ?? 50;

    const res = await api.get("/api/services", { params })
    return res.data
}
export const getServicesByName = async (name:string) => {
    const res = await api.get(`/api/services/${name}`)
    return res.data
}

export const changeStatus = async (url: string, status: string) => {
    const res = await api.put(url, {  status })

    return res.data
}   