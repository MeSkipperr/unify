"use client"

import { Button } from "@/components/ui/button"
import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"

import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"

import { ArrowUpRight, Badge, Ellipsis } from "lucide-react"
import Link from "next/link"
import { useEffect, useState } from "react"
import { changeStatus, getTraceSession } from "./trace.api"
import { getCompactRelativeTime } from "@/utils/time"
import TableRowSkeleton from "@/components/table/skeleton"
import { DeviceStatus } from "@/components/status"


export type TraceSessionType = {
    id: string
    status: string
    isReachable: boolean

    createdAt: Date
    lastRunAt: Date | null

    sourceIp: string
    destinationIp: string

    protocol: string
    port: number | null

    test: number
    note: string
    sendNotification: boolean
}


const TraceGroup = () => {
    const [datas, setDatas] = useState<TraceSessionType[]>([])
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                setIsLoading(true)
                const res = await getTraceSession({ page: 1, pageSize: 50 })
                const session: TraceSessionType[] = res.data.map((item: any) => ({
                    id: item.ID,
                    status: item.Status,
                    isReachable: item.IsReachable,

                    createdAt: new Date(item.CreatedAt),
                    lastRunAt: item.LastRunAt ? new Date(item.LastRunAt) : null,

                    sourceIp: item.SourceIP,
                    destinationIp: item.DestinationIP,

                    protocol: item.Protocol,
                    port: item.Port,

                    test: item.Test,
                    note: item.Note,
                    sendNotification: item.SendNotification,
                }))

                // console.log(session)
                setDatas(session)
            } catch (error) {
                console.log(error)
            } finally {
                setIsLoading(false)
            }
        }

        fetchData()
    }, [])
    return (
        <Table >
            <TableCaption>Overview of Trace Route Services Currently Running</TableCaption>
            <TableHeader>
                <TableRow>
                    <TableHead>Status</TableHead>
                    <TableHead>Source IP</TableHead>
                    <TableHead>Destiantion IP</TableHead>
                    <TableHead>Protocol</TableHead>
                    <TableHead>Port</TableHead>
                    <TableHead>LastRunAt</TableHead>
                    <TableHead>Is Reachable</TableHead>
                    <TableHead>Details</TableHead>
                    <TableHead className="text-right">Site</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {isLoading ? (
                    Array.from({ length: 10 }).map((_, i) => (
                        <TableRowSkeleton key={i} columns={8} />
                    ))
                ) : (
                    datas.map((data) => {
                        return (
                            <TableRow key={data.id} >

                                {/* Status */}
                                <TableCell>
                                    {data.status}
                                </TableCell>
                                <TableCell>{data.destinationIp}</TableCell>

                                {/* Destination IP */}
                                <TableCell>{data.destinationIp}</TableCell>

                                {/* Protocol */}
                                <TableCell className="uppercase">
                                    {data.protocol}
                                </TableCell>

                                {/* Port */}
                                <TableCell>
                                    {data.port ?? "-"}
                                </TableCell>

                                {/* Last Run At */}
                                <TableCell>
                                    {data.lastRunAt
                                        ? getCompactRelativeTime(data.lastRunAt)
                                        : "-"}
                                </TableCell>

                                {/* Is Reachable */}
                                <TableCell>
                                    <DeviceStatus isConnect={data.isReachable} />
                                </TableCell>

                                {/* Details */}
                                <TableCell >
                                    <Dialog>
                                        <DialogTrigger asChild>
                                            <Button variant="ghost">
                                                <Ellipsis />
                                            </Button>
                                        </DialogTrigger>

                                        <DialogContent>
                                            <DialogHeader>
                                                <DialogTitle>MTR Session Detail</DialogTitle>
                                            </DialogHeader>

                                            <div className="space-y-2 text-sm">
                                                <p><b>ID:</b> {data.id}</p>
                                                <p><b>Source IP:</b> {data.sourceIp || "-"}</p>
                                                <p><b>Destination IP:</b> {data.destinationIp}</p>
                                                <p><b>Protocol:</b> {data.protocol}</p>
                                                <p><b>Test:</b> {data.test}</p>
                                                <p><b>Note:</b> {data.note || "-"}</p>
                                                <p>
                                                    <b>Notification:</b>{" "}
                                                    {data.sendNotification ? "Enabled" : "Disabled"}
                                                </p>
                                            </div>

                                            <DialogFooter>
                                                <DialogClose asChild>
                                                    <Button variant="outline">Close</Button>
                                                </DialogClose>
                                            </DialogFooter>
                                        </DialogContent>
                                    </Dialog>
                                </TableCell>
                                <TableCell
                                    className="text-right flex justify-end"
                                >
                                    <Link href={`/services/trace-route/${data.id}`}>
                                        <ArrowUpRight />
                                    </Link>
                                </TableCell>
                            </TableRow>
                        )
                    }

                    ))
                }
            </TableBody>

        </Table >
    );
}

export default TraceGroup;

