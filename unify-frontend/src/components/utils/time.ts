/**
 * Returns compact relative time from a past date to now
 * Example: 1m, 1h, 1d, 1w, 1w1d, 1m3d
 */
export const getCompactRelativeTime = (past: Date): string => {
    const now = Date.now()
    const diff = now - new Date(past).getTime()

    if (diff <= 0) return "now"

    const minute = 1000 * 60
    const hour = minute * 60
    const day = hour * 24
    const week = day * 7
    const month = day * 30

    const months = Math.floor(diff / month)
    const weeks = Math.floor((diff % month) / week)
    const days = Math.floor((diff % week) / day)
    const hours = Math.floor((diff % day) / hour)
    const minutes = Math.floor((diff % hour) / minute)

    if (months > 0) {
        return `${months}m${days > 0 ? `${days}d` : ""}`
    }

    if (weeks > 0) {
        return `${weeks}w${days > 0 ? `${days}d` : ""}`
    }

    if (days > 0) return `${days}d`
    if (hours > 0) return `${hours}h`
    if (minutes > 0) return `${minutes}m`

    return "now"
}

/**
 * Format date to: HH:mm - dd/MM/yyyy
 * Example: 14:32 - 29/01/2026
 */
export const formatDateTime = (date: Date): string => {
    const d = new Date(date)

    const hours = String(d.getHours()).padStart(2, "0")
    const minutes = String(d.getMinutes()).padStart(2, "0")
    const day = String(d.getDate()).padStart(2, "0")
    const month = String(d.getMonth() + 1).padStart(2, "0")
    const year = d.getFullYear()

    return `${hours}:${minutes} - ${day}/${month}/${year}`
}