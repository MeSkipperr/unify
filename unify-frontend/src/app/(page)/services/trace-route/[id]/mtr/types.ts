// types/mtr.ts
export type Hub = {
    count: number
    host: string
    Dns: string
    "Loss%": number
    Snt: number
    Last: number
    Avg: number
    Best: number
    Wrst: number
    Stdev: number
}

export type MTRReport = {
    src: string
    dst: string
    tos: number
    tests: number
    psize: string
    bitpattern: string
    TotalHops: number
    Reachable: boolean
    AvgRTT: number
}

export type MTRMessage = {
    report: {
        mtr: MTRReport
        hubs: Hub[]
    }
}

export type MTRData = {
    time: string
    id: string
    message: MTRMessage
}
