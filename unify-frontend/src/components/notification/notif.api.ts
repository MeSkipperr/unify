import api from "@/api";
import { TableQuery } from "../table/types";


export const getNotification = async (data?:  Pick<TableQuery, "page" | "pageSize" >) => {
  const params: Record<
    string,
    string | number | boolean[] | string[] | undefined
  > = {};

  params.page = data?.page ?? 1;
  params.pageSize = data?.pageSize ?? 50;

  const res = await api.get("/api/notification", { params });
  return res.data;
};
