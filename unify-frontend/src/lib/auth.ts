"use client";

import { logOutRequest } from "@/services/auth.service";
import { useRouter } from "next/navigation";

export const useLogout = () => {
  const router = useRouter();

  const handleLogout = async () => {
    try {
      await logOutRequest();
    } catch (error) {
      console.error("Logout failed:", error);
    } finally {
      router.push("/login");
    }
  };

  return handleLogout;
};
