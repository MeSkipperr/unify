
'use client'

import { ColumnDef } from '@tanstack/react-table'
import { Device } from '../../types'
import { DeviceStatus } from '@/components/status'

import { Label } from '@radix-ui/react-label'

import { NotificationToggle } from '../notification-toogle'
import { getCompactRelativeTime } from '@/utils/time'


export const columns: ColumnDef<Device>[] =
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
            accessorKey: 'name',
            header: 'Name',
        },
        {
            accessorKey: "isConnect",
            header: "Status",
            cell: ({ row }) => (
                <DeviceStatus isConnect={row.original.isConnect} />
            ),
            meta: { label: "Status" },
        },

        {
            accessorKey: 'ipAddress',
            header: 'IP Address',
            meta: { label: "Ip Address" },
        },
        {
            accessorKey: 'roomNumber',
            header: 'Room Number',
            meta: { label: "Room Number" },
        },
        {
            accessorKey: 'macAddress',
            header: 'MAC Address',
            meta: { label: "Mac Address" },
        },
        {
            accessorKey: 'type',
            header: 'Type',
        },
        {
            accessorKey: 'deviceProduct',
            header: 'Product',
        },
        {
            accessorKey: 'statusUpdatedAt',
            header: 'Last Update',
            cell: ({ row }) => (
                <span>{getCompactRelativeTime(row.original.statusUpdatedAt)}</span>
            )
        },
        {
            accessorKey: 'notification',
            header: 'Notification',
            cell: ({ row }) => (
                <NotificationToggle
                    deviceId={row.original.id}
                    initialNotification={row.original.notification}
                />
            )
        },
    ]
