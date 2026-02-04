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
import { Button } from "../../../components/ui/button";
import { Textarea } from "../../../components/ui/textarea";
import { Label } from "../../../components/ui/label";
import { Input } from "../../../components/ui/input";
import { Plus } from "lucide-react";
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
    DialogTrigger,
} from "@/components/ui/dialog"
import { useState } from "react";
import { handleIPv4Input } from "@/utils/ipv4";
import { handleMacAddressInput } from "@/utils/macAddress";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { DeviceSchemas, type UserFormValues } from "../schemas/device.schema";
import { Controller } from "react-hook-form";
import { DeviceType } from "../types";
import { addDevice } from "../api/device.api";

const NewDataTable = () => {
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const [isUnsavedDialogOpen, setIsUnsavedDialogOpen] = useState(false);

    const {
        register,
        handleSubmit,
        setValue,
        reset,
        control,
        watch,
        formState: { errors, isSubmitting },
    } = useForm<UserFormValues>({
        resolver: zodResolver(DeviceSchemas),
        defaultValues: {
            name: "",
            description: "",
            ipAddress: "",
            macAddress: "",
            roomNumber: "",
            type: undefined,
        },
    });

    const watchedValues = watch();

    const hasAnyValue = (): boolean => {
        return Object.values(watchedValues).some(
            (value) => typeof value === "string" && value.trim() !== ""
        );
    };


    const handlerClose = () => {
        if (hasAnyValue()) {
            setIsUnsavedDialogOpen(true);
            return;
        }
        setIsOpen(false)
    }

    const handleDiscardChange = () => {
        reset();
        setIsUnsavedDialogOpen(false);
        setIsOpen(false);
    };


    const onSubmit = async (data: UserFormValues) => {
        try {
            await addDevice(data);
            reset();
            setIsOpen(false);
        } catch (err) {
            console.error(err);
        }
    };


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
                <Button onClick={() => setIsOpen(true)} >
                    <Plus />
                    Add New
                </Button>
            </SheetTrigger>
            <form onSubmit={handleSubmit(onSubmit)} >

                <SheetContent className="sm:max-w-md"
                    showCloseButton={false}
                    onInteractOutside={(event) => handlerClose()}>
                    <SheetHeader>
                        <SheetTitle>
                            Add New Device
                        </SheetTitle>
                        <SheetDescription>
                            Enter device details including name, type, IP, MAC, room number, and description.
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
                        <Button type="submit" onClick={handleSubmit(onSubmit)}>Save changes</Button>
                        <SheetClose asChild>
                            <Button
                                onClick={(event) => handlerClose()}
                                variant="outline">Close</Button>
                        </SheetClose>
                    </SheetFooter>
                </SheetContent>
            </form>
        </Sheet>
    );
}

export default NewDataTable;