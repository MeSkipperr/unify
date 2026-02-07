'use client'

import ChartTraceRoute from "./chart";
import { ChartDataItem, mapHubsToChartData } from "./mtr/mapHubsToChartData";
import { useMTR } from "./mtr/useMTR";
import TraceDetail from "./trace-detail";

const chartData = [
    { isConnect: true, time: new Date(Date.now() - 60_000).getTime(), ping: 10 },
    { isConnect: true, time: new Date(Date.now() - 50_000).getTime(), ping: 130 },
    { isConnect: true, time: new Date(Date.now() - 40_000).getTime(), ping: 20 },
    { isConnect: false, time: new Date(Date.now() - 30_000).getTime(), ping: 30 },
    { isConnect: false, time: new Date(Date.now() - 20_000).getTime(), ping: 30 },
    { isConnect: true, time: new Date(Date.now() - 10_000).getTime(), ping: 20 },
]

const TraceGroup = () => {
    const mtrData = useMTR("ws://localhost:8080/ws/mtr")

    const latestHubs: ChartDataItem[] = mtrData ? mapHubsToChartData(mtrData.message.report.hubs) : []

    return (
        <div className="w-[1000px] lg:w-full h-[85dvh] flex flex-col">
            <div className="w-full h-1/2 overflow-y-auto">
                <TraceDetail chartData={latestHubs}/>
            </div>
            {/* <div className="w-full h-1/2">
                <ChartTraceRoute chartData={chartData} />
            </div> */}
        </div>
    );
}

export default TraceGroup;