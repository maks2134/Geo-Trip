import type { Metadata } from "next";
import { QueryProvider } from "@/shared/lib/query-client";
import "./globals.css";

export const metadata: Metadata = {
  title: "AtlasFlow — NomadGlass Travel Scratch Map & Planner",
  description: "Personal travel scratch-map and day-by-day trip planner with adaptive glass UI.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" data-theme="dark">
      <body className="w-screen h-screen overflow-hidden">
        <QueryProvider>{children}</QueryProvider>
      </body>
    </html>
  );
}
