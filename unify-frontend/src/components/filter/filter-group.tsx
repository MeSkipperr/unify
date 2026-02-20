"use client"
import { FilterConfig, OptionValue } from "./types"
import FilterDropdown from "./filter-dropdown"
import FilterToggleMenu from "./filter-toggle-menu"

type FilterGroupProps = {
    data: FilterConfig[]
    onChange?: (filters: FilterConfig[]) => void
}

const FilterGroup = ({ data, onChange }: FilterGroupProps) => {

    const updateFilters = (
        updater: (prev: FilterConfig[]) => FilterConfig[]
    ) => {
        if (!onChange) return
        onChange(updater(data))
    }

    const toggleOption = (filterKey: string, optionValue: OptionValue) => {
        updateFilters((prev) =>
            prev.map((filter) =>
                filter.key !== filterKey
                    ? filter
                    : {
                        ...filter,
                        options: filter.options.map((opt) =>
                            opt.value === optionValue
                                ? { ...opt, isSelected: !opt.isSelected }
                                : opt
                        ),
                    }
            )
        )
    }

    const toggleFilter = (filterKey: string) => {
        updateFilters((prev) =>
            prev.map((filter) =>
                filter.key === filterKey
                    ? { ...filter, isEnabled: !filter.isEnabled }
                    : filter
            )
        )
    }

    return (
        <div className="flex gap-4">
            {data
                .filter((f) => f.isEnabled)
                .map((filter) => (
                    <FilterDropdown
                        key={filter.key}
                        filter={filter}
                        onToggleOption={toggleOption}
                    />
                ))}

            <FilterToggleMenu
                filters={data}
                onToggleFilter={toggleFilter}
            />
        </div>
    )
}

export default FilterGroup
