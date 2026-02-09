'use client'

import { useEffect, useRef, useState } from "react"
import ChartTraceRoute, { ChartData } from "./chart"
import { ChartDataItem, mapHubsToChartData } from "./mapHubsToChartData"
import { useMTR } from "./useMTR"
import TraceDetail from "./trace-detail"
import { getResultMTR } from "../mtr-result.api"

const MAX_POINTS = 256

const TraceDetailGroup = ({ traceId }: { traceId: string }) => {
    const mtrData = useMTR("ws://localhost:8080/ws/mtr")

    const [latestHubs, setLatestHubs] = useState<ChartDataItem[]>([]);


    const [chartLogs, setChartLogs] = useState<ChartData[]>([])
    const isInitialized = useRef(false)

    useEffect(() => {
        const fetchData = async () => {
            try {
                const res = await getResultMTR(traceId, { page: 1, pageSize: MAX_POINTS })

                const parsed: ChartData[] = res.data.map((item: any) => ({
                    isConnect: item.Reachable,
                    ping: item.AvgRTT ?? 0,
                    time: new Date(item.CreatedAt).getTime(),
                }))

                setChartLogs(parsed)
                isInitialized.current = true
            } catch (err) {
                console.error(err)
            }
        }

        fetchData()
    }, [traceId])

    useEffect(() => {
        if (!mtrData) return
        if (!isInitialized.current) return
        if (mtrData.id !== traceId) return

        setLatestHubs(mtrData ? mapHubsToChartData(mtrData.message.report.hubs) : [])

        setChartLogs(prev => {
            const next = [
                ...prev,
                {
                    isConnect: mtrData.message.report.mtr.Reachable ?? false,
                    ping: mtrData.message.report.mtr.AvgRTT ?? 0,
                    time: mtrData.time
                        ? new Date(mtrData.time).getTime()
                        : Date.now(),
                },
            ]

            return next.length > MAX_POINTS
                ? next.slice(-MAX_POINTS)
                : next
        })
    }, [mtrData])

    return (
        <div className="w-[1000px] lg:w-full h-[85dvh] flex flex-col">
            <div className="w-full h-1/2 overflow-y-auto">
                <TraceDetail chartData={latestHubs} />
            </div>
            <div className="w-full h-1/2">
                <ChartTraceRoute chartData={chartLogs} />
            </div>
        </div>
    )
}

export default TraceDetailGroup
