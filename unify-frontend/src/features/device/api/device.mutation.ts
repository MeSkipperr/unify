import api from "@/api"

export const updateDeviceNotification = async (
    deviceId: string,
    notification: boolean
) => {
    const res = await api.patch(`/api/devices/${deviceId}/notification`, { notification })
    return res.data
}
