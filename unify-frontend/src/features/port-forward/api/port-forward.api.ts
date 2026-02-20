import api from "@/api";
import { TableQuery } from "@/components/table/types";

export type PortForwardQuery = Pick<
  TableQuery,
  "page" | "pageSize" | "search"
> & {
  protocol?: string[];
  status?: string[];
  sort?: string[];
};

export const getPortForward = async (filter?: PortForwardQuery) => {
  const params: Record<
    string,
    string | number | boolean | string[] | undefined
  > = {};

  if (filter?.protocol?.length) {
    params.protocol = filter.protocol;
  }

  if (filter?.status?.length) {
    params.status = filter.status;
  }

  if (filter?.sort?.length) {
    params.sort = filter.sort.join(",");
  }

  if (filter?.search && filter.search.trim() !== "") {
    params.search = filter.search;
  }

  params.page = filter?.page ?? 1;
  params.pageSize = filter?.pageSize ?? 50;

  const res = await api.get("/api/services/port-forward", { params });
  return res.data;
};

export type createPortForwardPayload = {
  listenIp: string;
  expiresAt: string;
  destIp: string;
  protocol: string;
  destPort: number;
  ruleComment: string;
};

export const createPortForward = async (payload: createPortForwardPayload) => {
  const res = await api.post("/api/services/port-forward", payload);
  return res.data;
};

export const deactivatePortForward = async (id: string) => {
    const res = await api.patch(`/api/services/port-forward/${id}/deactivate`);
    return res.data;
};
  