import { FilterConfig } from "@/components/filter/types"

export function updateFilterOption(
    filters: FilterConfig[],      
    key: string,                 
    valueToSelect: string | boolean 
): FilterConfig[] {
    return filters.map(filter => {
        if (filter.key !== key) return filter

        return {
            ...filter,
            options: filter.options.map(option => ({
                ...option,
                isSelected: option.value === valueToSelect
            })),
            isEnabled: filter.options.some(opt => opt.value === valueToSelect) 
        }
    })
}
