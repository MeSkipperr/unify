"use client"

import { ColumnDef } from "@tanstack/react-table"
import { Device } from "./types"
import { Checkbox } from "@/components/ui/checkbox"
import { Button } from "@/components/ui/button"
import { Bell, BellOff } from "lucide-react"
import StatusBadge from "@/components/status-badge"

export const deviceColumns: ColumnDef<Device>[] = [
    {
        id: "select",
        header: ({ table }) => (
            <Checkbox
                checked={table.getIsAllPageRowsSelected()}
                onCheckedChange={(v) =>
                    table.toggleAllPageRowsSelected(!!v)
                }
            />
        ),
        cell: ({ row }) => (
            <Checkbox
                checked={row.getIsSelected()}
                onCheckedChange={(v) => row.toggleSelected(!!v)}
            />
        ),
        enableSorting: false,
        enableHiding: false,
        meta: { label: "Select" },
    },
    {
        accessorKey: "isConnect",
        header: "Status",
        cell: ({ row }) => (
            <StatusBadge isConnect={row.original.isConnect} />
        ),
        meta: { label: "Status" },
    },
    {
        accessorKey: "notification",
        header: "Notification",
        cell: ({ row }) => (
            <Button variant="ghost" size="icon">
                {row.original.notification ? (
                    <Bell className="size-4" />
                ) : (
                    <BellOff className="size-4" />
                )}
            </Button>
        ),
        meta: { label: "Notification" },
    },
    {
        accessorKey: "name",
        header: "Name",
        meta: { label: "Device Name" },
    },
    {
        accessorKey: "ipAddress",
        header: "IP Address",
        meta: { label: "IP Address" },
    },
    {
        accessorKey: "roomNumber",
        header: "Room",
        meta: { label: "Room Number" },
    },
    {
        accessorKey: "type",
        header: "Type",
        meta: { label: "Device Type" },
    },
]
