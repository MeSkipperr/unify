import api from "@/api"
import { TableQuery } from "@/components/table/types"
import { format } from "date-fns"

export type AdbQuery =
    Pick<TableQuery, "page" | "pageSize" | "search"> & {
        typeServices?: string[]
        sort?: string[]
        date?: Date | string 
    }



export const getAdbResults = async (filter?: AdbQuery) => {
    const params: Record<string, any> = {}

    if (filter?.typeServices?.length) {
        params.typeServices = filter.typeServices
    }

    if (filter?.sort?.length) {
        params.sort = filter.sort.join(",")
    }

    if (filter?.search && filter.search.trim() !== "") {
        params.search = filter.search
    }

    if (filter?.date) {
        params.date =
            filter.date instanceof Date
                ? format(filter.date, "yyyy-MM-dd")
                : filter.date
    }

    params.page = filter?.page ?? 1
    params.pageSize = filter?.pageSize ?? 50

    const res = await api.get("/api/services/adb", { params })
    return res.data
}


