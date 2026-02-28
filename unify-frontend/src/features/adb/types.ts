export type AdbResult = {
  id: string;
  index: number;
  status: string;
  startTime: Date;
  finishTime: Date;
  ipAddress: string;
  port: number;
  deviceName: string;
  result: string;
  serviceType: string;
};

export enum AdbCommand {
  Reboot = "reboot-device",
  RemoveYoutubeData = "remove-youtube-data",
  EnableYoutube = "enable-youtube",
  DisableYoutube = "disable-youtube",
  GetUptime = "get-uptime",
  GetProductModel = "get-product-model",
  GetSerialNumber = "get-serial-number",
  GetDeviceId = "get-device-id",
}
