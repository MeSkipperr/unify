"use client"
import { columns } from "./columns";
import { dataFilter, sortData } from "../filter-data";
import { Device } from "../types";
import DataTable from "@/components/table";
import React from "react";
import { SortBy } from "@/components/sort/types";
import { DeviceQuery, getDevices } from "../api/device.api";
import { FilterConfig } from "@/components/filter/types";
import { SearchBarProps } from "@/components/table/search";
import { useRouter, useSearchParams } from "next/navigation"

const DeviceTableData = () => {
    const router = useRouter()
    const searchParams = useSearchParams()

    const [data, setData] = React.useState<Device[]>([]);
    const [sortOptions, setSortOptions] = React.useState<SortBy[]>(sortData)
    const [filter, setFilter] = React.useState<FilterConfig[]>(dataFilter);
    const [isLoading, setIsLoading] = React.useState<boolean>(true);
    const [search, setSearch] = React.useState<string>(searchParams.get("search") || "");

    const searchParameter: SearchBarProps = {
        id: "device-search-bar",
        value: search,
        onChange: setSearch,
        description: "Search by name, MAC, IP, room, device type, or description.",
        label: "Search Device",
        placeholder: "DPSCY-..."
    }

    const updateParam = (key: string, value?: string) => {
        const params = new URLSearchParams(searchParams.toString())

        if (!value || value === "") {
            params.delete(key)
        } else {
            params.set(key, value)
        }
~
        router.replace(`?${params.toString()}`, { scroll: false })
    }

    const handleFetchData = React.useCallback(async () => {
        setIsLoading(true)

        const payload: DeviceQuery = {}

        filter.forEach((filter) => {
            const selectedValues = filter.options
                .filter(opt => opt.isSelected)
                .map(opt => opt.value)

            if (selectedValues.length > 0) {
                switch (filter.key) {
                    case "status":
                        payload.status = selectedValues as boolean[]
                        break
                    case "notification":
                        payload.notification = selectedValues as boolean[]
                        break
                    case "type":
                        payload.type = selectedValues as string[]
                        break
                }
            }
        })

        payload.sort = sortOptions
            .filter(s => s.value !== "none")
            .map(s => `${s.key}:${s.value}`)

        payload.search = search
        updateParam("search", search)

        try {
            const result = await getDevices(payload)
            setData(result)
        } catch (err) {
            console.table(err)
        } finally {
            setIsLoading(false)
        }
    }, [filter, sortOptions, search])



    React.useEffect(() => {
        const timer = setTimeout(() => {
            handleFetchData()
        }, 800)

        return () => clearTimeout(timer)
    }, [handleFetchData])

    return (
        <DataTable
            data={data}
            filter={filter}
            columns={columns}
            setFilter={setFilter}
            sort={sortOptions}
            setSort={setSortOptions}
            isLoading={isLoading}
            search={searchParameter}
            handleFetchData={handleFetchData}
        />
    )
}

export default DeviceTableData;