import DeviceTableData from "@/features/device/components"

export default async function DeviceType({
    params,
}: {
    params: Promise<{ type: string }>
}) {
    const { type } = await params
    return (
        <div className="h-[90dvh]">
            <DeviceTableData selectType={type} />
        </div>
    )
}