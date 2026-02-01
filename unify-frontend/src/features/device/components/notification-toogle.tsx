import React from "react"
import { Button } from "@/components/ui/button"
import { Bell, BellOff } from "lucide-react"
import { updateDeviceNotification } from '../api/device.mutation'
import { useDebouncedValue } from '@/hooks/useDebounce'

interface NotificationToggleProps {
    deviceId: string
    initialNotification: boolean
}

export const NotificationToggle: React.FC<NotificationToggleProps> = ({
    deviceId,
    initialNotification,
}) => {
    const { value, setDebouncedValue } = useDebouncedValue<boolean>({
        initialValue: initialNotification,
        delay: 500,
        onChange: async (val) => {
            try {
                await updateDeviceNotification(deviceId, val)
                return
            } catch (err) {
                return
            }
        },
    })

    return (
        <Button
            variant="ghost"
            size="icon"
            onClick={() => setDebouncedValue(!value)}
        >
            {value ? <Bell className="size-4" /> : <BellOff className="size-4" />}
        </Button>
    )
}
