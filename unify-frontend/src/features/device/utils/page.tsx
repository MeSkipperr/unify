export function getLastPage(totalData: number, perPage: number): number {
    if (perPage <= 0) return 0
    return Math.ceil(totalData / perPage)
}
export function generatePages(currentPage: number, lastPage: number) {
    const pages: (number | "ellipsis")[] = []

    if (lastPage <= 5) {
        for (let i = 1; i <= lastPage; i++) {
            pages.push(i)
        }
        return pages
    }

    // selalu halaman pertama
    pages.push(1)

    // ===== AWAL =====
    if (currentPage <= 3) {
        for (let i = 2; i <= 5; i++) {
            pages.push(i)
        }
        pages.push("ellipsis")
    }

    // ===== TENGAH =====
    else if (currentPage < lastPage - 2) {
        pages.push("ellipsis")
        pages.push(currentPage - 1, currentPage, currentPage + 1)
        pages.push("ellipsis")
    }

    // ===== AKHIR =====
    else {
        pages.push("ellipsis")
        for (let i = lastPage - 4; i < lastPage; i++) {
            pages.push(i)
        }
    }

    // selalu halaman terakhir
    pages.push(lastPage)

    return pages
}

