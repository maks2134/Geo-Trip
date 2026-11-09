"use client";

import React from "react";
import { Globe, MapPin, Moon, Sun } from "lucide-react";
import { useMapUIStore } from "@/stores/map-ui-store";
import { GlassPanel } from "./GlassPanel";
import { cn } from "@/shared/lib/utils";

export function Dock() {
  const { mode, setMode, theme, toggleTheme } = useMapUIStore();

  return (
    <div className="fixed bottom-6 left-1/2 -translate-x-1/2 z-40">
      <GlassPanel className="flex items-center gap-1.5 p-1.5 rounded-full shadow-2xl backdrop-blur-2xl">
        <button
          onClick={() => setMode("world")}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all duration-200",
            mode === "world"
              ? "bg-[var(--accent)] text-white shadow-md shadow-[var(--accent)]/30"
              : "text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-white/5"
          )}
        >
          <Globe className="w-4 h-4" />
          <span>World Map</span>
        </button>

        <button
          onClick={() => setMode("trip")}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all duration-200",
            mode === "trip"
              ? "bg-[var(--accent)] text-white shadow-md shadow-[var(--accent)]/30"
              : "text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-white/5"
          )}
        >
          <MapPin className="w-4 h-4" />
          <span>Trip Planner</span>
        </button>

        <div className="h-5 w-px bg-white/10 mx-1" />

        <button
          onClick={toggleTheme}
          title="Toggle Dark/Light Mode"
          aria-label="Toggle Dark/Light Mode"
          className="p-2 rounded-full text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-white/5 transition-colors"
        >
          {theme === "dark" ? <Sun className="w-4 h-4 text-amber-400" /> : <Moon className="w-4 h-4 text-blue-500" />}
        </button>
      </GlassPanel>
    </div>
  );
}