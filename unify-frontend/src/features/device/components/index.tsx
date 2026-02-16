"use client"
import { columns } from "./columns";
import { dataFilter, sortData } from "../filter-data";
import { Device } from "../types";
import DataTable from "@/components/table";
import React from "react";
import { SortBy } from "@/components/sort/types";
import { DeviceQuery, getDevices } from "../api/device.api";
import { FilterConfig } from "@/components/filter/types";
import { TableQuery } from "@/components/table/types";
import { updateFilterOption } from "../utils/select-options";
import NewDataTable from "./new-data";
import { ColumnDef } from "@tanstack/react-table";
import ActionsColums from "./columns/actions";
import { useSearchParams } from "next/navigation";
import ActionsColumns from "./columns/actions";

const DeviceTableData = (
    { selectType = "" }: { selectType?: string }
) => {
    const searchParams = useSearchParams()

    const defaultFilter = updateFilterOption(dataFilter, "type", selectType)
    const [data, setData] = React.useState<Device[]>([]);
    const [sortOptions, setSortOptions] = React.useState<SortBy[]>(sortData)
    const [filter, setFilter] = React.useState<FilterConfig[]>(defaultFilter);
    const [isLoading, setIsLoading] = React.useState<boolean>(true);
    const [totalData, setTotalData] = React.useState<number>(1);


    const searchParameter = {
        id: "device-search-bar",
        description: "Search by name, MAC, IP, room, device type, or description.",
        label: "Search Device",
        placeholder: "DPSCY-..."
    }


    const handleFetchData = async (payload?: TableQuery) => {
        setIsLoading(true)

        const page = payload?.page ?? Number(searchParams.get("page")) ?? 1
        const pageSize = payload?.pageSize ?? Number(searchParams.get("pageSize")) ?? 10
        const search = payload?.search ?? searchParams.get("search") ?? ""

        const dataPayload: DeviceQuery = {
            page,
            pageSize,
            search
        }

        console.log(dataPayload)

        filter.forEach((filter) => {
            const selectedValues = filter.options
                .filter(opt => opt.isSelected)
                .map(opt => opt.value)

            if (selectedValues.length > 0) {
                switch (filter.key) {
                    case "status":
                        dataPayload.status = selectedValues as boolean[]
                        break
                    case "notification":
                        dataPayload.notification = selectedValues as boolean[]
                        break
                    case "type":
                        dataPayload.type = selectedValues as string[]
                        break
                }
            }
        })

        dataPayload.sort = sortOptions
            .filter(s => s.value !== "none")
            .map(s => `${s.key}:${s.value}`)

        try {
            const result = await getDevices(dataPayload)
            setTotalData(result.total)
            const devices: Device[] = result.data.map((item: any) => ({
                id: item.ID,
                index: item.index,
                name: item.Name,
                ipAddress: item.IPAddress,
                macAddress: item.MacAddress,
                roomNumber: item.RoomNumber,
                isConnect: item.IsConnect,
                type: item.Type,
                description: item.Description,
                notification: item.Notification,
                statusUpdatedAt: item.Status_updated_at,
                errorCount: item.ErrorCount
            }))
            setData(devices)
        } catch (err) {
            console.table(err)
        } finally {
            setIsLoading(false)
        }
    }

    const finalColumns: ColumnDef<Device>[] = [
            ...columns,
            {
                id: 'actions',
                header: 'Actions',
                cell: ({ row }) => {
                    const deviceRow = row.original
                    return (
                        <ActionsColumns
                            row={deviceRow}
                            handleFetchData={() => handleFetchData()}
                        />
                    )
                },
                size: 150
            }
        ]


    return (
        <DataTable
            data={data}
            filter={filter}
            defaultFilter={defaultFilter}
            columns={finalColumns}
            setFilter={setFilter}
            sort={sortOptions}
            defaultSort={sortData}
            setSort={setSortOptions}
            isLoading={isLoading}
            setIsLoading={setIsLoading}
            searchProps={searchParameter}
            handleFetchData={handleFetchData}
            totalData={totalData}
            addNewData={<NewDataTable handleFetchData={handleFetchData} />}
        />
    )
}

export default DeviceTableData;