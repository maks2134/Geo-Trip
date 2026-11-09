"use client";

import React, { useState } from "react";
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import {
  Bookmark,
  Download,
  Footprints,
  Car,
  Route,
  ChevronDown,
  Wand2,
  Undo2,
  Navigation,
  Printer,
  Compass,
  Check,
  X,
  ArrowRight,
  Loader2,
} from "lucide-react";
import { Trip, PlaceRef } from "@/entities/trip";
import { useMapUIStore, NavAppPreference } from "@/stores/map-ui-store";
import { DraggablePlaceItem } from "./dnd/DraggablePlaceItem";
import { BacklogDrawer } from "./BacklogDrawer";
import { apiClient, OptimizeDayRouteResponse } from "@/shared/lib/api";
import { cn } from "@/shared/lib/utils";
import { useTranslation } from "@/shared/i18n/i18n-context";
import {
  buildNavigationRouteUrl,
  estimateTransit,
  calculateHaversineDistanceKm,
} from "@/shared/lib/navigation";

interface DayTimelineProps {
  trip: Trip;
  onTripUpdate: (updated: Trip) => void;
}

function triggerDownload(filename: string, content: string, mimeType: string) {
  const blob = new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

export function DayTimeline({ trip, onTripUpdate }: DayTimelineProps) {
  const {
    selectedDayId,
    setSelectedDayId,
    isBacklogOpen,
    toggleBacklog,
    setDossierOpen,
    navAppPreference,
    setNavAppPreference,
  } = useMapUIStore();
  const { t, plural } = useTranslation();
  const [isExportOpen, setIsExportOpen] = useState(false);

  // Transit leg navigation popover state (FEAT-01)
  const [activeLegNav, setActiveLegNav] = useState<{
    origin: PlaceRef;
    dest: PlaceRef;
    distanceKm: number;
    walkingMin: number;
    drivingMin: number;
  } | null>(null);

  // Route Optimization state (FEAT-04)
  const [isOptimizing, setIsOptimizing] = useState(false);
  const [optimizationPreview, setOptimizationPreview] = useState<OptimizeDayRouteResponse | null>(null);
  const [lockStart, setLockStart] = useState(true);
  const [lockEnd, setLockEnd] = useState(false);
  const [undoSnapshot, setUndoSnapshot] = useState<PlaceRef[] | null>(null);
  const [undoTimer, setUndoTimer] = useState<NodeJS.Timeout | null>(null);

  const activeDay = trip.days.find((d) => d.id === selectedDayId) || trip.days[0];

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 4 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
  );

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;

    const activeId = String(active.id);
    const overId = String(over.id);

    let targetDayId = activeDay?.id || null;
    let newIndex = 0;

    if (overId === "droppable-backlog") {
      targetDayId = null;
      newIndex = trip.backlog.length;
    } else if (activeDay) {
      const overIndex = activeDay.places.findIndex((p) => p.id === overId);
      newIndex = overIndex >= 0 ? overIndex : activeDay.places.length;
    }

    const updated = await apiClient.reorderPlace(trip.id, activeId, targetDayId, newIndex);
    onTripUpdate(updated);
  };

  // Compute leg distances and total day distance
  const places = activeDay?.places || [];
  const legs: {
    origin: PlaceRef;
    dest: PlaceRef;
    km: number;
    durationText: string;
    isWalk: boolean;
    walkingMin: number;
    drivingMin: number;
  }[] = [];
  let totalDistanceKm = 0;

  for (let i = 0; i < places.length - 1; i++) {
    const p1 = places[i];
    const p2 = places[i + 1];
    const km = calculateHaversineDistanceKm(p1.lat, p1.lng, p2.lat, p2.lng);
    totalDistanceKm += km;
    const isWalk = km < 2.5;
    const walkingMin = Math.max(3, Math.round(km * 12));
    const drivingMin = Math.max(5, Math.round(km * 2.2) + 4);
    const modeLabel = isWalk ? t("distance.walk") : t("distance.drive");
    const durationMin = isWalk ? walkingMin : drivingMin;

    legs.push({
      origin: p1,
      dest: p2,
      km,
      durationText: `${km.toFixed(1)} ${t("distance.km")} • ~${durationMin} min ${modeLabel}`,
      isWalk,
      walkingMin,
      drivingMin,
    });
  }

  // Trigger route optimization preview (FEAT-04)
  const handleRequestOptimization = async () => {
    if (!activeDay || places.length < 3) return;
    setIsOptimizing(true);
    try {
      const preview = await apiClient.optimizeDayRoute(trip.id, activeDay.id, lockStart, lockEnd);
      setOptimizationPreview(preview);
    } catch {
      // fallback handled in api client
    } finally {
      setIsOptimizing(false);
    }
  };

  // Apply route optimization with 10s undo window
  const handleApplyOptimization = () => {
    if (!optimizationPreview || !activeDay) return;

    // Save snapshot for undo
    const snapshot = [...activeDay.places];
    setUndoSnapshot(snapshot);

    // Map new place sequence
    const orderMap = new Map(optimizationPreview.places.map((p) => [p.itemId, p.sequenceOrder]));
    const reordered = [...activeDay.places].sort((a, b) => {
      const orderA = orderMap.get(a.id) ?? 999;
      const orderB = orderMap.get(b.id) ?? 999;
      return orderA - orderB;
    });

    activeDay.places = reordered;
    onTripUpdate({ ...trip });
    setOptimizationPreview(null);

    // Set 10s undo timeout
    if (undoTimer) clearTimeout(undoTimer);
    const timer = setTimeout(() => {
      setUndoSnapshot(null);
    }, 10000);
    setUndoTimer(timer);
  };

  // Restore previous order (Undo)
  const handleUndoOptimization = () => {
    if (!undoSnapshot || !activeDay) return;
    activeDay.places = undoSnapshot;
    onTripUpdate({ ...trip });
    setUndoSnapshot(null);
    if (undoTimer) clearTimeout(undoTimer);
  };

  // Route Export Handlers
  const exportGPX = () => {
    if (!activeDay || places.length === 0) return;
    const wpts = places
      .map(
        (p) =>
          `  <wpt lat="${p.lat}" lon="${p.lng}">\n    <name>${p.name.replace(/&/g, "&amp;")}</name>\n    <desc>${(p.notes || "").replace(/&/g, "&amp;")}</desc>\n  </wpt>`
      )
      .join("\n");
    const rtepts = places
      .map(
        (p) =>
          `    <rtept lat="${p.lat}" lon="${p.lng}">\n      <name>${p.name.replace(/&/g, "&amp;")}</name>\n    </rtept>`
      )
      .join("\n");

    const gpx = `<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="AtlasFlow NomadGlass" xmlns="http://www.topografix.com/GPX/1/1">
  <metadata>
    <name>${trip.title} - Day ${activeDay.dayNumber}</name>
  </metadata>
${wpts}
  <rte>
    <name>Day ${activeDay.dayNumber} Route</name>
${rtepts}
  </rte>
</gpx>`;
    triggerDownload(`atlasflow-${trip.slug || trip.id}-day-${activeDay.dayNumber}.gpx`, gpx, "application/gpx+xml");
    setIsExportOpen(false);
  };

  const exportGeoJSON = () => {
    if (!activeDay || places.length === 0) return;
    const features = [
      ...places.map((p, idx) => ({
        type: "Feature",
        properties: {
          id: p.id,
          order: idx + 1,
          name: p.name,
          category: p.category,
          notes: p.notes,
        },
        geometry: {
          type: "Point",
          coordinates: [p.lng, p.lat],
        },
      })),
      {
        type: "Feature",
        properties: {
          name: `Day ${activeDay.dayNumber} Itinerary Route`,
          totalDistanceKm: Number(totalDistanceKm.toFixed(2)),
        },
        geometry: {
          type: "LineString",
          coordinates: places.map((p) => [p.lng, p.lat]),
        },
      },
    ];

    const geojson = JSON.stringify(
      {
        type: "FeatureCollection",
        features,
      },
      null,
      2
    );
    triggerDownload(`atlasflow-${trip.slug || trip.id}-day-${activeDay.dayNumber}.geojson`, geojson, "application/geo+json");
    setIsExportOpen(false);
  };

  return (
    <div className="space-y-4">
      {/* Day Selector Dock */}
      <div className="flex items-center gap-1.5 overflow-x-auto pb-1 scrollbar-none">
        {trip.days.map((day) => {
          const isActive = day.id === (activeDay?.id ?? "");
          return (
            <button
              key={day.id}
              onClick={() => setSelectedDayId(day.id)}
              className={cn(
                "px-3.5 py-2 rounded-xl text-xs font-semibold whitespace-nowrap transition-all flex flex-col items-center",
                isActive
                  ? "bg-[var(--accent)] text-white shadow-lg shadow-[var(--accent)]/20"
                  : "bg-white/[0.04] text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-white/[0.08]"
              )}
            >
              <span>{t("trip.day")} {day.dayNumber}</span>
              <span className="text-[10px] opacity-75 font-normal">
                {plural(day.places.length, "точка", "точки", "точек")}
              </span>
            </button>
          );
        })}

        <button
          onClick={toggleBacklog}
          className={cn(
            "px-3 py-2 rounded-xl text-xs font-medium whitespace-nowrap transition-all flex items-center gap-1.5",
            isBacklogOpen
              ? "bg-sky-500/20 text-sky-300 border border-sky-500/40"
              : "bg-white/[0.03] text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-white/[0.07]"
          )}
        >
          <Bookmark className="w-3.5 h-3.5" />
          <span>{t("trip.backlog")} ({trip.backlog.length})</span>
        </button>
      </div>

      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <BacklogDrawer items={trip.backlog} isOpen={isBacklogOpen} onClose={toggleBacklog} />

        <div className="space-y-2">
          {/* Day Header with Distance & Actions */}
          <div className="flex items-center justify-between text-xs px-1">
            <div>
              <span className="font-semibold text-[var(--text)]">{activeDay?.title}</span>
              {totalDistanceKm > 0 && (
                <span className="ml-2 px-2 py-0.5 rounded-full bg-sky-500/10 text-sky-400 border border-sky-500/20 text-[10px] font-bold">
                  {t("trip.totalRouteDistance")}: {totalDistanceKm.toFixed(1)} {t("distance.km")}
                </span>
              )}
            </div>

            <div className="flex items-center gap-1.5">
              {/* Route Optimization Button (FEAT-04) */}
              {places.length >= 3 && (
                <button
                  onClick={handleRequestOptimization}
                  disabled={isOptimizing}
                  title="Optimize route sequence to minimize transit distance (2-opt TSP)"
                  className="flex items-center gap-1 px-2.5 py-1 rounded-lg bg-sky-500/20 hover:bg-sky-500/30 text-sky-300 border border-sky-500/40 transition-colors text-[11px] font-semibold"
                >
                  {isOptimizing ? (
                    <Loader2 className="w-3.5 h-3.5 animate-spin" />
                  ) : (
                    <Wand2 className="w-3.5 h-3.5 text-sky-300" />
                  )}
                  <span>{t("route.optimize")}</span>
                </button>
              )}

              {/* Dossier Print Trigger (FEAT-05) */}
              <button
                onClick={() => setDossierOpen(true)}
                title={t("dossier.printTooltip")}
                className="flex items-center gap-1 px-2.5 py-1 rounded-lg bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-300 border border-emerald-500/40 transition-colors text-[11px] font-semibold"
              >
                <Printer className="w-3.5 h-3.5" />
                <span className="hidden sm:inline">{t("dossier.title")}</span>
              </button>

              {/* Export Menu Dropdown */}
              <div className="relative">
                <button
                  onClick={() => setIsExportOpen(!isExportOpen)}
                  className="flex items-center gap-1 px-2.5 py-1 rounded-lg bg-white/[0.04] hover:bg-white/10 text-[var(--text-muted)] hover:text-[var(--text)] border border-white/10 transition-colors text-[11px]"
                >
                  <Download className="w-3.5 h-3.5" />
                  <span>{t("trip.exportRoute")}</span>
                  <ChevronDown className="w-3 h-3" />
                </button>

                {isExportOpen && (
                  <div className="absolute right-0 top-8 z-30 w-48 rounded-xl bg-slate-900/95 backdrop-blur-xl border border-white/15 shadow-2xl p-1.5 space-y-1">
                    <button
                      onClick={exportGPX}
                      className="w-full text-left px-3 py-1.5 rounded-lg text-xs text-slate-200 hover:bg-sky-500/20 hover:text-sky-300 font-medium transition-colors"
                    >
                      {t("trip.exportGPX")}
                    </button>
                    <button
                      onClick={exportGeoJSON}
                      className="w-full text-left px-3 py-1.5 rounded-lg text-xs text-slate-200 hover:bg-sky-500/20 hover:text-sky-300 font-medium transition-colors"
                    >
                      {t("trip.exportGeoJSON")}
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Undo Toast Notification (FEAT-04) */}
          {undoSnapshot && (
            <div className="p-2.5 rounded-xl bg-sky-500/25 border border-sky-400/50 backdrop-blur-md flex items-center justify-between text-xs text-sky-100 shadow-xl animate-in slide-in-from-top-2">
              <div className="flex items-center gap-2">
                <Check className="w-4 h-4 text-emerald-400" />
                <span>{t("route.optimizedSuccess")}</span>
              </div>
              <button
                onClick={handleUndoOptimization}
                className="flex items-center gap-1 px-2.5 py-1 rounded-lg bg-white/20 hover:bg-white/30 text-white font-bold transition-colors"
              >
                <Undo2 className="w-3.5 h-3.5" />
                <span>{t("common.undo")} (10s)</span>
              </button>
            </div>
          )}

          {activeDay && activeDay.places.length === 0 ? (
            <div className="text-center py-10 px-4 rounded-xl border border-dashed border-white/10 text-xs text-[var(--text-muted)]">
              {t("trip.noPlacesScheduled")}
            </div>
          ) : (
            <SortableContext
              items={activeDay ? activeDay.places.map((p) => p.id) : []}
              strategy={verticalListSortingStrategy}
            >
              <div className="space-y-2">
                {activeDay?.places.map((item, index) => (
                  <React.Fragment key={item.id}>
                    <DraggablePlaceItem item={item} orderNumber={index + 1} />

                    {/* Transit Leg Connector Badge with Click Popover (FEAT-01) */}
                    {index < legs.length && (
                      <div className="flex items-center justify-center my-0.5">
                        <button
                          type="button"
                          onClick={() =>
                            setActiveLegNav({
                              origin: legs[index].origin,
                              dest: legs[index].dest,
                              distanceKm: legs[index].km,
                              walkingMin: legs[index].walkingMin,
                              drivingMin: legs[index].drivingMin,
                            })
                          }
                          title="Click to open route navigation in Maps"
                          className="group flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/[0.04] hover:bg-sky-500/15 border border-white/10 hover:border-sky-400/40 backdrop-blur-md text-[10px] text-sky-400 font-medium shadow-sm transition-all"
                        >
                          {legs[index].isWalk ? (
                            <Footprints className="w-3 h-3 text-sky-300 group-hover:scale-110 transition-transform" />
                          ) : (
                            <Car className="w-3 h-3 text-emerald-300 group-hover:scale-110 transition-transform" />
                          )}
                          <span>{legs[index].durationText}</span>
                          <Navigation className="w-3 h-3 opacity-60 group-hover:opacity-100 transition-opacity ml-0.5" />
                        </button>
                      </div>
                    )}
                  </React.Fragment>
                ))}
              </div>
            </SortableContext>
          )}
        </div>
      </DndContext>

      {/* Transit Leg Navigation Modal (FEAT-01) */}
      {activeLegNav && (
        <div
          onClick={() => setActiveLegNav(null)}
          className="fixed inset-0 z-50 bg-black/75 backdrop-blur-md flex items-center justify-center p-4 animate-in fade-in"
        >
          <div
            onClick={(e) => e.stopPropagation()}
            className="w-full max-w-sm rounded-2xl bg-slate-900 border border-white/15 p-5 shadow-2xl space-y-4 text-xs"
          >
            <div className="flex items-center justify-between border-b border-white/10 pb-3">
              <div className="font-bold text-sm text-white flex items-center gap-2">
                <Navigation className="w-4 h-4 text-sky-400" />
                <span>{t("nav.transitDirections")}</span>
              </div>
              <button
                onClick={() => setActiveLegNav(null)}
                className="p-1 rounded-full text-slate-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="space-y-2 bg-white/[0.03] p-3 rounded-xl border border-white/5">
              <div className="flex items-center gap-2 font-medium text-slate-200 truncate">
                <span className="w-5 h-5 rounded-full bg-sky-500 text-white flex items-center justify-center text-[10px] font-bold shrink-0">
                  A
                </span>
                <span className="truncate">{activeLegNav.origin.name}</span>
              </div>
              <div className="flex items-center justify-center text-[11px] text-sky-400 font-semibold gap-2 py-0.5">
                <ArrowRight className="w-3.5 h-3.5" />
                <span>
                  {activeLegNav.distanceKm.toFixed(1)} km (~{activeLegNav.drivingMin}m drive, ~{activeLegNav.walkingMin}m walk)
                </span>
              </div>
              <div className="flex items-center gap-2 font-medium text-slate-200 truncate">
                <span className="w-5 h-5 rounded-full bg-emerald-500 text-white flex items-center justify-center text-[10px] font-bold shrink-0">
                  B
                </span>
                <span className="truncate">{activeLegNav.dest.name}</span>
              </div>
            </div>

            {/* Preferred Navigation App Selector */}
            <div className="space-y-1.5">
              <label className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">
                {t("nav.preferredApp")}
              </label>
              <div className="grid grid-cols-4 gap-1.5">
                {(["auto", "google", "apple", "yandex"] as NavAppPreference[]).map((app) => (
                  <button
                    key={app}
                    onClick={() => setNavAppPreference(app)}
                    className={cn(
                      "py-1.5 px-2 rounded-xl text-center capitalize font-semibold border text-[11px] transition-all",
                      navAppPreference === app
                        ? "bg-sky-500 text-white border-sky-500 shadow-md"
                        : "bg-white/[0.03] text-slate-400 border-white/10 hover:bg-white/[0.07]"
                    )}
                  >
                    {app}
                  </button>
                ))}
              </div>
            </div>

            {/* Launch Action Buttons */}
            <div className="space-y-2 pt-2">
              <a
                href={buildNavigationRouteUrl(
                  { latitude: activeLegNav.origin.lat, longitude: activeLegNav.origin.lng },
                  { latitude: activeLegNav.dest.lat, longitude: activeLegNav.dest.lng },
                  navAppPreference,
                  "driving"
                )}
                target="_blank"
                rel="noreferrer"
                className="w-full py-2.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold flex items-center justify-center gap-2 shadow-lg shadow-sky-500/25 transition-all text-sm"
              >
                <Car className="w-4 h-4" />
                <span>Drive with {navAppPreference.toUpperCase()}</span>
              </a>

              <a
                href={buildNavigationRouteUrl(
                  { latitude: activeLegNav.origin.lat, longitude: activeLegNav.origin.lng },
                  { latitude: activeLegNav.dest.lat, longitude: activeLegNav.dest.lng },
                  navAppPreference,
                  "walking"
                )}
                target="_blank"
                rel="noreferrer"
                className="w-full py-2.5 rounded-xl bg-white/10 hover:bg-white/15 text-white font-semibold flex items-center justify-center gap-2 border border-white/10 transition-colors"
              >
                <Footprints className="w-4 h-4 text-emerald-400" />
                <span>Walk route (~{activeLegNav.walkingMin} min)</span>
              </a>
            </div>
          </div>
        </div>
      )}

      {/* Route Optimization Preview Modal (FEAT-04) */}
      {optimizationPreview && (
        <div
          onClick={() => setOptimizationPreview(null)}
          className="fixed inset-0 z-50 bg-black/75 backdrop-blur-md flex items-center justify-center p-4 animate-in fade-in"
        >
          <div
            onClick={(e) => e.stopPropagation()}
            className="w-full max-w-md rounded-2xl bg-slate-900 border border-sky-500/30 p-5 shadow-2xl space-y-4 text-xs"
          >
            <div className="flex items-center justify-between border-b border-white/10 pb-3">
              <div className="font-bold text-sm text-white flex items-center gap-2">
                <Wand2 className="w-4 h-4 text-sky-400" />
                <span>{t("route.optimizationPreview")}</span>
              </div>
              <button
                onClick={() => setOptimizationPreview(null)}
                className="p-1 rounded-full text-slate-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            {/* Savings Callout */}
            <div className="p-4 rounded-xl bg-gradient-to-r from-sky-500/20 to-emerald-500/20 border border-sky-400/30 text-center space-y-1">
              <div className="text-2xl font-black text-white">
                -{optimizationPreview.savingsPercentage.toFixed(1)}% {t("route.distanceSaved")}
              </div>
              <p className="text-[11px] text-sky-200">
                Reduced from {optimizationPreview.originalDistanceKm.toFixed(1)} km to{" "}
                <strong className="text-emerald-300">
                  {optimizationPreview.optimizedDistanceKm.toFixed(1)} km
                </strong>{" "}
                (saving {optimizationPreview.distanceSavedKm.toFixed(1)} km of transit).
              </p>
            </div>

            {/* Lock Controls */}
            <div className="space-y-2 bg-white/[0.03] p-3 rounded-xl border border-white/5">
              <div className="text-[11px] font-semibold text-slate-300 uppercase tracking-wider mb-1">
                {t("route.constraints")}
              </div>
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={lockStart}
                  onChange={(e) => setLockStart(e.target.checked)}
                  className="rounded text-sky-500 focus:ring-0"
                />
                <span className="text-slate-200">
                  {t("route.lockFirstStop")} (e.g. Hotel / Departure)
                </span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={lockEnd}
                  onChange={(e) => setLockEnd(e.target.checked)}
                  className="rounded text-sky-500 focus:ring-0"
                />
                <span className="text-slate-200">
                  {t("route.lockLastStop")} (e.g. Dinner reservation)
                </span>
              </label>
            </div>

            <div className="flex gap-2 pt-2 border-t border-white/10">
              <button
                type="button"
                onClick={() => setOptimizationPreview(null)}
                className="flex-1 py-2.5 rounded-xl bg-white/5 hover:bg-white/10 text-slate-300 font-medium transition-colors"
              >
                {t("common.cancel")}
              </button>
              <button
                type="button"
                onClick={handleApplyOptimization}
                className="flex-1 py-2.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold shadow-lg shadow-sky-500/30 transition-all flex items-center justify-center gap-1.5"
              >
                <Check className="w-4 h-4" />
                <span>{t("route.applySort")}</span>
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
