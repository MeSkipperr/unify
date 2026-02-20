import api from "@/api";

export const getSummaryConnect = async () => {
    const res = await api.get("/api/devices/summary");
    return res.data;
}