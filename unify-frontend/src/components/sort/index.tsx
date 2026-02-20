import {
    Popover,
    PopoverContent,
    PopoverDescription,
    PopoverHeader,
    PopoverTitle,
    PopoverTrigger,
} from "@/components/ui/popover"

import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"

import { Button } from "../ui/button";
import { Label } from "@radix-ui/react-label";
import { ArrowDownUp, ChevronDown } from "lucide-react";
import { SortBy } from "./types";

interface SortGroupProps {
    sortOptions: SortBy[]
    onChange: (updated: SortBy[]) => void
}

export const SortGroup: React.FC<SortGroupProps> = ({ sortOptions, onChange }) => {
    const handleChange = (key: string, newValue: "ascending" | "descending" | "none") => {
        const updated = sortOptions.map(item =>
            item.key === key ? { ...item, value: newValue } : item
        )
        onChange(updated)
    }
    return (
        <Popover>
            <PopoverTrigger asChild>
                <Button variant="default"><ArrowDownUp />Sort<ChevronDown /></Button>
            </PopoverTrigger>
            <PopoverContent>
                <div className="flex flex-col items-center justify-between w-full mt-2 gap-2">
                    {sortOptions.map(item => (
                        <div key={item.key} className="w-full flex items-center justify-between">
                            <Label className="text-sm">{item.label} :</Label>
                            <Select value={item.value} onValueChange={(val) => handleChange(item.key, val as SortBy["value"])}>
                                <SelectTrigger className="w-[150px]">
                                    <SelectValue placeholder="Sort"   />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="ascending">Ascending</SelectItem>
                                    <SelectItem value="descending">Descending</SelectItem>
                                    <SelectItem value="none">None</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                    ))}

                </div>
            </PopoverContent>
        </Popover>
    );
}

export default SortGroup;