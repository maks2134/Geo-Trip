"use client";

import React from "react";
import { Globe, Award } from "lucide-react";
import { GeoStats } from "@/entities/place";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { ProgressBar } from "@/shared/ui/ProgressBar";
import { useTranslation } from "@/shared/i18n/i18n-context";

interface WorldStatsWidgetProps {
  stats: GeoStats;
}

export function WorldStatsWidget({ stats }: WorldStatsWidgetProps) {
  const { t, plural } = useTranslation();

  return (
    <GlassPanel className="p-5 w-80 space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="p-2 rounded-xl bg-sky-500/10 text-sky-400">
            <Globe className="w-5 h-5" />
          </div>
          <div>
            <h3 className="font-semibold text-sm">{t("stats.title")}</h3>
            <p className="text-xs text-[var(--text-muted)]">
              {stats.unlockedCountries} {plural(stats.unlockedCountries, "страна", "страны", "стран").split(" ")[1] || "countries"}
            </p>
          </div>
        </div>
        <div className="flex items-center gap-1 text-xs font-semibold px-2.5 py-1 rounded-full bg-amber-500/15 text-amber-400 border border-amber-500/30">
          <Award className="w-3.5 h-3.5" />
          <span>Tier 3</span>
        </div>
      </div>

      <div className="space-y-1.5">
        <div className="flex justify-between text-xs">
          <span className="text-[var(--text-muted)]">{t("stats.worldExplored")}</span>
          <span className="font-bold text-sky-400">{stats.worldPercentage}%</span>
        </div>
        <ProgressBar value={stats.worldPercentage} />
      </div>

      <div className="grid grid-cols-3 gap-2 pt-2 border-t border-white/5">
        <div className="p-2 rounded-xl bg-white/[0.03] text-center">
          <div className="text-lg font-extrabold text-[var(--text)]">{stats.unlockedCountries}</div>
          <div className="text-[10px] text-[var(--text-muted)] uppercase tracking-wider">
            {t("stats.countriesUnlocked")}
          </div>
        </div>
        <div className="p-2 rounded-xl bg-white/[0.03] text-center">
          <div className="text-lg font-extrabold text-[var(--text)]">{stats.unlockedCities}</div>
          <div className="text-[10px] text-[var(--text-muted)] uppercase tracking-wider">
            {t("stats.citiesVisited")}
          </div>
        </div>
        <div className="p-2 rounded-xl bg-white/[0.03] text-center">
          <div className="text-lg font-extrabold text-[var(--text)]">{(stats.totalKm / 1000).toFixed(0)}k</div>
          <div className="text-[10px] text-[var(--text-muted)] uppercase tracking-wider">
            {t("stats.totalDistance")}
          </div>
        </div>
      </div>
    </GlassPanel>
  );
}
