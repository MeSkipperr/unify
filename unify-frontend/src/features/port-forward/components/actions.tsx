
import { Ban } from "lucide-react";
import { deactivatePortForward } from "../api/port-forward.api";
import { toast } from "sonner";
import { PortForwardResult } from "../types";

type ActionsProps = {
    row: PortForwardResult
    handleFetchData: () => void
}

const DisablePortForward = ({ row, handleFetchData }: ActionsProps) => {
    const handleDisable = async () => {
        try {
            await deactivatePortForward(row.id)
            toast.success("Device has been deleted successfully!", { position: "bottom-right" })
        } catch (error) {
            toast.error("Failed to delete device. Please try again.", { position: "bottom-right" })
        } finally {
            handleFetchData()
        }
    }

    return (
        <Ban className='text-destructive cursor-pointer' onClick={() => handleDisable()} />
    );
}

export default DisablePortForward;