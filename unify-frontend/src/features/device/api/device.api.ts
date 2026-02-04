import api from "@/api"
import { TableQuery } from "@/components/table/types"
import { normalizeIPv4 } from "@/utils/ipv4"
import { normalizeMacAddress } from "@/utils/macAddress"
import { DeviceType } from "../types"

export type DeviceQuery =
    Pick<TableQuery, "page" | "pageSize" | "search"> & {
        status?: boolean[]
        notification?: boolean[],
        type?: String[]
        sort?: string[]
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

export type CreateDevicePayload = {
    name: string;
    ipAddress: string;
    macAddress: string;
    roomNumber?: string;
    description: string;
    type: DeviceType
};

export const addDevice = async (
    payload: CreateDevicePayload
) => {
    const normalizedPayload: CreateDevicePayload = {
        ...payload,
        name: payload.name.trim(),
        description: payload.description.trim(),
        ipAddress: normalizeIPv4(payload.ipAddress),
        macAddress: normalizeMacAddress(payload.macAddress),
        roomNumber: payload.roomNumber?.trim(),
        type: payload.type,
    };
    const res = await api.post("/api/devices", { normalizedPayload })

    return res.data;
};


export const deleteDevice = async (deviceId: string) => {
    const res = await api.delete(`/api/devices/${deviceId}`)
    return res.data
};
