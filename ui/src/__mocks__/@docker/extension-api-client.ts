import { vi } from "vitest";

const mockService = {
  get: vi.fn().mockResolvedValue({ checks: [] }),
  post: vi.fn().mockResolvedValue({
    checks: [],
    startedAt: new Date().toISOString(),
  }),
};

const mockClient = {
  extension: {
    vm: {
      service: mockService,
    },
  },
  desktopUI: {
    toast: {
      success: vi.fn(),
      error: vi.fn(),
      warning: vi.fn(),
    },
  },
  docker: {
    cli: { exec: vi.fn() },
  },
};

export const createDockerDesktopClient = () => mockClient;
