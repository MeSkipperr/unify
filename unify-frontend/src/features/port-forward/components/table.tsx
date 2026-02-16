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
import { ColumnDef } from "@tanstack/react-table"
import DisablePortForward from "./actions"

const searchParameter = {
    id: "port-forward-search-bar",
    label: "Search Session",
    description: "Search devices by IP address",
    placeholder: "e.g. 192.168.1.10 "
}

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
                defaultFilter={dataFilter}
                setFilter={setFilter}
                handleFetchData={handleFetchData}
                isLoading={isLoading}
                searchProps={searchParameter}
                setIsLoading={setIsLoading}
                setSort={setSort}
                sort={sort}
                defaultSort={sortData}
                totalData={totalData}
                addNewData={<NewDataTable handleFetchData={handleFetchData} />}

            />
        </div>
    )
}

export default PortForwardTable