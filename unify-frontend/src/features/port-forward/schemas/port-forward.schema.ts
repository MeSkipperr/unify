// schemas/user.schema.ts
import { isValidIPv4 } from "@/utils/ipv4";
import { z } from "zod";
import { EXPIRE_OPTIONS } from "../types";

export const PortForwardSchemas = z.object({
  listenIp: z.string().nonempty("Listen IP is required").refine(isValidIPv4, {
    message: "Invalid IPv4 address",
  }),

  destIp: z
    .string()
    .nonempty("Destination IP is required")
    .refine(isValidIPv4, {
      message: "Invalid IPv4 address",
    }),

  destPort: z
    .number()
    .int("Destination port must be an integer")
    .min(1, "Destination port must be between 1 and 65535")
    .max(65535, "Destination port must be between 1 and 65535"),

  ruleComment: z
    .string()
    .nonempty("Description is required")
    .min(5, "Description must be at least 5 characters"),

  protocol: z.enum(["tcp", "udp"]).refine((val) => !!val, {
    message: "Protocol is required",
  }),
  expiresAt: z.enum(EXPIRE_OPTIONS, {
    message: "Expiration time is required",
  }),
  
});

export type PortForwardFormValues = z.infer<typeof PortForwardSchemas>;
