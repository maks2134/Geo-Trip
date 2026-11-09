"use client";

import React, { useState, useRef, useEffect } from "react";
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
  Upload,
  Lock,
  Globe,
  Maximize2,
  Navigation,
  MapPin,
  Plus,
  Loader2,
  Image as ImageIcon,
} from "lucide-react";
import { Place, PinCategory, PinStatus, PlaceAttachment } from "@/entities/place";
import { useMapUIStore } from "@/stores/map-ui-store";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { Badge } from "@/shared/ui/Badge";
import { cn } from "@/shared/lib/utils";
import { useTranslation } from "@/shared/i18n/i18n-context";
import { apiClient } from "@/shared/lib/api";
import { buildSinglePlaceUrl } from "@/shared/lib/navigation";

interface PlaceDrawerProps {
  places: Place[];
}

export function PlaceDrawer({ places }: PlaceDrawerProps) {
  const {
    drawerPlaceId,
    isDrawerOpen,
    closePlaceDrawer,
    selectedTripId,
    selectedDayId,
    droppedPin,
    setDroppedPin,
    navAppPreference,
  } = useMapUIStore();
  const { t } = useTranslation();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const place = places.find((p) => p.id === drawerPlaceId);

  // Quick Pin Form State (FEAT-02)
  const [quickName, setQuickName] = useState("");
  const [quickCategory, setQuickCategory] = useState<PinCategory>("sight");
  const [quickDayId, setQuickDayId] = useState<string | null>(selectedDayId || "day-1");
  const [isAddingQuick, setIsAddingQuick] = useState(false);

  useEffect(() => {
    if (droppedPin) {
      setQuickName(
        droppedPin.name ||
          `${t("quickPin.droppedPin")} (${droppedPin.lat.toFixed(4)}, ${droppedPin.lng.toFixed(4)})`
      );
      setQuickDayId(selectedDayId || "day-1");
    }
  }, [droppedPin, selectedDayId, t]);

  const [status, setStatus] = useState<PinStatus>(place?.status || "planned");
  const [notes, setNotes] = useState<string>(place?.notes || "");
  const [attachments, setAttachments] = useState<PlaceAttachment[]>(place?.attachments || []);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [lightboxUrl, setLightboxUrl] = useState<string | null>(null);
  const [isPrivateUpload, setIsPrivateUpload] = useState<boolean>(true);

  useEffect(() => {
    if (place) {
      setStatus(place.status);
      setNotes(place.notes || "");
      setAttachments(place.attachments || []);
    }
  }, [place]);

  if (!isDrawerOpen) return null;

  // Handle Quick Add Pin Submission (FEAT-02)
  const handleQuickAdd = async () => {
    if (!droppedPin) return;
    setIsAddingQuick(true);
    try {
      await apiClient.addPlaceToItinerary(selectedTripId, {
        name: quickName.trim() || `Place (${droppedPin.lat.toFixed(4)}, ${droppedPin.lng.toFixed(4)})`,
        category: quickCategory,
        lat: droppedPin.lat,
        lng: droppedPin.lng,
        dayId: quickDayId,
        city: droppedPin.city,
        country: droppedPin.country,
        address: droppedPin.address,
      });
      setDroppedPin(null);
      closePlaceDrawer();
    } catch {
      setIsAddingQuick(false);
    }
  };

  // If Drawer opened via Quick Pin drop (FEAT-02)
  if (droppedPin) {
    return (
      <div className="fixed top-0 right-0 h-full w-[380px] max-w-full z-40 p-4 pointer-events-none">
        <GlassPanel className="h-full flex flex-col overflow-hidden pointer-events-auto rounded-3xl shadow-2xl border border-amber-500/30">
          <div className="relative p-5 bg-gradient-to-br from-amber-500/20 via-sky-500/10 to-transparent border-b border-white/10">
            <button
              onClick={() => {
                setDroppedPin(null);
                closePlaceDrawer();
              }}
              aria-label={t("place.close")}
              className="absolute top-4 right-4 p-1.5 rounded-full bg-black/40 text-white hover:bg-black/60 transition-colors"
            >
              <X className="w-4 h-4" />
            </button>

            <div className="flex items-center gap-2 text-amber-400 font-bold text-xs uppercase tracking-wider mb-1">
              <MapPin className="w-4 h-4 animate-bounce" />
              <span>{t("quickPin.header")}</span>
            </div>
            <h3 className="text-base font-bold text-white leading-tight">
              {droppedPin.name || `${droppedPin.lat.toFixed(4)}, ${droppedPin.lng.toFixed(4)}`}
            </h3>
            <p className="text-xs text-[var(--text-muted)] mt-1 truncate">
              {droppedPin.address || `${droppedPin.lat.toFixed(5)}, ${droppedPin.lng.toFixed(5)}`}
            </p>
          </div>

          <div className="flex-1 overflow-y-auto p-5 space-y-4 text-xs">
            {/* Place Name Edit */}
            <div className="space-y-1.5">
              <label className="text-[11px] font-semibold text-[var(--text-muted)] uppercase tracking-wider">
                {t("quickPin.placeName")}
              </label>
              <input
                type="text"
                value={quickName}
                onChange={(e) => setQuickName(e.target.value)}
                className="w-full px-3 py-2 rounded-xl bg-black/[0.03] dark:bg-white/[0.04] border border-[var(--glass-border)] text-[var(--text)] font-medium focus:outline-none focus:ring-1 focus:ring-[var(--accent)]"
                placeholder="E.g. Shinjuku Gyoen Park"
              />
            </div>

            {/* Category Selector (Fix: flex-wrap + whitespace-nowrap to prevent button text collisions) */}
            <div className="space-y-1.5">
              <label className="text-[11px] font-semibold text-[var(--text-muted)] uppercase tracking-wider">
                {t("filters.category")}
              </label>
              <div className="flex flex-wrap gap-1.5">
                {(["sight", "food", "stay", "nature", "transport"] as PinCategory[]).map((cat) => (
                  <button
                    key={cat}
                    type="button"
                    onClick={() => setQuickCategory(cat)}
                    className={cn(
                      "py-1.5 px-3 rounded-xl text-center text-xs font-medium border transition-all whitespace-nowrap",
                      quickCategory === cat
                        ? "bg-[var(--accent)] text-white border-[var(--accent)] font-semibold shadow-sm"
                        : "bg-black/[0.03] dark:bg-white/[0.03] text-[var(--text-muted)] border-[var(--glass-border)] hover:bg-black/5 dark:hover:bg-white/[0.06] hover:text-[var(--text)]"
                    )}
                  >
                    {t(`filters.${cat}`)}
                  </button>
                ))}
              </div>
            </div>

            {/* Destination Day Selection */}
            <div className="space-y-1.5">
              <label className="text-[11px] font-semibold text-[var(--text-muted)] uppercase tracking-wider">
                {t("quickPin.addTo")}
              </label>
              <div className="grid grid-cols-2 gap-2">
                <button
                  type="button"
                  onClick={() => setQuickDayId("day-1")}
                  className={cn(
                    "py-2 px-3 rounded-xl border text-left font-medium transition-all min-w-0",
                    quickDayId === "day-1"
                      ? "bg-sky-500/20 text-sky-400 dark:text-sky-300 border-sky-500/40"
                      : "bg-black/[0.03] dark:bg-white/[0.03] text-[var(--text-muted)] border-[var(--glass-border)] hover:bg-black/5 dark:hover:bg-white/[0.06]"
                  )}
                >
                  <div className="font-semibold truncate">{t("quickPin.dayOne")}</div>
                  <div className="text-[10px] opacity-70 truncate">{t("quickPin.activeItinerary")}</div>
                </button>
                <button
                  type="button"
                  onClick={() => setQuickDayId(null)}
                  className={cn(
                    "py-2 px-3 rounded-xl border text-left font-medium transition-all min-w-0",
                    quickDayId === null
                      ? "bg-sky-500/20 text-sky-400 dark:text-sky-300 border-sky-500/40"
                      : "bg-black/[0.03] dark:bg-white/[0.03] text-[var(--text-muted)] border-[var(--glass-border)] hover:bg-black/5 dark:hover:bg-white/[0.06]"
                  )}
                >
                  <div className="font-semibold truncate">{t("quickPin.backlog")}</div>
                  <div className="text-[10px] opacity-70 truncate">{t("quickPin.unassignedPool")}</div>
                </button>
              </div>
            </div>

            {/* External Navigation Link Preview */}
            <div className="pt-2">
              <a
                href={buildSinglePlaceUrl(
                  { latitude: droppedPin.lat, longitude: droppedPin.lng, name: quickName },
                  navAppPreference
                )}
                target="_blank"
                rel="noreferrer"
                className="flex items-center justify-center gap-1.5 py-2 px-3 rounded-xl bg-black/[0.03] dark:bg-white/[0.04] text-[var(--text-muted)] hover:text-[var(--text)] border border-[var(--glass-border)] transition-colors"
              >
                <Navigation className="w-3.5 h-3.5 text-sky-400" />
                <span>{t("nav.openInMaps")}</span>
              </a>
            </div>
          </div>

          <div className="p-4 border-t border-[var(--glass-border)] flex items-center gap-2">
            <button
              onClick={() => {
                setDroppedPin(null);
                closePlaceDrawer();
              }}
              className="flex-1 py-2 rounded-xl text-center font-medium bg-black/5 dark:bg-white/5 text-[var(--text-muted)] hover:bg-black/10 dark:hover:bg-white/10 hover:text-[var(--text)] transition-colors"
            >
              {t("common.cancel")}
            </button>
            <button
              onClick={handleQuickAdd}
              disabled={isAddingQuick}
              className="flex-1 py-2 rounded-xl text-center font-semibold bg-amber-500 hover:bg-amber-400 text-slate-950 shadow-lg shadow-amber-500/30 flex items-center justify-center gap-1.5 transition-all"
            >
              {isAddingQuick ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <>
                  <Plus className="w-4 h-4" />
                  <span>{t("quickPin.addBtn")}</span>
                </>
              )}
            </button>
          </div>
        </GlassPanel>
      </div>
    );
  }

  if (!place) return null;

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploadError(null);

    const isImage = file.type.startsWith("image/");
    const isPdf = file.type === "application/pdf";

    if (!isImage && !isPdf) {
      setUploadError("Only JPG, PNG, WEBP and PDF files are allowed.");
      return;
    }

    if (isImage && file.size > 10 * 1024 * 1024) {
      setUploadError("Photos must not exceed 10 MB.");
      return;
    }
    if (isPdf && file.size > 25 * 1024 * 1024) {
      setUploadError("PDF documents must not exceed 25 MB.");
      return;
    }

    setIsUploading(true);
    try {
      const presign = await apiClient.createPresignedUpload(selectedTripId, {
        filename: file.name,
        contentType: file.type,
        fileSize: file.size,
        visibility: isPrivateUpload ? "private" : "public",
        placeId: place.id,
        ownerType: "place",
      });

      await apiClient.uploadFileToMinIO(presign.uploadUrl, file);
      await apiClient.confirmUpload(selectedTripId, presign.attachmentId);

      const newAtt: PlaceAttachment = {
        id: presign.attachmentId,
        name: file.name,
        fileType: isImage ? "image" : "pdf",
        url: presign.downloadUrl,
        sizeBytes: file.size,
        isPrivate: isPrivateUpload,
        uploadedAt: new Date().toISOString(),
      };

      setAttachments((prev) => [newAtt, ...prev]);
    } catch (err: any) {
      setUploadError(err.message || "Upload failed. Please try again.");
    } finally {
      setIsUploading(false);
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  };

  const imageAttachments = attachments.filter((a) => a.fileType === "image");
  const docAttachments = attachments.filter((a) => a.fileType === "pdf");

  return (
    <>
      <div className="fixed top-0 right-0 h-full w-[380px] max-w-full z-40 p-4 pointer-events-none">
        <GlassPanel className="h-full flex flex-col overflow-hidden pointer-events-auto rounded-3xl shadow-2xl">
          {/* Header Image & Actions */}
          <div className="relative h-44 w-full bg-slate-800 shrink-0 overflow-hidden">
            {place.imageUrl ? (
              <img
                src={place.imageUrl}
                alt={place.name}
                className="w-full h-full object-cover cursor-pointer"
                onClick={() => setLightboxUrl(place.imageUrl || null)}
              />
            ) : (
              <div className="w-full h-full bg-gradient-to-tr from-sky-900 to-indigo-900 flex items-center justify-center text-sky-300">
                <FileText className="w-12 h-12 opacity-50" />
              </div>
            )}
            <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent pointer-events-none" />

            <button
              onClick={closePlaceDrawer}
              aria-label={t("place.close")}
              className="absolute top-3 right-3 p-1.5 rounded-full bg-black/50 text-white hover:bg-black/80 transition-colors z-10"
            >
              <X className="w-4 h-4" />
            </button>

            <div className="absolute bottom-3 left-4 right-4 pointer-events-none">
              <Badge category={place.category}>{t(`filters.${place.category}`)}</Badge>
              <h3 className="text-lg font-bold text-white mt-1 leading-tight drop-shadow-md">
                {place.name}
              </h3>
              <p className="text-xs text-white/80">
                {place.city}, {place.country}
              </p>
            </div>
          </div>

          <div className="flex-1 overflow-y-auto p-4 space-y-4 text-xs">
            {/* Travel Status */}
            <div className="space-y-1.5">
              <span className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                {t("filters.status")}
              </span>
              <div className="flex flex-wrap gap-1.5">
                {(["visited", "planned", "wishlist"] as PinStatus[]).map((s) => (
                  <button
                    key={s}
                    onClick={() => setStatus(s)}
                    className={cn(
                      "py-1.5 px-3 rounded-xl capitalize text-xs font-medium border transition-all flex items-center justify-center gap-1.5 whitespace-nowrap",
                      status === s
                        ? "bg-[var(--accent)] text-white border-[var(--accent)] font-semibold shadow-sm"
                        : "bg-black/[0.03] dark:bg-white/[0.04] text-[var(--text-muted)] border-[var(--glass-border)] hover:bg-black/5 dark:hover:bg-white/[0.08] hover:text-[var(--text)]"
                    )}
                  >
                    {s === "visited" && <CheckCircle2 className="w-3.5 h-3.5 text-amber-500 dark:text-amber-300" />}
                    {s === "planned" && <Calendar className="w-3.5 h-3.5 text-blue-500 dark:text-blue-300" />}
                    {s === "wishlist" && <Bookmark className="w-3.5 h-3.5 text-pink-500 dark:text-pink-300" />}
                    <span>{t(`filters.${s}`)}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Quick Metrics & Navigation Action (FEAT-01) */}
            <div className="grid grid-cols-3 gap-2 py-2 border-y border-[var(--glass-border)]">
              {place.rating && (
                <div className="p-2 rounded-xl bg-black/[0.03] dark:bg-white/[0.03] text-center">
                  <div className="flex items-center justify-center gap-1 font-bold text-[var(--text)]">
                    <Star className="w-3.5 h-3.5 text-amber-400 fill-amber-400" />
                    <span>{place.rating}</span>
                  </div>
                  <div className="text-[10px] text-[var(--text-muted)]">{t("place.rating")}</div>
                </div>
              )}
              {place.durationMinutes && (
                <div className="p-2 rounded-xl bg-black/[0.03] dark:bg-white/[0.03] text-center">
                  <div className="flex items-center justify-center gap-1 font-bold text-[var(--text)]">
                    <Clock className="w-3.5 h-3.5 text-sky-400" />
                    <span>{place.durationMinutes}m</span>
                  </div>
                  <div className="text-[10px] text-[var(--text-muted)]">{t("place.duration")}</div>
                </div>
              )}
              <div className="p-2 rounded-xl bg-black/[0.03] dark:bg-white/[0.03] text-center">
                <div className="flex items-center justify-center gap-1 font-bold text-[var(--text)]">
                  <DollarSign className="w-3.5 h-3.5 text-emerald-400" />
                  <span>${place.estimatedCost ?? 0}</span>
                </div>
                <div className="text-[10px] text-[var(--text-muted)]">{t("place.estimatedCost")}</div>
              </div>
            </div>

            {/* 1-Tap Navigation Button (FEAT-01) */}
            <div className="flex gap-2">
              <a
                href={buildSinglePlaceUrl(
                  { latitude: place.lat, longitude: place.lng, name: place.name },
                  navAppPreference
                )}
                target="_blank"
                rel="noreferrer"
                className="flex-1 flex items-center justify-center gap-1.5 py-2 px-3 rounded-xl bg-sky-500/20 hover:bg-sky-500/30 text-sky-300 border border-sky-500/30 font-semibold transition-colors"
              >
                <Navigation className="w-3.5 h-3.5" />
                <span>{t("nav.openInMaps")}</span>
              </a>
              {place.externalLink && (
                <a
                  href={place.externalLink}
                  target="_blank"
                  rel="noreferrer"
                  className="p-2 rounded-xl bg-black/[0.04] dark:bg-white/[0.04] text-[var(--text-muted)] hover:text-[var(--text)] border border-[var(--glass-border)] transition-colors flex items-center justify-center"
                  title="Official Website"
                >
                  <ExternalLink className="w-4 h-4" />
                </a>
              )}
            </div>

            {/* Photo Gallery Strip */}
            {imageAttachments.length > 0 && (
              <div className="space-y-1.5">
                <span className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                  {t("place.photos")} ({imageAttachments.length})
                </span>
                <div className="flex gap-2 overflow-x-auto pb-1 scrollbar-none">
                  {imageAttachments.map((img) => (
                    <div
                      key={img.id}
                      onClick={() => setLightboxUrl(img.url)}
                      className="relative w-20 h-20 rounded-xl overflow-hidden shrink-0 border border-[var(--glass-border)] group cursor-pointer hover:scale-105 transition-transform"
                    >
                      <img src={img.url} alt={img.name} className="w-full h-full object-cover" />
                      <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity">
                        <Maximize2 className="w-4 h-4 text-white" />
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Travel Notes */}
            <div className="space-y-1.5">
              <span className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                {t("place.notes")}
              </span>
              <textarea
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Add personal notes, reservations, or subway tips..."
                rows={3}
                className="w-full p-2.5 rounded-xl bg-black/[0.03] dark:bg-white/[0.03] border border-[var(--glass-border)] text-[var(--text)] placeholder:text-[var(--text-muted)] focus:outline-none focus:ring-1 focus:ring-[var(--accent)] resize-none"
              />
            </div>

            {/* MinIO Upload Area */}
            <div className="space-y-2 pt-1 border-t border-[var(--glass-border)]">
              <div className="flex items-center justify-between">
                <span className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                  {t("place.attachments")}
                </span>
                <button
                  type="button"
                  onClick={() => setIsPrivateUpload(!isPrivateUpload)}
                  className="flex items-center gap-1 text-[10px] text-[var(--text-muted)] hover:text-[var(--text)] transition-colors"
                >
                  {isPrivateUpload ? <Lock className="w-3 h-3 text-amber-500 dark:text-amber-400" /> : <Globe className="w-3 h-3 text-sky-500 dark:text-sky-400" />}
                  <span>{isPrivateUpload ? t("place.privateBadge") : t("place.publicBadge")}</span>
                </button>
              </div>

              <input
                ref={fileInputRef}
                type="file"
                accept="image/*,application/pdf"
                className="hidden"
                onChange={handleFileUpload}
              />

              <button
                type="button"
                onClick={() => fileInputRef.current?.click()}
                disabled={isUploading}
                className="w-full py-2.5 px-3 rounded-xl border border-dashed border-[var(--glass-border)] hover:border-[var(--accent)] bg-black/[0.02] dark:bg-white/[0.02] hover:bg-black/[0.04] dark:hover:bg-white/[0.04] transition-all flex items-center justify-center gap-2 text-[var(--text-muted)] hover:text-[var(--text)] disabled:opacity-50"
              >
                {isUploading ? (
                  <>
                    <div className="w-3.5 h-3.5 border-2 border-[var(--accent)] border-t-transparent rounded-full animate-spin" />
                    <span>{t("place.uploading")}</span>
                  </>
                ) : (
                  <>
                    <Upload className="w-3.5 h-3.5" />
                    <span>{t("place.uploadDropzone")}</span>
                  </>
                )}
              </button>

              {uploadError && <div className="text-red-400 text-[10px] leading-tight">{uploadError}</div>}

              {docAttachments.length > 0 && (
                <div className="space-y-1.5 pt-1">
                  {docAttachments.map((att) => (
                    <div
                      key={att.id}
                      className="flex items-center justify-between p-2 rounded-xl bg-black/[0.03] dark:bg-white/[0.03] border border-[var(--glass-border)]"
                    >
                      <div className="flex items-center gap-2 truncate">
                        <Paperclip className="w-3.5 h-3.5 text-sky-400 shrink-0" />
                        <span className="truncate text-[var(--text)]">{att.name}</span>
                      </div>
                      <a
                        href={att.url}
                        target="_blank"
                        rel="noreferrer"
                        className="text-xs text-sky-400 hover:underline shrink-0 ml-2 font-medium"
                      >
                        PDF
                      </a>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </GlassPanel>
      </div>

      {lightboxUrl && (
        <div
          onClick={() => setLightboxUrl(null)}
          className="fixed inset-0 z-50 bg-black/80 backdrop-blur-md flex items-center justify-center p-6 animate-in fade-in duration-200"
        >
          <button
            onClick={() => setLightboxUrl(null)}
            className="absolute top-6 right-6 p-2 rounded-full bg-white/10 text-white hover:bg-white/20 transition-colors"
          >
            <X className="w-6 h-6" />
          </button>
          <img
            src={lightboxUrl}
            alt="Enlarged view"
            className="max-w-4xl max-h-[85vh] rounded-2xl object-contain shadow-2xl border border-white/10"
          />
        </div>
      )}
    </>
  );
}
