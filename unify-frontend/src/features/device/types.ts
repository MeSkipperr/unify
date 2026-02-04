export type Device = {
    id: string
    index: number
    name: string
    ipAddress: string
    macAddress: string
    roomNumber: string
    isConnect: boolean
    type: string
    description?: string
    statusUpdatedAt: Date
    notification: boolean
}