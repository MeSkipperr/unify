import api from "@/api";
import { TableQuery } from "@/components/table/types";
import { format } from "date-fns";
import { AdbCommand } from "../types";
import { normalizeIPv4 } from "@/utils/ipv4";

export type AdbQuery = Pick<TableQuery, "page" | "pageSize" | "search"> & {
  typeServices?: string[];
  sort?: string[];
  date?: Date | string;
};

export const getAdbResults = async (filter?: AdbQuery) => {
  const params: Record<
    string,
    string | number | boolean[] | string[] | undefined
  > = {};

  if (filter?.typeServices?.length) {
    params.typeServices = filter.typeServices;
  }

  if (filter?.sort?.length) {
    params.sort = filter.sort.join(",");
  }

  if (filter?.search && filter.search.trim() !== "") {
    params.search = filter.search;
  }

  if (filter?.date) {
    params.date =
      filter.date instanceof Date
        ? format(filter.date, "yyyy-MM-dd")
        : filter.date;
  }

  params.page = filter?.page ?? 1;
  params.pageSize = filter?.pageSize ?? 50;

  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const res = await api.get("/api/services/adb", {
    headers: {
      "X-Timezone": timezone,
    },
    params,
  });
  return res.data;
};

type createRunningAdbProps = {
  ipAddress: string;
  port: number;
  command: AdbCommand;
  name: string;
};

export const createRunningAdb = async (payload: createRunningAdbProps) => {
  const normalizedPayload: createRunningAdbProps = {
    ...payload,
    name: payload.name.trim(),
    ipAddress: normalizeIPv4(payload.ipAddress),
  };
  const res = await api.post("/api/services/adb", normalizedPayload);

  return res.data;
};
