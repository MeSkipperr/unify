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
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
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
import { useEffect, useRef, useState } from "react";
import { handleIPv4Input } from "@/utils/ipv4";
import { Controller, useForm } from "react-hook-form";
import { toast } from "sonner";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { AdbSchemas, UserFormValues } from "../schemas/adb.schema";
import { AdbCommand } from "../types";
import { createRunningAdb } from "../api/adb-result.api";

import {
    CodeBlock,
    CodeBlockBody,
    CodeBlockContent,
    CodeBlockItem,
} from "@/components/kibo-ui/code-block";
import type { BundledLanguage } from "shiki"

type NewDataProps = {
    handleFetchData: () => Promise<void>
}


const NewDataTable = ({ handleFetchData }: NewDataProps) => {
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const [isUnsavedDialogOpen, setIsUnsavedDialogOpen] = useState(false);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [showResult, setShowResult] = useState<boolean>(false)
    const [codes, setCodes] = useState([
        {
            language: "bash",
            filename: "",
            code: "",
        },
    ]);
    const isOpenRef = useRef(isOpen);
    const eventSourceRef = useRef<EventSource | null>(null);
    useEffect(() => {
        isOpenRef.current = isOpen;
    }, [isOpen]);


    const handleChangeCode = (index: number, value: string) => {
        setCodes((prev) =>
            prev.map((item, i) =>
                i === index ? { ...item, code: value } : item
            )
        );
    };


    const {
        register,
        handleSubmit,
        setValue,
        reset,
        control,
        watch,
        formState: { errors, isSubmitting },
    } = useForm<UserFormValues>({
        resolver: zodResolver(AdbSchemas),
        defaultValues: {
            name: "",
            ipAddress: "",
            command: "",
            port: 5555
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
        setIsLoading(true);

        // close SSE lama jika ada
        if (eventSourceRef.current) {
            eventSourceRef.current.close();
        }

        try {
            const res = await createRunningAdb(data);
            const jobId = res.data.id;

            const eventSource = new EventSource(
                `http://localhost:8080/events/services`
            );

            eventSourceRef.current = eventSource;

            eventSource.onmessage = (event) => {
                const result = JSON.parse(event.data);
                if (result.type !== "running-adb" || !isOpenRef.current) return;

                if (result.data.id === jobId) {
                    handleChangeCode(0, result.data.result);
                    setShowResult(true);
                    setIsLoading(false);
                    eventSource.close();
                    eventSourceRef.current = null; 
                }
            };

            eventSource.onerror = () => {
                eventSource.close();
                eventSourceRef.current = null;
                toast.error("SSE connection failed", { position: "bottom-right" });
                setIsLoading(false);
            };
        } catch (err) {
            toast.error("Failed to running adb. Please try again.", { position: "bottom-right" });
            setIsLoading(false);
        }
    };



    const handleCloseResult = async () => {
        reset()

        setShowResult(false)
        setIsUnsavedDialogOpen(false)
        setIsOpen(false)
        await handleFetchData()
    }


    return (
        <Sheet open={isOpen}>
            <Dialog open={showResult}>
                <DialogContent showCloseButton={false}>
                    <DialogHeader>
                        <DialogTitle>Result</DialogTitle>
                        <DialogDescription>
                        </DialogDescription>
                    </DialogHeader>
                    <CodeBlock data={codes} defaultValue={codes[0].language}>
                        <CodeBlockBody>
                            {(item) => (
                                <CodeBlockItem
                                    key={item.language}
                                    lineNumbers={false}
                                    value={item.language}
                                >
                                    <CodeBlockContent language={item.language as BundledLanguage}>
                                        {item.code}
                                    </CodeBlockContent>
                                </CodeBlockItem>
                            )}
                        </CodeBlockBody>
                    </CodeBlock>
                    <DialogFooter>
                        <DialogClose asChild>
                            <Button variant="default"
                                onClick={() => handleCloseResult()}>
                                Close
                            </Button>
                        </DialogClose>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <Dialog open={isUnsavedDialogOpen && !showResult}>
                <DialogContent showCloseButton={false}>
                    <DialogHeader>
                        <DialogTitle>Unsaved Changes</DialogTitle>
                        <DialogDescription>
                            You have unsaved changes. If you leave now, your changes will be lost.
                            Do you want to continue editing or discard your changes?
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
                            Running Adb
                        </SheetTitle>
                        <SheetDescription>
                            Execute the selected ADB command on the target device.
                            Please ensure the device is connected before proceeding.
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
                                    Port
                                </Label>
                                <Controller
                                    name="port"
                                    control={control}
                                    render={({ field }) => (
                                        <Input
                                            {...field}
                                            type="number"
                                            placeholder="xxxx"
                                            className="h-9 text-sm"
                                            onChange={(e) => {
                                                const value = e.target.value;
                                                field.onChange(value === "" ? "" : Number(value));
                                            }}
                                        />
                                    )}
                                />
                                {errors.port && (
                                    <p className="text-red-500 text-sm">
                                        {errors.port.message}
                                    </p>
                                )}
                            </div>
                        </div>
                        <div className="space-y-1">
                            <Controller
                                control={control}
                                name="command"
                                render={({ field }) => (
                                    <>
                                        <Label className="text-xs text-muted-foreground">
                                            Command
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
                                                    {Object.values(AdbCommand).map((type) => (
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
                            {errors.command && (
                                <p className="text-red-500 text-sm">
                                    {errors.command.message}
                                </p>
                            )}
                        </div>

                    </div>
                    <SheetFooter>
                        <Button
                            type="button"
                            onClick={handleSubmit(onSubmit)}
                            disabled={isLoading}
                        >
                            {isLoading ? (
                                <span className="flex gap-2 items-center">
                                    <Spinner />
                                    Running...
                                </span>
                            ) : (
                                "Run Command"
                            )}
                        </Button>

                        <SheetClose asChild>
                            <Button
                                onClick={handlerClose}
                                variant="outline"
                                disabled={isLoading}
                            >
                                Cancel
                            </Button>
                        </SheetClose>
                    </SheetFooter>
                </SheetContent>
            </form>
        </Sheet>
    );
}

export default NewDataTable;