"use client";

import React, { useState, useRef } from "react";
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
  Image as ImageIcon,
} from "lucide-react";
import { Place, PinStatus, PlaceAttachment } from "@/entities/place";
import { useMapUIStore } from "@/stores/map-ui-store";
import { GlassPanel } from "@/shared/ui/GlassPanel";
import { Badge } from "@/shared/ui/Badge";
import { cn } from "@/shared/lib/utils";
import { useTranslation } from "@/shared/i18n/i18n-context";
import { apiClient } from "@/shared/lib/api";

interface PlaceDrawerProps {
  places: Place[];
}

export function PlaceDrawer({ places }: PlaceDrawerProps) {
  const { drawerPlaceId, isDrawerOpen, closePlaceDrawer, selectedTripId } = useMapUIStore();
  const { t } = useTranslation();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const place = places.find((p) => p.id === drawerPlaceId);

  const [status, setStatus] = useState<PinStatus>(place?.status || "planned");
  const [notes, setNotes] = useState<string>(place?.notes || "");
  const [attachments, setAttachments] = useState<PlaceAttachment[]>(place?.attachments || []);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [lightboxUrl, setLightboxUrl] = useState<string | null>(null);
  const [isPrivateUpload, setIsPrivateUpload] = useState<boolean>(true);

  if (!isDrawerOpen || !place) return null;

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploadError(null);

    // Validate type and size
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
      // 1. Request presigned upload URL from trip-core
      const presign = await apiClient.createPresignedUpload(selectedTripId, {
        filename: file.name,
        contentType: file.type,
        fileSize: file.size,
        visibility: isPrivateUpload ? "private" : "public",
        placeId: place.id,
        ownerType: "place",
      });

      // 2. Upload binary directly to MinIO/S3
      await apiClient.uploadFileToMinIO(presign.uploadUrl, file);

      // 3. Confirm upload with trip-core
      await apiClient.confirmUpload(selectedTripId, presign.attachmentId);

      // 4. Update UI
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
      setUploadError(err?.message || "Failed to upload file");
    } finally {
      setIsUploading(false);
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  };

  const imageAttachments = attachments.filter((a) => a.fileType === "image");
  const docAttachments = attachments.filter((a) => a.fileType === "pdf");

  return (
    <>
      <div className="fixed top-6 right-6 bottom-24 w-96 z-40 animate-in slide-in-from-right-8 duration-200">
        <GlassPanel className="h-full flex flex-col overflow-hidden shadow-2xl border-white/15">
          {/* Header Image */}
          <div className="relative h-44 w-full bg-slate-900 shrink-0">
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
                    <span>{t(`filters.${s}`)}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Quick Metrics */}
            <div className="grid grid-cols-3 gap-2 py-2 border-y border-white/5">
              {place.rating && (
                <div className="p-2 rounded-xl bg-white/[0.03] text-center">
                  <div className="flex items-center justify-center gap-1 font-bold text-[var(--text)]">
                    <Star className="w-3.5 h-3.5 text-amber-400 fill-amber-400" />
                    <span>{place.rating}</span>
                  </div>
                  <div className="text-[10px] text-[var(--text-muted)]">{t("place.rating")}</div>
                </div>
              )}
              <div className="p-2 rounded-xl bg-white/[0.03] text-center">
                <div className="flex items-center justify-center gap-1 font-bold text-[var(--text)]">
                  <Clock className="w-3.5 h-3.5 text-sky-400" />
                  <span>{place.durationMinutes || 60}{t("place.minutes")}</span>
                </div>
                <div className="text-[10px] text-[var(--text-muted)]">{t("place.duration")}</div>
              </div>
              <div className="p-2 rounded-xl bg-white/[0.03] text-center">
                <div className="flex items-center justify-center gap-1 font-bold text-[var(--text)]">
                  <DollarSign className="w-3.5 h-3.5 text-emerald-400" />
                  <span>${place.estimatedCost ?? 0}</span>
                </div>
                <div className="text-[10px] text-[var(--text-muted)]">{t("place.estimatedCost")}</div>
              </div>
            </div>

            {/* Photo Gallery Strip (MH-C) */}
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
                      className="relative w-20 h-20 rounded-xl overflow-hidden shrink-0 border border-white/10 group cursor-pointer hover:scale-105 transition-transform"
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
                className="w-full p-2.5 rounded-xl bg-white/[0.03] border border-white/10 text-[var(--text)] placeholder:text-[var(--text-muted)] focus:outline-none focus:ring-1 focus:ring-[var(--accent)] resize-none"
              />
            </div>

            {/* Official Link */}
            {place.externalLink && (
              <div className="space-y-1.5">
                <span className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                  {t("place.links")}
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

            {/* MinIO Upload Area (MH-C) */}
            <div className="space-y-2 pt-2 border-t border-white/5">
              <div className="flex items-center justify-between">
                <span className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                  {t("place.documents")} & {t("place.photos")}
                </span>
                <div className="flex items-center gap-1.5 text-[10px]">
                  <button
                    type="button"
                    onClick={() => setIsPrivateUpload(!isPrivateUpload)}
                    className={cn(
                      "flex items-center gap-1 px-2 py-0.5 rounded-md border transition-all",
                      isPrivateUpload
                        ? "bg-amber-500/10 border-amber-500/30 text-amber-400"
                        : "bg-sky-500/10 border-sky-500/30 text-sky-400"
                    )}
                  >
                    {isPrivateUpload ? <Lock className="w-3 h-3" /> : <Globe className="w-3 h-3" />}
                    <span>{isPrivateUpload ? t("place.privateBadge") : t("place.publicBadge")}</span>
                  </button>
                </div>
              </div>

              {/* Upload Dropzone */}
              <div
                onClick={() => fileInputRef.current?.click()}
                className="p-3.5 rounded-xl bg-white/[0.02] border border-dashed border-white/15 hover:border-[var(--accent)] cursor-pointer text-center space-y-1.5 transition-colors group"
              >
                <input
                  ref={fileInputRef}
                  type="file"
                  accept="image/jpeg,image/png,image/webp,application/pdf"
                  onChange={handleFileUpload}
                  className="hidden"
                />
                <div className="flex justify-center">
                  <Upload className="w-5 h-5 text-[var(--text-muted)] group-hover:text-[var(--accent)] transition-colors" />
                </div>
                <p className="text-[11px] text-[var(--text)] font-medium">
                  {isUploading ? t("place.uploading") : t("place.uploadDropzone")}
                </p>
                <p className="text-[10px] text-[var(--text-muted)]">
                  {t("place.uploadSizeLimit")}
                </p>
              </div>

              {uploadError && (
                <p className="text-[11px] text-rose-400 bg-rose-500/10 p-2 rounded-lg border border-rose-500/20">
                  {uploadError}
                </p>
              )}

              {/* Document List */}
              {docAttachments.length > 0 && (
                <div className="space-y-1.5 pt-1">
                  {docAttachments.map((att) => (
                    <div
                      key={att.id}
                      className="flex items-center justify-between p-2.5 rounded-xl bg-white/[0.03] border border-white/5"
                    >
                      <div className="flex items-center gap-2 min-w-0">
                        <FileText className="w-4 h-4 text-rose-400 shrink-0" />
                        <span className="truncate font-medium">{att.name}</span>
                        {att.isPrivate && (
                          <span title="Hidden from public share" className="p-0.5 rounded bg-amber-500/10 text-amber-400 shrink-0">
                            <Lock className="w-3 h-3" />
                          </span>
                        )}
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

      {/* Lightbox Modal (MH-C) */}
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
