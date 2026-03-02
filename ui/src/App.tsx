import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import Alert from "@mui/material/Alert";
import RefreshIcon from "@mui/icons-material/Refresh";
import { ContainerTable } from "./components/ContainerTable";
import { useBackend } from "./hooks/useBackend";

export default function App() {
  const { checks, checkAll, loading, error, fetchChecks } = useBackend();

  const handleCheckAll = () => {
    checkAll();
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

      {error != null && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => fetchChecks()}>
          {error}
        </Alert>
      )}

      <ContainerTable checks={checks} loading={loading} />

      <Typography variant="caption" color="text.disabled" sx={{ mt: "auto", pt: 2, textAlign: "center" }}>
        DockPulse v{__APP_VERSION__}
      </Typography>
    </Box>
  );
}
