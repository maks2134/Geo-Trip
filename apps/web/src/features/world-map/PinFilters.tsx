"use client";

import React from "react";
import { Utensils, Hotel, Eye, Trees, Train, CheckCircle2, Calendar, Bookmark } from "lucide-react";
import { PinCategory, PinStatus } from "@/entities/place";
import { useMapUIStore } from "@/stores/map-ui-store";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { cn } from "@/shared/lib/utils";
import { useTranslation } from "@/shared/i18n/i18n-context";

export function PinFilters() {
  const { statusFilters, categoryFilters, toggleStatusFilter, toggleCategoryFilter, resetFilters } =
    useMapUIStore();
  const { t } = useTranslation();

  const statuses: { id: PinStatus; labelKey: string; icon: React.ReactNode; color: string }[] = [
    { id: "visited", labelKey: "filters.visited", icon: <CheckCircle2 className="w-3.5 h-3.5" />, color: "text-amber-500 dark:text-amber-400" },
    { id: "planned", labelKey: "filters.planned", icon: <Calendar className="w-3.5 h-3.5" />, color: "text-blue-500 dark:text-blue-400" },
    { id: "wishlist", labelKey: "filters.wishlist", icon: <Bookmark className="w-3.5 h-3.5" />, color: "text-pink-500 dark:text-pink-400" },
  ];

  const categories: { id: PinCategory; labelKey: string; icon: React.ReactNode }[] = [
    { id: "sight", labelKey: "filters.sight", icon: <Eye className="w-3.5 h-3.5" /> },
    { id: "food", labelKey: "filters.food", icon: <Utensils className="w-3.5 h-3.5" /> },
    { id: "stay", labelKey: "filters.stay", icon: <Hotel className="w-3.5 h-3.5" /> },
    { id: "nature", labelKey: "filters.nature", icon: <Trees className="w-3.5 h-3.5" /> },
    { id: "transport", labelKey: "filters.transport", icon: <Train className="w-3.5 h-3.5" /> },
  ];

  return (
    <GlassPanel className="p-4 w-80 space-y-3.5">
      <div className="flex items-center justify-between text-xs">
        <span className="font-semibold uppercase tracking-wider text-[var(--text-muted)]">{t("filters.status")}</span>
        <button
          onClick={resetFilters}
          className="text-[11px] text-[var(--accent)] hover:underline"
        >
          {t("filters.reset")}
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
                  ? "bg-[var(--accent)]/15 border-[var(--accent)] text-[var(--accent)] font-bold shadow-sm"
                  : "bg-transparent border-transparent text-[var(--text-muted)] hover:bg-black/5 dark:hover:bg-white/5"
              )}
            >
              <span className={s.color}>{s.icon}</span>
              <span>{t(s.labelKey)}</span>
            </button>
          );
        })}
      </div>

      {/* Category filters */}
      <div className="flex flex-wrap gap-1.5 pt-2 border-t border-[var(--glass-border)]">
        {categories.map((c) => {
          const active = categoryFilters.includes(c.id);
          return (
            <button
              key={c.id}
              onClick={() => toggleCategoryFilter(c.id)}
              className={cn(
                "flex items-center gap-1.5 py-1 px-2.5 rounded-lg text-xs font-medium border transition-all",
                active
                  ? "bg-[var(--accent)]/15 border-[var(--accent)] text-[var(--accent)] font-bold shadow-sm"
                  : "bg-black/[0.03] dark:bg-white/[0.02] border-[var(--glass-border)] text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-black/5 dark:hover:bg-white/5"
              )}
            >
              {c.icon}
              <span>{t(c.labelKey)}</span>
            </button>
          );
        })}
      </div>
    </GlassPanel>
  );
}
