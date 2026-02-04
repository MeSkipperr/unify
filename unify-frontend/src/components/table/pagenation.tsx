"use client"

import {
    Pagination,
    PaginationContent,
    PaginationEllipsis,
    PaginationItem,
    PaginationLink,
    PaginationNext,
    PaginationPrevious,
} from "@/components/ui/pagination"
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { Field, FieldLabel } from "@/components/ui/field"
import { generatePages, getLastPage } from "@/features/device/utils/page"
import React from "react"

type PagenationTableProps = {
    totalData: number,
    pageQuery: number
    setPageQuery: React.Dispatch<React.SetStateAction<number>>
    pageSizeQuery: number
    setPageSizeQuery: React.Dispatch<React.SetStateAction<number>>
    rowPageList?: number[]
}

const PagenationTable = ({
    totalData = 0,
    pageQuery,
    setPageQuery,
    pageSizeQuery,
    setPageSizeQuery,
    rowPageList = [15, 25, 50, 100]
}: PagenationTableProps) => {

    const lastPage = getLastPage(totalData, pageSizeQuery)
    const pages = generatePages(pageQuery, lastPage)

    return (
        <div className="flex w-full shrink-0 items-center justify-between">
            <Pagination className="justify-start">
                <PaginationContent>

                    {/* PREVIOUS */}
                    <PaginationItem>
                        <PaginationPrevious
                            onClick={() => {
                                if (pageQuery > 1) setPageQuery(pageQuery - 1)
                            }}
                        />
                    </PaginationItem>

                    {/* PAGE NUMBERS */}
                    {pages.map((page, index) => (
                        <PaginationItem key={index}>
                            {page === "ellipsis" ? (
                                <PaginationEllipsis />
                            ) : (
                                <PaginationLink
                                    isActive={page === pageQuery}
                                    onClick={() => setPageQuery(page)}
                                >
                                    {page}
                                </PaginationLink>
                            )}
                        </PaginationItem>
                    ))}

                    {/* NEXT */}
                    <PaginationItem>
                        <PaginationNext
                            onClick={() => {
                                if (pageQuery < lastPage) setPageQuery(pageQuery + 1)
                            }}
                        />
                    </PaginationItem>

                </PaginationContent>
            </Pagination>

            {/* ROWS PER PAGE */}
            <Field orientation="horizontal" className="w-fit items-end">
                <FieldLabel htmlFor="select-rows-per-page">
                    Rows per page
                </FieldLabel>
                <Select
                    value={pageSizeQuery.toString()}
                    onValueChange={(val) => {
                        setPageSizeQuery(Number(val))
                        setPageQuery(1) // reset ke page 1
                    }}
                >
                    <SelectTrigger className="w-20" id="select-rows-per-page">
                        <SelectValue />
                    </SelectTrigger>
                    <SelectContent align="start">
                        <SelectGroup>
                            {rowPageList.map((val) => (
                                <SelectItem
                                    key={val}
                                    value={val.toString()}
                                >
                                    {val}
                                </SelectItem>
                            ))}
                        </SelectGroup>
                    </SelectContent>
                </Select>
            </Field>
        </div>
    )
}

export default PagenationTable;