// schemas/user.schema.ts
import { isValidIPv4 } from "@/utils/ipv4";
import { z } from "zod";
import { AdbCommand } from "../types";

export const AdbSchemas = z.object({
  name: z.string().min(3, "Name must be at least 3 characters"),

  ipAddress: z.string().refine(isValidIPv4, "Invalid IPv4 address"),

  port: z
    .number()
    .int("Port must be a whole number")
    .min(1, "Port must be between 1 and 65535")
    .max(65535, "Port must be between 1 and 65535"),

  command: z.nativeEnum(AdbCommand, {
    message: "Command is required",
  }),
});

export type UserFormValues = z.infer<typeof AdbSchemas>;
