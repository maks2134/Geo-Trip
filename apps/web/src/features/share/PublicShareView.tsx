"use client";

import React, { useState } from "react";
import { Calendar, Copy, Check, Clock, Globe } from "lucide-react";
import { PublicTripView } from "@/entities/trip";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { Badge } from "@/shared/ui/Badge";

interface PublicShareViewProps {
  trip: PublicTripView;
}

export function PublicShareView({ trip }: PublicShareViewProps) {
  const [selectedDayId, setSelectedDayId] = useState(trip.days[0]?.id || "");
  const [copied, setCopied] = useState(false);

  const activeDay = trip.days.find((d) => d.id === selectedDayId) || trip.days[0];

  const handleCopyItinerary = () => {
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative w-full h-full overflow-hidden">
      <div className="absolute top-6 left-6 bottom-6 w-96 z-20 flex flex-col pointer-events-none">
        <GlassPanel className="h-full flex flex-col pointer-events-auto shadow-2xl p-4 space-y-4 overflow-hidden">
          <div className="space-y-2 pb-3 border-b border-white/10">
            <div className="flex items-center gap-2 text-xs font-semibold text-sky-400">
              <Globe className="w-4 h-4" />
              <span>Public Guest Itinerary</span>
            </div>
            <h1 className="text-xl font-black text-[var(--text)]">{trip.title}</h1>
            <p className="text-xs text-[var(--text-muted)] line-clamp-2">{trip.description}</p>
            <div className="flex items-center gap-2 text-xs text-[var(--text-muted)] pt-1">
              <Calendar className="w-3.5 h-3.5" />
              <span>
                {trip.startDate} — {trip.endDate}
              </span>
            </div>
          </div>

          <button
            onClick={handleCopyItinerary}
            className="w-full flex items-center justify-center gap-2 py-2.5 px-4 rounded-xl bg-[var(--accent)] text-white font-semibold text-xs shadow-lg shadow-[var(--accent)]/20 hover:opacity-95 transition-all"
          >
            {copied ? <Check className="w-4 h-4 text-emerald-300" /> : <Copy className="w-4 h-4" />}
            <span>{copied ? "Itinerary Copied to Your Account!" : "Copy Route to My Trips"}</span>
          </button>

          <div className="flex items-center gap-1.5 overflow-x-auto pb-1 scrollbar-none">
            {trip.days.map((day) => {
              const active = day.id === activeDay?.id;
              return (
                <button
                  key={day.id}
                  onClick={() => setSelectedDayId(day.id)}
                  className={`px-3 py-1.5 rounded-xl text-xs font-semibold whitespace-nowrap transition-all ${
                    active
                      ? "bg-white/15 text-white shadow-sm"
                      : "bg-white/[0.03] text-[var(--text-muted)] hover:text-white"
                  }`}
                >
                  Day {day.dayNumber}
                </button>
              );
            })}
          </div>

          <div className="flex-1 overflow-y-auto space-y-2 pr-1">
            <div className="text-xs font-semibold text-[var(--text)] px-1">{activeDay?.title}</div>

            {activeDay && activeDay.places.length === 0 ? (
              <div className="text-center py-8 text-xs text-[var(--text-muted)]">No places listed for this day.</div>
            ) : (
              activeDay?.places.map((place) => (
                <div
                  key={place.id}
                  className="flex items-start gap-3 p-3 rounded-xl bg-white/[0.03] border border-white/5"
                >
                  <div className="w-6 h-6 rounded-full bg-sky-500 text-white font-bold text-xs flex items-center justify-center shrink-0 mt-0.5">
                    {place.orderNumber}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <h4 className="font-semibold text-sm truncate text-[var(--text)]">{place.name}</h4>
                      <Badge category={place.category}>{place.category}</Badge>
                    </div>
                    {place.durationMinutes && (
                      <div className="flex items-center gap-1 text-[11px] text-[var(--text-muted)] mt-1">
                        <Clock className="w-3 h-3" />
                        <span>{place.durationMinutes} min</span>
                      </div>
                    )}
                    {place.notes && <p className="text-xs text-[var(--text-muted)] mt-1.5 italic">"{place.notes}"</p>}
                  </div>
                </div>
              ))
            )}
          </div>
        </GlassPanel>
      </div>
    </div>
  );
}
