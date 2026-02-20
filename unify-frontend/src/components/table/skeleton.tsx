import { Skeleton } from "../ui/skeleton";
import { TableCell, TableRow } from "../ui/table";

const TableRowSkeleton = ({ columns }: { columns: number }) => (
    <TableRow>
        {Array.from({ length: columns }).map((_, i) => (
            <TableCell key={i}>
                <Skeleton className="h-4 w-full" />
            </TableCell>
        ))}
    </TableRow>
)

export default TableRowSkeleton;