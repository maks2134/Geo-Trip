"use client";

import React from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { MapCanvas } from "@/shared/map/MapCanvas";
import { Dock } from "@/shared/ui/Dock";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { WorldStatsWidget } from "@/features/world-map/WorldStatsWidget";
import { PinFilters } from "@/features/world-map/PinFilters";
import { TripSwitcher } from "@/features/trip-workspace/TripSwitcher";
import { DayTimeline } from "@/features/trip-workspace/DayTimeline";
import { PlaceDrawer } from "@/features/place-drawer/PlaceDrawer";
import { BudgetTab } from "@/features/budget/BudgetTab";
import { ChecklistTab } from "@/features/checklists/ChecklistTab";
import { PrintDossier } from "@/features/trip-workspace/PrintDossier";
import { useMapUIStore } from "@/stores/map-ui-store";
import { useTranslation } from "@/shared/i18n/i18n-context";
import { apiClient } from "@/shared/lib/api";
import { Trip } from "@/entities/trip";
import { DollarSign, CheckSquare, Layers } from "lucide-react";
import { cn } from "@/shared/lib/utils";

export default function HomePage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const {
    mode,
    selectedTripId,
    activeWorkspaceTab,
    setActiveWorkspaceTab,
  } = useMapUIStore();

  const { data: places = [] } = useQuery({
    queryKey: ["places"],
    queryFn: apiClient.getPlaces,
  });

  const { data: geoStats } = useQuery({
    queryKey: ["geo-stats"],
    queryFn: apiClient.getGeoStats,
  });

  const { data: trips = [] } = useQuery({
    queryKey: ["trips"],
    queryFn: apiClient.getTrips,
  });

  const activeTrip = trips.find((t) => t.id === selectedTripId) || trips[0] || null;

  const handleTripUpdate = (updatedTrip: Trip) => {
    queryClient.setQueryData(["trips"], (old: Trip[] | undefined) =>
      old ? old.map((t) => (t.id === updatedTrip.id ? updatedTrip : t)) : [updatedTrip]
    );
  };

  return (
    <main className="relative w-screen h-screen overflow-hidden">
      <MapCanvas places={places} trip={activeTrip} />
      <Dock />

      {mode === "world" && geoStats && (
        <div className="absolute top-6 left-6 z-20 flex flex-col gap-4 pointer-events-none">
          <div className="pointer-events-auto">
            <WorldStatsWidget stats={geoStats} />
          </div>
          <div className="pointer-events-auto">
            <PinFilters />
          </div>
        </div>
      )}

      {mode === "trip" && activeTrip && (
        <div className="absolute top-6 left-6 bottom-24 w-[420px] z-20 pointer-events-none">
          <GlassPanel className="h-full flex flex-col pointer-events-auto shadow-2xl p-4 space-y-4 overflow-hidden">
            <TripSwitcher trips={trips} activeTrip={activeTrip} onTripUpdate={handleTripUpdate} />

            <div className="flex items-center gap-1 p-1 bg-white/[0.03] rounded-xl border border-white/5 shrink-0">
              <button
                onClick={() => setActiveWorkspaceTab("itinerary")}
                className={cn(
                  "flex-1 flex items-center justify-center gap-1.5 py-1.5 rounded-lg text-xs font-semibold transition-all",
                  activeWorkspaceTab === "itinerary"
                    ? "bg-[var(--accent)] text-white shadow-sm"
                    : "text-[var(--text-muted)] hover:text-[var(--text)]"
                )}
              >
                <Layers className="w-3.5 h-3.5" />
                <span>{t("tabs.itinerary")}</span>
              </button>

              <button
                onClick={() => setActiveWorkspaceTab("budget")}
                className={cn(
                  "flex-1 flex items-center justify-center gap-1.5 py-1.5 rounded-lg text-xs font-semibold transition-all",
                  activeWorkspaceTab === "budget"
                    ? "bg-[var(--accent)] text-white shadow-sm"
                    : "text-[var(--text-muted)] hover:text-[var(--text)]"
                )}
              >
                <DollarSign className="w-3.5 h-3.5" />
                <span>{t("tabs.budget")}</span>
              </button>

              <button
                onClick={() => setActiveWorkspaceTab("checklists")}
                className={cn(
                  "flex-1 flex items-center justify-center gap-1.5 py-1.5 rounded-lg text-xs font-semibold transition-all",
                  activeWorkspaceTab === "checklists"
                    ? "bg-[var(--accent)] text-white shadow-sm"
                    : "text-[var(--text-muted)] hover:text-[var(--text)]"
                )}
              >
                <CheckSquare className="w-3.5 h-3.5" />
                <span>{t("tabs.checklists")}</span>
              </button>
            </div>

            <div className="flex-1 overflow-y-auto pr-1">
              {activeWorkspaceTab === "itinerary" && (
                <DayTimeline trip={activeTrip} onTripUpdate={handleTripUpdate} />
              )}
              {activeWorkspaceTab === "budget" && (
                <BudgetTab trip={activeTrip} onTripUpdate={handleTripUpdate} />
              )}
              {activeWorkspaceTab === "checklists" && (
                <ChecklistTab trip={activeTrip} onTripUpdate={handleTripUpdate} />
              )}
            </div>
          </GlassPanel>
        </div>
      )}

      <PlaceDrawer places={places} />
      {activeTrip && <PrintDossier trip={activeTrip} />}
    </main>
  );
}
