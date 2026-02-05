import ServicesTable from "./services-table";

const ServicesPage = () => {
    return (
        <div className="w-full p-4">
            <h1 className="sm:text-xl text-2xl font-bold">Service</h1>
            <ServicesTable/>
        </div>
    );
}

export default ServicesPage;