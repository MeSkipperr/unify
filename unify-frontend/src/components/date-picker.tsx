"use client"

import * as React from "react"
import { Calendar } from "@/components/ui/calendar"
import { Button, buttonVariants } from "@/components/ui/button"
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"
import { format } from "date-fns"
import { cn } from "@/lib/utils"
import { VariantProps } from "class-variance-authority"
import { CalendarIcon } from "lucide-react"

type DatePickerProps = {
    value?: Date
    onChange?: (date: Date) => void
    onAfterChange?: (date: Date) => void

    className?: string
    variant?: VariantProps<typeof buttonVariants>["variant"]
    minDate?: Date
    maxDate?: Date
}

export function DatePicker({
    value,
    onChange,
    className,
    variant,
    minDate,
    maxDate,
    onAfterChange
}: DatePickerProps) {
    const [open, setOpen] = React.useState(false)

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button
                    variant={variant ?? "default"}
                    className={cn(" justify-start", className)}
                >
                    <CalendarIcon />
                    {value ? format(value, "dd MMM yyyy") : "Pick a date"}
                </Button>
            </PopoverTrigger>

            <PopoverContent className="w-auto p-0" align="start">
                <Calendar
                    mode="single"
                    selected={value}
                    onSelect={(selected) => {
                        if (!selected) return

                        if (minDate && selected < minDate) return
                        if (maxDate && selected > maxDate) return

                        onChange?.(selected)
                        onAfterChange?.(selected)
                        setOpen(false)
                    }}

                    disabled={(date) => {
                        if (minDate && date < minDate) return true
                        if (maxDate && date > maxDate) return true
                        return false
                    }}
                />

            </PopoverContent>
        </Popover>
    )
}
