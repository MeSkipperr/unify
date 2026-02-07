// hooks/useMTR.ts
import { useEffect, useState } from "react"
import type { MTRData } from "./types"

export const useMTR = (url: string) => {
    const [data, setData] = useState<MTRData>()

    useEffect(() => {
        const ws = new WebSocket(url)

        ws.onopen = () => {
            console.log("WebSocket connected")
        }

        ws.onmessage = (event) => {
            try {
                const parsed: MTRData = JSON.parse(event.data)
                setData(parsed)
            } catch (err) {
                console.error("Failed to parse message", err)
            }
        }

        ws.onclose = () => {
            console.log("WebSocket disconnected")
        }

        ws.onerror = (err) => {
            console.error("WebSocket error", err)
        }

        return () => ws.close()
    }, [url])

    return data
}
