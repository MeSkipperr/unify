export type OptionValue = boolean | string

export type OptionStruct = {
    value: OptionValue
    label: string
    isSelected: boolean
}

export type FilterConfig = {
    key: string
    label: string
    type: "boolean" | "select"
    isEnabled: boolean
    options: OptionStruct[]
}
export type SelectedOptionSummary = {
    isEmpty: boolean
    isFull: boolean
    values: OptionValue[]
    labels: string[]
    text: string
}
