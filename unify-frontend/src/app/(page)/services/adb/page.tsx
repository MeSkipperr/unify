"use client"

import AdbTable from "@/features/adb/components/table"
import { Suspense } from "react";

const AdbPage = () => {


    return (
        <div className="h-[85dvh]">
            <Suspense fallback={<div>Loading...</div>}>
                <AdbTable serviceType="manual" addNewData />
            </Suspense>
        </div>

    );
}

export default AdbPage;