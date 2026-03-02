import Chip from "@mui/material/Chip";
import type { CheckStatus } from "../types";

interface StatusChipProps {
  status: CheckStatus;
}

const STATUS_CONFIG: Record<CheckStatus, { label: string; color: "success" | "warning" | "error" | "default" | "info" }> = {
  "up-to-date": { label: "Up to Date", color: "success" },
  "update-available": { label: "Update Available", color: "warning" },
  "check-failed": { label: "Check Failed", color: "error" },
  "unknown": { label: "Unknown", color: "default" },
  "checking": { label: "Checking...", color: "info" },
};

export function StatusChip({ status }: StatusChipProps) {
  const config = STATUS_CONFIG[status];

  return (
    <Chip
      label={config.label}
      color={config.color}
      size="small"
      variant="filled"
    />
  );
}
