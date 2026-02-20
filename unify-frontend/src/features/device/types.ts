export type Device = {
  id: string;
  index: number;
  name: string;
  ipAddress: string;
  macAddress: string;
  roomNumber: string;
  isConnect: boolean;
  type: string;
  deviceProduct: string;
  description?: string;
  statusUpdatedAt: Date;
  notification: boolean;
};

export enum DeviceType {
  AccessPoint = "access-point",
  IPTV = "iptv",
  CCTV = "cctv",
  Switch = "sw",
}
