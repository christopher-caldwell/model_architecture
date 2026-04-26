import { execFileSync } from "node:child_process";
import { describe, expect, it } from "vitest";

describe("import boundary checker", () => {
  it("passes for the current workspace", () => {
    const output = execFileSync("node", ["./scripts/check-boundaries.mjs"], {
      cwd: process.cwd(),
      encoding: "utf8"
    });

    expect(output).toContain("Import boundaries ok");
  });
});
