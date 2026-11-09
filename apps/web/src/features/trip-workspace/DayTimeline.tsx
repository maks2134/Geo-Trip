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
import { Bookmark, Download, Footprints, Car, Route, ChevronDown } from "lucide-react";
import { Trip, PlaceRef } from "@/entities/trip";
import { useMapUIStore } from "@/stores/map-ui-store";
import { DraggablePlaceItem } from "./dnd/DraggablePlaceItem";
import { BacklogDrawer } from "./BacklogDrawer";
import { apiClient } from "@/shared/lib/api";
import { cn } from "@/shared/lib/utils";
import { useTranslation } from "@/shared/i18n/i18n-context";

interface DayTimelineProps {
  trip: Trip;
  onTripUpdate: (updated: Trip) => void;
}

function calculateHaversineKm(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const R = 6371; // Earth's radius in km
  const dLat = (lat2 - lat1) * (Math.PI / 180);
  const dLon = (lon2 - lon1) * (Math.PI / 180);
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(lat1 * (Math.PI / 180)) *
      Math.cos(lat2 * (Math.PI / 180)) *
      Math.sin(dLon / 2) *
      Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
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
  const { selectedDayId, setSelectedDayId, isBacklogOpen, toggleBacklog } = useMapUIStore();
  const { t, plural } = useTranslation();
  const [isExportOpen, setIsExportOpen] = useState(false);

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

  // Compute leg distances and total day distance (CAND-01)
  const places = activeDay?.places || [];
  const legs: { km: number; durationText: string; isWalk: boolean }[] = [];
  let totalDistanceKm = 0;

  for (let i = 0; i < places.length - 1; i++) {
    const p1 = places[i];
    const p2 = places[i + 1];
    const km = calculateHaversineKm(p1.lat, p1.lng, p2.lat, p2.lng);
    totalDistanceKm += km;
    const isWalk = km < 2.5;
    const durationMin = isWalk ? Math.max(5, Math.round(km * 12)) : Math.max(8, Math.round(km * 2.2) + 5);
    const modeLabel = isWalk ? t("distance.walk") : t("distance.drive");
    legs.push({
      km,
      durationText: `${km.toFixed(1)} ${t("distance.km")} • ~${durationMin} min ${modeLabel}`,
      isWalk,
    });
  }

  // Route Export Handlers (CAND-08)
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
          {/* Day Header with Distance & Export Controls */}
          <div className="flex items-center justify-between text-xs px-1">
            <div>
              <span className="font-semibold text-[var(--text)]">{activeDay?.title}</span>
              {totalDistanceKm > 0 && (
                <span className="ml-2 px-2 py-0.5 rounded-full bg-sky-500/10 text-sky-400 border border-sky-500/20 text-[10px] font-bold">
                  {t("trip.totalRouteDistance")}: {totalDistanceKm.toFixed(1)} {t("distance.km")}
                </span>
              )}
            </div>

            {/* Export Menu Dropdown (CAND-08) */}
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
                    <DraggablePlaceItem
                      item={item}
                      orderNumber={index + 1}
                    />

                    {/* Transit Leg Connector Badge between consecutive stops (CAND-01) */}
                    {index < legs.length && (
                      <div className="flex items-center justify-center my-0.5">
                        <div className="flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/[0.04] border border-white/10 backdrop-blur-md text-[10px] text-sky-400 font-medium shadow-sm">
                          {legs[index].isWalk ? (
                            <Footprints className="w-3 h-3 text-sky-300" />
                          ) : (
                            <Car className="w-3 h-3 text-emerald-300" />
                          )}
                          <span>{legs[index].durationText}</span>
                        </div>
                      </div>
                    )}
                  </React.Fragment>
                ))}
              </div>
            </SortableContext>
          )}
        </div>
      </DndContext>
    </div>
  );
}
