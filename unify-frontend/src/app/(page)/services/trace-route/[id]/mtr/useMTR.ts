// hooks/useMTR.ts
import { useEffect, useRef, useState } from "react"
import type { MTRData } from "./types"

export const useMTR = (url: string) => {
    const [data, setData] = useState<MTRData>()
    const wsRef = useRef<WebSocket | null>(null)

    useEffect(() => {
        if (wsRef.current) return // ⛔ cegah double connect

        const ws = new WebSocket(url)
        wsRef.current = ws

        ws.onopen = () => {
            console.log("WebSocket connected")
        }

        ws.onmessage = (event) => {
            try {
                setData(JSON.parse(event.data))
            } catch (err) {
                console.error("Parse error", err)
            }
        }

        ws.onerror = () => {
            console.warn("WebSocket error")
        }

        ws.onclose = (event) => {
            console.log("WebSocket closed", event)
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
