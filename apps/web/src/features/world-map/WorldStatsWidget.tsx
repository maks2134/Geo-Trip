"use client";

import React, { useState, useEffect } from "react";
import {
  Globe,
  Award,
  Compass,
  MapPin,
  Navigation,
  Send,
  Flag,
  Shield,
  Star,
  Crown,
  ChevronRight,
  Sparkles,
  X,
} from "lucide-react";
import { GeoStats } from "@/entities/place";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { ProgressBar } from "@/shared/ui/ProgressBar";
import { useTranslation } from "@/shared/i18n/i18n-context";
import { useMapUIStore } from "@/stores/map-ui-store";
import { apiClient, RegionalMilestonesResponse, MilestoneBadge, RegionProgress } from "@/shared/lib/api";
import { cn } from "@/shared/lib/utils";

interface WorldStatsWidgetProps {
  stats: GeoStats;
}

export function WorldStatsWidget({ stats }: WorldStatsWidgetProps) {
  const { t, plural } = useTranslation();
  const { activeRegionFilter, setActiveRegionFilter } = useMapUIStore();
  const [activeTab, setActiveTab] = useState<"regions" | "milestones">("regions");
  const [data, setData] = useState<RegionalMilestonesResponse | null>(null);
  const [selectedMilestone, setSelectedMilestone] = useState<MilestoneBadge | null>(null);

  useEffect(() => {
    let isMounted = true;
    apiClient.getRegionalMilestones().then((res) => {
      if (isMounted) setData(res);
    });
    return () => {
      isMounted = false;
    };
  }, []);

  const getBadgeIcon = (iconName: string, unlocked: boolean) => {
    const className = cn("w-4 h-4", unlocked ? "text-amber-300" : "text-slate-500");
    switch (iconName) {
      case "compass": return <Compass className={className} />;
      case "map-pin": return <MapPin className={className} />;
      case "navigation": return <Navigation className={className} />;
      case "send": return <Send className={className} />;
      case "globe": return <Globe className={className} />;
      case "flag": return <Flag className={className} />;
      case "shield": return <Shield className={className} />;
      case "star": return <Star className={className} />;
      case "crown": return <Crown className={className} />;
      default: return <Award className={className} />;
    }
  };

  const unlockedCount = data?.milestones.filter((m) => m.unlocked).length || 5;

  return (
    <>
      <GlassPanel className="p-4 w-84 space-y-3.5 shadow-2xl rounded-3xl">
        {/* Header with Title and Badges Tab */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="p-2 rounded-xl bg-sky-500/10 text-sky-400">
              <Globe className="w-5 h-5" />
            </div>
            <div>
              <h3 className="font-semibold text-sm">{t("stats.title")}</h3>
              <p className="text-xs text-[var(--text-muted)]">
                {stats.unlockedCountries}{" "}
                {plural(stats.unlockedCountries, "страна", "страны", "стран").split(" ")[1] || "countries"}
              </p>
            </div>
          </div>

          <button
            onClick={() => setActiveTab(activeTab === "milestones" ? "regions" : "milestones")}
            className="flex items-center gap-1.5 text-xs font-semibold px-2.5 py-1 rounded-full bg-amber-500/15 text-amber-300 border border-amber-500/30 hover:bg-amber-500/25 transition-all"
          >
            <Award className="w-3.5 h-3.5" />
            <span>{unlockedCount}/10 Badges</span>
          </button>
        </div>

        {/* Global Progress */}
        <div className="space-y-1">
          <div className="flex justify-between text-xs">
            <span className="text-[var(--text-muted)]">{t("stats.worldExplored")}</span>
            <span className="font-bold text-sky-400">{stats.worldPercentage}%</span>
          </div>
          <ProgressBar value={stats.worldPercentage} />
        </div>

        {/* 3-metric strip */}
        <div className="grid grid-cols-3 gap-1.5 py-2 border-y border-white/5">
          <div className="p-1.5 rounded-xl bg-white/[0.03] text-center">
            <div className="text-base font-black text-[var(--text)]">{stats.unlockedCountries}</div>
            <div className="text-[9px] text-[var(--text-muted)] uppercase tracking-wider">
              {t("stats.countriesUnlocked")}
            </div>
          </div>
          <div className="p-1.5 rounded-xl bg-white/[0.03] text-center">
            <div className="text-base font-black text-[var(--text)]">{stats.unlockedCities}</div>
            <div className="text-[9px] text-[var(--text-muted)] uppercase tracking-wider">
              {t("stats.citiesVisited")}
            </div>
          </div>
          <div className="p-1.5 rounded-xl bg-white/[0.03] text-center">
            <div className="text-base font-black text-[var(--text)]">{(stats.totalKm / 1000).toFixed(0)}k</div>
            <div className="text-[9px] text-[var(--text-muted)] uppercase tracking-wider">
              {t("stats.totalDistance")}
            </div>
          </div>
        </div>

        {/* Sub-tabs: UN Geoscheme Regions vs Milestones (FEAT-06) */}
        <div className="flex items-center gap-1 p-0.5 rounded-xl bg-white/[0.03] border border-white/5 text-xs">
          <button
            onClick={() => setActiveTab("regions")}
            className={cn(
              "flex-1 py-1 rounded-lg font-semibold transition-all text-center",
              activeTab === "regions"
                ? "bg-[var(--accent)] text-white shadow-sm"
                : "text-[var(--text-muted)] hover:text-white"
            )}
          >
            {t("milestones.regionsTab")}
          </button>
          <button
            onClick={() => setActiveTab("milestones")}
            className={cn(
              "flex-1 py-1 rounded-lg font-semibold transition-all text-center",
              activeTab === "milestones"
                ? "bg-[var(--accent)] text-white shadow-sm"
                : "text-[var(--text-muted)] hover:text-white"
            )}
          >
            {t("milestones.badgesTab")}
          </button>
        </div>

        {/* Tab 1: UN Geoscheme Regional Completion Bars */}
        {activeTab === "regions" && (
          <div className="space-y-2 max-h-52 overflow-y-auto pr-1">
            <div className="text-[10px] text-[var(--text-muted)] italic px-0.5">
              {t("milestones.clickToFly")}
            </div>
            {data?.regions.map((reg) => {
              const isSelected = activeRegionFilter === reg.subregion;
              return (
                <button
                  key={`${reg.region}-${reg.subregion}`}
                  onClick={() =>
                    setActiveRegionFilter(isSelected ? null : reg.subregion)
                  }
                  className={cn(
                    "w-full text-left p-2 rounded-xl border transition-all text-xs space-y-1 block",
                    isSelected
                      ? "bg-sky-500/20 border-sky-400/50 shadow-md shadow-sky-500/10"
                      : "bg-white/[0.02] border-white/5 hover:bg-white/[0.05]"
                  )}
                >
                  <div className="flex items-center justify-between">
                    <span className="font-semibold text-slate-200 truncate">
                      {reg.subregion}
                    </span>
                    <span className="text-[11px] font-bold text-sky-400">
                      {reg.visitedCountries}/{reg.totalCountries} ({reg.completionPercentage.toFixed(0)}%)
                    </span>
                  </div>
                  <ProgressBar
                    value={reg.completionPercentage}
                    barColor={isSelected ? "bg-sky-400" : "bg-sky-500/70"}
                  />
                </button>
              );
            })}
          </div>
        )}

        {/* Tab 2: Milestone Badges Showcase */}
        {activeTab === "milestones" && (
          <div className="grid grid-cols-2 gap-1.5 max-h-52 overflow-y-auto pr-1">
            {data?.milestones.map((m) => (
              <button
                key={m.id}
                onClick={() => setSelectedMilestone(m)}
                className={cn(
                  "p-2 rounded-xl border text-left flex items-center gap-2 transition-all",
                  m.unlocked
                    ? "bg-amber-500/10 border-amber-500/30 hover:bg-amber-500/20"
                    : "bg-white/[0.02] border-white/5 opacity-50"
                )}
              >
                <div
                  className={cn(
                    "w-8 h-8 rounded-lg flex items-center justify-center shrink-0",
                    m.unlocked ? "bg-amber-500/20 shadow-sm" : "bg-white/5"
                  )}
                >
                  {getBadgeIcon(m.icon, m.unlocked)}
                </div>
                <div className="min-w-0">
                  <div className="font-bold text-[11px] text-slate-200 truncate">{m.name}</div>
                  <div className="text-[9px] text-slate-400">
                    {m.unlocked ? "Unlocked" : `${m.threshold} countries`}
                  </div>
                </div>
              </button>
            ))}
          </div>
        )}
      </GlassPanel>

      {/* Milestone Celebration Modal (FEAT-06) */}
      {selectedMilestone && (
        <div
          onClick={() => setSelectedMilestone(null)}
          className="fixed inset-0 z-50 bg-black/80 backdrop-blur-md flex items-center justify-center p-4 animate-in fade-in"
        >
          <div
            onClick={(e) => e.stopPropagation()}
            className="w-full max-w-sm rounded-3xl bg-slate-900 border border-amber-500/40 p-6 shadow-2xl text-center space-y-4"
          >
            <div className="w-16 h-16 rounded-2xl bg-amber-500/20 border border-amber-500/40 mx-auto flex items-center justify-center shadow-lg shadow-amber-500/20">
              {getBadgeIcon(selectedMilestone.icon, true)}
            </div>

            <div>
              <div className="flex items-center justify-center gap-1.5 text-xs font-bold text-amber-400 uppercase tracking-widest">
                <Sparkles className="w-4 h-4" />
                <span>EXPEDITION MILESTONE</span>
              </div>
              <h2 className="text-xl font-black text-white mt-1">
                {selectedMilestone.name}
              </h2>
              <p className="text-xs text-slate-300 mt-2 leading-relaxed">
                {selectedMilestone.description}
              </p>
              <div className="mt-3 px-3 py-1 rounded-full bg-white/5 border border-white/10 text-xs font-semibold text-slate-400 inline-block">
                Threshold: {selectedMilestone.threshold} Countries Visited
              </div>
            </div>

            <button
              onClick={() => setSelectedMilestone(null)}
              className="w-full py-2.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-lg shadow-amber-500/25 transition-all"
            >
              Continue Exploring
            </button>
          </div>
        </div>
      )}
    </>
  );
}
