"use client"
import { Button } from "@/components/ui/button";
import {
    CardContent,
    CardDescription,
    CardFooter,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { loginAction } from "./actions";
import { useState } from "react";


const LoginForm = () => {
    const property = process.env.NEXT_PUBLIC_PROPERTY;
    const date = new Date();

    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [isLoading, setIsLoading] = useState<boolean>(false);

    const submitHandler = async () => {
        setIsLoading(true)

        if (username.trim() === "" || password.trim() === "") {
            return
        }

        const payload = {
            username,
            password
        }
        try {
            await loginAction(payload)
        } catch (error) {

        } finally {
            setIsLoading(false)
        }
    }
    return (
        <>
            <CardContent >
                <form onSubmit={submitHandler}>
                    <div className="flex flex-col gap-6">
                        <div className="grid gap-2">
                            <Label htmlFor="username">Username</Label>
                            <Input
                                id="username"
                                type="username"
                                required
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                            />
                        </div>
                        <div className="grid gap-2">
                            <Label htmlFor="password">Password</Label>
                            <Input
                                id="password"
                                type="password"
                                required
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                            />
                        </div>
                    </div>
                </form>
            </CardContent>
            <CardFooter className="flex-col gap-2">
                <Button type="submit" disabled={isLoading} className="w-full" onClick={submitHandler}>
                    Sign In
                </Button>
                <CardDescription className="text-center flex flex-col">
                    <span>
                        {property}
                    </span>
                    <span>
                        © Unify {date.getFullYear()}
                    </span>
                </CardDescription>
            </CardFooter>
        </>
    );
}

export default LoginForm;