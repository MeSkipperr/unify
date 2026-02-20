import AdbTable from "@/features/adb/components/table"

export default async function DeviceType({
    params,
}: {
    params: Promise<{ type: string }>
}) {
    const { type } = await params
    return (
        <div className="h-[85dvh]">
            <AdbTable serviceType={type} hasDefaultValue/>
        </div>
    )
}