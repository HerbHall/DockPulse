import { createDockerDesktopClient } from "@docker/extension-api-client";
import { useState, useCallback, useEffect } from "react";
import type { ImageCheck, CheckAllResponse } from "../types";

const ddClient = createDockerDesktopClient();

export function useBackend() {
  const [checks, setChecks] = useState<ImageCheck[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchChecks = useCallback(async () => {
    try {
      const result = await ddClient.extension.vm?.service?.get("/api/checks");
      const data = result as { checks: ImageCheck[] };
      setChecks(data.checks ?? []);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to fetch checks");
    }
  }, []);

  const checkAll = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await ddClient.extension.vm?.service?.post("/api/check-all", {});
      const data = result as CheckAllResponse;
      setChecks(data.checks ?? []);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Check failed");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchChecks();
  }, [fetchChecks]);

  return { checks, checkAll, loading, error, fetchChecks };
}
