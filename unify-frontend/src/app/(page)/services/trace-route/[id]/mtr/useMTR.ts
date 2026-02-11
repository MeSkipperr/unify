// hooks/useMTR.ts
import { useEffect, useRef, useState } from "react"
import type { MTRData } from "./types"
import { toast } from "sonner"

export const useMTR = (url: string) => {
    const [data, setData] = useState<MTRData>()
    const wsRef = useRef<WebSocket | null>(null)

    useEffect(() => {
        if (wsRef.current) return 

        const ws = new WebSocket(url)
        wsRef.current = ws

        ws.onopen = () => {
            toast.success("WebSocket connected", { position: "bottom-right" })
        }

        ws.onmessage = (event) => {
            try {
                setData(JSON.parse(event.data))
            } catch (err) {
                toast.error("Parse error", { position: "bottom-right" })
            }
        }

        ws.onerror = () => {
            
            toast.error("WebSocket error", { position: "bottom-right" })
        }

        ws.onclose = (event) => {
            toast.success("WebSocket closed", { position: "bottom-right" })
            wsRef.current = null
        }

        return () => {
            if (ws.readyState === WebSocket.OPEN) {
                ws.close()
            }
        }
    }, [url])

    return data
}
