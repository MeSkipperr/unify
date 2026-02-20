import { FilterConfig } from "@/components/filter/types";
import { SortBy } from "@/components/sort/types";

export const dataFilter: FilterConfig[] = [
    {
        key: "status",
        label: "Status",
        type: "boolean",
        isEnabled: false,
        options: [
            { value: false, label: "DOWN", isSelected: false },
            { value: true, label: "UP", isSelected: false },
        ],
    },
    {
        key: "notification",
        label: "Notification",
        type: "boolean",
        isEnabled: false,
        options: [
            { value: false, label: "Off", isSelected: false },
            { value: true, label: "On", isSelected: false },
        ],
    },
    {
        key: "type",
        label: "Types",
        type: "select",
        isEnabled: false,
        options: [
            { value: "cctv", label: "CCTV", isSelected: false },
            { value: "iptv", label: "IPTV", isSelected: false },
            { value: "access-point", label: "Access Point", isSelected: false },
            { value: "sw", label: "Switch", isSelected: false },
        ],
    },
]

export const sortData: SortBy[] = [
    { key: "roomNumber", label: "Room Number", value: "ascending" },
    { key: "lastUpdate", label: "Last Update", value: "none" },
    { key: "created_at", label: "Created At", value: "ascending" },
]