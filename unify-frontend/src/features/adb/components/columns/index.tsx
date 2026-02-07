
'use client'

import { ColumnDef } from '@tanstack/react-table'
import { AdbResult } from '../../types'
import { formatDateTime } from '@/utils/time'


import { Label } from '@radix-ui/react-label'


export const columns: ColumnDef<AdbResult>[] =
    [
        {
            id: "number",
            header: () => (
                <Label className='w-full flex justify-center items-center'>No</Label>
            ),
            cell: ({ row }) => (
                <Label className='w-full flex justify-center items-center'>{row.original.index}</Label>
            ),
            size: 50,
        },
        {
            accessorKey: 'status',
            header: 'Status',
        },
        {
            accessorKey: 'ipAddress',
            header: 'IP Address',
        },
        {
            accessorKey: 'port',
            header: 'Port',
        },
        {
            accessorKey: 'deviceName',
            header: 'Device Name',
        },

        {
            accessorKey: 'serviceType',
            header: 'Service Type',
        },
        {
            accessorKey: 'startTime',
            header: 'Start At',
            cell: ({ row }) => (
                <span>{formatDateTime(row.original.startTime)}</span>
            )
        },
        {
            accessorKey: 'finishTime',
            header: 'Finnish At',
            cell: ({ row }) => (
                <span>{formatDateTime(row.original.finishTime)}</span>
            )
        },
    ]
