"use client";

import React from "react";
import { cn } from "@/shared/lib/utils";

interface GlassPanelProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  className?: string;
  variant?: "default" | "subtle" | "active";
}

export function GlassPanel({
  children,
  className,
  variant = "default",
  ...props
}: GlassPanelProps) {
  return (
    <div
      className={cn(
        "rounded-2xl backdrop-blur-glass transition-all duration-200 border",
        "bg-[var(--glass)] border-[var(--glass-border)] text-[var(--text)]",
        "shadow-glass",
        variant === "active" && "ring-1 ring-[var(--accent)]",
        variant === "subtle" && "bg-opacity-60",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}