"use client";

import React from "react";
import { cn } from "@/shared/lib/utils";
import { PinCategory, PinStatus } from "@/entities/place";

interface BadgeProps {
  children: React.ReactNode;
  category?: PinCategory;
  status?: PinStatus;
  variant?: "outline" | "filled";
  className?: string;
}

export function Badge({ children, category, status, variant = "filled", className }: BadgeProps) {
  let colorStyles = "bg-neutral-800/60 text-neutral-300 border-neutral-700/50";

  if (category === "food") {
    colorStyles = "bg-orange-500/15 text-orange-400 border-orange-500/30";
  } else if (category === "stay") {
    colorStyles = "bg-purple-500/15 text-purple-400 border-purple-500/30";
  } else if (category === "sight") {
    colorStyles = "bg-emerald-500/15 text-emerald-400 border-emerald-500/30";
  } else if (category === "nature") {
    colorStyles = "bg-teal-500/15 text-teal-400 border-teal-500/30";
  } else if (category === "transport") {
    colorStyles = "bg-sky-500/15 text-sky-400 border-sky-500/30";
  }

  if (status === "visited") {
    colorStyles = "bg-yellow-500/15 text-yellow-400 border-yellow-500/30";
  } else if (status === "planned") {
    colorStyles = "bg-blue-500/15 text-blue-400 border-blue-500/30";
  } else if (status === "wishlist") {
    colorStyles = "bg-pink-500/15 text-pink-400 border-pink-500/30";
  }

  return (
    <span
      className={cn(
        "inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border uppercase tracking-wider",
        colorStyles,
        className
      )}
    >
      {children}
    </span>
  );
}