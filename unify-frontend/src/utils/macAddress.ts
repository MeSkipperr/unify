/**
 * Membatasi input agar hanya bisa format MAC Address
 * Format: AA:BB:CC:DD:EE:FF
 */
export function handleMacAddressInput(
    value: string,
    setValue: (val: string) => void
) {
    // hanya hex dan colon
    if (!/^[0-9a-fA-F:]*$/.test(value)) return;

    // ubah ke uppercase biar konsisten
    const upper = value.toUpperCase();

    const parts = upper.split(":");

    // maksimal 6 oktet
    if (parts.length > 6) return;

    for (const part of parts) {
        // tiap oktet max 2 hex
        if (part.length > 2) return;
    }

    setValue(upper);
}

/**
 * Validasi MAC Address
 * Accept:
 *  - AA:BB:CC:DD:EE:FF
 */
export function isValidMacAddress(value: string): boolean {
    const macRegex = /^([0-9A-F]{2}:){5}[0-9A-F]{2}$/;
    return macRegex.test(value.toUpperCase());
}

export function normalizeMacAddress(value: string): string {
    if (!value) return "";

    // Trim & uppercase
    const cleaned = value.trim().toUpperCase();

    // Ambil hanya hex
    const hex = cleaned.replace(/[^0-9A-F]/g, "");

    // Harus 12 karakter hex
    if (hex.length !== 12) return "";

    // Format ke AA:BB:CC:DD:EE:FF
    const parts = hex.match(/.{2}/g);
    if (!parts) return "";

    return parts.join(":");
}