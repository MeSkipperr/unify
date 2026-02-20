// schemas/user.schema.ts
import { DeviceType } from "../types";
import { isValidIPv4 } from "@/utils/ipv4";
import { isValidMacAddress } from "@/utils/macAddress";
import { z } from "zod";

export const DeviceSchemas = z.object({
    name: z
        .string()
        .min(3, "Name must be at least 3 characters"),

    ipAddress: z
        .string()
        .refine(
            isValidIPv4,
            "Invalid IPv4 address"
        ),

    macAddress: z
        .string()
        .refine(
            isValidMacAddress,
            "Invalid MAC address format"
        ),

    roomNumber: z
        .string()
        .optional()
        .refine(
            val => !val || val.length >= 2,
            "Room number must be at least 2 characters"
        ),

    description: z
        .string()
        .min(5, "Description must be at least 5 characters"),

    deviceProduct: z
        .string()
        .min(5, "Devices Product must be at least 3 characters"),

    type: z
        .nativeEnum(DeviceType)
        .refine(val => val !== undefined, {
            message: "Device type is required",
        }),

});

export type UserFormValues = z.infer<typeof DeviceSchemas>;
