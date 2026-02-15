import { FilterConfig } from "@/components/filter/types";
import { SortBy } from "@/components/sort/types";

export const dataFilter: FilterConfig[] = [
    {
        key: "serviceType",
        label: "Service Type",
        type: "select",
        isEnabled: true,
        options: [
            { value: "remove-youtube-data", label: "Remove Youtube Data", isSelected: false },
            { value: "get-uptime-adb", label: "Get Uptime", isSelected: false },
            { value: "manual", label: "Manual", isSelected: false },
        ],
    },
]

export const sortData: SortBy[] = [
    { key: "startTime", label: "Start At", value: "none" },
    { key: "finishTime", label: "Finish At", value: "descending" },
]