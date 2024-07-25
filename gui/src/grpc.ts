import { ChannelCredentials, GrpcObject } from "@grpc/grpc-js"

export enum ServingStatus {
	UNKNOWN = "UNKNOWN",
	SERVING = "SERVING",
	NOT_SERVING = "NOT_SERVING",
	SERVICE_UNKNOWN = "SERVICE_UNKNOWN",
}

export type HealthCheckResponse = {
	status: ServingStatus
}

export type HealthResponse = {
	Check: (request: { service: string }, callback: (err: any, response: HealthCheckResponse) => void) => void
}

export type HealthProtobufTypeDefinition = GrpcObject & {
	grpc: {
		health: {
			v1: {
				Health: {
					new (address: string, credentials: ChannelCredentials, options?: any): HealthResponse
				}
			}
		}
	}
}
