"use client";

import React, { useState } from "react";
import { CheckCircle2, Circle, Plane, Luggage, ShoppingBag, Sparkles, X } from "lucide-react";
import { Trip, ChecklistItem } from "@/entities/trip";
import { ProgressBar } from "@/shared/ui/ProgressBar";
import { apiClient } from "@/shared/lib/api";
import { useTranslation } from "@/shared/i18n/i18n-context";

interface ChecklistTabProps {
  trip: Trip;
  onTripUpdate: (updated: Trip) => void;
}

export function ChecklistTab({ trip, onTripUpdate }: ChecklistTabProps) {
  const { t } = useTranslation();
  const [isTemplateModalOpen, setIsTemplateModalOpen] = useState(false);

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

  const handleApplyTemplate = async (templateKey: string) => {
    const updatedChecklists = await apiClient.applyChecklistTemplate(trip.id, templateKey);
    onTripUpdate({
      ...trip,
      checklists: updatedChecklists,
    });
    setIsTemplateModalOpen(false);
  };

  const sections: {
    key: ChecklistItem["section"];
    labelKey: string;
    icon: React.ReactNode;
  }[] = [
    { key: "pre-flight", labelKey: "checklist.preFlight", icon: <Plane className="w-3.5 h-3.5 text-sky-400" /> },
    { key: "packing", labelKey: "checklist.packing", icon: <Luggage className="w-3.5 h-3.5 text-amber-400" /> },
    { key: "buy-on-site", labelKey: "checklist.buyOnSite", icon: <ShoppingBag className="w-3.5 h-3.5 text-emerald-400" /> },
  ];

  const starterTemplates = [
    { key: "digital_nomad", titleKey: "checklist.templateNomad", descKey: "checklist.templateNomadDesc", icon: "💻" },
    { key: "beach_sun", titleKey: "checklist.templateBeach", descKey: "checklist.templateBeachDesc", icon: "🏖️" },
    { key: "mountain_trekking", titleKey: "checklist.templateTrek", descKey: "checklist.templateTrekDesc", icon: "🥾" },
    { key: "winter_city", titleKey: "checklist.templateWinter", descKey: "checklist.templateWinterDesc", icon: "❄️" },
  ];

  return (
    <div className="space-y-4 text-xs">
      {/* Departure Readiness Summary */}
      <div className="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-2">
        <div className="flex items-center justify-between">
          <span className="font-semibold text-xs text-[var(--text)]">{t("checklist.readiness")}</span>
          <span className="font-bold text-sky-400">{progressPercent}%</span>
        </div>
        <ProgressBar value={progressPercent} />
        <div className="flex items-center justify-between text-[11px] text-[var(--text-muted)] pt-1">
          <span>{completedCount} / {totalCount}</span>
          <button
            onClick={() => setIsTemplateModalOpen(true)}
            className="flex items-center gap-1 text-[var(--accent)] hover:underline font-semibold"
          >
            <Sparkles className="w-3.5 h-3.5 text-amber-400" />
            <span>{t("checklist.applyTemplate")}</span>
          </button>
        </div>
      </div>

      {/* Sections & Items */}
      <div className="space-y-4">
        {sections.map((sec) => {
          const items = trip.checklists.filter((c) => c.section === sec.key);
          return (
            <div key={sec.key} className="space-y-2">
              <div className="flex items-center gap-2 font-semibold text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                {sec.icon}
                <span>{t(sec.labelKey)}</span>
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

      {/* Starter Templates Modal (CAND-06) */}
      {isTemplateModalOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="w-full max-w-md p-5 rounded-2xl bg-slate-900/95 border border-white/15 shadow-2xl space-y-4 text-slate-200">
            <div className="flex items-center justify-between border-b border-white/10 pb-3">
              <div className="flex items-center gap-2">
                <Sparkles className="w-4 h-4 text-amber-400" />
                <h3 className="font-bold text-sm text-white">{t("checklist.templatesTitle")}</h3>
              </div>
              <button
                onClick={() => setIsTemplateModalOpen(false)}
                className="p-1 rounded-full text-white/60 hover:text-white hover:bg-white/10"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="space-y-2">
              {starterTemplates.map((tpl) => (
                <button
                  key={tpl.key}
                  onClick={() => handleApplyTemplate(tpl.key)}
                  className="w-full text-left p-3 rounded-xl bg-white/[0.03] hover:bg-sky-500/15 border border-white/10 hover:border-sky-500/30 transition-all group flex items-start gap-3"
                >
                  <span className="text-2xl">{tpl.icon}</span>
                  <div>
                    <h4 className="font-semibold text-xs text-white group-hover:text-sky-300">
                      {t(tpl.titleKey)}
                    </h4>
                    <p className="text-[11px] text-slate-400 mt-0.5">
                      {t(tpl.descKey)}
                    </p>
                  </div>
                </button>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
