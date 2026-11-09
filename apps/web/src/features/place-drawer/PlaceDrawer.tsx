"use client";

import React, { useState } from "react";
import {
  X,
  Star,
  Clock,
  DollarSign,
  ExternalLink,
  FileText,
  Paperclip,
  CheckCircle2,
  Calendar,
  Bookmark,
} from "lucide-react";
import { Place, PinStatus } from "@/entities/place";
import { useMapUIStore } from "@/stores/map-ui-store";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { Badge } from "@/shared/ui/Badge";
import { cn } from "@/shared/lib/utils";

interface PlaceDrawerProps {
  places: Place[];
}

export function PlaceDrawer({ places }: PlaceDrawerProps) {
  const { drawerPlaceId, isDrawerOpen, closePlaceDrawer } = useMapUIStore();
  const place = places.find((p) => p.id === drawerPlaceId);

  const [status, setStatus] = useState<PinStatus>(place?.status || "planned");
  const [notes, setNotes] = useState<string>(place?.notes || "");

  if (!isDrawerOpen || !place) return null;

  return (
    <div className="fixed top-6 right-6 bottom-24 w-96 z-40 animate-in slide-in-from-right-8 duration-200">
      <GlassPanel className="h-full flex flex-col overflow-hidden shadow-2xl border-white/15">
        <div className="relative h-44 w-full bg-slate-900 shrink-0">
          {place.imageUrl ? (
            <img
              src={place.imageUrl}
              alt={place.name}
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full bg-gradient-to-tr from-sky-900 to-indigo-900 flex items-center justify-center text-sky-300">
              <FileText className="w-12 h-12 opacity-50" />
            </div>
          )}
          <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent" />

          <button
            onClick={closePlaceDrawer}
            aria-label="Close place details"
            className="absolute top-3 right-3 p-1.5 rounded-full bg-black/50 text-white hover:bg-black/80 transition-colors"
          >
            <X className="w-4 h-4" />
          </button>

          <div className="absolute bottom-3 left-4 right-4">
            <Badge category={place.category}>{place.category}</Badge>
            <h3 className="text-lg font-bold text-white mt-1 leading-tight drop-shadow-md">
              {place.name}
            </h3>
            <p className="text-xs text-white/80">
              {place.city}, {place.country}
            </p>
          </div>
        </div>

        <div className="flex-1 overflow-y-auto p-4 space-y-4 text-xs">
          <div className="space-y-1.5">
            <span className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
              Travel Status
            </span>
            <div className="grid grid-cols-3 gap-1.5">
              {(["visited", "planned", "wishlist"] as PinStatus[]).map((s) => (
                <button
                  key={s}
                  onClick={() => setStatus(s)}
                  className={cn(
                    "py-1.5 rounded-xl capitalize font-medium border transition-all flex items-center justify-center gap-1",
                    status === s
                      ? "bg-[var(--accent)] text-white border-[var(--accent)]"
                      : "bg-white/[0.04] text-[var(--text-muted)] border-white/5 hover:bg-white/[0.08]"
                  )}
                >
                  {s === "visited" && <CheckCircle2 className="w-3 h-3 text-amber-300" />}
                  {s === "planned" && <Calendar className="w-3 h-3 text-blue-300" />}
                  {s === "wishlist" && <Bookmark className="w-3 h-3 text-pink-300" />}
                  <span>{s}</span>
                </button>
              ))}
            </div>
          </div>

          <div className="grid grid-cols-3 gap-2 py-2 border-y border-white/5">
            {place.rating && (
              <div className="p-2 rounded-xl bg-white/[0.03] text-center">
                <div className="flex items-center justify-center gap-1 font-bold text-[var(--text)]">
                  <Star className="w-3.5 h-3.5 text-amber-400 fill-amber-400" />
                  <span>{place.rating}</span>
                </div>
                <div className="text-[10px] text-[var(--text-muted)]">Rating</div>
              </div>
            )}
            <div className="p-2 rounded-xl bg-white/[0.03] text-center">
              <div className="flex items-center justify-center gap-1 font-bold text-[var(--text)]">
                <Clock className="w-3.5 h-3.5 text-sky-400" />
                <span>{place.durationMinutes || 60}m</span>
              </div>
              <div className="text-[10px] text-[var(--text-muted)]">Duration</div>
            </div>
            <div className="p-2 rounded-xl bg-white/[0.03] text-center">
              <div className="flex items-center justify-center gap-1 font-bold text-[var(--text)]">
                <DollarSign className="w-3.5 h-3.5 text-emerald-400" />
                <span>${place.estimatedCost ?? 0}</span>
              </div>
              <div className="text-[10px] text-[var(--text-muted)]">Cost</div>
            </div>
          </div>

          <div className="space-y-1.5">
            <span className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
              Travel Notes
            </span>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Add personal notes, reservations, or subway tips..."
              rows={3}
              className="w-full p-2.5 rounded-xl bg-white/[0.03] border border-white/10 text-[var(--text)] placeholder:text-[var(--text-muted)] focus:outline-none focus:ring-1 focus:ring-[var(--accent)] resize-none"
            />
          </div>

          {place.externalLink && (
            <div className="space-y-1.5">
              <span className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                Official Guide
              </span>
              <a
                href={place.externalLink}
                target="_blank"
                rel="noreferrer"
                className="flex items-center justify-between p-2.5 rounded-xl bg-white/[0.03] border border-white/10 text-sky-400 hover:bg-white/[0.06] transition-colors"
              >
                <span className="truncate">{place.externalLink}</span>
                <ExternalLink className="w-3.5 h-3.5 shrink-0 ml-2" />
              </a>
            </div>
          )}

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                Attachments & Vouchers
              </span>
              <span className="text-[10px] text-[var(--text-muted)]">
                {place.attachments?.length || 0} files
              </span>
            </div>

            {place.attachments && place.attachments.length > 0 ? (
              <div className="space-y-1.5">
                {place.attachments.map((att) => (
                  <div
                    key={att.id}
                    className="flex items-center justify-between p-2.5 rounded-xl bg-white/[0.03] border border-white/5"
                  >
                    <div className="flex items-center gap-2 min-w-0">
                      <Paperclip className="w-3.5 h-3.5 text-sky-400 shrink-0" />
                      <span className="truncate font-medium">{att.name}</span>
                    </div>
                    <a
                      href={att.url}
                      target="_blank"
                      rel="noreferrer"
                      className="text-xs text-sky-400 hover:underline shrink-0 ml-2 font-medium"
                    >
                      View
                    </a>
                  </div>
                ))}
              </div>
            ) : (
              <div className="p-3 rounded-xl bg-white/[0.02] border border-dashed border-white/10 text-center text-[var(--text-muted)]">
                No tickets attached yet
              </div>
            )}
          </div>
        </div>
      </GlassPanel>
    </div>
  );
}
