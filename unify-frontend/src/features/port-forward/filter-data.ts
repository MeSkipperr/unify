import { FilterConfig } from "@/components/filter/types";
import { SortBy } from "@/components/sort/types";

export const dataFilter: FilterConfig[] = [
    {
        key: "protocol",
        label: "Protocol",
        type: "select",
        isEnabled: false,
        options: [
            { value: "tcp", label: "TCP", isSelected: false },
            { value: "udp", label: "UDP", isSelected: false },
        ],
    },
    {
        key: "status",
        label: "Status",
        type: "select",
        isEnabled: true,
        options: [
            { value: "active", label: "Active", isSelected: true },
            { value: "inactive", label: "Inactive", isSelected: false },
            { value: "expired", label: "Expried", isSelected: false },
            { value: "pending", label: "Pending", isSelected: true },
        ],
    },
]

export const sortData: SortBy[] = [
    { key: "startTime", label: "Start At", value: "descending" },
    { key: "finishTime", label: "Finish At", value: "none" },
]