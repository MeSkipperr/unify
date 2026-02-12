import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from "@/components/ui/sheet"

import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog"

import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Button } from "@/components/ui/button";
import { EllipsisVertical } from "lucide-react";
import { Label } from "@radix-ui/react-label";
import { DeviceStatus } from "@/components/status";
import Link from "next/link";
import { deleteDevice } from "../../api/device.api";
import { Device } from "../../types";
import { formatDateTime } from "@/utils/time";
import { toast } from "sonner";
import ChangeData from "./change-data";

type ActionsProps = {
    row: Device
    handleFetchData: () => void
}

const ActionsColumns = ({ row, handleFetchData }: ActionsProps) => {
    const handleDeleteDevice = async () => {
        try {
            await deleteDevice(row.id)
            toast.success("Device has been deleted successfully!", { position: "bottom-right" })
        } catch (error) {
            toast.error("Failed to delete device. Please try again.", { position: "bottom-right" })
        } finally {
            handleFetchData()
        }
    }

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" >
                    <EllipsisVertical className="h-4 w-4" />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-40" align="start">
                <DropdownMenuLabel className=' font-bold'>Actions</DropdownMenuLabel>
                <Sheet>
                    <SheetTrigger asChild>
                        <Button variant="ghost" className='w-full flex justify-start gap-0 py-0 px-2'>View Details</Button>
                    </SheetTrigger>
                    <SheetContent className="sm:max-w-md">
                        <SheetHeader>
                            <SheetTitle className="flex items-center gap-4">
                                Device Details
                                <DeviceStatus isConnect={row.isConnect} />
                            </SheetTitle>

                            <SheetDescription>
                                Detailed information about the selected network device.
                            </SheetDescription>
                            <SheetDescription className='text-primary'>
                                Last status update {formatDateTime(row.statusUpdatedAt)}
                            </SheetDescription>
                        </SheetHeader>

                        <div className="mt-2 space-y-6 px-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div className="space-y-1">
                                    <Label className="text-xs text-muted-foreground">
                                        Device Name
                                    </Label>
                                    <p className="text-sm font-medium">
                                        {row.name}
                                    </p>
                                </div>

                                <div className="space-y-1">
                                    <Label className="text-xs text-muted-foreground">
                                        Device Type
                                    </Label>
                                    <p className="text-sm font-medium">
                                        {row.type}
                                    </p>
                                </div>
                                <div className="space-y-1">
                                    <Label className="text-xs text-muted-foreground">
                                        IP Address
                                    </Label>
                                    <p className="font-mono text-sm">
                                        {row.ipAddress}
                                    </p>
                                </div>

                                <div className="space-y-1">
                                    <Label className="text- text-muted-foreground">
                                        MAC Address
                                    </Label>
                                    <p className="font-mono text-sm">
                                        {row.macAddress}
                                    </p>
                                </div>
                                <div className="space-y-1">
                                    <Label className="text-xs text-muted-foreground">
                                        Room Number
                                    </Label>
                                    <p className="text-sm font-medium">
                                        {row.roomNumber}
                                    </p>
                                </div>
                                <div className="space-y-1">
                                    <Label className="text-xs text-muted-foreground">
                                        Room Number
                                    </Label>
                                    <p className="text-sm font-medium">
                                        {row.roomNumber}
                                    </p>
                                </div>
                            </div>

                        </div>
                        {row.description?.trim() && (
                            <div className="space-y-1 px-4">
                                <Label className="text-xs text-muted-foreground">
                                    Description
                                </Label>
                                <p className="text-sm font-medium leading-relaxed">
                                    {row.description}
                                </p>
                            </div>
                        )}
                    </SheetContent>
                </Sheet>
                <ChangeData row={row} handleFetchData={handleFetchData} />
                <DropdownMenuItem >
                    <Link href={`/port-forward?listen-ip=${row.ipAddress}`} className='size-full'>
                        Port Forward
                    </Link>
                </DropdownMenuItem>
                <AlertDialog>
                    <AlertDialogTrigger asChild>
                        <Button variant="ghost" className='text-destructive w-full flex justify-start gap-0 py-0 px-2 hover:text-destructive'>
                            Delete
                        </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                        <AlertDialogHeader>
                            <AlertDialogTitle>
                                Delete device {row.name}?
                            </AlertDialogTitle>
                            <AlertDialogDescription>
                                This action cannot be undone. This will permanently remove this device from the system and all related configurations may be lost.
                            </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                            <AlertDialogCancel>Keep Device</AlertDialogCancel>
                            <AlertDialogAction variant="destructive" onClick={() => handleDeleteDevice()}>Delete Device</AlertDialogAction>
                        </AlertDialogFooter>
                    </AlertDialogContent>
                </AlertDialog>
            </DropdownMenuContent>
        </DropdownMenu>
    );
}

export default ActionsColumns;