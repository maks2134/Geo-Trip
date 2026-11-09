"use client";

import React, { useId } from "react";
import { Printer, X, MapPin, Calendar, Compass, Clock } from "lucide-react";
import { Trip, Day, PlaceRef } from "@/entities/trip";
import { useMapUIStore } from "@/stores/map-ui-store";
import { generateQRCodeSVG } from "@/shared/lib/qrcode";
import { useTranslation } from "@/shared/i18n/i18n-context";

interface PrintDossierProps {
  trip: Trip;
  isPublicShare?: boolean;
}

export function PrintDossier({ trip, isPublicShare = false }: PrintDossierProps) {
  const { dossierOpen, setDossierOpen } = useMapUIStore();
  const { t } = useTranslation();

  if (!dossierOpen) return null;

  const handlePrint = () => {
    if (typeof window !== "undefined") {
      window.print();
    }
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto bg-black/80 backdrop-blur-md flex flex-col items-center p-4 sm:p-6 print:p-0 print:bg-white print:static print:inset-auto">
      <style jsx global>{`
        @media print {
          @page {
            margin: 1.2cm;
            size: A4 portrait;
          }
          body {
            background-color: #ffffff !important;
            color: #000000 !important;
          }
          .no-print {
            display: none !important;
          }
          .dossier-page {
            page-break-after: always !important;
            break-after: page !important;
            min-height: 100vh;
            background: #ffffff !important;
            color: #111827 !important;
            padding: 0 !important;
            box-shadow: none !important;
            border: none !important;
          }
        }
      `}</style>

      {/* Floating Action Bar (hidden when printing) */}
      <div className="no-print sticky top-4 z-50 w-full max-w-4xl flex items-center justify-between p-3 rounded-2xl bg-slate-900/90 border border-white/20 shadow-2xl backdrop-blur-xl mb-6">
        <div className="flex items-center gap-2 text-white font-bold text-sm px-2">
          <Printer className="w-4 h-4 text-emerald-400" />
          <span>{t("dossier.previewTitle")}</span>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={handlePrint}
            className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold text-xs shadow-lg shadow-emerald-500/25 transition-all"
          >
            <Printer className="w-4 h-4" />
            <span>{t("dossier.printBtn")}</span>
          </button>
          <button
            onClick={() => setDossierOpen(false)}
            className="p-2 rounded-xl bg-white/10 hover:bg-white/20 text-white transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>
      </div>

      {/* Dossier Document Container: 1 Page Per Day */}
      <div className="w-full max-w-4xl space-y-8 print:space-y-0">
        {trip.days.map((day, dayIdx) => (
          <div
            key={day.id}
            className="dossier-page rounded-3xl bg-white text-slate-900 p-8 sm:p-10 shadow-2xl border border-slate-200"
          >
            {/* Day Page Header */}
            <div className="border-b-2 border-slate-900 pb-5 mb-6 flex items-start justify-between">
              <div>
                <div className="flex items-center gap-2 text-xs font-black tracking-widest uppercase text-sky-600 mb-1">
                  <span>ATLASFLOW EXPEDITION DOSSIER</span>
                  <span>•</span>
                  <span>DAY {day.dayNumber} OF {trip.days.length}</span>
                </div>
                <h1 className="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900">
                  {day.title || `Day ${day.dayNumber}`}
                </h1>
                <div className="flex items-center gap-4 text-xs text-slate-600 mt-1 font-medium">
                  <span className="flex items-center gap-1">
                    <Calendar className="w-3.5 h-3.5" />
                    {day.date || trip.startDate}
                  </span>
                  <span>•</span>
                  <span>{trip.title}</span>
                </div>
              </div>

              <div className="text-right">
                <div className="text-xs font-bold text-slate-500 uppercase tracking-wider">
                  Total Stops
                </div>
                <div className="text-2xl font-black text-slate-900">
                  {day.places.length}
                </div>
              </div>
            </div>

            {/* Day Itinerary Stops */}
            {day.places.length === 0 ? (
              <div className="py-20 text-center text-slate-400 font-medium text-sm">
                No stops scheduled for this day.
              </div>
            ) : (
              <div className="space-y-4">
                {day.places.map((place, idx) => {
                  const mapUrl = `https://maps.google.com/?q=${place.lat},${place.lng}`;
                  const qrSvg = generateQRCodeSVG(mapUrl, 64, "#0f172a", "#ffffff");

                  return (
                    <div
                      key={place.id}
                      className="flex items-start justify-between gap-4 p-4 rounded-2xl border border-slate-200 bg-slate-50/70"
                    >
                      {/* Left: Stop Index & Info */}
                      <div className="flex items-start gap-3.5 flex-1 min-w-0">
                        <div className="w-8 h-8 rounded-full bg-slate-900 text-white font-black text-sm flex items-center justify-center shrink-0 mt-0.5">
                          {idx + 1}
                        </div>

                        <div className="space-y-1 flex-1 min-w-0">
                          <div className="flex items-center gap-2">
                            <h3 className="font-bold text-base text-slate-900 leading-tight truncate">
                              {place.name}
                            </h3>
                            <span className="px-2 py-0.5 rounded-md bg-slate-200 text-slate-700 text-[10px] font-bold uppercase tracking-wider">
                              {place.category}
                            </span>
                          </div>

                          <div className="flex items-center gap-3 text-xs text-slate-500 font-mono">
                            <span>{place.lat.toFixed(5)}, {place.lng.toFixed(5)}</span>
                            {place.durationMinutes && (
                              <span className="flex items-center gap-1 font-sans font-medium text-slate-600">
                                <Clock className="w-3 h-3" />
                                {place.durationMinutes} min
                              </span>
                            )}
                          </div>

                          {place.notes && (
                            <p className="text-xs text-slate-700 mt-1 italic leading-relaxed">
                              "{place.notes}"
                            </p>
                          )}
                        </div>
                      </div>

                      {/* Right: Sharp SVG QR Code (Scenario 5.2) */}
                      <div className="shrink-0 flex flex-col items-center gap-1">
                        <div
                          className="w-16 h-16 p-1 bg-white rounded-lg border border-slate-300 shadow-sm flex items-center justify-center"
                          dangerouslySetInnerHTML={{ __html: qrSvg }}
                        />
                        <span className="text-[9px] font-mono text-slate-500">Scan GPS</span>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}

            {/* Page Footer */}
            <div className="mt-8 pt-4 border-t border-slate-200 flex items-center justify-between text-[11px] text-slate-500">
              <div>AtlasFlow / NomadGlass • Personal Travel Itinerary</div>
              <div>Page {dayIdx + 1} of {trip.days.length}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
