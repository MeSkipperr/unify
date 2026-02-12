'use client'

import DataTable from "@/components/table"
import { columns } from "./columns"
import { useState } from "react"
import { dataFilter, sortData } from "../filter-data"
import { SortBy } from "@/components/sort/types"
import { FilterConfig } from "@/components/filter/types"
import { PortForwardResult } from "../types"
import NewDataTable from "./new-data"
import { TableQuery } from "@/components/table/types"
import { useSearchParams } from "next/navigation"
import { getPortForward, PortForwardQuery } from "../api/port-forward.api"
import ActionsColumns from "@/features/device/components/columns/actions"
import { ColumnDef } from "@tanstack/react-table"
import DisablePortForward from "./actions"

const searchParameter = {
    id: "port-forward-search-bar",
    label: "Search Session",
    description: "Search devices by IP address",
    placeholder: "e.g. 192.168.1.10 "
}

// export const portForwardDummy: PortForwardResult[] = [
//     {
//         id: "as",
//         index: 1,
//         listenIp: "0.0.0.0",
//         listenPort: 8080,
//         destIp: "192.168.1.10",
//         destPort: 80,
//         protocol: "TCP",
//         startTime: new Date("2026-02-10T08:00:00Z"),
//         finishTime: new Date("2026-02-10T10:00:00Z"),
//         status: "running",
//         createdAt: new Date("2026-02-10T07:55:00Z"),
//         expiredAt: new Date("2026-02-10T12:00:00Z"),
//     },
//     {
//         id: "as",
//         index: 2,
//         listenIp: "127.0.0.1",
//         listenPort: 2222,
//         destIp: "192.168.1.20",
//         destPort: 22,
//         protocol: "TCP",
//         startTime: new Date("2026-02-10T06:30:00Z"),
//         finishTime: new Date("2026-02-10T07:45:00Z"),
//         status: "finished",
//         createdAt: new Date("2026-02-10T06:25:00Z"),
//         expiredAt: new Date("2026-02-10T08:00:00Z"),
//     },
//     {
//         id: "as",

//         index: 3,
//         listenIp: "0.0.0.0",
//         listenPort: 3307,
//         destIp: "10.0.0.5",
//         destPort: 3306,
//         protocol: "TCP",
//         startTime: new Date("2026-02-10T09:00:00Z"),
//         finishTime: new Date("2026-02-10T09:20:00Z"),
//         status: "failed",
//         createdAt: new Date("2026-02-10T08:58:00Z"),
//         expiredAt: new Date("2026-02-10T09:30:00Z"),
//     },
//     {
//         id: "as",

//         index: 4,
//         listenIp: "0.0.0.0",
//         listenPort: 5353,
//         destIp: "224.0.0.251",
//         destPort: 5353,
//         protocol: "UDP",
//         startTime: new Date("2026-02-10T05:00:00Z"),
//         finishTime: new Date("2026-02-10T11:00:00Z"),
//         status: "running",
//         createdAt: new Date("2026-02-10T04:55:00Z"),
//         expiredAt: new Date("2026-02-10T13:00:00Z"),
//     },
//     {
//         id: "as",

//         index: 5,
//         listenIp: "192.168.1.1",
//         listenPort: 9000,
//         destIp: "192.168.1.50",
//         destPort: 9000,
//         protocol: "TCP",
//         startTime: new Date("2026-02-10T01:00:00Z"),
//         finishTime: new Date("2026-02-10T03:00:00Z"),
//         status: "expired",
//         createdAt: new Date("2026-02-10T00:55:00Z"),
//         expiredAt: new Date("2026-02-10T03:00:00Z"),
//     },
// ]


const PortForwardTable = () => {
    const searchParams = useSearchParams()

    const [data, setData] = useState<PortForwardResult[]>([]);
    const [isLoading, setIsLoading] = useState<boolean>(false)
    const [filter, setFilter] = useState<FilterConfig[]>(dataFilter);
    const [sort, setSort] = useState<SortBy[]>(sortData);
    const [totalData, setTotalData] = useState<number>(1);

    const handleFetchData = async (payload?: TableQuery) => {
        setIsLoading(true)

        const page = payload?.page ?? Number(searchParams.get("page")) ?? 1
        const pageSize = payload?.pageSize ?? Number(searchParams.get("pageSize")) ?? 10
        const search = payload?.search ?? searchParams.get("search") ?? ""

        const dataPayload: PortForwardQuery = {
            page,
            pageSize,
            search
        }


        filter.forEach((filter) => {
            const selectedValues = filter.options
                .filter(opt => opt.isSelected)
                .map(opt => opt.value)

            if (selectedValues.length > 0) {
                switch (filter.key) {
                    case "protocol":
                        dataPayload.protocol = selectedValues as string[]
                        break
                    case "status":
                        dataPayload.status = selectedValues as string[]
                        break
                }
            }
        })

        dataPayload.sort = sort
            .filter(s => s.value !== "none")
            .map(s => `${s.key}:${s.value}`)

        try {
            const result = await getPortForward(dataPayload)
            setTotalData(result.total)
            const portForwards: PortForwardResult[] = result.data.map(
                (item: PortForwardResult) => ({
                    ...item,
                    createdAt: new Date(item.createdAt),
                    expiresAt: new Date(item.expiresAt),
                    lastAppliedAt: item.lastAppliedAt
                        ? new Date(item.lastAppliedAt)
                        : null,
                })
            )
            setData(portForwards)
        } catch (err) {
            console.table(err)
        } finally {
            setIsLoading(false)
        }
    }
    const finalColumns: ColumnDef<PortForwardResult>[] = [
        ...columns,
        {
            id: 'actions',
            header: 'Actions',
            cell: ({ row }) => {
                const deviceRow = row.original
                return (
                    <DisablePortForward
                        row={deviceRow}
                        handleFetchData={() => handleFetchData()}
                    />
                )
            },
            size: 150
        }
    ]



    return (
        <div className="h-[85dvh] ">
            <DataTable
                columns={finalColumns}
                data={data}
                filter={filter}
                setFilter={setFilter}
                handleFetchData={handleFetchData}
                isLoading={isLoading}
                searchProps={searchParameter}
                setIsLoading={setIsLoading}
                setSort={setSort}
                sort={sort}
                totalData={totalData}
                addNewData={<NewDataTable handleFetchData={handleFetchData} />}

            />
        </div>
    )
}

export default PortForwardTable