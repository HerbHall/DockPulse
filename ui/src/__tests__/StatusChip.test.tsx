import { render, screen } from "@testing-library/react";
import { StatusChip } from "../components/StatusChip";
import type { CheckStatus } from "../types";

describe("StatusChip", () => {
  const cases: { status: CheckStatus; expectedLabel: string }[] = [
    { status: "up-to-date", expectedLabel: "Up to Date" },
    { status: "update-available", expectedLabel: "Update Available" },
    { status: "check-failed", expectedLabel: "Check Failed" },
    { status: "unknown", expectedLabel: "Unknown" },
    { status: "checking", expectedLabel: "Checking..." },
  ];

  it.each(cases)(
    "renders '$expectedLabel' for status '$status'",
    ({ status, expectedLabel }) => {
      render(<StatusChip status={status} />);
      expect(screen.getByText(expectedLabel)).toBeInTheDocument();
    },
  );
});
