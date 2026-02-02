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
import { useRef, useState } from "react";
import { useRouter } from "next/navigation";


const LoginForm = () => {
    const router = useRouter();

    const property = process.env.NEXT_PUBLIC_PROPERTY;
    const date = new Date();

    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [errorMsg, setErrorMsg] = useState<string>("");
    const errorTimerRef = useRef<NodeJS.Timeout | null>(null);


    const submitHandler = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setErrorMsg("");

        if (!username.trim() || !password.trim()) {
            setErrorMsg("Username and password are required.");
            setIsLoading(false);
            return;
        }

        try {
            await loginAction({ username, password });
            // success handling / redirect
            router.push("/")
        } catch (error: any) {
            if (error.status === 401) {
                setErrorMsg("Invalid username or password.");
            } else {
                setErrorMsg("Something went wrong. Please try again.");
            }

            // clear previous timer
            if (errorTimerRef.current) {
                clearTimeout(errorTimerRef.current);
            }

            // auto clear after 3 seconds
            errorTimerRef.current = setTimeout(() => {
                setErrorMsg("");
            }, 3000);
        } finally {
            setIsLoading(false);
        }
    };


    return (
        <form onSubmit={submitHandler} className="space-y-2">
            <CardContent >
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
                        <Label className="text-destructive">
                            {errorMsg}
                        </Label>
                    </div>
                </div>
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
        </form>
    );
}

export default LoginForm;