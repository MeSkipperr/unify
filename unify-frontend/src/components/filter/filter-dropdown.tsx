import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Button } from "@/components/ui/button"
import { ChevronDown } from "lucide-react"
import { FilterConfig, OptionValue } from "./types"
import { getSelectedOptionSummary } from "./utils"

type Props = {
    filter: FilterConfig
    onToggleOption: (filterKey: string, optionValue: OptionValue) => void
}

const FilterDropdown = ({ filter, onToggleOption }: Props) => {
    const summary = getSelectedOptionSummary(filter.options)

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button
                    size="sm"
                    variant={summary.isEmpty ? "outline" : "default"}
                >
                    {filter.label}
                    {!summary.isEmpty && `: ${summary.text}`}
                    <ChevronDown className="ml-1 h-4 w-4" />
                </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent>
                <DropdownMenuGroup>
                    {filter.options.map((option) => (
                        <DropdownMenuCheckboxItem
                            key={`${filter.key}-${option.value}`}
                            checked={option.isSelected}
                            onCheckedChange={() =>
                                onToggleOption(filter.key, option.value)
                            }
                        >
                            {option.label}
                        </DropdownMenuCheckboxItem>
                    ))}
                </DropdownMenuGroup>
            </DropdownMenuContent>
        </DropdownMenu>
    )
}

export default FilterDropdown
