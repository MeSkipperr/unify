import { Label } from "@/components/ui/label";
import TraceDetailGroup from "./mtr/trace-group";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";

const TraceId = async ({
    params,
}: {
    params: Promise<{ id: string }>
}) => {
    const { id } = await params
    return (
        <div className="">
            <Link href="/services/trace-route" >
                <Label className="size-8 cursor-pointer">
                    <ArrowLeft />
                </Label>
            </Link>
            <TraceDetailGroup traceId={id} />
        </div>
    );
}

export default TraceId;