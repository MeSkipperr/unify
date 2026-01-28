
'use client'

import { ColumnDef } from '@tanstack/react-table'
import { Checkbox } from '@/components/ui/checkbox'
import { Button } from '@/components/ui/button'
import { Pencil } from 'lucide-react'
import { Device } from './data-table'

export const columns: ColumnDef<Device>[] = [
    {
        id: 'select',
        header: ({ table }) => (
            <Checkbox
                checked={table.getIsAllPageRowsSelected()}
                onCheckedChange={(value) =>
                    table.toggleAllPageRowsSelected(!!value)
                }
                aria-label="Select all"
            />
        ),
        cell: ({ row }) => (
            <Checkbox
                checked={row.getIsSelected()}
                onCheckedChange={(value) =>
                    row.toggleSelected(!!value)
                }
                aria-label="Select row"
            />
        ),
        enableSorting: false,
        enableHiding: false,
    },
    {
        accessorKey: 'name',
        header: 'Name',
    },
    {
        accessorKey: 'ipAddress',
        header: 'IP Address',
    },
    {
        accessorKey: 'roomNumber',
        header: 'Room Number',
    },
    {
        accessorKey: 'macAddress',
        header: 'MAC Address',
    },
    
    {
        id: 'actions',
        header: 'Edit',
        cell: () => (
            <Button variant="ghost" size="icon">
                <Pencil className="h-4 w-4" />
            </Button>
        ),
        enableSorting: false,
        enableHiding: false,
    },
]
