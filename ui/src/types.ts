export type CheckStatus =
  | "up-to-date"
  | "update-available"
  | "check-failed"
  | "unknown"
  | "checking";

export interface ImageCheck {
  id: number;
  containerName: string;
  containerId: string;
  imageRef: string;
  localDigest: string;
  remoteDigest: string;
  status: CheckStatus;
  checkedAt: string;
  registry: string;
}

export interface CheckAllResponse {
  checks: ImageCheck[];
  startedAt: string;
}

export interface StatusResponse {
  healthy: boolean;
  version: string;
}
