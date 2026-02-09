"use client"

import {
    CartesianGrid,
    Line,
    LineChart,
    ReferenceArea,
    XAxis,
    YAxis,
} from "recharts"
import {
    ChartContainer,
    ChartTooltip,
    type ChartConfig,
} from "@/components/ui/chart"
import { formatDateTime } from "@/utils/time"

/* ================= CONFIG ================= */

const YELLOW_LIMIT = 50
const RED_LIMIT = 100

const chartConfig = {
    ping: {
        label: "Ping",
        color: "#fff",
    },
} satisfies ChartConfig

/* ================= TYPES ================= */

export type ChartData = {
    isConnect: boolean
    time: number
    ping: number | null
}

type ChartProps = {
    chartData: ChartData[]
}

/* ================= COMPONENT ================= */

const ChartTraceRoute = ({ chartData }: ChartProps) => {
    /**
     * 1️⃣ Tambahkan index agar XAxis jaraknya konsisten
     */
    const indexedData = chartData.map((item, index) => ({
        ...item,
        index,
        ping: item.isConnect ? item.ping : 0,
    }))

    /**
     * 2️⃣ Hitung ping valid
     */
    const validPing = indexedData
        .map(d => d.ping)
        .filter((p): p is number => p !== 0)

    const maxPing = validPing.length > 0 ? Math.max(...validPing) : 0
    const maxY = maxPing > 0 ? Math.ceil(maxPing * 1.2) : 10

    const hasYellow = validPing.some(p => p >= YELLOW_LIMIT)
    const hasRed = validPing.some(p => p >= RED_LIMIT)

    /**
     * 3️⃣ Hitung range disconnect (index-based)
     */
    const disconnectRanges: { start: number; end: number }[] = []

    for (let i = 0; i < indexedData.length; i++) {
        if (!indexedData[i].isConnect) {
            const start = indexedData[i].index

            let j = i
            while (j + 1 < indexedData.length && !indexedData[j + 1].isConnect) {
                j++
            }

            const end = indexedData[j + 1]?.index ?? j
            disconnectRanges.push({ start, end })
            i = j
        }
    }

    return (
        <ChartContainer config={chartConfig} className="h-full w-full">
            <LineChart
                data={indexedData}
                margin={{ left: 12, right: 12 }}
            >
                <CartesianGrid vertical={false} />

                {/* ================= X AXIS (INDEX BASED) ================= */}
                <XAxis
                    dataKey="index"
                    type="number"
                    domain={[0, "dataMax"]}
                    tickLine={false}
                    axisLine={false}
                    tickMargin={8}
                    tickFormatter={(i) =>
                        indexedData[i]
                            ? formatDateTime(new Date(indexedData[i].time))
                            : ""
                    }
                />

                {/* ================= TOOLTIP ================= */}
                <ChartTooltip
                    content={({ payload }) => {
                        if (!payload || payload.length === 0) return null
                        const { time, ping, isConnect } = payload[0].payload

                        return (
                            <div className="rounded-md border bg-background px-3 py-2 text-sm shadow">
                                <div className="font-medium">
                                    {formatDateTime(new Date(time))}
                                </div>
                                <div>Ping: {ping} ms</div>
                                {!isConnect && (
                                    <div className="text-red-500 font-semibold">
                                        Disconnected
                                    </div>
                                )}
                            </div>
                        )
                    }}
                />

                {/* ================= Y AXIS ================= */}
                <YAxis
                    dataKey="ping"
                    domain={[0, maxY]}
                />

                {/* ================= BACKGROUND ZONES ================= */}
                <ReferenceArea
                    y1={0}
                    y2={hasYellow ? YELLOW_LIMIT : maxY}
                    fill="green"
                    fillOpacity={0.4}
                />

                {hasYellow && (
                    <ReferenceArea
                        y1={YELLOW_LIMIT}
                        y2={Math.min(RED_LIMIT, maxY)}
                        fill="yellow"
                        fillOpacity={0.4}
                    />
                )}

                {hasRed && (
                    <ReferenceArea
                        y1={RED_LIMIT}
                        y2={maxY}
                        fill="red"
                        fillOpacity={0.4}
                    />
                )}

                {/* ================= DISCONNECT AREA ================= */}
                {disconnectRanges.map((range, i) => (
                    <ReferenceArea
                        key={i}
                        x1={range.start}
                        x2={range.end}
                        y1={0}
                        y2={maxY}
                        fill="red"
                        fillOpacity={0.9}
                        isFront
                    />
                ))}

                {/* ================= LINE ================= */}
                <Line
                    dataKey="ping"
                    type="step"
                    stroke="var(--color-ping)"
                    strokeWidth={2}
                    dot={false}
                    connectNulls={false}
                />
            </LineChart>
        </ChartContainer>
    )
}

export default ChartTraceRoute
