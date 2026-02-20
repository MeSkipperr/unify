import api from "@/api"

type ServicesQueryProps = {
    page?: number
    pageSize?: number
}
export const getResultMTR = async (id: string,query: ServicesQueryProps) => {
    const params: Record<string, any> = {}

    params.page = query?.page ?? 1;
    params.pageSize = query?.pageSize ?? 50;

    const res = await api.get(`/api/services/mtr-sessions/result/${id}`, { params })
    return res.data
}
