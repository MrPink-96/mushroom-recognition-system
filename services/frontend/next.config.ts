import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',
  images: {
    remotePatterns: [
      {
        protocol: "http",
        hostname: process.env.API_HOST || "localhost",
        port: "8080",
      },
    ],
  },
};

export default nextConfig;
