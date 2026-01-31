"use client"
import * as React from "react"
import { FilterConfig, OptionValue } from "./types"
import FilterDropdown from "./filter-dropdown"
import FilterToggleMenu from "./filter-toggle-menu"

type FilterGroupProps = {
    data: FilterConfig[]
    onChange?: (filters: FilterConfig[]) => void
}


const FilterGroup = ({ data, onChange }: FilterGroupProps) => {
    const [filterData, setFilterData] = React.useState<FilterConfig[]>(data)

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


    const updateFilters = (
        updater: (prev: FilterConfig[]) => FilterConfig[]
    ) => {
        setFilterData((prev) => {
            const next = updater(prev)
            onChange?.(next)
            return next
        })
    }


    return (
        <div className="flex gap-4">
            {filterData
                .filter((f) => f.isEnabled)
                .map((filter) => (
                    <FilterDropdown
                        key={filter.key}
                        filter={filter}
                        onToggleOption={toggleOption}
                    />
                ))}

            <FilterToggleMenu
                filters={filterData}
                onToggleFilter={toggleFilter}
            />
        </div>
    )
}

export default FilterGroup