import { describe, expect, it, vi } from "vitest";

describe("graphql client", () => {
  it("sends a no-store JSON POST request and returns response data", async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ data: { health: "ok" } }),
    });
    vi.stubGlobal("fetch", fetchMock);

    const { graphql } = await import("./graphql");
    await expect(graphql<{ health: string }>("query { health }")).resolves.toEqual({ health: "ok" });
    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining("/graphql"),
      expect.objectContaining({ method: "POST", cache: "no-store" }),
    );
  });
});
