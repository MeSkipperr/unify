import { Field, FieldDescription, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"

const SearchBar = ({className}:{className ?: string}) => {
    return (
        <Field className={className}>
            <FieldLabel htmlFor="">Search Device</FieldLabel>
            <Input id="" type="search" placeholder="DPSCY-..." />
            <FieldDescription>
                Search by name, MAC, IP, room, device type, or description.
            </FieldDescription>
        </Field>
    );
}

export default SearchBar;