
'use client'

import { ColumnDef } from '@tanstack/react-table'
import { PortForwardResult } from '../types'
import { formatDateTime } from '@/utils/time'


import { Label } from '@radix-ui/react-label'


export const columns: ColumnDef<PortForwardResult>[] =
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
            accessorKey: 'listen',
            header: 'Listen',
            cell: (({ row }) => (
                <div className="flex justify-between w-full">
                    <span>{row.original.listenIp}</span>
                    <span className='w-2/4'>:{row.original.listenPort}</span>
                </div>
            ))
        },
        {
            accessorKey: 'dest',
            header: 'Destination',
            cell: ({ row }) => (
                <div className="flex justify-between w-full">
                    <span>{row.original.destIp}</span>
                    <span className='w-2/4'>:{row.original.destPort}</span>
                </div>
            )
        },
        {
            accessorKey: 'protocol',
            header: 'Protocol',
        },
        {
            accessorKey: 'startTime',
            header: 'Start At',
            cell: ({ row }) => (
                <span>{formatDateTime(row.original.createdAt)}</span>
            )
        },
        {
            accessorKey: 'finishTime',
            header: 'Finnish At',
            cell: ({ row }) => (
                <span>{formatDateTime(row.original.expiresAt)}</span>
            )
        },
    ]
