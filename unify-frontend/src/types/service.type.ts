export type Services<TConfig> = {
    id: string
    serviceName: string
    displayName: string
    description: string
    version: string
    type: string
    config: TConfig
    updatedAt: Date
    status:string
}
