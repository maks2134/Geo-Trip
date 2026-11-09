"use client";

import React from "react";
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
import { Bookmark } from "lucide-react";
import { Trip } from "@/entities/trip";
import { useMapUIStore } from "@/stores/map-ui-store";
import { DraggablePlaceItem } from "./dnd/DraggablePlaceItem";
import { BacklogDrawer } from "./BacklogDrawer";
import { apiClient } from "@/shared/lib/api";
import { cn } from "@/shared/lib/utils";

interface DayTimelineProps {
  trip: Trip;
  onTripUpdate: (updated: Trip) => void;
}

export function DayTimeline({ trip, onTripUpdate }: DayTimelineProps) {
  const { selectedDayId, setSelectedDayId, isBacklogOpen, toggleBacklog } = useMapUIStore();

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

  return (
    <div className="space-y-4">
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
              <span>Day {day.dayNumber}</span>
              <span className="text-[10px] opacity-75 font-normal">
                {day.places.length} stops
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
          <span>Backlog ({trip.backlog.length})</span>
        </button>
      </div>

      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <BacklogDrawer items={trip.backlog} isOpen={isBacklogOpen} onClose={toggleBacklog} />

        <div className="space-y-2">
          <div className="flex items-center justify-between text-xs px-1">
            <span className="font-semibold text-[var(--text)]">{activeDay?.title}</span>
            <span className="text-[var(--text-muted)]">{activeDay?.date}</span>
          </div>

          {activeDay && activeDay.places.length === 0 ? (
            <div className="text-center py-10 px-4 rounded-xl border border-dashed border-white/10 text-xs text-[var(--text-muted)]">
              No places scheduled for Day {activeDay.dayNumber}. Drag pins from backlog or click pins on map to add!
            </div>
          ) : (
            <SortableContext
              items={activeDay ? activeDay.places.map((p) => p.id) : []}
              strategy={verticalListSortingStrategy}
            >
              <div className="space-y-2">
                {activeDay?.places.map((item, index) => (
                  <DraggablePlaceItem
                    key={item.id}
                    item={item}
                    orderNumber={index + 1}
                  />
                ))}
              </div>
            </SortableContext>
          )}
        </div>
      </DndContext>
    </div>
  );
}
