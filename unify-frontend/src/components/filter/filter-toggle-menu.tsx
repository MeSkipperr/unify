import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Button } from "@/components/ui/button"
import { Plus } from "lucide-react"
import { FilterConfig } from "./types"

type Props = {
    filters: FilterConfig[]
    onToggleFilter: (filterKey: string) => void
}

const FilterToggleMenu = ({ filters, onToggleFilter }: Props) => {
    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                    <Plus className="mr-1 h-4 w-4" />
                    Filter
                </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent>
                <DropdownMenuGroup>
                    {filters.map((filter) => (
                        <DropdownMenuCheckboxItem
                            key={`enable-${filter.key}`}
                            checked={filter.isEnabled}
                            onCheckedChange={() => onToggleFilter(filter.key)}
                        >
                            {filter.label}
                        </DropdownMenuCheckboxItem>
                    ))}
                </DropdownMenuGroup>
            </DropdownMenuContent>
        </DropdownMenu>
    )
}

export default FilterToggleMenu