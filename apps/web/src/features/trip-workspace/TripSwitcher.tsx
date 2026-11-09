"use client";

import React, { useState, useRef } from "react";
import { ChevronDown, Share2, Calendar, Check, Camera, Loader2, Image as ImageIcon } from "lucide-react";
import { Trip } from "@/entities/trip";
import { useMapUIStore } from "@/stores/map-ui-store";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { cn } from "@/shared/lib/utils";
import { useTranslation } from "@/shared/i18n/i18n-context";
import { apiClient } from "@/shared/lib/api";

interface TripSwitcherProps {
  trips: Trip[];
  activeTrip: Trip;
  onTripUpdate?: (trip: Trip) => void;
}

export function TripSwitcher({ trips, activeTrip, onTripUpdate }: TripSwitcherProps) {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  const [isUploadingCover, setIsUploadingCover] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { setSelectedTripId } = useMapUIStore();

  const handleCopyShare = () => {
    if (typeof window !== "undefined") {
      const shareUrl = `${window.location.origin}/share/${activeTrip.slug}`;
      navigator.clipboard.writeText(shareUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > 10 * 1024 * 1024) {
      alert("Cover image must be 10MB or less");
      return;
    }

    setIsUploadingCover(true);
    try {
      const presign = await apiClient.createPresignedUpload(activeTrip.id, {
        filename: file.name,
        contentType: file.type,
        fileSize: file.size,
        visibility: "public",
        ownerType: "trip_cover",
      });

      await apiClient.uploadFileToMinIO(presign.uploadUrl, file);
      await apiClient.confirmUpload(activeTrip.id, presign.attachmentId);
      await apiClient.updateTripCover(activeTrip.id, presign.downloadUrl);

      const updated = { ...activeTrip, coverImageUrl: presign.downloadUrl };
      onTripUpdate?.(updated);
    } catch (err: any) {
      console.error("Cover upload error:", err);
    } finally {
      setIsUploadingCover(false);
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  };

  return (
    <div className="relative space-y-2">
      <input
        ref={fileInputRef}
        type="file"
        accept="image/jpeg,image/png,image/webp"
        className="hidden"
        onChange={handleFileChange}
      />

      {/* Cover Image Banner if exists */}
      {activeTrip.coverImageUrl && (
        <div className="relative w-full h-24 rounded-2xl overflow-hidden group shadow-inner border border-white/10">
          <img
            src={activeTrip.coverImageUrl}
            alt={activeTrip.title}
            className="w-full h-full object-cover transition-transform group-hover:scale-105 duration-500"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent" />
          <button
            onClick={() => fileInputRef.current?.click()}
            disabled={isUploadingCover}
            title={t("trip.changeCover")}
            className="absolute bottom-2 right-2 px-2.5 py-1 rounded-lg bg-black/60 hover:bg-black/80 text-white text-[11px] font-medium backdrop-blur-md border border-white/15 flex items-center gap-1.5 transition-all opacity-85 hover:opacity-100"
          >
            {isUploadingCover ? (
              <Loader2 className="w-3.5 h-3.5 animate-spin" />
            ) : (
              <Camera className="w-3.5 h-3.5" />
            )}
            <span>{isUploadingCover ? t("place.uploading") : t("trip.changeCover")}</span>
          </button>
        </div>
      )}

      {/* Trip Switcher Header */}
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

        <div className="flex items-center gap-1">
          {!activeTrip.coverImageUrl && (
            <button
              onClick={() => fileInputRef.current?.click()}
              disabled={isUploadingCover}
              title={t("trip.uploadCover")}
              className="p-2 rounded-xl bg-white/5 hover:bg-white/10 text-[var(--text)] transition-colors flex items-center gap-1 text-xs"
            >
              {isUploadingCover ? (
                <Loader2 className="w-3.5 h-3.5 animate-spin text-sky-400" />
              ) : (
                <Camera className="w-3.5 h-3.5" />
              )}
            </button>
          )}

          <button
            onClick={handleCopyShare}
            title={t("share.title")}
            className="p-2 rounded-xl bg-white/5 hover:bg-white/10 text-[var(--text)] transition-colors flex items-center gap-1.5 text-xs"
          >
            {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Share2 className="w-3.5 h-3.5" />}
            <span>{copied ? t("trip.shareLinkCopied") : t("nav.share")}</span>
          </button>
        </div>
      </div>

      {/* Switcher Dropdown */}
      {isOpen && (
        <GlassPanel className="absolute top-full left-0 right-0 mt-2 p-2 z-50 space-y-1 shadow-2xl">
          <div className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)] px-3 py-1">
            {t("trip.myTrips")}
          </div>
          {trips.map((tItem) => (
            <button
              key={tItem.id}
              onClick={() => {
                setSelectedTripId(tItem.id);
                setIsOpen(false);
              }}
              className={cn(
                "w-full text-left p-2.5 rounded-xl text-xs font-medium transition-all flex items-center justify-between",
                tItem.id === activeTrip.id
                  ? "bg-[var(--accent)] text-white"
                  : "text-[var(--text)] hover:bg-white/5"
              )}
            >
              <div>
                <div className="font-semibold">{tItem.title}</div>
                <div className="text-[10px] opacity-75">{tItem.cities.join(", ")}</div>
              </div>
              <span className="text-[10px] uppercase px-2 py-0.5 rounded-full bg-black/20">
                {tItem.status}
              </span>
            </button>
          ))}
        </GlassPanel>
      )}
    </div>
  );
}
