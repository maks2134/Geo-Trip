import { create } from "zustand";
import { PinCategory, PinStatus } from "@/entities/place";

export type WorkspaceTab = "itinerary" | "budget" | "checklists";
export type ThemeMode = "system" | "dark" | "light";
export type BasemapMode = "streets" | "satellite";

interface MapUIState {
  mode: "world" | "trip";
  themeMode: ThemeMode;
  theme: "dark" | "light"; // computed active theme
  basemap: BasemapMode;
  selectedTripId: string;
  selectedDayId: string | null;
  drawerPlaceId: string | null;
  isDrawerOpen: boolean;
  activeWorkspaceTab: WorkspaceTab;
  isBacklogOpen: boolean;
  statusFilters: PinStatus[];
  categoryFilters: PinCategory[];
  satelliteErrorToast: boolean;

  setMode: (mode: "world" | "trip") => void;
  setThemeMode: (mode: ThemeMode) => void;
  setBasemap: (basemap: BasemapMode) => void;
  toggleBasemap: () => void;
  setSatelliteErrorToast: (show: boolean) => void;
  setSelectedTripId: (tripId: string) => void;
  setSelectedDayId: (dayId: string | null) => void;
  openPlaceDrawer: (placeId: string) => void;
  closePlaceDrawer: () => void;
  setActiveWorkspaceTab: (tab: WorkspaceTab) => void;
  toggleBacklog: () => void;
  toggleStatusFilter: (status: PinStatus) => void;
  toggleCategoryFilter: (category: PinCategory) => void;
  resetFilters: () => void;
}

function resolveActiveTheme(mode: ThemeMode): "dark" | "light" {
  if (mode === "dark" || mode === "light") return mode;
  if (typeof window !== "undefined" && window.matchMedia) {
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }
  return "dark";
}

export const useMapUIStore = create<MapUIState>((set) => ({
  mode: "world",
  themeMode: "system",
  theme: "dark",
  basemap: "streets",
  selectedTripId: "trip-japan-2026",
  selectedDayId: "day-1",
  drawerPlaceId: null,
  isDrawerOpen: false,
  activeWorkspaceTab: "itinerary",
  isBacklogOpen: false,
  statusFilters: ["visited", "planned", "wishlist"],
  categoryFilters: ["food", "stay", "sight", "nature", "transport"],
  satelliteErrorToast: false,

  setMode: (mode) => set({ mode }),
  setThemeMode: (themeMode) => {
    const active = resolveActiveTheme(themeMode);
    if (typeof document !== "undefined") {
      document.documentElement.setAttribute("data-theme", active);
      localStorage.setItem("atlasflow_theme", themeMode);
    }
    set({ themeMode, theme: active });
  },
  setBasemap: (basemap) => {
    if (typeof localStorage !== "undefined") {
      localStorage.setItem("atlasflow_basemap", basemap);
    }
    set({ basemap });
  },
  toggleBasemap: () =>
    set((s) => {
      const next = s.basemap === "streets" ? "satellite" : "streets";
      if (typeof localStorage !== "undefined") {
        localStorage.setItem("atlasflow_basemap", next);
      }
      return { basemap: next };
    }),
  setSatelliteErrorToast: (satelliteErrorToast) => set({ satelliteErrorToast }),
  setSelectedTripId: (selectedTripId) =>
    set({ selectedTripId, selectedDayId: "day-1" }),
  setSelectedDayId: (selectedDayId) => set({ selectedDayId }),
  openPlaceDrawer: (drawerPlaceId) =>
    set({ drawerPlaceId, isDrawerOpen: true }),
  closePlaceDrawer: () => set({ drawerPlaceId: null, isDrawerOpen: false }),
  setActiveWorkspaceTab: (activeWorkspaceTab) => set({ activeWorkspaceTab }),
  toggleBacklog: () => set((s) => ({ isBacklogOpen: !s.isBacklogOpen })),
  toggleStatusFilter: (status) =>
    set((state) => {
      const exists = state.statusFilters.includes(status);
      const next = exists
        ? state.statusFilters.filter((s) => s !== status)
        : [...state.statusFilters, status];
      return { statusFilters: next.length === 0 ? [status] : next };
    }),
  toggleCategoryFilter: (category) =>
    set((state) => {
      const exists = state.categoryFilters.includes(category);
      const next = exists
        ? state.categoryFilters.filter((c) => c !== category)
        : [...state.categoryFilters, category];
      return { categoryFilters: next.length === 0 ? [category] : next };
    }),
  resetFilters: () =>
    set({
      statusFilters: ["visited", "planned", "wishlist"],
      categoryFilters: ["food", "stay", "sight", "nature", "transport"],
    }),
}));
