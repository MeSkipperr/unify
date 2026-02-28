'use client'

import { useEffect, useRef, useState } from "react"
import ChartTraceRoute, { ChartData } from "./chart"
import { ChartDataItem, mapHubsToChartData } from "./mapHubsToChartData"
import { useMTR } from "./useMTR"
import TraceDetail from "./trace-detail"
import { getResultMTR } from "../mtr-result.api"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { clamp } from "@/utils/clamp"
import { useSSE } from "@/hooks/useSSE"
import { MTRData } from "./types"

const MAX_POINTS = 256

const TraceDetailGroup = ({ traceId }: { traceId: string }) => {
    // const {
    //     data: mtrData,
    //     start,
    //     stop,
    //     isRunning
    // } = useMTR("http://localhost:8080/events/mtr")

    const {
        start,
        stop,
        isRunning
    } = useSSE<MTRData>({
        url: "/events/mtr",
        onMessage: (msg) => {
            if (!msg) return
            if (!isInitialized.current) return
            if (msg.id !== traceId) return

            console.log(msg)

            setLatestHubs(msg ? mapHubsToChartData(msg.message.report.hubs) : [])

            setChartLogs(prev => {
                const next = [
                    ...prev,
                    {
                        isConnect: msg.message.report.mtr.Reachable ?? false,
                        ping: msg.message.report.mtr.AvgRTT ?? 0,
                        time: msg.time
                            ? new Date(msg.time).getTime()
                            : Date.now(),
                    },
                ]

                return next.length > MAX_POINTS
                    ? next.slice(-MAX_POINTS)
                    : next
            })
        }
    })

    const [page, setPage] = useState<number>(1)
    const [maxPage, setMaxPage] = useState<number>(1)
    const [latestHubs, setLatestHubs] = useState<ChartDataItem[]>([]);


    const [chartLogs, setChartLogs] = useState<ChartData[]>([])
    const isInitialized = useRef(false)

    const fetchData = async () => {
        if (page === 1) start()
        try {
            const res = await getResultMTR(traceId, { page: page, pageSize: MAX_POINTS })

            setMaxPage(res.totalPages)
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

    useEffect(() => {
        const load = async () => {
            await fetchData()
        }

        load()
    }, [traceId, page])


    const handleNextChart = () => {
        setPage(prev => {
            const newPage = clamp(prev - 1, 0, maxPage)

            if (newPage !== prev) {
                stop()
            }

            return newPage
        })
    }

    const handlePreviewsChart = () => {
        setPage(prev => {
            const newPage = clamp(prev + 1, 0, maxPage)

            stop()

            return newPage
        })
    }


    return (
        <div className="w-[1000px] lg:w-full h-[85dvh] flex flex-col">
            <div className="w-full h-1/2 overflow-y-auto">
                <TraceDetail chartData={latestHubs} />
            </div>
            <div className="w-full h-1/2 flex gap-2">
                <Button disabled={page === maxPage} className="h-full flex justify-center items-center" variant="ghost" onClick={() => handlePreviewsChart()}>
                    <ChevronLeft />
                </Button>
                <ChartTraceRoute chartData={chartLogs} />
                <Button disabled={isRunning} className="h-full flex justify-center items-center" variant="ghost" onClick={() => handleNextChart()}>
                    <ChevronRight />
                </Button>
            </div>
        </div>
    )
}

export default TraceDetailGroup
