import type { Metadata } from "next";

export const metadata: Metadata = {
    title: "Devices | Unify",
    description: "Manage and monitor all connected devices in the Unify system.",
};

export default function DeviceLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <div>
            {children}
        </div>
    )
    ;
}
