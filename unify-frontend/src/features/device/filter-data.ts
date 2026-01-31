import { FilterConfig } from "@/components/filter/types";

export const dataFilter: FilterConfig[] = [
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