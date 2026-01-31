"use client"

import { FilterConfig } from "@/components/filter/types"
import DeviceTable from "@/features/device/table"
import { Device } from "@/features/device/types"



const data: Device[] = [
    {
        id: "1",
        name: "Access Point Lobby",
        ipAddress: "192.168.1.10",
        macAddress: "00:00:00:00:00:00",
        roomNumber: "1001",
        isConnect: true,
        type: "access-point",
        statusUpdatedAt: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000),
        notification: false,
        description: "Server: NVR4,Channel: 29,Coverage: C-108 Room 3208 - 3201,Distribution: In Front of Room 3207 & 3208 pointed to Elevator LP-03"
    },
    {
        id: "2",
        name: "CCTV Gate",
        ipAddress: "192.168.1.20",
        macAddress: "00:00:00:00:00:00",
        roomNumber: "1002",
        isConnect: false,
        type: "cctv",
        notification: true,
        statusUpdatedAt: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000),
    },
]

const dataFilter: FilterConfig[] = [
    {
        key: "status",
        label: "Status",
        type: "boolean",
        isEnabled: true,
        options: [
            { value: false, label: "DOWN", isSelected: false },
            { value: true, label: "UP", isSelected: false },
        ],
    },
    {
        key: "notification",
        label: "Notification",
        type: "boolean",
        isEnabled: true,
        options: [
            { value: false, label: "Off", isSelected: true },
            { value: true, label: "On", isSelected: false },
        ],
    },
    {
        key: "type",
        label: "Types",
        type: "select",
        isEnabled: true,
        options: [
            { value: "cctv", label: "CCTV", isSelected: false },
            { value: "iptv", label: "IPTV", isSelected: true },
            { value: "access-point", label: "Access Point", isSelected: false },
            { value: "sw", label: "Switch", isSelected: true },
        ],
    },
]

export default function Table() {

    return (
        <div className="w-full space-y-4">

            <DeviceTable data={data} />
        </div>
    )
}
