"use client";

import React from "react";
import { Utensils, Hotel, Eye, Trees, Train, CheckCircle2, Calendar, Bookmark } from "lucide-react";
import { PinCategory, PinStatus } from "@/entities/place";
import { useMapUIStore } from "@/stores/map-ui-store";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { cn } from "@/shared/lib/utils";

export function PinFilters() {
  const { statusFilters, categoryFilters, toggleStatusFilter, toggleCategoryFilter, resetFilters } =
    useMapUIStore();

  const statuses: { id: PinStatus; label: string; icon: React.ReactNode; color: string }[] = [
    { id: "visited", label: "Visited", icon: <CheckCircle2 className="w-3.5 h-3.5" />, color: "text-amber-400" },
    { id: "planned", label: "Planned", icon: <Calendar className="w-3.5 h-3.5" />, color: "text-blue-400" },
    { id: "wishlist", label: "Wishlist", icon: <Bookmark className="w-3.5 h-3.5" />, color: "text-pink-400" },
  ];

  const categories: { id: PinCategory; label: string; icon: React.ReactNode }[] = [
    { id: "sight", label: "Sights", icon: <Eye className="w-3.5 h-3.5" /> },
    { id: "food", label: "Food", icon: <Utensils className="w-3.5 h-3.5" /> },
    { id: "stay", label: "Stays", icon: <Hotel className="w-3.5 h-3.5" /> },
    { id: "nature", label: "Nature", icon: <Trees className="w-3.5 h-3.5" /> },
    { id: "transport", label: "Transit", icon: <Train className="w-3.5 h-3.5" /> },
  ];

  return (
    <GlassPanel className="p-4 w-80 space-y-3.5">
      <div className="flex items-center justify-between text-xs">
        <span className="font-semibold uppercase tracking-wider text-[var(--text-muted)]">Pin Filters</span>
        <button
          onClick={resetFilters}
          className="text-[11px] text-[var(--accent)] hover:underline"
        >
          Reset All
        </button>
      </div>

      {/* Status filters */}
      <div className="flex gap-1.5">
        {statuses.map((s) => {
          const active = statusFilters.includes(s.id);
          return (
            <button
              key={s.id}
              onClick={() => toggleStatusFilter(s.id)}
              className={cn(
                "flex-1 flex items-center justify-center gap-1.5 py-1.5 px-2 rounded-xl text-xs font-medium border transition-all",
                active
                  ? "bg-white/10 border-white/20 text-white shadow-sm"
                  : "bg-transparent border-transparent text-[var(--text-muted)] hover:bg-white/5"
              )}
            >
              <span className={s.color}>{s.icon}</span>
              <span>{s.label}</span>
            </button>
          );
        })}
      </div>

      {/* Category filters */}
      <div className="flex flex-wrap gap-1.5 pt-2 border-t border-white/5">
        {categories.map((c) => {
          const active = categoryFilters.includes(c.id);
          return (
            <button
              key={c.id}
              onClick={() => toggleCategoryFilter(c.id)}
              className={cn(
                "flex items-center gap-1.5 py-1 px-2.5 rounded-lg text-xs font-medium border transition-all",
                active
                  ? "bg-[var(--accent)]/15 border-[var(--accent)]/30 text-[var(--accent)]"
                  : "bg-white/[0.02] border-white/5 text-[var(--text-muted)] hover:text-[var(--text)]"
              )}
            >
              {c.icon}
              <span>{c.label}</span>
            </button>
          );
        })}
      </div>
    </GlassPanel>
  );
}