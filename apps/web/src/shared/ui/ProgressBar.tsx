"use client";

import React from "react";
import { cn } from "@/shared/lib/utils";

interface ProgressBarProps {
  value: number; // 0 to 100
  className?: string;
  barColor?: string;
}

export function ProgressBar({ value, className, barColor }: ProgressBarProps) {
  const clamped = Math.min(100, Math.max(0, value));
  return (
    <div className={cn("w-full bg-white/10 rounded-full h-2 overflow-hidden", className)}>
      <div
        className={cn(
          "h-full rounded-full transition-all duration-500 ease-out",
          barColor || "bg-[var(--accent)]"
        )}
        style={{ width: `${clamped}%` }}
      />
    </div>
  );
}