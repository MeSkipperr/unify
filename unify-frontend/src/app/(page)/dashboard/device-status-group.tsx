"use client"

import { DeviceStatus, DeviceStatusProps } from "@/components/dashboard/device-status"
import { DeviceType } from "@/features/device/types"
import { useSSE } from "@/hooks/useSSE"
import { getSummaryConnect } from "@/services/device.service"
import { formatDateTime } from "@/utils/time"
import { Cctv, EthernetPort, LucideIcon, Monitor, Wifi } from "lucide-react"
import { useEffect, useState } from "react"
import { toast } from "sonner"


const deviceType: DeviceTypeStruct[] = [
    {
        key: DeviceType.AccessPoint,
        label: "Access Point",
        icon: Wifi
    },
    {
        key: DeviceType.CCTV,
        label: "CCTV",
        icon: Cctv
    },
    {
        key: DeviceType.IPTV,
        label: "IPTV",
        icon: Monitor
    },
    {
        key: DeviceType.Switch,
        label: "Switch",
        icon: EthernetPort
    },
]

export type DeviceTypeStruct = {
    key: string
    label: string
    icon: LucideIcon
}

type DeviceSummary = {
    type: string
    total: number
    online: number
    offline: number
}


const StatusGroup = () => {
    const [data, setData] = useState<DeviceSummary[]>([])

    const {start } = useSSE<DeviceSummary>({
        url: "/events/dashboard",
        onMessage: (msg) => {
            try {
                const updated = msg as DeviceSummary

                setData((prev) =>
                    prev.map((d) =>
                        d.type === updated.type 
                            ? {
                                ...d,
                                online: updated.online,
                                offline: updated.offline,
                                total: updated.total,
                            }
                            : d
                    )
                )
            } catch (error) {
                console.error(error)
                toast.error(`Failed to get data dashboard`)
            }
        },
    })

    useEffect(() => {
        const loadData = async () => {
            try {
                const result = await getSummaryConnect()
                setData(result?.data ?? [])
            } catch (error) {
                console.error(error)
                toast.error(`Failed to get data dashboard`)
            }
        }

        loadData()
        start()
    }, [])

    return (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 place-items-center">
            {deviceType.map((device) => {
                const summary = data.find((d) => d.type === device.key)

                const props: DeviceStatusProps = {
                    name: device.label,
                    icon: device.icon,
                    total: summary?.total ?? 0,
                    online: summary?.online ?? 0,
                    offline: summary?.offline ?? 0,
                }

                return (
                    <DeviceStatus
                        key={device.key}
                        {...props}
                    />
                )
            })}
        </div>
    )
}

export default StatusGroup
