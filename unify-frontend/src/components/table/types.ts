import { ColumnDef, ColumnFiltersState, RowSelectionState, SortingState, VisibilityState } from "@tanstack/react-table"
import { FilterConfig } from "../filter/types"
import { SortBy } from "../sort/types"
import { SearchBarProps } from "./search"
import React from "react"

export type TableQuery = {
    search?: string
    page?: number
    pageSize?: number
    date?: Date
}

export type TableProps<TData> = {
    data: TData[]
    columns: ColumnDef<TData, any>[]
    filter: FilterConfig[]
    defaultFilter?: FilterConfig[]
    setFilter: React.Dispatch<React.SetStateAction<FilterConfig[]>>
    sort: SortBy[]
    setSort: React.Dispatch<React.SetStateAction<SortBy[]>>
    isLoading: boolean
    setIsLoading: React.Dispatch<React.SetStateAction<boolean>>
    searchProps: Omit<SearchBarProps, "value" | "onChange">
    useObserver?: boolean
    handleFetchData: (payload: TableQuery) => Promise<void>
    totalData: number
    addNewData?: React.ReactNode
}

export type DataTableStateProps = {

    sorting: SortingState
    setSorting: React.Dispatch<React.SetStateAction<SortingState>>

    columnFilters: ColumnFiltersState
    setColumnFilters: React.Dispatch<React.SetStateAction<ColumnFiltersState>>

    columnVisibility: VisibilityState
    setColumnVisibility: React.Dispatch<React.SetStateAction<VisibilityState>>

    rowSelection: RowSelectionState
    setRowSelection: React.Dispatch<React.SetStateAction<RowSelectionState>>
}
