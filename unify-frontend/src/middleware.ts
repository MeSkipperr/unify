// import { cookies } from "next/headers"
import { cookies } from "next/headers"
import { NextRequest, NextResponse } from "next/server"

const PUBLIC_PATHS = ["/login", "/register"]

export async function middleware(req: NextRequest) {
    const pathname = req.nextUrl.pathname

    if (PUBLIC_PATHS.includes(pathname)) {
        return NextResponse.next()
    }

    const cookieStore = await cookies()
    const cookieHeader = cookieStore
        .getAll()
        .map(c => `${c.name}=${c.value}`)
        .join("; ")
    const origin = req.nextUrl.origin;
    const baseUrl =
        process.env.NEXT_PUBLIC_API_BASE_URL?.replace(/\/$/, "") ||
        origin;
    const res = await fetch(
        `${baseUrl}/auth/me`,
        {
            method: "POST",
            headers: {
                Cookie: cookieHeader,
                "Content-Type": "application/json",
            },
            credentials: "include",
            cache: "no-store",
        }
    )

    if (!res.ok) {
        return NextResponse.redirect(new URL("/login", req.url))
    }

    return NextResponse.next()
}

export const config = {
    matcher: [
        "/((?!login|register|_next/static|_next/image|favicon.ico).*)",
    ],
}
