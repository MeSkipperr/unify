import { ColumnDef, ColumnFiltersState, RowSelectionState, SortingState, VisibilityState } from "@tanstack/react-table"
import { FilterConfig } from "../filter/types"

export type TableProps<TData> = {
    data: TData[]
    columns: ColumnDef<TData, any>[]
    filter : FilterConfig[]
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
