import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Skeleton from "@mui/material/Skeleton";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import { StatusChip } from "./StatusChip";
import type { ImageCheck } from "../types";

interface ContainerTableProps {
  checks: ImageCheck[];
  loading: boolean;
}

function formatCheckedAt(checkedAt: string): string {
  if (!checkedAt) return "Never";

  const date = new Date(checkedAt);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSec = Math.floor(diffMs / 1000);

  if (diffSec < 60) return "Just now";
  if (diffSec < 3600) return `${Math.floor(diffSec / 60)} min ago`;
  if (diffSec < 86400) return `${Math.floor(diffSec / 3600)} hr ago`;

  return date.toLocaleDateString();
}

function LoadingSkeleton() {
  return (
    <>
      {[1, 2, 3].map((row) => (
        <TableRow key={row}>
          <TableCell><Skeleton variant="text" /></TableCell>
          <TableCell><Skeleton variant="text" /></TableCell>
          <TableCell><Skeleton variant="rounded" width={100} height={24} /></TableCell>
          <TableCell><Skeleton variant="text" width={80} /></TableCell>
        </TableRow>
      ))}
    </>
  );
}

export function ContainerTable({ checks, loading }: ContainerTableProps) {
  if (!loading && checks.length === 0) {
    return (
      <Box sx={{ textAlign: "center", py: 6 }}>
        <Typography variant="body1" color="text.secondary">
          No containers found. Start some containers and check again.
        </Typography>
      </Box>
    );
  }

  return (
    <TableContainer sx={{ flex: 1 }}>
      <Table size="small" stickyHeader>
        <TableHead>
          <TableRow>
            <TableCell>Container</TableCell>
            <TableCell>Image</TableCell>
            <TableCell>Status</TableCell>
            <TableCell>Last Checked</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {loading ? (
            <LoadingSkeleton />
          ) : (
            checks.map((check) => (
              <TableRow key={check.id} hover>
                <TableCell>
                  <Typography variant="body2" fontWeight={500}>
                    {check.containerName}
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    {check.containerId.substring(0, 12)}
                  </Typography>
                </TableCell>
                <TableCell>
                  <Typography variant="body2">{check.imageRef}</Typography>
                </TableCell>
                <TableCell>
                  <StatusChip status={check.status} />
                </TableCell>
                <TableCell>
                  <Typography variant="body2" color="text.secondary">
                    {formatCheckedAt(check.checkedAt)}
                  </Typography>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
