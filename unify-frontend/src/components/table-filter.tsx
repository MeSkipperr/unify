'use client'

import * as React from 'react'
import { Button } from '@/components/ui/button'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { Filter } from 'lucide-react'

type FilterType = 'status' | null

export function TableFilter({
    onChange,
}: {
    onChange: (filter: { type: string; value: string }) => void
}) {
    const [filterType, setFilterType] = React.useState<FilterType>(null)
    const [filterValue, setFilterValue] = React.useState<string>('')

    React.useEffect(() => {
        if (filterType && filterValue) {
            onChange({ type: filterType, value: filterValue })
        }
    }, [filterType, filterValue, onChange])

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                    <Filter className="mr-2 size-4" />
                    Filter
                </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent className="w-64 p-4 space-y-3">
                {/* FILTER BY */}
                <div className="space-y-1">
                    <p className="text-xs text-muted-foreground">Filter by</p>
                    <Select onValueChange={(v) => setFilterType(v as FilterType)}>
                        <SelectTrigger>
                            <SelectValue placeholder="Select field" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="status">Status</SelectItem>
                        </SelectContent>
                    </Select>
                </div>

                {/* FILTER VALUE */}
                {filterType === 'status' && (
                    <div className="space-y-1">
                        <p className="text-xs text-muted-foreground">Status</p>
                        <Select onValueChange={setFilterValue}>
                            <SelectTrigger>
                                <SelectValue placeholder="Select status" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="online">Online</SelectItem>
                                <SelectItem value="offline">Offline</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>
                )}
            </DropdownMenuContent>
        </DropdownMenu>
    )
}
