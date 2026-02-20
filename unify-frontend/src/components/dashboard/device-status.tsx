"use client"

import { Wifi, Clock, Server, CheckCircle, XCircle, LucideIcon } from "lucide-react"
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"

export type DeviceStatusProps = {
    online?: number
    offline?: number
    total?: number
    name?: string
    icon?: LucideIcon
}

export const DeviceStatus = ({
    online = 60,
    offline = 20,
    total = 80,
    name = "Device",
    icon: Icon,
}: DeviceStatusProps) => {

    // Validasi agar tidak melebihi total
    const safeOnline = Math.min(online, total)
    const safeOffline = Math.min(offline, total - safeOnline)

    const onlinePercentage = Math.round((safeOnline / total) * 100)
    const offlinePercentage = Math.round((safeOffline / total) * 100)

    const getStatus = () => {
        if (offlinePercentage > 30)
            return { label: "Critical", variant: "destructive" as const }
        if (offlinePercentage > 10)
            return { label: "Warning", variant: "secondary" as const }
        return { label: "Stable", variant: "default" as const }
    }

    const status = getStatus()

    return (
        <Card className="w-full rounded-2xl shadow-md hover:shadow-xl transition-all duration-300">
            <CardHeader className="pb-2">
                <div className="flex items-center justify-between">
                    <CardTitle className="flex items-center gap-2 text-lg">
                        {Icon && <Icon className="w-5 h-5 text-primary" />}
                        {name}
                    </CardTitle>

                    <Badge variant={status.variant}>{status.label}</Badge>
                </div>

            </CardHeader>

            <CardContent className="space-y-4">

                {/* Category Breakdown */}
                <div className="grid grid-cols-2 gap-3 text-sm font-medium">

                <div className="flex items-center justify-between bg-green-500/10 border border-green-500/20 rounded-xl px-3 py-2 hover:bg-green-500/15 transition-colors">
                        <div className="flex items-center gap-2 text-green-600 dark:text-green-400">
                            <CheckCircle className="w-4 h-4" />
                            <span className="text-sm font-medium">
                                Online
                            </span>
                        </div>

                        <div className="text-right">
                            <div className="font-semibold text-green-600 ">
                                {safeOnline}
                            </div>
                            <div className="text-xs text-muted-foreground">
                                {onlinePercentage}%
                            </div>
                        </div>
                    </div>

                    <div className="flex items-center justify-between bg-red-500/10 border border-red-500/20 rounded-xl px-3 py-2 hover:bg-red-500/15 transition-colors">
                        <div className="flex items-center gap-2 text-red-600 dark:text-red-400">
                            <XCircle className="w-4 h-4" />
                            <span className="text-sm font-medium">
                                Offline
                            </span>
                        </div>

                        <div className="text-right">
                            <div className="font-semibold text-red-600">
                                {safeOffline}
                            </div>
                            <div className="text-xs text-muted-foreground">
                                {offlinePercentage}%
                            </div>
                        </div>
                    </div>


                </div>

                {/* Online Progress */}
                <div className="space-y-2">
                    <div className="flex justify-between text-xs text-muted-foreground">
                        <span>Online Ratio</span>
                        <span>{onlinePercentage}%</span>
                    </div>
                    <Progress value={onlinePercentage} />
                </div>

                <Separator />

                {/* Summary */}
                <div className="grid grid-cols-3 gap-3 text-center">

                    <div>
                        <div className="text-lg font-bold text-primary">{total}</div>
                        <div className="text-xs text-muted-foreground">Total</div>
                    </div>

                    <div>
                        <div className="text-lg font-bold text-green-600">{safeOnline}</div>
                        <div className="text-xs text-muted-foreground">Online</div>
                    </div>

                    <div>
                        <div className="text-lg font-bold text-red-600">{safeOffline}</div>
                        <div className="text-xs text-muted-foreground">Offline</div>
                    </div>

                </div>

                <Separator />

                {/* Capacity Info */}
                <div className="flex items-center justify-between text-xs text-muted-foreground">
                    <div className="flex items-center gap-2">
                        <Server className="w-4 h-4" />
                        Device Availability
                    </div>
                    <span>
                        {safeOnline} / {total} Active Devices
                    </span>
                </div>

            </CardContent>
        </Card>
    )
}
