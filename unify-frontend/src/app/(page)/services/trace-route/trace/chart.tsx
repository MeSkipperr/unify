"use client"
import { CartesianGrid, Line, LineChart, ReferenceArea, ReferenceLine, XAxis, YAxis } from "recharts"
import {
    ChartContainer,
    ChartTooltip,
    ChartTooltipContent,
    type ChartConfig,
} from "@/components/ui/chart"


const YELLOW_LIMIT = 50
const RED_LIMIT = 100

const chartConfig = {
    ping: {
        label: "Ping",
        color: "#fff",
    },
} satisfies ChartConfig

type ChartData = {
    isConnect: boolean
    time: number,
    ping: number | null
}

type ChartProps = {
    chartData: ChartData[]
}


const ChartTraceRoute = ({ chartData }: ChartProps) => {
    const mappedData = chartData.map(item => ({
        ...item,
        ping: item.isConnect ? item.ping : null,
    }));

    const validPing = mappedData
        .map(d => d.ping)
        .filter((ping): ping is number => ping !== null);


    const hasYellow = validPing.some(p => p >= YELLOW_LIMIT);
    const hasRed = validPing.some(p => p >= RED_LIMIT);

    const maxY = validPing.length > 0 ? Math.max(...validPing) : 0;

    const disconnectRanges: { start: number; end: number }[] = [];

    for (let i = 0; i < mappedData.length; i++) {
        if (!mappedData[i].isConnect) {
            const start = mappedData[i].time;

            let j = i;
            while (j + 1 < mappedData.length && !mappedData[j + 1].isConnect) {
                j++;
            }

            const end = mappedData[j + 1]?.time;
            if (end) disconnectRanges.push({ start, end });

            i = j;
        }
    }

    return (
        <ChartContainer config={chartConfig} className="h-full w-full">
            <LineChart
                accessibilityLayer
                data={mappedData}
                margin={{
                    left: 12,
                    right: 12,
                }}
            >

                <CartesianGrid vertical={false} />
                <XAxis
                    dataKey="time"
                    type="number"
                    domain={["dataMin", "dataMax"]}
                    tickLine={false}
                    axisLine={false}
                    tickMargin={8}
                    tickFormatter={(value) =>
                        new Date(value).toLocaleTimeString("id-ID", {
                            hour: "2-digit",
                            minute: "2-digit",
                            second: "2-digit",
                        })
                    }
                />


                <YAxis domain={[0, "dataMax"]} />

                <ReferenceArea
                    y1={0}
                    y2={hasYellow ? 50 : maxY}
                    fill="green"
                    fillOpacity={0.5}
                />
                {hasYellow && (
                    <ReferenceArea
                        y1={YELLOW_LIMIT}
                        y2={Math.min(RED_LIMIT, maxY)}
                        fill="yellow"
                        fillOpacity={0.5}
                    />
                )}

                {hasRed && (
                    <ReferenceArea
                        y1={RED_LIMIT}
                        y2={maxY}
                        fill="red"
                        fillOpacity={0.5}
                    />
                )}
                {disconnectRanges.map((range, i) => (
                    <ReferenceArea
                        key={i}
                        x1={range.start}
                        x2={range.end}
                        y1={0}
                        y2={maxY}
                        fill="red"
                        fillOpacity={1}
                        isFront
                    />
                ))}

                <ChartTooltip
                    cursor={false}
                    content={<ChartTooltipContent hideLabel />}
                />
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
    );
}

export default ChartTraceRoute;