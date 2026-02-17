"use client"

import { Button } from "@/components/ui/button"
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from "@/components/ui/sheet"
import { Bell, ChevronLeft, ChevronRight } from "lucide-react"
import { NotificationType } from "@/types/notification";
import { NotificationList } from "./notification-list";
import { useEffect, useState } from "react";
import { useSSE } from "@/hooks/useSSE";
import { getNotification } from "./notif.api";
import { TableQuery } from "../table/types";
import { Skeleton } from "@/components/ui/skeleton";

const PAGESIZE = 25;
const MAX_NOTIFICATION = 100;

export function Notification() {
    const [open, setOpen] = useState(false)
    const [hasNewNotification, setHasNewNotification] = useState<boolean>(false)
    const [dataNotification, setDataNotification] = useState<NotificationType[]>([])
    const [page, setPage] = useState<number>(1)
    const [maxPage, setMaxPage] = useState<number>(1)
    const [isLoading, setIsLoading] = useState<boolean>(true)

    // SSE Hook
    const { start, stop } = useSSE<NotificationType>({
        url: "/events/notification",
        onMessage: (msg) => {
            setHasNewNotification(true)
            setDataNotification((prev) => {
                const updated = [msg, ...prev]
                // Limit jumlah notification di client
                if (updated.length > MAX_NOTIFICATION) {
                    updated.pop()
                }
                return updated
            })
        },
    });

    // Fetch initial data
    const fetchNotification = async () => {
        setIsLoading(true)
        try {
            const payload: Pick<TableQuery, "page" | "pageSize"> = {
                page,
                pageSize: PAGESIZE
            }
            const notif = await getNotification(payload)
            console.log(notif)
            setDataNotification(notif.data)
            setMaxPage(notif.totalPages)
        } catch (err) {
            console.error("Failed to fetch notifications:", err)
        } finally {
            setIsLoading(false)
        }
    }

    useEffect(() => {
        fetchNotification()
        start()
        return () => stop()
    }, [page])

    // Pagination handlers
    const handlePrevPage = () => {
        if (page > 1) {
            setPage(page - 1)
            stop()
        } else {
            start()
        }
    }

    const handleNextPage = () => {
        if (page <= maxPage) {
            setPage(page + 1)
            stop()
        }

    }

    return (
        <Sheet open={open} onOpenChange={setOpen}>
            <SheetTrigger asChild>
                <Button
                    variant="ghost"
                    size="icon"
                    className="relative p-2"
                    onClick={() => setHasNewNotification(false)}
                >
                    <Bell className="size-5" />
                    {hasNewNotification && (
                        <span className="absolute top-1 right-1 h-2 w-2 rounded-full bg-red-500" />
                    )}
                </Button>
            </SheetTrigger>

            <SheetContent
                side="right"
                className="
          w-full 
          sm:max-w-sm 
          md:max-w-md 
          lg:max-w-lg
          p-0
        "
            >
                <SheetHeader className="px-6 pt-6 pb-4 border-b flex flex-row">
                    <div className="w-2/3">
                        <SheetTitle>Notifications</SheetTitle>
                        <SheetDescription>
                            Latest system updates and alerts
                        </SheetDescription>
                    </div>
                    <div className="w-full flex gap-2 justify-end p-2">
                        <Button variant="ghost" onClick={handlePrevPage} disabled={page === 1}>
                            <ChevronLeft />
                        </Button>
                        <Button variant="ghost" onClick={handleNextPage} disabled={page === maxPage}>
                            <ChevronRight />
                        </Button>
                    </div>
                </SheetHeader>

                <div className="h-[calc(100vh-120px)] overflow-y-auto px-4 py-4">
                    {isLoading ? (
                        <div className="space-y-3">
                            {Array.from({ length: PAGESIZE }).map((_, idx) => (
                                <Skeleton key={idx} className="h-16 w-full rounded-md" />
                            ))}
                        </div>
                    ) : (
                        <NotificationList notifications={dataNotification ?? []} />
                    )}
                </div>
            </SheetContent>
        </Sheet>
    )
}
