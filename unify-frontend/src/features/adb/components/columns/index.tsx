
'use client'

import { ColumnDef } from '@tanstack/react-table'
import { AdbResult } from '../../types'
import { formatDateTime } from '@/utils/time'


import { Label } from '@radix-ui/react-label'
import {  MessageSquareText } from 'lucide-react'

import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from '@/components/ui/button'

import {
    CodeBlock,
    CodeBlockBody,
    CodeBlockContent,
    CodeBlockItem,
} from "@/components/kibo-ui/code-block";
import type { BundledLanguage } from "shiki"

export const columns: ColumnDef<AdbResult>[] =
    [
        {
            id: "number",
            header: () => (
                <Label className='w-full flex justify-center items-center'>No</Label>
            ),
            cell: ({ row }) => (
                <Label className='w-full flex justify-center items-center'>{row.original.index}</Label>
            ),
            size: 50,
        },
        {
            accessorKey: 'status',
            header: 'Status',
        },
        {
            accessorKey: 'ipAddress',
            header: 'IP Address',
        },
        {
            accessorKey: 'port',
            header: 'Port',
        },
        {
            accessorKey: 'deviceName',
            header: 'Device Name',
        },

        {
            accessorKey: 'serviceType',
            header: 'Service Type',
        },
        {
            accessorKey: 'startTime',
            header: 'Start At',
            cell: ({ row }) => (
                <span>{formatDateTime(row.original.startTime)}</span>
            )
        },
        {
            accessorKey: 'finishTime',
            header: 'Finnish At',
            cell: ({ row }) => (
                <span>{formatDateTime(row.original.finishTime)}</span>
            )
        },
        {
            accessorKey: 'result',
            header: 'Result',
            cell: ({ row }) => {
                const code = [
                    {
                        language: "bash",
                        filename: "",
                        code: row.original.result,
                    },
                ];

                return (
                    <Dialog>
                        <DialogTrigger asChild>
                            <Button variant="ghost">
                                <MessageSquareText />
                            </Button>
                        </DialogTrigger>
                        <DialogContent >
                            <DialogHeader>
                                <DialogTitle>{row.original.deviceName}</DialogTitle>
                            </DialogHeader>
                            <DialogDescription>
                                Detailed command output from the selected device.
                            </DialogDescription>
                            <CodeBlock data={code} defaultValue={code[0].language}>
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
                        </DialogContent>
                    </Dialog>
                )
            }
        },
    ]
