"use client"

import * as React from "react"
import { Area, AreaChart, CartesianGrid, XAxis } from "recharts"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    ChartContainer,
    ChartLegend,
    ChartLegendContent,
    ChartTooltip,
    ChartTooltipContent,
    type ChartConfig,
} from "@/components/ui/chart"

import { formatDateTime } from "@/utils/time"
import { StackId } from "recharts/types/util/ChartUtils"
import { getSpeedTestResult } from "./speedtest.api"


const chartConfig = {
    visitors: {
        label: "Visitors",
    },
    upload: {
        label: "Upload",
        color: "#00e1e2",
    },
    download: {
        label: "Download",
        color: "#A020F0",
    },
    ping: {
        label: "Ping",
        color: "#00b809",
    },
} satisfies ChartConfig

type DataSpeedTestResultProps = {
    date: Date
    ping: number
    upload: number
    download: number
}
export type SpeedTestInformation = {
    name: string
    ipAddress: string
    interface: string
    serverId: number
}
type SpeedNetworkChartProps = {
    information: SpeedTestInformation
}

const SpeedNetworkChart = ({ information }: SpeedNetworkChartProps) => {
    const [data, setData] = React.useState<DataSpeedTestResultProps[]>([]);
    const [serverName, setServerName] = React.useState<string>("");

    React.useEffect(() => {
        const fetchSpeedTest = async () => {
            const query = {
                internalIp: information.ipAddress,
                serverId: information.serverId,
            }

            try {
                const result = await getSpeedTestResult(query)
                const mappedData: DataSpeedTestResultProps[] = result.map((item: any) => ({
                    date: item.testedAt,
                    ping: item.pingMs,
                    upload: item.uploadMbps,
                    download: item.downloadMbps,
                }))
                setServerName(result[0].serverName)
                console.log(mappedData)
                setData(mappedData)

            } catch (error) {
                console.error(error)
            }
        }

        fetchSpeedTest()
    }, [])


    return (
        <Card className="pt-0">
            <CardHeader className="flex flex-col gap-3 border-b py-5 sm:flex-row sm:items-center">
                <div className="flex flex-1 flex-col gap-1">
                    <CardTitle className="text-base font-semibold">
                        {information.name}
                    </CardTitle>
                    <CardDescription className="text-sm">
                        Upload, download, and ping performance overview
                    </CardDescription>
                </div>

                <div className="flex flex-wrap  gap-4 text-xs text-muted-foreground">
                    <div className="flex items-center gap-1">
                        <span className="font-medium">Source:</span>
                        <span>{information.ipAddress}</span>
                    </div>
                    <div className="flex items-center gap-1">
                        <span className="font-medium">Interface:</span>
                        <span>{information.interface}</span>
                    </div>
                    <div className="flex items-center gap-1">
                        <span className="font-medium">Server:</span>
                        <span>{serverName}</span>
                    </div>
                </div>
            </CardHeader>

            <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6 ">
                <ChartContainer
                    config={chartConfig}
                    className="aspect-auto h-[250px]"
                >
                    <AreaChart data={data}>
                        <defs>
                            <linearGradient id="fillDownload" x1="0" y1="0" x2="0" y2="1">
                                <stop
                                    offset="5%"
                                    stopColor="var(--color-download)"
                                    stopOpacity={1}
                                />
                                <stop
                                    offset="95%"
                                    stopColor="var(--color-download)"
                                    stopOpacity={0.5}
                                />
                            </linearGradient>
                            <linearGradient id="fillUpload" x1="0" y1="0" x2="0" y2="1">
                                <stop
                                    offset="5%"
                                    stopColor="var(--color-upload)"
                                    stopOpacity={1}
                                />
                                <stop
                                    offset="95%"
                                    stopColor="var(--color-upload)"
                                    stopOpacity={0.5}
                                />
                            </linearGradient>
                            <linearGradient id="fillPing" x1="0" y1="0" x2="0" y2="1">
                                <stop
                                    offset="5%"
                                    stopColor="var(--color-ping)"
                                    stopOpacity={1}
                                />
                                <stop
                                    offset="95%"
                                    stopColor="var(--color-ping)"
                                    stopOpacity={0.5}
                                />
                            </linearGradient>
                        </defs>
                        <CartesianGrid vertical={false} />
                        <XAxis
                            dataKey="date"
                            tickLine={false}
                            axisLine={false}
                            interval={0}
                            tickMargin={8}
                            minTickGap={32}
                            tickFormatter={(value) => {
                                const date = new Date(value)
                                return date.toLocaleString("id-ID", {
                                    hour: "2-digit",
                                    minute: "2-digit",
                                })
                            }}
                        />
                        <ChartTooltip
                            cursor={false}
                            content={
                                <ChartTooltipContent
                                    labelFormatter={(value) => {
                                        const date = new Date(value)
                                        return formatDateTime(date)
                                    }}
                                    formatter={(value, name, item: any) => {
                                        let unit = ""

                                        if (name === "download" || name === "upload") {
                                            unit = " Mbps"
                                        }

                                        if (name === "ping") {
                                            unit = " ms"
                                        }

                                        return [
                                            <span key="data-speedtes-result" className="flex w-full justify-between gap-2">
                                                <span className="capitalize">
                                                    {name}
                                                </span>
                                                <span style={{ color: item?.color }}>
                                                    {value}
                                                    {unit}
                                                </span>
                                            </span>

                                        ]
                                    }}
                                    indicator="dot"
                                />
                            }
                        />

                        <Area
                            dataKey="download"
                            type="natural"
                            fill="url(#fillDownload)"
                            stroke="var(--color-download)"
                            stackId="a"
                        />
                        <Area
                            dataKey="upload"
                            type="natural"
                            fill="url(#fillUpload)"
                            stroke="var(--color-upload)"
                            stackId="a"
                        />
                        <Area
                            dataKey="ping"
                            type="natural"
                            fill="url(#fillPing)"
                            stroke="var(--color-ping)"
                            stackId="a"
                        />
                        <ChartLegend content={<ChartLegendContent />} />
                    </AreaChart>
                </ChartContainer>
            </CardContent>
        </Card>
    );
}

export default SpeedNetworkChart;