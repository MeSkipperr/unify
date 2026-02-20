import { NotificationType } from "@/types/notification";
import NotificationCard from "./notification-card";

interface ListProps {
    notifications: NotificationType[];
}

export function NotificationList({ notifications }: ListProps) {
    if (!notifications || notifications.length === 0) {
        return (
            <div className="text-sm text-muted-foreground text-center py-6">
                No notifications available
            </div>
        );
    }

    return (
        <div className="space-y-3">
            {notifications.map((item) => (
                <NotificationCard key={item.id} data={item} />
            ))}
        </div>
    );
}
