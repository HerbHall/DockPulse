import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import App from "../App";
import { createDockerDesktopClient } from "@docker/extension-api-client";
import type { ImageCheck } from "../types";

const mockChecks: ImageCheck[] = [
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
];

function getMockService() {
  const client = createDockerDesktopClient();
  return client.extension.vm!.service!;
}

describe("App", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders the heading and check button", async () => {
    render(<App />);
    expect(screen.getByText("DockPulse")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /check now/i })).toBeInTheDocument();
  });

  it("fetches checks on mount and displays them", async () => {
    const service = getMockService();
    vi.mocked(service.get).mockResolvedValueOnce({ checks: mockChecks });

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText("web-server")).toBeInTheDocument();
    });
    expect(screen.getByText("database")).toBeInTheDocument();
  });

  it("calls check-all when button is clicked", async () => {
    const service = getMockService();
    vi.mocked(service.get).mockResolvedValueOnce({ checks: [] });
    vi.mocked(service.post).mockResolvedValueOnce({
      checks: mockChecks,
      startedAt: new Date().toISOString(),
    });

    render(<App />);

    const button = screen.getByRole("button", { name: /check now/i });
    fireEvent.click(button);

    await waitFor(() => {
      expect(service.post).toHaveBeenCalledWith("/api/check-all", {});
    });

    await waitFor(() => {
      expect(screen.getByText("web-server")).toBeInTheDocument();
    });
  });

  it("displays error alert when fetch fails", async () => {
    const service = getMockService();
    vi.mocked(service.get).mockRejectedValueOnce(new Error("Connection refused"));

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText("Connection refused")).toBeInTheDocument();
    });
  });

  it("shows empty state when no containers", async () => {
    const service = getMockService();
    vi.mocked(service.get).mockResolvedValueOnce({ checks: [] });

    render(<App />);

    await waitFor(() => {
      expect(
        screen.getByText("No containers found. Start some containers and check again."),
      ).toBeInTheDocument();
    });
  });
});
