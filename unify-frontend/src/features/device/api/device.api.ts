import api from "@/api"

export type DeviceQuery = {
    status?: boolean[]
    notification?: boolean[]
    type?: string[]
    sort?: string[]
    search?: string
    page?: number
    pageSize?: number
}

export const getDevices = async (filter?: DeviceQuery) => {
    const params: Record<string, any> = {}

    if (filter?.status?.length) {
        params.status = filter.status
    }

    if (filter?.notification?.length) {
        params.notification = filter.notification
    }

    if (filter?.type?.length) {
        params.type = filter.type
    }

    if (filter?.sort?.length) {
        params.sort = filter.sort.join(",")
    }

    if (filter?.search && filter.search.trim() !== "") {
        params.search = filter.search
    }

    params.page = filter?.page ?? 1;
    params.pageSize = filter?.pageSize ?? 50;

    const res = await api.get("/api/devices", { params })
    console.log(res)
    return res.data
}
