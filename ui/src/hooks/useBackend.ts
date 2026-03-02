import { useState, useCallback } from "react";
import type { ImageCheck } from "../types";

// TODO: Replace with real ddClient.extension.vm.service calls in Wave 3
export function useBackend() {
  const [loading, setLoading] = useState(false);

  const getChecks = useCallback(async (): Promise<ImageCheck[]> => {
    // Placeholder -- will call ddClient.extension.vm.service.get("/api/checks")
    return [];
  }, []);

  const checkAll = useCallback(async (): Promise<ImageCheck[]> => {
    setLoading(true);
    try {
      // Placeholder -- will call ddClient.extension.vm.service.post("/api/check-all")
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  return { getChecks, checkAll, loading };
}
