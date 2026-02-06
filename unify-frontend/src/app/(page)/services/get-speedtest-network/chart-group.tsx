"use client"

import { Services } from "@/types/service.type";
import { useEffect, useState } from "react";
import { getServicesByName } from "../service.api";
import SpeedNetworkChart from "./speed-network-chart";

type NetworkDetailProps = {
    name: string,
    interface: string,
    ipAddress: string
}

type ConfigProps = {
    cron: string
    network: NetworkDetailProps[]
    serverId: string[]
}

const ChartGroup = ({ data }: { data: Services<ConfigProps> }) => {
    const [isLoading, setIsLoading] = useState<boolean>(false);



    return (
        <div className="flex flex-col gap-8">
            {data.config.network.map((net) => (
                <div key={net.name} className="mb-4">
                    <p>Network: {net.name}</p>
                    <p>Interface: {net.interface}</p>
                    <p>IP: {net.ipAddress}</p>

                    <ul className="ml-4 list-disc">
                        {data.config.serverId.map((serverId) => (
                            <li key={serverId}>{serverId}</li>
                            // <SpeedNetworkChart key={}/>
                        ))}
                    </ul>
                </div>
            ))}
        </div>
    );
}

export default ChartGroup;