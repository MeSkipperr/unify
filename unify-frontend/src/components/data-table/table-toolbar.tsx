"use client"

import { ChevronDown } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuContent,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { DataTableToolbarProps } from "./types"

const DataTableToolbar = <TData,>({
    table,
    searchColumn,
    searchPlaceholder = "Search...",
    rightSlot,
}: DataTableToolbarProps<TData>) => {
    return (
        <div className="flex items-end justify-between gap-4">
            <div className="flex gap-4">
                {/* SEARCH */}
                {searchColumn && (
                    <Input
                        className="w-64"
                        placeholder={searchPlaceholder}
                        value={
                            (table
                                .getColumn(searchColumn)
                                ?.getFilterValue() as string) ?? ""
                        }
                        onChange={(e) =>
                            table
                                .getColumn(searchColumn)
                                ?.setFilterValue(e.target.value)
                        }
                    />
                )}

                {/* COLUMN TOGGLE */}
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="outline">
                            Columns
                            <ChevronDown className="ml-2 h-4 w-4" />
                        </Button>
                    </DropdownMenuTrigger>

                    <DropdownMenuContent align="start">
                        {table
                            .getAllColumns()
                            .filter((col) => col.getCanHide())
                            .map((col) => (
                                <DropdownMenuCheckboxItem
                                    key={col.id}
                                    checked={col.getIsVisible()}
                                    onCheckedChange={(checked) =>
                                        col.toggleVisibility(!!checked)
                                    }
                                    
                                    onSelect={(e) => e.preventDefault()}
                                >
                                    {col.columnDef.meta?.label ?? col.id}
                                </DropdownMenuCheckboxItem>
                            ))}
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>

            {rightSlot}
        </div>
    )
}

export default DataTableToolbar
