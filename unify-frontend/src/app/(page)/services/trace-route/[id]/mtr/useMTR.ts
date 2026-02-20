import { useEffect, useRef, useState, useCallback } from "react";
import type { MTRData } from "./types";
import { toast } from "sonner";

export const useMTR = (url: string) => {
  const [data, setData] = useState<MTRData>();
  const [enabled, setEnabled] = useState(false);
  const eventSourceRef = useRef<EventSource | null>(null);

  const start = useCallback(() => {
    if (eventSourceRef.current) return;

    const es = new EventSource(url);
    eventSourceRef.current = es;

    es.onmessage = (event) => {
      try {
        const parsed = JSON.parse(event.data);
        setData(parsed);
      } catch {
        toast.error("SSE parse error", { position: "bottom-right" });
      }
    };

    es.onerror = () => {
      console.error("SSE error");
    };
  }, [url]);

  const stop = useCallback(() => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }
  }, []);

  useEffect(() => {
    if (enabled) {
      start();
    } else {
      stop();
    }

    return () => stop();
  }, [enabled, start, stop]);

  return {
    data,
    start: () => setEnabled(true),
    stop: () => setEnabled(false),
    isRunning: enabled,
  };
};
