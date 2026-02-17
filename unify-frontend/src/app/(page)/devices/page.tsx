import DeviceTableData from "@/features/device/components";
import { Suspense } from "react";

const DevicesList = () => {
    return (
        <div className="h-[90dvh]">
            <Suspense fallback={<div>Loading...</div>}>
                <DeviceTableData />
            </Suspense>
        </div>
    );
}

export default DevicesList;