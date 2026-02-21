import { Services } from "@/types/service.type"
import { cookies, headers } from "next/headers"
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

    const headersList = await headers() 
    const host = headersList.get("host")

    const protocol = "http"

    const baseUrl =
        process.env.NEXT_PUBLIC_API_BASE_URL ||
        `${protocol}://${host}`

    const response = await fetch(
        `${baseUrl }/api/services/${serviceName}`,
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
        id: res.id,
        serviceName: res.serviceName,
        displayName: res.displayName,
        description: res.description,
        version: res.version,
        type: res.type,
        config: {
            cron: res.config.cron,
            network: res.config.network.map((net: NetworkDetailProps) => ({
                name: net.name,
                interface: net.interface,
                ipAddress: net.ipAddress,
            })),
            serverId: res.config.server_id,
        },
        updatedAt: new Date(res.updatedAt),
        status: res.status
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
                            <SpeedNetworkChart information={info} key={serverId} />
                        )
                    })}
                </div>
            ))}
        </div>
    )
}
