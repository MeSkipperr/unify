import api from "@/api"
import { TableQuery } from "@/components/table/types"

export type DeviceQuery =
    Pick<TableQuery, "page" | "pageSize" | "search"> & {
        status?: boolean[]
        notification?: boolean[],
        type?:String[]
        sort?:string[]
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
    return res.data
}
