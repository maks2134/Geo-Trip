"use client";

import React, { useState } from "react";
import { ChevronDown, Share2, Calendar, Check } from "lucide-react";
import { Trip } from "@/entities/trip";
import { useMapUIStore } from "@/stores/map-ui-store";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { cn } from "@/shared/lib/utils";

interface TripSwitcherProps {
  trips: Trip[];
  activeTrip: Trip;
}

export function TripSwitcher({ trips, activeTrip }: TripSwitcherProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  const { setSelectedTripId } = useMapUIStore();

  const handleCopyShare = () => {
    if (typeof window !== "undefined") {
      const shareUrl = `${window.location.origin}/share/${activeTrip.slug}`;
      navigator.clipboard.writeText(shareUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div className="relative">
      <div className="flex items-center justify-between gap-2 p-3 bg-white/[0.03] rounded-2xl border border-white/5">
        <button
          onClick={() => setIsOpen(!isOpen)}
          className="flex-1 flex items-center justify-between text-left hover:opacity-85 transition-opacity"
        >
          <div>
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
              <h2 className="font-bold text-sm text-[var(--text)]">{activeTrip.title}</h2>
            </div>
            <div className="flex items-center gap-2 text-xs text-[var(--text-muted)] mt-0.5">
              <Calendar className="w-3 h-3" />
              <span>
                {activeTrip.startDate} — {activeTrip.endDate}
              </span>
            </div>
          </div>
          <ChevronDown className={cn("w-4 h-4 transition-transform text-[var(--text-muted)]", isOpen && "rotate-180")} />
        </button>

        <button
          onClick={handleCopyShare}
          title="Copy Public Share Link"
          className="p-2 rounded-xl bg-white/5 hover:bg-white/10 text-[var(--text)] transition-colors flex items-center gap-1.5 text-xs"
        >
          {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Share2 className="w-3.5 h-3.5" />}
          <span>{copied ? "Copied!" : "Share"}</span>
        </button>
      </div>

      {isOpen && (
        <GlassPanel className="absolute top-full left-0 right-0 mt-2 p-2 z-50 space-y-1 shadow-2xl">
          <div className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)] px-3 py-1">
            My Trips
          </div>
          {trips.map((t) => (
            <button
              key={t.id}
              onClick={() => {
                setSelectedTripId(t.id);
                setIsOpen(false);
              }}
              className={cn(
                "w-full text-left p-2.5 rounded-xl text-xs font-medium transition-all flex items-center justify-between",
                t.id === activeTrip.id
                  ? "bg-[var(--accent)] text-white"
                  : "text-[var(--text)] hover:bg-white/5"
              )}
            >
              <div>
                <div className="font-semibold">{t.title}</div>
                <div className="text-[10px] opacity-75">{t.cities.join(", ")}</div>
              </div>
              <span className="text-[10px] uppercase px-2 py-0.5 rounded-full bg-black/20">
                {t.status}
              </span>
            </button>
          ))}
        </GlassPanel>
      )}
    </div>
  );
}
