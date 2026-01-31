import { columns } from "./columns";
import { dataFilter } from "./filter-data";
import { Device } from "./types";
import DataTable from "@/components/table";

type DeviceTableDataProps = {
    data : Device[]
}

const DeviceTableData = ({data}: DeviceTableDataProps) => {
    return ( 
        <DataTable data={data} filter={dataFilter} columns={columns}/>
    );
    
}

export default DeviceTableData;