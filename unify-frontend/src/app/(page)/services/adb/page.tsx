"use client"

import AdbTable from "@/features/adb/components/table"

const AdbPage = () => {
    

    return (
        <div className="h-[85dvh]">
            <AdbTable serviceType="manual"  addNewData />
        </div>
        
    );
}

export default AdbPage;