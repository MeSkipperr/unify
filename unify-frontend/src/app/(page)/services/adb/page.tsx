"use client"

import AdbTable from "@/features/adb/components/table"
import { useEffect } from "react"

const AdbPage = () => {
    
    useEffect(() => {

        const eventSource = new EventSource(
            "http://localhost:8080/events/adb"
        )

        eventSource.onmessage = (event) => {
            console.log("SSE Message:", event.data)
        }

        eventSource.onerror = (error) => {
            console.error("SSE Error:", error)
        }

        return () => {
            eventSource.close()
        }

    }, [])

    return (
        <div className="h-[85dvh]">
            <AdbTable serviceType="manual" />
        </div>
        
    );
}

export default AdbPage;