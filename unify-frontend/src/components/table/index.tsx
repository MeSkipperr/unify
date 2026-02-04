"use client"

import React from "react";
import { TableProps, TableQuery } from "./types";
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
import { ChevronDown, Plus, RotateCw, TriangleAlert } from "lucide-react"
import FilterGroup from "../filter/filter-group";
import { SearchBar } from "./search";
import SortGroup from "../sort";
import { Empty, EmptyContent, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from "../ui/empty";
import PagenationTable from "./pagenation";
import { useRouter, useSearchParams } from "next/navigation";
import { Label } from "@radix-ui/react-label";
import TableRowSkeleton from "./skeleton";

const DataTable = <TData,>({
    data,
    columns,
    filter,
    defaultFilter,
    setFilter,
    sort,
    setSort,
    isLoading,
    setIsLoading,
    searchProps,
    useObserver = false,
    handleFetchData,
    totalData,
    addNewData
}: TableProps<TData>) => {
    const router = useRouter();
    const searchParams = useSearchParams()
    const [sorting, setSorting] = React.useState<SortingState>([])
    const [columnFilters, setColumnFilters] =
        React.useState<ColumnFiltersState>([])
    const [columnVisibility, setColumnVisibility] =
        React.useState<VisibilityState>({})
    const [rowSelection, setRowSelection] = React.useState({})

    const targetRef = React.useRef<HTMLDivElement | null>(null);
    const [search, setSearch] = React.useState<string>(searchParams.get("search") || "");
    const [pageSizeQuery, setPageSizeQuery] = React.useState<number>(25);
    const [pageQuery, setPageQuery] = React.useState<number>(1);

    React.useEffect(() => {
        if (!targetRef.current) return;

        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting && (!isLoading || table.getRowModel().rows.length > 0)) {
                        setPageQuery?.(prev => {
                            const newPage = prev + 1;
                            return newPage;
                        });
                    }
                });
            },
            { root: null, rootMargin: "0px", threshold: 0 }
        );

        observer.observe(targetRef.current);

        return () => {
            if (targetRef.current) observer.unobserve(targetRef.current);
        };
    }, [setPageQuery]);


    const fetchData = React.useCallback(async () => {
        const payload: TableQuery = {
            page: pageQuery,
            pageSize: pageSizeQuery,
            search,
        }
        const params = new URLSearchParams(searchParams.toString())

        if (search.trim() !== "") {
            params.set("search", search.toString())
        } else {
            params.delete("search") // hapus key jika kosong
        }

        router.push(`${window.location.pathname}?${params.toString()}`)

        handleFetchData(payload)
    }, [filter, sort, search, pageQuery, pageSizeQuery])

    React.useEffect(() => {
        setPageQuery(1)
    }, [filter, sort, search, pageSizeQuery])

    React.useEffect(() => {
        setIsLoading(true)
        const timer = setTimeout(() => {
            fetchData()
        }, 800)

        return () => clearTimeout(timer)
    }, [fetchData])

    const handleResetFilter = () => {
        setFilter(defaultFilter ?? [])
        setSort([])
        setSearch("")
    }

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
        // getPaginationRowModel: getPaginationRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
    })
    return (
        <div className="w-full h-full flex flex-col gap-4">

            <div className="flex justify-end items-center gap-4 shrink-0">
                <SearchBar
                    id={searchProps.id}
                    label={searchProps.label}
                    value={search}
                    onChange={setSearch}
                    placeholder={searchProps.placeholder}
                    description={searchProps.description}
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
                <Button variant="destructive" onClick={() => handleResetFilter()}>
                    <RotateCw /> Reset Filter
                </Button>
                {addNewData}
            </div>
            <div className="flex gap-4 items-center shrink-0 justify-between">
                <div className="flex gap-4 items-center">
                    <SortGroup sortOptions={sort} onChange={setSort} />
                    <FilterGroup data={filter} onChange={setFilter} />
                </div>
                {!isLoading &&
                    <Label>Total Items: {totalData}</Label>
                }

            </div>
            <div className="flex-1 min-h-0">
                <div className="h-full rounded-md border overflow-auto">
                    <Table>
                        <TableHeader >
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
                                Array.from({ length: 10 }).map((_, i) => (
                                    <TableRowSkeleton key={i} columns={columns.length} />
                                ))
                            ) : table.getRowModel().rows.length > 0 ? (
                                table.getRowModel().rows.map(row => (
                                    <TableRow key={row.id}>
                                        {row.getVisibleCells().map(cell => (
                                            <TableCell key={cell.id}>
                                                {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                            </TableCell>
                                        ))}
                                    </TableRow>
                                ))
                            ) :
                                !isLoading &&
                                <TableRow>
                                    <TableCell colSpan={columns.length} className="h-32 text-center">
                                        <Empty className="">
                                            <EmptyHeader className="py-0 my-0 gap-0">
                                                <EmptyMedia variant="icon">
                                                    <TriangleAlert />
                                                </EmptyMedia>
                                                <EmptyTitle >Data dot found</EmptyTitle>
                                                <EmptyDescription>Try adjusting your search or filters.</EmptyDescription>
                                            </EmptyHeader>
                                            <EmptyContent>
                                                <Button onClick={async () => fetchData()}>
                                                    Reload Data
                                                </Button>
                                            </EmptyContent>
                                        </Empty>
                                    </TableCell>
                                </TableRow>
                            }

                        </TableBody>
                    </Table>
                </div>
            </div>

            {
                useObserver ?
                    <div
                        ref={targetRef}
                        style={{
                            height: "1",
                            marginTop: "1",
                            backgroundColor: "transparent",
                        }}
                    />
                    :
                    <PagenationTable
                        pageQuery={pageQuery}
                        setPageQuery={setPageQuery}
                        pageSizeQuery={pageSizeQuery}
                        setPageSizeQuery={setPageSizeQuery}
                        totalData={totalData}
                    />
            }
        </div>
    )
}

export default DataTable
