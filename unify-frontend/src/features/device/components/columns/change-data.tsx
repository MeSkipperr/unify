"use client"
import {
    Sheet,
    SheetClose,
    SheetContent,
    SheetDescription,
    SheetFooter,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from "@/components/ui/sheet"
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { useState } from "react";
import { handleIPv4Input } from "@/utils/ipv4";
import { handleMacAddressInput } from "@/utils/macAddress";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller } from "react-hook-form";
import { toast } from "sonner";
import { Spinner } from "@/components/ui/spinner";
import { Device, DeviceType } from "../../types";
import { DeviceSchemas, UserFormValues } from "../../schemas/device.schema";
import { changeData } from "../../api/device.api";
import { Button } from "@/components/ui/button";
import { Label } from "@radix-ui/react-label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { DeviceStatus } from "@/components/status";
import { formatDateTime } from "@/utils/time";

type ChangeDataProps = {
    handleFetchData: () => void
    row: Device
}


const NewDataTable = ({ row, handleFetchData }: ChangeDataProps) => {
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const [isUnsavedDialogOpen, setIsUnsavedDialogOpen] = useState(false);
    const [isLoading, setIsLoading] = useState<boolean>(false);

    const {
        register,
        handleSubmit,
        reset,
        control,
        watch,
        formState: { errors },
    } = useForm<UserFormValues>({
        resolver: zodResolver(DeviceSchemas),
        defaultValues: {
            name: row.name,
            description: row.description,
            ipAddress: row.ipAddress,
            macAddress: row.macAddress,
            roomNumber: row.roomNumber,
            type: (row.type as DeviceType) ?? undefined,
            deviceProduct: row.deviceProduct
        },
    });

    const watchedValues = watch();

    const hasAnyValue = (): boolean => {
        return Object.values(watchedValues).some(
            (value) => typeof value === "string" && value.trim() !== ""
        );
    };

    const hasChanges = (): boolean => {
        return Object.keys(watchedValues).some((key) => {
            const field = key as keyof typeof watchedValues
            return watchedValues[field] !== row[field]
        })
    }


    const handlerClose = () => {
        if (hasAnyValue() && hasChanges()) {
            setIsUnsavedDialogOpen(true);
            return;
        }
        setIsOpen(false)
    }

    const handleDiscardChange = () => {
        setIsUnsavedDialogOpen(false);
        setIsOpen(false);
        reset()
    };


    const onSubmit = async (data: UserFormValues) => {
        setIsLoading(true)

        try {
            await changeData(row.id, data)

            reset()

            setIsOpen(false)

            await handleFetchData()

            toast.success("Device has been updated successfully!", { position: "bottom-right" })
        } catch (err) {
            toast.error("Failed to update device. Please try again.", { position: "bottom-right" })
        } finally {
            setIsLoading(false)
        }
    }




    return (
        <Sheet open={isOpen}>
            <Dialog open={isUnsavedDialogOpen}>
                <DialogContent showCloseButton={false}>
                    <DialogHeader>
                        <DialogTitle>Unsaved Changes</DialogTitle>
                        <DialogDescription>
                            You have unsaved changes in this form.
                            If you close it now, all changes will be lost.
                        </DialogDescription>
                    </DialogHeader>

                    <DialogFooter className="gap-2">
                        <DialogClose asChild>
                            <Button variant="default"
                                onClick={() => setIsUnsavedDialogOpen(false)}>
                                Continue Editing
                            </Button>
                        </DialogClose>

                        <DialogClose asChild>
                            <Button variant="destructive"
                                onClick={() => handleDiscardChange()}
                            >
                                Discard Changes
                            </Button>
                        </DialogClose>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <SheetTrigger asChild>
                <Button onClick={() => setIsOpen(true)} variant="ghost" className='w-full flex justify-start gap-0 py-0 px-2'>Edit</Button>
            </SheetTrigger>
            <form onSubmit={handleSubmit(onSubmit)} >

                <SheetContent className="sm:max-w-md"
                    showCloseButton={false}
                    onInteractOutside={(event) => handlerClose()}>
                    <SheetHeader>
                        <SheetTitle className="flex items-center gap-4">
                            Device Details
                            <DeviceStatus isConnect={row.isConnect} />
                        </SheetTitle>

                        <SheetDescription>
                            View and manage configuration details, network identity, and operational status of this device.
                        </SheetDescription>

                        <SheetDescription className="text-primary">
                            Last status update: {formatDateTime(row.statusUpdatedAt)}
                        </SheetDescription>
                    </SheetHeader>
                    <div className="mt-2 space-y-6 px-4">
                        <div className="space-y-1">
                            <Label className="text-xs text-muted-foreground">
                                Device Name
                            </Label>
                            <Input
                                placeholder="DPSCY-..."
                                className="h-9 text-sm"
                                required
                                {...register("name")}
                            />
                            {errors.name && (
                                <p className="text-red-500 text-sm">
                                    {errors.name.message}
                                </p>
                            )}
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-1">
                                <Controller
                                    control={control}
                                    name="type"
                                    render={({ field }) => (
                                        <>
                                            <Label className="text-xs text-muted-foreground">
                                                Device Type
                                            </Label>

                                            <Select
                                                onValueChange={field.onChange}
                                                value={field.value}
                                            >
                                                <SelectTrigger className="text-sm w-full">
                                                    <SelectValue placeholder="Select device type" />
                                                </SelectTrigger>

                                                <SelectContent>
                                                    <SelectGroup>
                                                        {Object.values(DeviceType).map((type) => (
                                                            <SelectItem key={type} value={type}>
                                                                {type.replace("-", " ").toUpperCase()}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectGroup>
                                                </SelectContent>
                                            </Select>
                                        </>
                                    )}
                                />
                                {errors.type && (
                                    <p className="text-red-500 text-sm">
                                        {errors.type.message}
                                    </p>
                                )}
                            </div>

                            <div className="space-y-1">
                                <Label className="text-xs text-muted-foreground">
                                    IP Address
                                </Label>
                                <Controller
                                    name="ipAddress"
                                    control={control}
                                    render={({ field }) => (
                                        <Input
                                            required
                                            {...field}
                                            placeholder="xxx.xxx.xxx.xxx"
                                            className="h-9 text-sm font-mono"
                                            onChange={(e) =>
                                                handleIPv4Input(e.target.value, field.onChange)
                                            }
                                        />
                                    )}
                                />
                                {errors.ipAddress && (
                                    <p className="text-red-500 text-sm">
                                        {errors.ipAddress.message}
                                    </p>
                                )}
                            </div>

                            <div className="space-y-1">
                                <Label className="text-xs text-muted-foreground">
                                    MAC Address
                                </Label>
                                <Controller
                                    name="macAddress"
                                    control={control}
                                    render={({ field }) => (
                                        <Input
                                            {...field}
                                            placeholder="00:00:00:00:00:00"
                                            className="h-9 text-sm font-mono"
                                            required
                                            onChange={(e) =>
                                                handleMacAddressInput(e.target.value, field.onChange)
                                            }
                                        />
                                    )}
                                />
                                {errors.macAddress && (
                                    <p className="text-red-500 text-sm">
                                        {errors.macAddress.message}
                                    </p>
                                )}
                            </div>

                            <div className="space-y-1">
                                <Label className="text-xs text-muted-foreground">
                                    Room Number
                                </Label>
                                <Controller
                                    name="roomNumber"
                                    control={control}
                                    render={({ field }) => (
                                        <Input
                                            {...field}
                                            placeholder="xxxx"
                                            className="h-9 text-sm"
                                            onChange={(e) =>
                                                field.onChange(
                                                    e.target.value.toUpperCase().trim()
                                                )
                                            }
                                        />
                                    )}
                                />
                            </div>
                            <div className="space-y-1">
                                <Label className="text-xs text-muted-foreground">
                                    Product
                                </Label>
                                <Controller
                                    name="deviceProduct"
                                    control={control}
                                    render={({ field }) => (
                                        <Input
                                            {...field}
                                            placeholder="Product..."
                                            className="h-9 text-sm"
                                        />
                                    )}
                                />
                            </div>
                        </div>
                    </div>
                    <div className="space-y-1 px-4">
                        <Label className="text-xs text-muted-foreground">
                            Description
                        </Label>

                        <Textarea
                            className="text-sm min-h-[100px]"
                            placeholder="Location at .."
                            {...register("description")}
                        />
                        {errors.description && (
                            <p className="text-red-500 text-sm">
                                {errors.description.message}
                            </p>
                        )}
                    </div>

                    <SheetFooter>
                        <Button
                            type="submit"
                            onClick={
                                handleSubmit(onSubmit)
                            }

                            disabled={isLoading || !hasChanges()}
                        >
                            {isLoading ?
                                <Label className="flex gap-4 items-center">
                                    <Spinner /> Saving...
                                </Label>
                                :
                                "Save changes"
                            }
                        </Button>
                        <SheetClose asChild>
                            <Button
                                onClick={(event) => handlerClose()}
                                variant="outline"
                                disabled={isLoading}
                            >
                                Close
                            </Button>
                        </SheetClose>
                    </SheetFooter>
                </SheetContent>
            </form>
        </Sheet>
    );
}

export default NewDataTable;