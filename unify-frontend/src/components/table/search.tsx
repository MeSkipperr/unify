import { Field, FieldDescription, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"

export type SearchBarProps = {
    value: string
    onChange: (value: string) => void
    className?: string
    label?: string
    id: string
    placeholder?: string
    description?: string
}

export const SearchBar = ({
    value,
    onChange,
    className,
    label = "Search",
    id,
    placeholder,
    description,
}: SearchBarProps) => {
    return (
        <Field className={className}>
            <FieldLabel htmlFor={id}>{label}</FieldLabel>

            <Input
                id={id}
                type="search"
                placeholder={placeholder}
                value={value}
                onChange={(e) => onChange(e.target.value)}
            />

            {description || description?.trim() == "" && (
                <FieldDescription>{description}</FieldDescription>
            )}
        </Field>
    )
}