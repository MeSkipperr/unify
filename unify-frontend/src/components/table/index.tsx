"use client"

import React from "react";
import { TableProps } from "./types";
import { ColumnFiltersState, flexRender, getCoreRowModel, getFilteredRowModel, getPaginationRowModel, SortingState, useReactTable, VisibilityState } from "@tanstack/react-table";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuContent,
    DropdownMenuTrigger,
} from "../ui/dropdown-menu"
import { Button } from "../ui/button"
import { ChevronDown, TriangleAlert } from "lucide-react"
import FilterGroup from "../filter/filter-group";
import { SearchBar } from "./search";
import SortGroup from "../sort";
import TableRowSkeleton from "./skeleton";
import { Empty, EmptyContent, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from "../ui/empty";

const DataTable = <TData,>({
    data,
    columns,
    filter,
    setFilter,
    sort,
    setSort,
    isLoading,
    search,
    handleFetchData
}: TableProps<TData>) => {
    const [sorting, setSorting] = React.useState<SortingState>([])
    const [columnFilters, setColumnFilters] =
        React.useState<ColumnFiltersState>([])
    const [columnVisibility, setColumnVisibility] =
        React.useState<VisibilityState>({})
    const [rowSelection, setRowSelection] = React.useState({})

    const table = useReactTable({
        data,
        columns,
        state: {
            sorting,
            columnFilters,
            columnVisibility,
            rowSelection,
        },
        onSortingChange: setSorting,
        onColumnFiltersChange: setColumnFilters,
        onColumnVisibilityChange: setColumnVisibility,
        onRowSelectionChange: setRowSelection,
        getCoreRowModel: getCoreRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
    })
    return (
        <div className="w-full space-y-4">
            <div className="flex justify-end items-center gap-4">
                <SearchBar
                    id={search.id}
                    label={search.label}
                    value={search.value}
                    onChange={search.onChange}
                    placeholder={search.placeholder}
                    description={search.description}
                />
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="outline">
                            Columns <ChevronDown />
                        </Button>
                    </DropdownMenuTrigger>

                    <DropdownMenuContent align="end">
                        {table
                            .getAllColumns()
                            .filter((column) => column.getCanHide())
                            .map((column) => (
                                <DropdownMenuCheckboxItem
                                    key={column.id}
                                    checked={column.getIsVisible()}
                                    onCheckedChange={(value) =>
                                        column.toggleVisibility(!!value)
                                    }
                                    className="capitalize"
                                >
                                    {column.columnDef.meta?.label ?? column.id}
                                </DropdownMenuCheckboxItem>
                            ))}
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
            <div className="flex gap-4 items-center ">
                <SortGroup sortOptions={sort} onChange={setSort} />
                <FilterGroup data={filter} onChange={setFilter} />
            </div>
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        {table.getHeaderGroups().map(headerGroup => (
                            <TableRow key={headerGroup.id}>
                                {headerGroup.headers.map(header => (
                                    <TableHead key={header.id}>
                                        {flexRender(
                                            header.column.columnDef.header,
                                            header.getContext()
                                        )}
                                    </TableHead>
                                ))}
                            </TableRow>
                        ))}
                    </TableHeader>

                    <TableBody>
                        {isLoading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRowSkeleton
                                    key={i}
                                    columns={columns.length}
                                />
                            ))
                        ) : table.getRowModel().rows.length ? (
                            table.getRowModel().rows.map(row => (
                                <TableRow key={row.id}>
                                    {row.getVisibleCells().map(cell => (
                                        <TableCell key={cell.id}>
                                            {flexRender(
                                                cell.column.columnDef.cell,
                                                cell.getContext()
                                            )}
                                        </TableCell>
                                    ))}
                                </TableRow>
                            ))
                        ) : (
                            <TableRow>
                                <TableCell colSpan={columns.length} className="h-32 text-center">
                                    <Empty>
                                        <EmptyHeader className="py-0 my-0 gap-0">
                                            <EmptyMedia variant="icon">
                                                <TriangleAlert />
                                            </EmptyMedia>
                                            <EmptyTitle >Data dot found</EmptyTitle>
                                            <EmptyDescription>Try adjusting your search or filters.</EmptyDescription>
                                        </EmptyHeader>
                                        {handleFetchData && (
                                            <EmptyContent>
                                                <Button onClick={() => handleFetchData()}>
                                                    Reload Data
                                                </Button>
                                            </EmptyContent>
                                        )}
                                    </Empty>
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>

            </div>
        </div>
    )
}

export default DataTable
