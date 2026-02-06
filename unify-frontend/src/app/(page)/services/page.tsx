import { Metadata } from "next";
import ServicesTable from "./services-table";


export const metadata: Metadata = {
    title: "Services | Unify",
    description: "Manage and monitor all services in the Unify system.",
};


const ServicesPage = () => {
    return (
        <div className="w-full p-4">
            <h1 className="sm:text-xl text-2xl font-bold">Service</h1>
            <ServicesTable/>
        </div>
    );
}

export default ServicesPage;