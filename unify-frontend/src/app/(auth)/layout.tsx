import AppNavbar from "@/components/navbar";
import AppSidebar from "@/components/sidebar";
import { SidebarProvider } from "@/components/ui/sidebar";

export default function AuthLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <main >
            {children}
        </main>
    );
}
