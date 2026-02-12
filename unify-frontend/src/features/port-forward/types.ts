export type PortForwardResult = {
  id: string;
  createdAt: Date;
  expiresAt: Date;
  status: string;
  listenIp: string;
  listenPort: number;
  destIp: string;
  destPort: number;
  protocol: "tcp" | "udp";
  ruleComment: string;
  lastAppliedAt: Date | null;
  index: number;
};


export const EXPIRE_OPTIONS = ["5m", "15m", "1h", "30h"] as const;
