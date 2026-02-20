import {
    Card,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import LoginForm from "./login-form";

const LoginPage = () => {
    const hostname = process.env.HOSTNAME;
    return (
        <div className="h-dvh w-full flex justify-center items-center ">
            <Card className="w-full max-w-sm">
                <CardHeader>
                    <div className="w-full flex justify-between">
                    <CardTitle>UNIFY</CardTitle>
                    <CardDescription>{hostname}</CardDescription>
                    </div>
                    <CardDescription>
                        Sign in to continue to the system
                    </CardDescription>
                </CardHeader>
                <LoginForm/>
            </Card>
        </div>
    );
}

export default LoginPage;