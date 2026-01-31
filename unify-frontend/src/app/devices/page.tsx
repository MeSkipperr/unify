import DeviceTableData from "@/features/device";
import { Device } from "@/features/device/types";

const data: Device[] = [
    {
        id: "1",
        name: "Access Point Lobby",
        ipAddress: "192.168.1.10",
        macAddress: "00:00:00:00:00:00",
        roomNumber: "1001",
        isConnect: true,
        type: "access-point",
        statusUpdatedAt: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000),
        notification: false,
        description: "Server: NVR4,Channel: 29,Coverage: C-108 Room 3208 - 3201,Distribution: In Front of Room 3207 & 3208 pointed to Elevator LP-03"
    },
    {
        id: "2",
        name: "CCTV Gate",
        ipAddress: "192.168.1.20",
        macAddress: "00:00:00:00:00:00",
        roomNumber: "1002",
        isConnect: false,
        type: "cctv",
        notification: true,
        statusUpdatedAt: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000),
    },
]

const DevicesList = () => {
    return (
        <div className="">
            <DeviceTableData data={data}/>
        </div>
    );
}

export default DevicesList;