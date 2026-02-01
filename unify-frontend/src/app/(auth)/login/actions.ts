"use server"

import { loginRequest } from "@/services/auth.service"

export async function loginAction(payload: {
    username: string
    password: string
}) {
    console.log('satu')
    if (payload.username.trim() === "" || payload.password.trim() === "") {
        throw new Error("Invalid credentials")
    }
    console.log('dua')

    const result = await loginRequest(payload)

    // contoh: set cookie / session
    // cookies().set("token", result.token)

    return result
}
