import {
    AiOutlineInfoCircle,
    AiOutlineCheckCircle,
    AiOutlineWarning,
    AiOutlineCloseCircle
} from "react-icons/ai"

import { NotificationType } from "@/types/notification"
import { getCompactRelativeTime } from "@/utils/time"

interface Props {
    data: NotificationType
}

export default function NotificationCard({ data }: Props) {
    const levelConfig = {
        info: {
            icon: <AiOutlineInfoCircle size={18} />,
            iconColor: "text-primary",
            bg: "bg-primary/10",
        },
        success: {
            icon: <AiOutlineCheckCircle size={18} />,
            iconColor: "text-emerald-500 dark:text-emerald-400",
            bg: "bg-emerald-500/10",
        },
        alert: {
            icon: <AiOutlineWarning size={18} />,
            iconColor: "text-amber-500 dark:text-amber-400",
            bg: "bg-amber-500/10",
        },
        error: {
            icon: <AiOutlineCloseCircle size={18} />,
            iconColor: "text-destructive",
            bg: "bg-destructive/10",
        },
    }

    const config = levelConfig[data.level]

    return (
        <div
            className="
                w-full 
                rounded-xl
                bg-card
                p-3
                transition-all
                shadow-sm
                hover:bg-accent
            "
        >
            <div className="flex items-start gap-2">
                <div
                    className={`
                        flex h-7 w-7 items-center justify-center rounded-lg
                        ${config.bg}
                        ${config.iconColor}
                    `}
                >
                    {config.icon}
                </div>

                <div className="flex-1 space-y-1">
                    <div className="flex items-center justify-between gap-2">
                        <h3 className="text-sm font-medium text-foreground leading-none">
                            {data.title}
                        </h3>

                        <span className="text-xs text-muted-foreground whitespace-nowrap">
                            {getCompactRelativeTime(new Date(data.createdAt))}
                        </span>
                    </div>

                    {data.detail && (
                        <p className="text-xs text-muted-foreground leading-snug">
                            {data.detail}
                        </p>
                    )}

                    {data.url && (
                        <a
                            href={data.url}
                            className="inline-block text-xs font-medium text-primary hover:underline"
                        >
                            View Details
                        </a>
                    )}
                </div>
            </div>
        </div>
    )
}
