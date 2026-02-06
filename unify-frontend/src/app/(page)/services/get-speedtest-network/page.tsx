import { Services } from "@/types/service.type"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import SpeedNetworkChart, { SpeedTestInformation } from "./speed-network-chart"
import { Metadata } from "next"

type NetworkDetailProps = {
    name: string,
    interface: string,
    ipAddress: string
}

type ConfigProps = {
    cron: string
    network: NetworkDetailProps[]
    serverId: number[]
}

export const metadata: Metadata = {
    title: "Speed Test - Services | Unify",
    description: "Monitor all speed network in the Unify system.",
};


export default async function Page() {
    const serviceName = "get-speedtest-network"

    const cookieStore = await cookies()
    const cookieHeader = cookieStore
        .getAll()
        .map(c => `${c.name}=${c.value}`)
        .join("; ")

    const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/services/${serviceName}`,
        {
            method: "GET",
            headers: {
                Cookie: cookieHeader,
                "Content-Type": "application/json",
            },
            credentials: "include",
            cache: "no-store",
        }
    )
    if (!response.ok) {
        redirect("/login")
    }

    const res = await response.json()

    const service: Services<ConfigProps> = {
        id: res.ID,
        serviceName: res.ServiceName,
        displayName: res.DisplayName,
        description: res.Description,
        version: res.Version,
        type: res.Type,
        config: {
            cron: res.Config.cron,
            network: res.Config.network.map((net: any) => ({
                name: net.name,
                interface: net.interface,
                ipAddress: net.ip_address,
            })),
            serverId: res.Config.server_id,
        },
        updatedAt: new Date(res.UpdatedAt),
        status: res.Status
    }

    return (
        <div>
            {service.config.network.map((net) => (
                <div key={net.name} className="mb-4">
                        {service.config.serverId.map((serverId) => {
                            const info: SpeedTestInformation = {
                                interface: net.interface,
                                ipAddress: net.ipAddress,
                                name: net.name,
                                serverId: serverId
                            }
                            return (
                                <SpeedNetworkChart information={info} key={serverId}/>
                            )
                        })}
                </div>
            ))}
        </div>
    )
}
