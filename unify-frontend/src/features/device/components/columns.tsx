
'use client'

import { ColumnDef } from '@tanstack/react-table'
import { Button } from '@/components/ui/button'
import { EllipsisVertical } from 'lucide-react'
import { Device } from '../types'
import { DeviceStatus } from '@/components/status'
import { formatDateTime, getCompactRelativeTime } from '@/utils/time'

import {
    Sheet,
    SheetClose,
    SheetContent,
    SheetDescription,
    SheetFooter,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from "@/components/ui/sheet"

import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog"

import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Label } from '@radix-ui/react-label'
import Link from 'next/link'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'


import { NotificationToggle } from './notification-toogle'

export const columns: ColumnDef<Device>[] = [
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
    {
        id: 'actions',
        header: 'Actions',
        cell: ({ row }) => (
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon" >
                        <EllipsisVertical className="h-4 w-4" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-40" align="start">
                    <DropdownMenuLabel className=' font-bold'>Actions</DropdownMenuLabel>
                    <Sheet>
                        <SheetTrigger asChild>
                            <Button variant="ghost" className='w-full flex justify-start gap-0 py-0 px-2'>View Details</Button>
                        </SheetTrigger>
                        <SheetContent className="sm:max-w-md">
                            <SheetHeader>
                                <SheetTitle className="flex items-center gap-4">
                                    Device Details
                                    <DeviceStatus isConnect={row.original.isConnect} />
                                </SheetTitle>

                                <SheetDescription>
                                    Detailed information about the selected network device.
                                </SheetDescription>
                                <SheetDescription className='text-primary'>
                                    Last status update {formatDateTime(row.original.statusUpdatedAt)}
                                </SheetDescription>
                            </SheetHeader>

                            <div className="mt-2 space-y-6 px-4">
                                <div className="grid grid-cols-2 gap-4">
                                    <div className="space-y-1">
                                        <Label className="text-xs text-muted-foreground">
                                            Device Name
                                        </Label>
                                        <p className="text-sm font-medium">
                                            {row.original.name}
                                        </p>
                                    </div>

                                    <div className="space-y-1">
                                        <Label className="text-xs text-muted-foreground">
                                            Device Type
                                        </Label>
                                        <p className="text-sm font-medium">
                                            {row.original.type}
                                        </p>
                                    </div>
                                    <div className="space-y-1">
                                        <Label className="text-xs text-muted-foreground">
                                            IP Address
                                        </Label>
                                        <p className="font-mono text-sm">
                                            {row.original.ipAddress}
                                        </p>
                                    </div>

                                    <div className="space-y-1">
                                        <Label className="text- text-muted-foreground">
                                            MAC Address
                                        </Label>
                                        <p className="font-mono text-sm">
                                            {row.original.macAddress}
                                        </p>
                                    </div>
                                    <div className="space-y-1">
                                        <Label className="text-xs text-muted-foreground">
                                            Room Number
                                        </Label>
                                        <p className="text-sm font-medium">
                                            {row.original.roomNumber}
                                        </p>
                                    </div>
                                    <div className="space-y-1">
                                        <Label className="text-xs text-muted-foreground">
                                            Room Number
                                        </Label>
                                        <p className="text-sm font-medium">
                                            {row.original.roomNumber}
                                        </p>
                                    </div>
                                </div>

                            </div>
                            {row.original.description?.trim() && (
                                <div className="space-y-1 px-4">
                                    <Label className="text-xs text-muted-foreground">
                                        Description
                                    </Label>
                                    <p className="text-sm font-medium leading-relaxed">
                                        {row.original.description}
                                    </p>
                                </div>
                            )}
                        </SheetContent>

                    </Sheet>
                    <Sheet>
                        <SheetTrigger asChild>
                            <Button variant="ghost" className='w-full flex justify-start gap-0 py-0 px-2'>Edit</Button>
                        </SheetTrigger>
                        <SheetContent className="sm:max-w-md">
                            <SheetHeader>
                                <SheetTitle className="flex items-center gap-4">
                                    Device Details
                                    <DeviceStatus isConnect={row.original.isConnect} />
                                </SheetTitle>

                                <SheetDescription>
                                    View and manage configuration details, network identity, and operational status of this device.
                                </SheetDescription>

                                <SheetDescription className="text-primary">
                                    Last status update: {formatDateTime(row.original.statusUpdatedAt)}
                                </SheetDescription>
                            </SheetHeader>


                            <div className="mt-2 space-y-6 px-4">
                                <div className="grid grid-cols-2 gap-4">

                                    <div className="space-y-1">
                                        <Label className="text-xs text-muted-foreground">
                                            Device Name
                                        </Label>
                                        <Input
                                            defaultValue={row.original.name}
                                            className="h-8 text-sm"
                                        />
                                    </div>

                                    <div className="space-y-1">
                                        <Label className="text-xs text-muted-foreground">
                                            Device Type
                                        </Label>
                                        <Input
                                            defaultValue={row.original.type}
                                            className="h-8 text-sm"
                                        />
                                    </div>

                                    <div className="space-y-1">
                                        <Label className="text-xs text-muted-foreground">
                                            IP Address
                                        </Label>
                                        <Input
                                            defaultValue={row.original.ipAddress}
                                            className="h-8 text-sm font-mono"
                                        />
                                    </div>

                                    <div className="space-y-1">
                                        <Label className="text-xs text-muted-foreground">
                                            MAC Address
                                        </Label>
                                        <Input
                                            defaultValue={row.original.macAddress}
                                            className="h-8 text-sm font-mono"
                                        />
                                    </div>

                                    <div className="space-y-1">
                                        <Label className="text-xs text-muted-foreground">
                                            Room Number
                                        </Label>
                                        <Input
                                            defaultValue={row.original.roomNumber}
                                            className="h-8 text-sm"
                                        />
                                    </div>

                                </div>
                            </div>
                            <div className="space-y-1 px-4">
                                <Label className="text-xs text-muted-foreground">
                                    Description
                                </Label>

                                <Textarea
                                    defaultValue={row.original.description || ""}
                                    className="text-sm min-h-[100px]"
                                    placeholder="Enter device description..."
                                />
                            </div>

                            <SheetFooter>
                                <Button type="submit">Save changes</Button>
                                <SheetClose asChild>
                                    <Button variant="outline">Close</Button>
                                </SheetClose>
                            </SheetFooter>
                        </SheetContent>

                    </Sheet>
                    <DropdownMenuItem >
                        <Link href={`/port-forward?listen-ip=${row.original.ipAddress}`} className='size-full'>
                            Port Forward
                        </Link>
                    </DropdownMenuItem>
                    <AlertDialog>
                        <AlertDialogTrigger asChild>
                            <Button variant="ghost" className='text-destructive w-full flex justify-start gap-0 py-0 px-2 hover:text-destructive'>
                                Delete
                            </Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                            <AlertDialogHeader>
                                <AlertDialogTitle>
                                    Delete device {row.original.name}?
                                </AlertDialogTitle>
                                <AlertDialogDescription>
                                    This action cannot be undone. This will permanently remove this device from the system and all related configurations may be lost.
                                </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                                <AlertDialogCancel>Keep Device</AlertDialogCancel>
                                <AlertDialogAction variant="destructive">Delete Device</AlertDialogAction>
                            </AlertDialogFooter>
                        </AlertDialogContent>
                    </AlertDialog>
                </DropdownMenuContent>
            </DropdownMenu>
        ),
        enableSorting: false,
        enableHiding: false,
    },
]
