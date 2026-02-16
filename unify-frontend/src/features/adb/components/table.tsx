"use client"
import DataTable from "@/components/table";
import { AdbResult } from "../types";
import { columns } from "./columns";
import { dataFilter, sortData } from "../filter-data";
import { useEffect, useState } from "react";
import { FilterConfig } from "@/components/filter/types";
import { SortBy } from "@/components/sort/types";
import { TableQuery } from "@/components/table/types";
import { DatePicker } from "@/components/date-picker";
import { updateFilterOption } from "@/features/device/utils/select-options";
import { useSearchParams } from "next/navigation";
import { AdbQuery, getAdbResults } from "../api/adb-result.api";
import NewDataTable from "./new-data";

const searchParameter = {
    id: "adb-search-bar",
    label: "Search ADB Results",
    description: "Search devices by IP address, device name, or connection status.",
    placeholder: "e.g. 192.168.1.10 or device-name"
}
type AdbTableProps = {
    serviceType: string
    hasDefaultValue?: boolean
    addNewData?: boolean
}

const AdbTable = ({ serviceType, hasDefaultValue = false, addNewData = false }: AdbTableProps) => {
    const searchParams = useSearchParams()

    const defaultFilter = updateFilterOption(dataFilter, "serviceType", serviceType)

    const [data, setData] = useState<AdbResult[]>([]);
    const [filter, setFilter] = useState<FilterConfig[]>(defaultFilter);
    const [sort, setSort] = useState<SortBy[]>(sortData);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [date, setDate] = useState<Date | undefined>(
        hasDefaultValue ? new Date() : undefined
    );


    const [totalData, setTotalData] = useState<number>(1);

    const handleFetchData = async (payload?: TableQuery) => {
        setIsLoading(false)

        const page = payload?.page ?? Number(searchParams.get("page")) ?? 1
        const pageSize = payload?.pageSize ?? Number(searchParams.get("pageSize")) ?? 10
        const search = payload?.search ?? searchParams.get("search") ?? ""

        const dataPayload: AdbQuery = {
            page,
            pageSize,
            search,
            date
        }


        filter.forEach((filter) => {
            const selectedValues = filter.options
                .filter(opt => opt.isSelected)
                .map(opt => opt.value)

            if (selectedValues.length > 0) {
                switch (filter.key) {
                    case "serviceType":
                        dataPayload.typeServices = selectedValues as string[]
                        break
                }
            }
        })

        dataPayload.sort = sort
            .filter(s => s.value !== "none")
            .map(s => `${s.key}:${s.value}`)

        try {
            const result = await getAdbResults(dataPayload)
            setTotalData(result.total)
            const adbRes: AdbResult[] = result.data.map((item: any) => ({
                id: item.ID,
                index: item.index,
                status: item.Status,
                ipAddress: item.IPAddress,
                finishTime: item.FinishTime,
                startTime: item.StartTime,
                port: item.Port,
                deviceName: item.NameDevice,
                result: item.Result,
                serviceType: item.TypeServices
            }))
            setData(adbRes)
            console.log(adbRes)
        } catch (err) {
            console.log(err)
        } finally {
            setIsLoading(false)
        }
    }

    useEffect(() => {
        if (!date) return

        handleFetchData({
            date: new Date(date),
        })
    }, [date])


    return (
        <DataTable
            data={data}
            columns={columns}
            filter={filter}
            setFilter={setFilter}
            sort={sort}
            setSort={setSort}
            isLoading={isLoading}
            setIsLoading={setIsLoading}
            searchProps={searchParameter}
            handleFetchData={handleFetchData}
            totalData={totalData}
            addNewData={
                <>
                    <DatePicker
                        value={date}
                        onChange={setDate}
                        maxDate={new Date()}
                    />
                    {addNewData &&
                        <NewDataTable handleFetchData={handleFetchData} />
                    }
                </>

            }
        />
    );
}

export default AdbTable;