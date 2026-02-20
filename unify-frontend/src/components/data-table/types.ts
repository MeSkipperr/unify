import { Table } from "@tanstack/react-table"

export type DataTableProps<TData> = {
    table: Table<TData>
    emptyMessage?: string
}

export type DataTableToolbarProps<TData> = {
    table: Table<TData>
    searchColumn?: string
    searchPlaceholder?: string
    rightSlot?: React.ReactNode
}
