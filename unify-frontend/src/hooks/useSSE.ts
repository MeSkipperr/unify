import { useEffect, useRef, useState } from "react";

interface UseSSEOptions<T> {
  url: string;
  onMessage?: (data: T) => void;
  withCredentials?: boolean;
}

export function useSSE<T = string>({
  url,
  onMessage,
  withCredentials = true,
}: UseSSEOptions<T>) {
  const [data, setData] = useState<T | null>(null);
  const [enabled, setEnabled] = useState(false);

  const eventSourceRef = useRef<EventSource | null>(null);

  const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "";
  const fullUrl = url.startsWith("http") ? url : `${baseURL}${url}`;

  useEffect(() => {
    if (!enabled) {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
        eventSourceRef.current = null;
      }
      return;
    }

    if (eventSourceRef.current) return;

    const es = new EventSource(fullUrl, { withCredentials });

    eventSourceRef.current = es;

    es.onmessage = (event) => {
      try {
        const parsed = JSON.parse(event.data);
        setData(parsed);
        onMessage?.(parsed);
      } catch (err) {
        console.error("SSE parse error:", err);
      }
    };

    es.onerror = (err) => {
      console.error("SSE error:", err);
      es.close();
      eventSourceRef.current = null;
      setEnabled(false);
    };

    return () => {
      es.close();
      eventSourceRef.current = null;
    };
  }, [enabled, fullUrl, onMessage, withCredentials]);

  return {
    data,
    start: () => setEnabled(true),
    stop: () => setEnabled(false),
    isRunning: enabled,
  };
}
