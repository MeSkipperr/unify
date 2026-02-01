import api from "@/api"

export const updateDeviceNotification = async (
    deviceId: string,
    notification: boolean
) => {
    const res = await api.patch(`/devices/${deviceId}`, { notification })
    return res.data
}
