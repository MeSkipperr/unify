"use client"

import React from "react"
import {
    getCoreRowModel,
    getFilteredRowModel,
    getPaginationRowModel,
    useReactTable,
} from "@tanstack/react-table"
import { Device } from "./types"
import { deviceColumns } from "./columns"
import DataTable from "@/components/data-table/table"
import DataTableToolbar from "@/components/data-table/table-toolbar"


const DeviceTable = ({ data }: { data: Device[] }) => {
    const table = useReactTable<Device>({
        data,
        columns: deviceColumns,
        getCoreRowModel: getCoreRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
    })

    return (
        <div className="space-y-4" >
            <DataTableToolbar
                table={table}
                searchColumn="name"
                searchPlaceholder="Search device..."
            />

            <DataTable table={table} />
        </div>
    )
}


export default DeviceTable
