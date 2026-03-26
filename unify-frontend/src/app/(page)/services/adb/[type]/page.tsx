import AdbTable from "@/features/adb/components/table"
import { Suspense } from "react"

export default async function DeviceType({
    params,
}: {
    params: Promise<{ type: string }>
}) {
    const { type } = await params
    return (
        <div className="h-[85dvh]">
            <Suspense fallback={<div>Loading...</div>}>
                <AdbTable serviceType={type} hasDefaultValue />
            </Suspense>
        </div>
    )
}