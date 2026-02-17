import PortForwardTable from "@/features/port-forward/components/table"
import { Suspense } from "react"

const PortForwardPage = () => {
    return (
        <div >
            <Suspense fallback={<div>Loading...</div>}>
                <PortForwardTable />
            </Suspense>
        </div>
    )
}

export default PortForwardPage  