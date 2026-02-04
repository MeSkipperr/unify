/**
 * Handle input IPv4:
 * - Hanya angka & titik
 * - Maksimal 4 oktet
 * - Range 0–255
 * - Menghapus leading zero (008 → 8)
 */
export function handleIPv4Input(
    value: string,
    setValue: (val: string) => void
) {
    // hanya angka dan titik
    if (!/^[0-9.]*$/.test(value)) return;

    const rawParts = value.split(".");

    // maksimal 4 oktet
    if (rawParts.length > 4) return;

    const normalizedParts: string[] = [];

    for (const part of rawParts) {
        // max 3 digit per oktet
        if (part.length > 3) return;

        if (part === "") {
            normalizedParts.push("");
            continue;
        }

        const num = Number(part);

        // range IPv4
        if (Number.isNaN(num) || num > 255) return;

        // 🔥 hapus leading zero
        normalizedParts.push(String(num));
    }

    setValue(normalizedParts.join("."));
}

/**
 * Validasi final IPv4
 */
export function isValidIPv4(value: string): boolean {
    const ipv4Regex =
        /^(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}$/;

    return ipv4Regex.test(value);
}

export function normalizeIPv4(value: string): string {
    if (!value) return "";

    const trimmed = value.trim();
    const parts = trimmed.split(".");

    // harus 4 oktet
    if (parts.length !== 4) return "";

    const normalized = parts.map((part) => {
        // kosong atau bukan angka
        if (!/^\d+$/.test(part)) return "";

        // parse untuk menghilangkan leading zero
        const num = Number(part);

        // range IPv4
        if (num < 0 || num > 255) return "";

        return String(num);
    });

    // jika ada oktet invalid
    if (normalized.some((p) => p === "")) return "";

    return normalized.join(".");
}