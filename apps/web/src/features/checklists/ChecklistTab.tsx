"use client";

import React from "react";
import { CheckCircle2, Circle, Plane, Luggage, ShoppingBag } from "lucide-react";
import { Trip, ChecklistItem } from "@/entities/trip";
import { ProgressBar } from "@/shared/ui/ProgressBar";
import { apiClient } from "@/shared/lib/api";

interface ChecklistTabProps {
  trip: Trip;
  onTripUpdate: (updated: Trip) => void;
}

export function ChecklistTab({ trip, onTripUpdate }: ChecklistTabProps) {
  const completedCount = trip.checklists.filter((c) => c.completed).length;
  const totalCount = trip.checklists.length;
  const progressPercent = totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0;

  const handleToggle = async (id: string) => {
    const updatedChecklists = await apiClient.toggleChecklistItem(trip.id, id);
    onTripUpdate({
      ...trip,
      checklists: updatedChecklists,
    });
  };

  const sections: {
    key: ChecklistItem["section"];
    label: string;
    icon: React.ReactNode;
  }[] = [
    { key: "pre-flight", label: "Pre-Flight Prep", icon: <Plane className="w-3.5 h-3.5 text-sky-400" /> },
    { key: "packing", label: "Luggage Packing", icon: <Luggage className="w-3.5 h-3.5 text-amber-400" /> },
    { key: "buy-on-site", label: "On-Site Arrivals", icon: <ShoppingBag className="w-3.5 h-3.5 text-emerald-400" /> },
  ];

  return (
    <div className="space-y-4 text-xs">
      <div className="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-2">
        <div className="flex items-center justify-between">
          <span className="font-semibold text-xs text-[var(--text)]">Departure Readiness</span>
          <span className="font-bold text-sky-400">{progressPercent}%</span>
        </div>
        <ProgressBar value={progressPercent} />
        <div className="text-[11px] text-[var(--text-muted)] text-right">
          {completedCount} of {totalCount} tasks completed
        </div>
      </div>

      <div className="space-y-4">
        {sections.map((sec) => {
          const items = trip.checklists.filter((c) => c.section === sec.key);
          return (
            <div key={sec.key} className="space-y-2">
              <div className="flex items-center gap-2 font-semibold text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                {sec.icon}
                <span>{sec.label}</span>
              </div>

              <div className="space-y-1.5">
                {items.map((item) => (
                  <button
                    key={item.id}
                    onClick={() => handleToggle(item.id)}
                    className="w-full flex items-center gap-2.5 p-2.5 rounded-xl bg-white/[0.02] border border-white/5 hover:bg-white/[0.05] transition-colors text-left"
                  >
                    {item.completed ? (
                      <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                    ) : (
                      <Circle className="w-4 h-4 text-[var(--text-muted)] shrink-0" />
                    )}
                    <span
                      className={`text-xs ${
                        item.completed ? "line-through text-[var(--text-muted)]" : "text-[var(--text)]"
                      }`}
                    >
                      {item.title}
                    </span>
                  </button>
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
