import { useState, useRef, useEffect, useCallback } from "react"

/**
 * useDebouncedValue
 * 
 * A reusable React hook for debouncing any value (boolean, string, number, object, etc.).
 * It triggers a callback only after the value stops changing for a specified delay.
 * If the value remains the same as the initial value, the callback will not run.
 *
 * @template T - Type of the value (boolean, string, number, etc.)
 *
 * @param {object} options - Hook configuration options
 * @param {T} options.initialValue - Initial value to track
 * @param {number} [options.delay=800] - Debounce delay in milliseconds
 * @param {(value: T) => void} [options.onChange] - Callback triggered after debounce
 *
 * @returns {{
 *   value: T, 
 *   setDebouncedValue: (newValue: T) => void
 * }} 
 * - value: Current value (state)
 * - setDebouncedValue: Function to update the value and trigger debounce
 *
 * @example
 * // For search input
 * const { value, setDebouncedValue } = useDebouncedValue<string>({
 *   initialValue: "",
 *   delay: 500,
 *   onChange: (val) => console.log("Search API call:", val),
 * })
 *
 * <input
 *   value={value}
 *   onChange={(e) => setDebouncedValue(e.target.value)}
 * />
 *
 * @example
 * // For boolean toggle like button
 * const { value, setDebouncedValue } = useDebouncedValue<boolean>({
 *   initialValue: false,
 *   delay: 1000,
 *   onChange: (val) => console.log("Update like:", val),
 * })
 *
 * <button onClick={() => setDebouncedValue(!value)}>
 *   {value ? "❤️ Liked" : "🤍 Like"}
 * </button>
 */
export function useDebouncedValue<T>({
    initialValue,
    delay = 800,
    onChange,
}: {
    initialValue: T
    delay?: number
    onChange?: (value: T) => void
}) {
    const [value, setValue] = useState<T>(initialValue)
    const initialRef = useRef<T>(initialValue)
    const debounceRef = useRef<NodeJS.Timeout | null>(null)

    const setDebouncedValue = useCallback(
        (newValue: T) => {
            setValue(newValue)

            if (debounceRef.current) clearTimeout(debounceRef.current)

            debounceRef.current = setTimeout(() => {
                if (newValue === initialRef.current) {
                    return
                }

                onChange?.(newValue) // Replace with fetch or API call

                initialRef.current = newValue
            }, delay)
        },
        [delay, onChange]
    )

    // Update value if initialValue changes externally
    useEffect(() => {
        setValue(initialValue)
        initialRef.current = initialValue
    }, [initialValue])

    // Cleanup debounce timer on unmount
    useEffect(() => {
        return () => {
            if (debounceRef.current) clearTimeout(debounceRef.current)
        }
    }, [])

    return { value, setDebouncedValue }
}
