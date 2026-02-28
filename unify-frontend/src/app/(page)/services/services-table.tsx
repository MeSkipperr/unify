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

import { ArrowUpRight, Ellipsis } from "lucide-react"
import Link from "next/link"
import { useEffect, useState } from "react"
import { changeStatus, getServices } from "./service.api"
import { Services } from "@/types/service.type"
import { getCompactRelativeTime } from "@/utils/time"
import TableRowSkeleton from "@/components/table/skeleton"


const ServicesTable = () => {


    const [datas, setDatas] = useState<Services<any>[]>([])
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                setIsLoading(true)
                const res = await getServices({ page: 1, pageSize: 50 })
                const services: Services<any>[] = res.data.map((item: Services<any>) => ({
                    id: item.id,
                    serviceName: item.serviceName,
                    displayName: item.displayName,
                    description: item.description,
                    version: item.version,
                    type: item.type,
                    config: item.config,
                    updatedAt: new Date(item.updatedAt),
                    status: item.status
                }))
                setDatas(services)
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
            <TableCaption>Overview of Backend Services Currently Running</TableCaption>
            <TableHeader>
                <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Version</TableHead>
                    <TableHead>Last Update</TableHead>
                    <TableHead>Description</TableHead>
                    <TableHead className="text-right">Details</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {isLoading ?
                    Array.from({ length: 10 }).map((_, i) => (
                        <TableRowSkeleton key={i} columns={6} />
                    ))
                    :

                    datas.map((data) => {
                        const url = () => {
                            if (data.serviceName === "monitoring-network") return "/devices"
                            if (data.serviceName === "remove-youtube-data-adb" || data.serviceName === "get-uptime-adb") return "/services/adb/" + data.serviceName
                            return "/services/" + data.serviceName
                        }
                        const isDisabled = data.status !== "RUNNING" && data.status !== "STOPPED"
                        return (
                            <TableRow key={"services-" + data.serviceName}>
                                <TableCell className="font-medium">{data.displayName}</TableCell>
                                <TableCell>
                                    <Select
                                        onValueChange={(val) => {
                                            if (!isDisabled) {
                                                changeStatus(`/api/services/${data.serviceName}/status`, val)
                                            }
                                        }}
                                        defaultValue={data.status}
                                        disabled={isDisabled}
                                    >
                                        <SelectTrigger className="w-1/2">
                                            <SelectValue placeholder="Status" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectGroup>
                                                <SelectItem value="RUNNING">Running</SelectItem>
                                                <SelectItem value="STOPPED">Stopped</SelectItem>
                                            </SelectGroup>
                                        </SelectContent>
                                    </Select>
                                </TableCell>
                                <TableCell>{data.version}</TableCell>
                                <TableCell>{getCompactRelativeTime(data.updatedAt)}</TableCell>
                                <TableCell>
                                    <Dialog>
                                        <DialogTrigger asChild>
                                            <Button variant="ghost" className="truncate">
                                                <Ellipsis />
                                            </Button>
                                        </DialogTrigger>
                                        <DialogContent>
                                            <DialogHeader>
                                                <DialogTitle>Description</DialogTitle>
                                            </DialogHeader>
                                            <div className="no-scrollbar -mx-4 max-h-[50vh] overflow-y-auto px-4">
                                                <p className="mb-4 leading-normal">
                                                    {data.description}
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
                                    <Link href={url()}>
                                        <ArrowUpRight />
                                    </Link>
                                </TableCell>
                            </TableRow>
                        )
                    })
                }

            </TableBody>
        </Table >
    );
}

export default ServicesTable;