import { Button } from "@/components/ui/button";
import { ArrowUpIcon } from "lucide-react"

const NewPage = () => {
    return (
        <div className="w-full h-dvh flex justify-center items-center">
            <Button variant="outline">Button</Button>
            <Button variant="outline" size="icon" aria-label="Submit">
                <ArrowUpIcon />
            </Button>
        </div>
    );
}

export default NewPage;