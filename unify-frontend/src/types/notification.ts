export type NotificationLevel = "info" | "success" | "alert" | "error";

export interface NotificationType {
  id: string;
  level: NotificationLevel;
  title: string;
  detail?: string;
  url?: string;
  createdAt: string;
}
