import { create } from "zustand";
import { PinCategory, PinStatus } from "@/entities/place";

export type WorkspaceTab = "itinerary" | "budget" | "checklists";

interface MapUIState {
  mode: "world" | "trip";
  theme: "dark" | "light";
  selectedTripId: string;
  selectedDayId: string | null;
  drawerPlaceId: string | null;
  isDrawerOpen: boolean;
  activeWorkspaceTab: WorkspaceTab;
  isBacklogOpen: boolean;
  statusFilters: PinStatus[];
  categoryFilters: PinCategory[];

  setMode: (mode: "world" | "trip") => void;
  setTheme: (theme: "dark" | "light") => void;
  toggleTheme: () => void;
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

export const useMapUIStore = create<MapUIState>((set) => ({
  mode: "world",
  theme: "dark",
  selectedTripId: "trip-japan-2026",
  selectedDayId: "day-1",
  drawerPlaceId: null,
  isDrawerOpen: false,
  activeWorkspaceTab: "itinerary",
  isBacklogOpen: false,
  statusFilters: ["visited", "planned", "wishlist"],
  categoryFilters: ["food", "stay", "sight", "nature", "transport"],

  setMode: (mode) => set({ mode }),
  setTheme: (theme) => {
    if (typeof document !== "undefined") {
      document.documentElement.setAttribute("data-theme", theme);
    }
    set({ theme });
  },
  toggleTheme: () =>
    set((state) => {
      const nextTheme = state.theme === "dark" ? "light" : "dark";
      if (typeof document !== "undefined") {
        document.documentElement.setAttribute("data-theme", nextTheme);
      }
      return { theme: nextTheme };
    }),
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
