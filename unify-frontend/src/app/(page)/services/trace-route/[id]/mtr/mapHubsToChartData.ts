import type { Hub } from "./types"

export type ChartDataItem = {
  hop: number
  ip: string
  name: string
  avg: number
  min: number
  max: number
  loss: number
}

export const mapHubsToChartData = (hubs: Hub[]): ChartDataItem[] => {
  return hubs.map((hub) => ({
    hop: hub.count,
    ip: hub.host,
    name: hub.Dns || hub.host,
    avg: hub.Avg,
    min: hub.Best,
    max: hub.Wrst,
    loss: hub["Loss%"],
  }))
}