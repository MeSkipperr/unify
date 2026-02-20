"use client"

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { ChartDataItem } from "./mapHubsToChartData"


// Height per row in px
const ROW_HEIGHT = 50

const TraceDetail = ({chartData}:{chartData:ChartDataItem[]}) => {

    return (
        <div className="w-full flex ">
            <Table >
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-[100px]">Hop</TableHead>
                        <TableHead>Ip Address</TableHead>
                        <TableHead>Name</TableHead>
                        <TableHead className="text-right w-[50px]">Avg</TableHead>
                        <TableHead className="text-right w-[50px]">Min</TableHead>
                        <TableHead className="text-right w-[50px]">Max</TableHead>
                        <TableHead className="text-right w-[50px]">Loss</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {chartData.map((row) => (
                        <TableRow key={row.hop} style={{ height: `${ROW_HEIGHT}px` }}>
                            <TableCell className="font-medium">{row.hop}</TableCell>
                            <TableCell>{row.ip}</TableCell>
                            <TableCell>{row.name}</TableCell>
                            <TableCell className="text-right">{row.avg}</TableCell>
                            <TableCell className="text-right">{row.min}</TableCell>
                            <TableCell className="text-right">{row.max}</TableCell>
                            <TableCell className="text-right">{row.loss}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    )
}

export default TraceDetail
