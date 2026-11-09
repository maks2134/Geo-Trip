"use client";

import React from "react";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical, Clock, DollarSign, Footprints, Train, Car, Navigation } from "lucide-react";
import { PlaceRef } from "@/entities/trip";
import { Badge } from "@/shared/ui/Badge";
import { useMapUIStore } from "@/stores/map-ui-store";
import { cn } from "@/shared/lib/utils";
import { useTranslation } from "@/shared/i18n/i18n-context";
import { buildSinglePlaceUrl } from "@/shared/lib/navigation";

interface DraggablePlaceItemProps {
  item: PlaceRef;
  orderNumber: number;
}

export function DraggablePlaceItem({ item, orderNumber }: DraggablePlaceItemProps) {
  const { t } = useTranslation();
  const { openPlaceDrawer, navAppPreference } = useMapUIStore();
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: item.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div ref={setNodeRef} style={style} className="group relative">
      <div
        className={cn(
          "flex items-center gap-3 p-3 rounded-xl border transition-all duration-200",
          "bg-white/[0.04] border-white/5 hover:border-white/15 hover:bg-white/[0.07]",
          isDragging && "opacity-50 border-sky-400 bg-sky-500/10 shadow-xl z-30"
        )}
      >
        <button
          {...attributes}
          {...listeners}
          aria-label={t("trip.reorder")}
          className="text-[var(--text-muted)] hover:text-[var(--text)] cursor-grab active:cursor-grabbing p-1"
        >
          <GripVertical className="w-4 h-4" />
        </button>

        <div className="w-6 h-6 rounded-full bg-sky-500 text-white font-bold text-xs flex items-center justify-center shrink-0">
          {orderNumber}
        </div>

        <div
          className="flex-1 min-w-0 cursor-pointer"
          onClick={() => openPlaceDrawer(item.placeId)}
        >
          <div className="flex items-center gap-2">
            <h4 className="font-semibold text-sm truncate text-[var(--text)] group-hover:text-sky-400 transition-colors">
              {item.name}
            </h4>
            <Badge category={item.category}>{t(`filters.${item.category}`)}</Badge>
          </div>

          <div className="flex items-center gap-3 mt-1 text-xs text-[var(--text-muted)]">
            {item.durationMinutes && (
              <span className="flex items-center gap-1">
                <Clock className="w-3 h-3" />
                {item.durationMinutes} {t("place.minutes")}
              </span>
            )}
            {item.cost !== undefined && item.cost > 0 && (
              <span className="flex items-center gap-1">
                <DollarSign className="w-3 h-3" />
                ${item.cost}
              </span>
            )}
            {item.notes && <span className="truncate max-w-[140px] italic">"{item.notes}"</span>}
          </div>
        </div>

        {/* 1-Tap Navigation Link (FEAT-01) */}
        <a
          href={buildSinglePlaceUrl(
            { latitude: item.lat, longitude: item.lng, name: item.name },
            navAppPreference
          )}
          target="_blank"
          rel="noreferrer"
          title={t("nav.openInMaps")}
          onClick={(e) => e.stopPropagation()}
          className="p-1.5 rounded-lg bg-white/5 hover:bg-sky-500/20 text-[var(--text-muted)] hover:text-sky-300 transition-colors shrink-0"
        >
          <Navigation className="w-3.5 h-3.5" />
        </a>
      </div>

      {item.transitToNext && (
        <div className="flex items-center gap-2 pl-9 py-1 text-[11px] text-[var(--text-muted)]">
          <div className="w-px h-4 bg-white/20" />
          {item.transitToNext.mode === "walk" && <Footprints className="w-3 h-3 text-emerald-400" />}
          {item.transitToNext.mode === "transit" && <Train className="w-3 h-3 text-sky-400" />}
          {item.transitToNext.mode === "drive" && <Car className="w-3 h-3 text-amber-400" />}
          <span>
            {item.transitToNext.durationMinutes} {t("place.minutes")} ({item.transitToNext.distanceKm} {t("distance.km")})
          </span>
        </div>
      )}
    </div>
  );
}
