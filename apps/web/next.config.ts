import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,
  // Avoid a stale locked .next directory from interrupted Windows dev sessions.
  distDir: process.env.NEXT_DIST_DIR || ".next-build",
};

export default nextConfig;
