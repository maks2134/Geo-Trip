"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { Calendar, Copy, Check, Clock, Globe, ShieldCheck, Image as ImageIcon, Loader2 } from "lucide-react";
import { PublicTripView } from "@/entities/trip";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { Badge } from "@/shared/ui/Badge";
import { useTranslation } from "@/shared/i18n/i18n-context";
import { apiClient } from "@/shared/lib/api";

interface PublicShareViewProps {
  trip: PublicTripView;
}

export function PublicShareView({ trip }: PublicShareViewProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const [selectedDayId, setSelectedDayId] = useState(trip.days[0]?.id || "");
  const [isCopying, setIsCopying] = useState(false);
  const [copied, setCopied] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const activeDay = trip.days.find((d) => d.id === selectedDayId) || trip.days[0];

  const handleCopyItinerary = async () => {
    setIsCopying(true);
    setErrorMessage(null);
    try {
      const cloned = await apiClient.clonePublicTrip(trip.slug);
      setCopied(true);
      setTimeout(() => {
        router.push(`/?tripId=${cloned.id}`);
      }, 1200);
    } catch (err: any) {
      setErrorMessage(err?.message || "Failed to copy itinerary");
      setIsCopying(false);
    }
  };

  return (
    <div className="relative w-full h-full overflow-hidden pointer-events-none">
      <div className="absolute top-6 left-6 bottom-6 w-[420px] z-20 flex flex-col pointer-events-none">
        <GlassPanel className="h-full flex flex-col pointer-events-auto shadow-2xl p-4 space-y-4 overflow-hidden">
          {/* Trip Header & Cover */}
          <div className="space-y-3 pb-3 border-b border-white/10 shrink-0">
            {trip.coverImageUrl && (
              <div className="relative w-full h-32 rounded-xl overflow-hidden shadow-inner group">
                <img
                  src={trip.coverImageUrl}
                  alt={trip.title}
                  className="w-full h-full object-cover"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent" />
                <div className="absolute top-2 right-2 px-2 py-0.5 rounded-full bg-black/50 backdrop-blur-md text-[10px] text-white/90 font-medium border border-white/10 flex items-center gap-1">
                  <ShieldCheck className="w-3 h-3 text-emerald-400" />
                  <span>{t("share.readOnlyBadge")}</span>
                </div>
              </div>
            )}

            {!trip.coverImageUrl && (
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-1.5 text-xs font-semibold text-sky-400">
                  <Globe className="w-4 h-4" />
                  <span>{t("share.title")}</span>
                </div>
                <div className="px-2 py-0.5 rounded-full bg-sky-500/15 border border-sky-400/20 text-[10px] text-sky-300 font-medium">
                  {t("share.readOnlyBadge")}
                </div>
              </div>
            )}

            <div>
              <h1 className="text-xl font-black text-[var(--text)] tracking-tight">{trip.title}</h1>
              {trip.description && (
                <p className="text-xs text-[var(--text-muted)] line-clamp-2 mt-1">{trip.description}</p>
              )}
            </div>

            <div className="flex items-center justify-between text-xs text-[var(--text-muted)] pt-0.5">
              <div className="flex items-center gap-1.5">
                <Calendar className="w-3.5 h-3.5" />
                <span>
                  {trip.startDate} — {trip.endDate}
                </span>
              </div>
              {trip.cities && trip.cities.length > 0 && (
                <span className="text-[11px] font-medium text-[var(--text-muted)]">
                  {trip.cities.join(", ")}
                </span>
              )}
            </div>
          </div>

          {/* Fork / Clone CTA button (CAND-10) */}
          <div className="shrink-0 space-y-1.5">
            <button
              onClick={handleCopyItinerary}
              disabled={isCopying || copied}
              className={`w-full flex items-center justify-center gap-2 py-2.5 px-4 rounded-xl font-semibold text-xs shadow-lg transition-all ${
                copied
                  ? "bg-emerald-500 text-white shadow-emerald-500/25"
                  : "bg-[var(--accent)] hover:opacity-95 text-white shadow-[var(--accent)]/20 active:scale-[0.99]"
              }`}
            >
              {isCopying && <Loader2 className="w-4 h-4 animate-spin" />}
              {!isCopying && copied && <Check className="w-4 h-4 text-white" />}
              {!isCopying && !copied && <Copy className="w-4 h-4" />}
              <span>
                {isCopying
                  ? t("share.copying")
                  : copied
                  ? t("share.copySuccess")
                  : t("share.copyToMyTrips")}
              </span>
            </button>
            {errorMessage && (
              <div className="text-[11px] text-rose-400 text-center font-medium">{errorMessage}</div>
            )}
          </div>

          {/* Day selection tabs */}
          <div className="flex items-center gap-1.5 overflow-x-auto pb-1 shrink-0 scrollbar-none">
            {trip.days.map((day) => {
              const active = day.id === activeDay?.id;
              return (
                <button
                  key={day.id}
                  onClick={() => setSelectedDayId(day.id)}
                  className={`px-3 py-1.5 rounded-xl text-xs font-semibold whitespace-nowrap transition-all ${
                    active
                      ? "bg-[var(--accent)] text-white shadow-sm"
                      : "bg-white/[0.04] text-[var(--text-muted)] hover:text-white hover:bg-white/[0.08]"
                  }`}
                >
                  {t("trip.day")} {day.dayNumber}
                </button>
              );
            })}
          </div>

          {/* Place itinerary list */}
          <div className="flex-1 overflow-y-auto space-y-2 pr-1">
            <div className="text-xs font-semibold text-[var(--text)] px-1">{activeDay?.title}</div>

            {activeDay && activeDay.places.length === 0 ? (
              <div className="text-center py-8 text-xs text-[var(--text-muted)]">
                {t("trip.noPlacesScheduled")}
              </div>
            ) : (
              activeDay?.places.map((place) => (
                <div
                  key={place.id}
                  className="flex flex-col gap-2 p-3 rounded-xl bg-white/[0.03] border border-white/5 hover:border-white/10 transition-colors"
                >
                  <div className="flex items-start gap-3">
                    <div className="w-6 h-6 rounded-full bg-sky-500 text-white font-bold text-xs flex items-center justify-center shrink-0 mt-0.5">
                      {place.orderNumber}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <h4 className="font-semibold text-sm truncate text-[var(--text)]">{place.name}</h4>
                        <Badge category={place.category}>{t(`filters.${place.category}`)}</Badge>
                      </div>
                      {place.durationMinutes && (
                        <div className="flex items-center gap-1 text-[11px] text-[var(--text-muted)] mt-1">
                          <Clock className="w-3 h-3" />
                          <span>
                            {place.durationMinutes} {t("place.minutes")}
                          </span>
                        </div>
                      )}
                      {place.notes && (
                        <p className="text-xs text-[var(--text-muted)] mt-1.5 italic">"{place.notes}"</p>
                      )}
                    </div>
                  </div>

                  {/* Public Photos if available */}
                  {place.attachments && place.attachments.length > 0 && (
                    <div className="flex items-center gap-2 pl-9 pt-1">
                      {place.attachments.map((att) => (
                        <a
                          key={att.id}
                          href={att.url}
                          target="_blank"
                          rel="noreferrer"
                          className="relative w-12 h-12 rounded-lg overflow-hidden border border-white/10 hover:border-white/30 transition-all shrink-0 group"
                        >
                          <img
                            src={att.url}
                            alt={att.name}
                            className="w-full h-full object-cover group-hover:scale-105 transition-transform"
                          />
                        </a>
                      ))}
                    </div>
                  )}
                </div>
              ))
            )}
          </div>

          {/* Privacy Guarantee footer */}
          <div className="shrink-0 pt-2 border-t border-white/5 flex items-center gap-2 text-[10px] text-[var(--text-muted)]">
            <ShieldCheck className="w-3.5 h-3.5 text-emerald-400 shrink-0" />
            <span>{t("share.privateNoteNotice")}</span>
          </div>
        </GlassPanel>
      </div>
    </div>
  );
}
