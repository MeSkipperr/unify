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
} from "@/components/ui/dialog"
import { useEffect, useState } from "react";
import { handleIPv4Input, isValidIPv4 } from "@/utils/ipv4";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller } from "react-hook-form";
import { toast } from "sonner";
import { Spinner } from "@/components/ui/spinner";
import { PortForwardFormValues, PortForwardSchemas } from "../schemas/port-forward.schema";
import { EXPIRE_OPTIONS } from "../types";
import { createPortForward, createPortForwardPayload } from "../api/port-forward.api";
import { convertToDate } from "../utils/date";
import { usePathname, useRouter, useSearchParams } from "next/navigation";

type NewDataProps = {
    handleFetchData: () => Promise<void>
}


const NewDataTable = ({ handleFetchData }: NewDataProps) => {
    const searchParams = useSearchParams();
    const router = useRouter();
    const pathname = usePathname();

    const [isOpen, setIsOpen] = useState<boolean>(false);
    const [isUnsavedDialogOpen, setIsUnsavedDialogOpen] = useState(false);
    const [isLoading, setIsLoading] = useState<boolean>(false);



    const {
        register,
        handleSubmit,
        reset,
        setValue,
        control,
        watch,
        formState: { errors },
    } = useForm<PortForwardFormValues>({
        resolver: zodResolver(PortForwardSchemas),
        defaultValues: {
            destIp: "",
            destPort: 80,
            listenIp: "", // ambil dari param ?listen-ip=192.168.160.184
            protocol: "tcp",
            ruleComment: ""
        },
    });

    useEffect(() => {
        const listenIpParam = searchParams.get("listen-ip");

        if (!listenIpParam) return;

        if (listenIpParam && isValidIPv4(listenIpParam)) {
            setValue("listenIp", listenIpParam); // isi form
            setIsOpen(true); // buka modal / section
        } else {
            const params = new URLSearchParams(searchParams.toString());
            params.delete("listen-ip");

            router.replace(`${pathname}?${params.toString()}`);
        }

    }, [searchParams,pathname,router,setValue]);

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


    const onSubmit = async (data: PortForwardFormValues) => {
        setIsLoading(true)
        const payload: createPortForwardPayload = {
            ...data,
            expiresAt: convertToDate(data.expiresAt).toISOString(),
        };
        try {
            await createPortForward(payload)
            reset()

            setIsOpen(false)

            await handleFetchData()

            toast.success("Port forward rule created successfully.", {
                position: "bottom-right",
            });

        } catch (err) {
            console.log(err)
            toast.success("Port forwarding rule is now active and ready to use.", {
                position: "bottom-right",
            });
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
                <Button onClick={() => setIsOpen(true)} >
                    <Plus />
                    Add New
                </Button>
            </SheetTrigger>
            <form onSubmit={handleSubmit(onSubmit)} >

                <SheetContent className="sm:max-w-md"
                    showCloseButton={false}
                    onInteractOutside={() => handlerClose()}>
                    <SheetHeader>
                        <SheetTitle>
                            Add New Device
                        </SheetTitle>
                        <SheetDescription>
                            Enter device details including name, type, IP, MAC, room number, and description.
                        </SheetDescription>
                    </SheetHeader>
                    <div className="mt-2 space-y-6 px-4">
                        <div className="grid grid-cols-2 gap-4">

                            <div className="space-y-1">
                                <Label className="text-xs text-muted-foreground">
                                    Listen Ip
                                </Label>
                                <Controller
                                    name="listenIp"
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
                                {errors.listenIp && (
                                    <p className="text-red-500 text-sm">
                                        {errors.listenIp.message}
                                    </p>
                                )}
                            </div>
                            <div className="space-y-1">

                                <Controller
                                    control={control}
                                    name="expiresAt"
                                    render={({ field }) => (
                                        <>
                                            <Label className="text-xs text-muted-foreground">
                                                Experis At
                                            </Label>

                                            <Select
                                                onValueChange={field.onChange}
                                                value={field.value}
                                                defaultValue={field.value}
                                            >
                                                <SelectTrigger className="text-sm w-full">
                                                    <SelectValue placeholder="Select expiration time" />
                                                </SelectTrigger>

                                                <SelectContent>
                                                    <SelectGroup>
                                                        {EXPIRE_OPTIONS.map((time) => (
                                                            <SelectItem key={time} value={time}>
                                                                {time}
                                                            </SelectItem>
                                                        ))}

                                                    </SelectGroup>
                                                </SelectContent>
                                            </Select>
                                        </>
                                    )}
                                />
                                {errors.expiresAt && (
                                    <p className="text-red-500 text-sm">
                                        {errors.expiresAt.message}
                                    </p>
                                )}
                            </div>
                            <div className="space-y-1">
                                <Label className="text-xs text-muted-foreground">
                                    Destination Ip
                                </Label>
                                <Controller
                                    name="destIp"
                                    control={control}
                                    render={({ field }) => (
                                        <Input
                                            {...field}
                                            placeholder="xxx:xxx:xxx:xxx"
                                            className="h-9 text-sm font-mono"
                                            required
                                            onChange={(e) =>
                                                handleIPv4Input(e.target.value, field.onChange)
                                            }
                                        />
                                    )}
                                />
                                {errors.destIp && (
                                    <p className="text-red-500 text-sm">
                                        {errors.destIp.message}
                                    </p>
                                )}
                            </div>

                            <div className="space-y-1">
                                <Label className="text-xs text-muted-foreground">
                                    Destination Port
                                </Label>
                                <Controller
                                    name="destPort"
                                    control={control}
                                    render={({ field }) => (
                                        <Input
                                            {...field}
                                            type="number"
                                            placeholder="xxxx"
                                            className="h-9 text-sm"
                                            onChange={(e) => field.onChange(Number(e.target.value))}
                                        />
                                    )}
                                />
                                {errors.destPort && (
                                    <p className="text-red-500 text-sm">
                                        {errors.destPort.message}
                                    </p>
                                )}
                            </div>
                        </div>
                    </div>
                    <div className="mt-2 space-y-6 px-4">
                        <div className="space-y-1">

                            <Controller
                                control={control}
                                name="protocol"
                                render={({ field }) => (
                                    <>
                                        <Label className="text-xs text-muted-foreground">
                                            Protocol
                                        </Label>

                                        <Select
                                            onValueChange={field.onChange}
                                            value={field.value}
                                            defaultValue={field.value}
                                        >
                                            <SelectTrigger className="text-sm w-full">
                                                <SelectValue placeholder="Select device type" />
                                            </SelectTrigger>

                                            <SelectContent>
                                                <SelectGroup>
                                                    {Object.values(["tcp", "udp"]).map((type) => (
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
                            {errors.protocol && (
                                <p className="text-red-500 text-sm">
                                    {errors.protocol.message}
                                </p>
                            )}
                        </div>
                    </div>

                    <div className="space-y-1 px-4">
                        <Label className="text-xs text-muted-foreground">
                            Comment
                        </Label>

                        <Textarea
                            className="text-sm min-h-[100px]"
                            placeholder="Device..."
                            {...register("ruleComment")}
                        />
                        {errors.ruleComment && (
                            <p className="text-red-500 text-sm">
                                {errors.ruleComment.message}
                            </p>
                        )}
                    </div>

                    <SheetFooter>
                        <Button
                            type="submit"
                            onClick={
                                handleSubmit(onSubmit)
                            }
                            disabled={isLoading}
                        >
                            {isLoading ?
                                <Label className="flex gap-4 items-center">
                                    <Spinner /> Saving...
                                </Label>
                                :
                                "Save Devices"
                            }
                        </Button>
                        <SheetClose asChild>
                            <Button
                                onClick={() => handlerClose()}
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