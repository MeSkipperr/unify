import api from "@/api"

export type SpeedTestQueryProps = {
    internalIp: string
    serverId: string | number
}
export const getSpeedTestResult = async (
    query: SpeedTestQueryProps
) => {
    const params: Record<string, any> = {}

    if (query.internalIp) {
        params.internalIp = query.internalIp
    }

    if (query.serverId) {
        params.serverId = query.serverId
    }

    const res = await api.get("/api/services/speedtest", { params })
    return res.data
}
