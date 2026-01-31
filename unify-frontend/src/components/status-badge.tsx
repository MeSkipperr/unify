import { CheckCircle2, XCircle } from "lucide-react"
import { Badge } from "./ui/badge";

const StatusBadge = ({ isConnect }: { isConnect: boolean }) => {
    return isConnect ? (
        <>
            <Badge className="flex items-center gap-1 bg-green-100 text-green-700 font-bold">
                <CheckCircle2 className="h-4 w-4 text-green-700 " />
                UP
            </Badge>
        </>

    ) : (
        <Badge className="flex items-center gap-1 bg-red-100 text-red-700 font-bold" >
            <XCircle className="h-4 w-4  text-red-700 " />
            DOWN
        </Badge>
    )
}

export default StatusBadge;