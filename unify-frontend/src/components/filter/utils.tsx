import { OptionStruct, SelectedOptionSummary, } from "./types"

export const getSelectedOptionSummary = (
    options: OptionStruct[]
): SelectedOptionSummary => {
    const selected = options.filter((o) => o.isSelected)

    const labels = selected.map((o) => o.label)
    const values = selected.map((o) => o.value)

    return {
        isEmpty: selected.length === 0,
        isFull: selected.length === options.length,
        values,
        labels,
        text: labels.join(", "),
    }
}
