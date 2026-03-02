import { useState } from "react";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import RefreshIcon from "@mui/icons-material/Refresh";
import { ContainerTable } from "./components/ContainerTable";
import type { ImageCheck } from "./types";

const MOCK_CHECKS: ImageCheck[] = [
  {
    id: 1,
    containerName: "web-server",
    containerId: "abc123def456",
    imageRef: "nginx:latest",
    localDigest: "sha256:aaa111",
    remoteDigest: "sha256:bbb222",
    status: "update-available",
    checkedAt: new Date().toISOString(),
    registry: "dockerhub",
  },
  {
    id: 2,
    containerName: "database",
    containerId: "def456ghi789",
    imageRef: "postgres:16",
    localDigest: "sha256:ccc333",
    remoteDigest: "sha256:ccc333",
    status: "up-to-date",
    checkedAt: new Date().toISOString(),
    registry: "dockerhub",
  },
  {
    id: 3,
    containerName: "cache",
    containerId: "ghi789jkl012",
    imageRef: "redis:7-alpine",
    localDigest: "",
    remoteDigest: "",
    status: "unknown",
    checkedAt: "",
    registry: "dockerhub",
  },
];

export default function App() {
  const [checks] = useState<ImageCheck[]>(MOCK_CHECKS);
  const [loading, setLoading] = useState(false);

  const handleCheckAll = () => {
    setLoading(true);
    // TODO: Replace with real backend call in Wave 3
    setTimeout(() => {
      setLoading(false);
    }, 1500);
  };

  return (
    <Box sx={{ p: 3, height: "100vh", display: "flex", flexDirection: "column" }}>
      <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 2 }}>
        <Typography variant="h3">DockPulse</Typography>
        <Button
          variant="contained"
          startIcon={<RefreshIcon />}
          onClick={handleCheckAll}
          disabled={loading}
        >
          {loading ? "Checking..." : "Check Now"}
        </Button>
      </Stack>

      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        Check if your running container images have newer versions available.
      </Typography>

      <ContainerTable checks={checks} loading={loading} />

      <Typography variant="caption" color="text.disabled" sx={{ mt: "auto", pt: 2, textAlign: "center" }}>
        DockPulse v{__APP_VERSION__}
      </Typography>
    </Box>
  );
}
