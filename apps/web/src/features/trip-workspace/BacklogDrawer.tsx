"use client";

import React from "react";
import { useDroppable } from "@dnd-kit/core";
import { SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";
import { Bookmark } from "lucide-react";
import { PlaceRef } from "@/entities/trip";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { DraggablePlaceItem } from "./dnd/DraggablePlaceItem";
import { useTranslation } from "@/shared/i18n/i18n-context";

interface BacklogDrawerProps {
  items: PlaceRef[];
  isOpen: boolean;
  onClose: () => void;
}

export function BacklogDrawer({ items, isOpen }: BacklogDrawerProps) {
  const { t, plural } = useTranslation();
  const { setNodeRef, isOver } = useDroppable({
    id: "droppable-backlog",
  });

  if (!isOpen) return null;

  return (
    <GlassPanel className="p-4 space-y-3 border-dashed border-sky-400/40 bg-sky-950/20">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Bookmark className="w-4 h-4 text-sky-400" />
          <h3 className="font-semibold text-xs uppercase tracking-wider text-sky-400">
            {t("trip.backlogTitle")}
          </h3>
        </div>
        <span className="text-xs px-2 py-0.5 rounded-full bg-sky-500/20 text-sky-300 font-semibold">
          {plural(items.length, {
            one: "place",
            few: "places",
            many: "places",
          })}
        </span>
      </div>

      <div
        ref={setNodeRef}
        className={`min-h-[90px] p-2 rounded-xl transition-colors space-y-2 ${
          isOver ? "bg-sky-500/20 ring-1 ring-sky-400" : "bg-black/10"
        }`}
      >
        {items.length === 0 ? (
          <div className="text-center py-6 text-xs text-[var(--text-muted)]">
            {t("trip.backlogEmpty")}
          </div>
        ) : (
          <SortableContext items={items.map((i) => i.id)} strategy={verticalListSortingStrategy}>
            {items.map((item, idx) => (
              <DraggablePlaceItem key={item.id} item={item} orderNumber={idx + 1} />
            ))}
          </SortableContext>
        )}
      </div>
    </GlassPanel>
  );
}
