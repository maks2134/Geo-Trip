"use client";

import React from "react";
import { Globe, MapPin, Moon, Sun, Laptop, Crosshair, Printer } from "lucide-react";
import { useMapUIStore, ThemeMode } from "@/stores/map-ui-store";
import { GlassPanel } from "./GlassPanel";
import { cn } from "@/shared/lib/utils";
import { useTranslation, Locale } from "@/shared/i18n/i18n-context";

export function Dock() {
  const {
    mode,
    setMode,
    themeMode,
    theme,
    setThemeMode,
    quickPinMode,
    toggleQuickPinMode,
    setDossierOpen,
  } = useMapUIStore();
  const { t, locale, setLocale } = useTranslation();

  const cycleTheme = () => {
    const next: Record<ThemeMode, ThemeMode> = {
      system: "light",
      light: "dark",
      dark: "system",
    };
    setThemeMode(next[themeMode]);
  };

  return (
    <div className="fixed bottom-6 left-1/2 -translate-x-1/2 z-40">
      <GlassPanel className="flex items-center gap-1.5 p-1.5 rounded-full shadow-2xl backdrop-blur-2xl">
        <button
          onClick={() => setMode("world")}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all duration-200",
            mode === "world"
              ? "bg-[var(--accent)] text-white shadow-md shadow-[var(--accent)]/30"
              : "text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-white/5"
          )}
        >
          <Globe className="w-4 h-4" />
          <span>{t("nav.worldMap")}</span>
        </button>

        <button
          onClick={() => setMode("trip")}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all duration-200",
            mode === "trip"
              ? "bg-[var(--accent)] text-white shadow-md shadow-[var(--accent)]/30"
              : "text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-white/5"
          )}
        >
          <MapPin className="w-4 h-4" />
          <span>{t("nav.trips")}</span>
        </button>

        {/* Quick Pin Tool (FEAT-02) */}
        <button
          onClick={toggleQuickPinMode}
          title={t("dock.quickPinTooltip")}
          className={cn(
            "flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-medium transition-all duration-200 border",
            quickPinMode
              ? "bg-amber-500/25 text-amber-300 border-amber-500/50 shadow-md shadow-amber-500/20"
              : "text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-white/5 border-transparent"
          )}
        >
          <Crosshair className={cn("w-3.5 h-3.5", quickPinMode && "animate-spin")} />
          <span className="hidden md:inline">
            {quickPinMode ? t("dock.quickPinActive") : t("dock.quickPin")}
          </span>
        </button>

        {/* Dossier Print Tool (FEAT-05) */}
        {mode === "trip" && (
          <button
            onClick={() => setDossierOpen(true)}
            title={t("dossier.printTooltip")}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-medium text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-white/5 transition-colors border border-transparent"
          >
            <Printer className="w-3.5 h-3.5 text-emerald-400" />
            <span className="hidden md:inline">{t("dossier.title")}</span>
          </button>
        )}

        <div className="h-5 w-px bg-white/10 mx-1" />

        {/* Language Switcher Pill (MH-A) */}
        <div className="flex items-center bg-white/5 rounded-full p-0.5 border border-white/10 text-xs font-semibold">
          <button
            onClick={() => setLocale("en")}
            className={cn(
              "px-2.5 py-1 rounded-full transition-all",
              locale === "en" ? "bg-[var(--accent)] text-white shadow" : "text-[var(--text-muted)] hover:text-[var(--text)]"
            )}
          >
            EN
          </button>
          <button
            onClick={() => setLocale("ru")}
            className={cn(
              "px-2.5 py-1 rounded-full transition-all",
              locale === "ru" ? "bg-[var(--accent)] text-white shadow" : "text-[var(--text-muted)] hover:text-[var(--text)]"
            )}
          >
            RU
          </button>
        </div>

        {/* Theme Mode Toggle (MH-B) */}
        <button
          onClick={cycleTheme}
          title={`Theme: ${themeMode} (${theme})`}
          aria-label="Cycle theme: System -> Light -> Dark"
          className="flex items-center gap-1 px-2.5 py-1.5 rounded-full text-xs font-medium text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-white/5 transition-colors"
        >
          {themeMode === "system" && <Laptop className="w-4 h-4 text-sky-400" />}
          {themeMode === "light" && <Sun className="w-4 h-4 text-amber-400" />}
          {themeMode === "dark" && <Moon className="w-4 h-4 text-indigo-400" />}
          <span className="capitalize">{t(`theme.${themeMode}`)}</span>
        </button>
      </GlassPanel>
    </div>
  );
}
