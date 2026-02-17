"use client";

import {
    Sheet,
    SheetClose,
    SheetContent,
    SheetDescription,
    SheetFooter,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from "@/components/ui/sheet";
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
} from "@/components/ui/select";
import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { useEffect, useRef, useState, useCallback } from "react";
import { handleIPv4Input } from "@/utils/ipv4";
import { Controller, useForm } from "react-hook-form";
import { toast } from "sonner";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { AdbSchemas, UserFormValues } from "../schemas/adb.schema";
import { AdbCommand } from "../types";
import { createRunningAdb } from "../api/adb-result.api";
import { useSSE } from "@/hooks/useSSE";

import {
    CodeBlock,
    CodeBlockBody,
    CodeBlockContent,
    CodeBlockItem,
} from "@/components/kibo-ui/code-block";
import type { BundledLanguage } from "shiki";

type NewDataProps = {
    handleFetchData: () => Promise<void>;
};

type RunningAdbSseEvent = {
    type: "running-adb";
    data: {
        id: string;
        result: string;
    };
};

const NewDataTable = ({ handleFetchData }: NewDataProps) => {
    const [isOpen, setIsOpen] = useState(false);
    const [isUnsavedDialogOpen, setIsUnsavedDialogOpen] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [showResult, setShowResult] = useState(false);
    const [jobId, setJobId] = useState<string | null>(null);

    const [codes, setCodes] = useState([
        { language: "bash", filename: "", code: "" },
    ]);

    const isOpenRef = useRef(isOpen);

    useEffect(() => {
        isOpenRef.current = isOpen;
    }, [isOpen]);

    const handleChangeCode = (index: number, value: string) => {
        setCodes((prev) =>
            prev.map((item, i) => (i === index ? { ...item, code: value } : item))
        );
    };

    const {
        register,
        handleSubmit,
        reset,
        control,
        watch,
        formState: { errors },
    } = useForm<UserFormValues>({
        resolver: zodResolver(AdbSchemas),
        defaultValues: {
            name: "",
            ipAddress: "",
            port: 5555,
            command: AdbCommand.Reboot,
        },
    });

    const watchedValues = watch();

    const hasAnyValue = () =>
        Object.values(watchedValues).some(
            (value) => typeof value === "string" && value.trim() !== ""
        );

    const handleClose = () => {
        if (hasAnyValue()) {
            setIsUnsavedDialogOpen(true);
            return;
        }
        setIsOpen(false);
    };

    const handleDiscardChange = () => {
        reset();
        setIsUnsavedDialogOpen(false);
        setIsOpen(false);
    };

    const handleCloseResult = async () => {
        reset();
        setShowResult(false);
        setIsUnsavedDialogOpen(false);
        setIsOpen(false);
        setJobId(null);
        await handleFetchData();
    };

    const handleSseMessage = useCallback(
        (result: RunningAdbSseEvent) => {
            if (result.type !== "running-adb") return;
            if (!isOpenRef.current) return;
            if (!jobId) return;

            if (result.data.id === jobId) {
                handleChangeCode(0, result.data.result);
                setShowResult(true);
                setIsLoading(false);
                stop();
            }
        },
        [jobId]
    );

    const { start, stop } = useSSE<RunningAdbSseEvent>({
        url: "/events/services",
        onMessage: handleSseMessage,
    });

    const onSubmit = async (data: UserFormValues) => {
        try {
            setIsLoading(true);
            stop();

            const res = await createRunningAdb(data);
            const newJobId = res.data.id;

            setJobId(newJobId);
            start();
        } catch  {
            toast.error("Failed to run adb. Please try again.", {
                position: "bottom-right",
            });
            setIsLoading(false);
        }
    };

    return (
        <Sheet open={isOpen}>
            <Dialog open={showResult}>
                <DialogContent showCloseButton={false}>
                    <DialogHeader>
                        <DialogTitle>Result</DialogTitle>
                        <DialogDescription />
                    </DialogHeader>

                    <CodeBlock data={codes} defaultValue={codes[0].language}>
                        <CodeBlockBody>
                            {(item) => (
                                <CodeBlockItem
                                    key={item.language}
                                    lineNumbers={false}
                                    value={item.language}
                                >
                                    <CodeBlockContent
                                        language={item.language as BundledLanguage}
                                    >
                                        {item.code}
                                    </CodeBlockContent>
                                </CodeBlockItem>
                            )}
                        </CodeBlockBody>
                    </CodeBlock>

                    <DialogFooter>
                        <DialogClose asChild>
                            <Button onClick={handleCloseResult}>Close</Button>
                        </DialogClose>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <Dialog open={isUnsavedDialogOpen && !showResult}>
                <DialogContent showCloseButton={false}>
                    <DialogHeader>
                        <DialogTitle>Unsaved Changes</DialogTitle>
                        <DialogDescription>
                            You have unsaved changes. If you leave now, your changes will be
                            lost.
                        </DialogDescription>
                    </DialogHeader>

                    <DialogFooter className="gap-2">
                        <DialogClose asChild>
                            <Button onClick={() => setIsUnsavedDialogOpen(false)}>
                                Continue Editing
                            </Button>
                        </DialogClose>
                        <DialogClose asChild>
                            <Button variant="destructive" onClick={handleDiscardChange}>
                                Discard Changes
                            </Button>
                        </DialogClose>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <SheetTrigger asChild>
                <Button onClick={() => setIsOpen(true)}>
                    <Plus />
                    Add New
                </Button>
            </SheetTrigger>

            <form onSubmit={handleSubmit(onSubmit)}>
                <SheetContent
                    className="sm:max-w-md"
                    showCloseButton={false}
                    onInteractOutside={handleClose}
                >
                    <SheetHeader>
                        <SheetTitle>Running Adb</SheetTitle>
                        <SheetDescription>
                            Execute the selected ADB command on the target device.
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
                                <Label className="text-xs text-muted-foreground">Port</Label>
                                <Controller
                                    name="port"
                                    control={control}
                                    render={({ field }) => (
                                        <Input
                                            {...field}
                                            type="number"
                                            className="h-9 text-sm"
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
                                                <SelectValue placeholder="Select command" />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectGroup>
                                                    {Object.values(AdbCommand).map((cmd) => (
                                                        <SelectItem key={cmd} value={cmd}>
                                                            {cmd.replace("-", " ").toUpperCase()}
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
                        <Button type="submit" disabled={isLoading}>
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
                                onClick={handleClose}
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
};

export default NewDataTable;
