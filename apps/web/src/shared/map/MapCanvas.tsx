"use client";

import React, { useEffect, useRef, useState, useCallback } from "react";
import maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";
import { useMapUIStore } from "@/stores/map-ui-store";
import { Place } from "@/entities/place";
import { Trip } from "@/entities/trip";
import { VISITED_COUNTRIES_GEOJSON } from "@/features/world-map/world-geojson";
import { fitBoundsPlaces } from "./flyTo";
import { useTranslation } from "@/shared/i18n/i18n-context";
import { apiClient } from "@/shared/lib/api";

interface MapCanvasProps {
  places?: Place[];
  trip?: Trip | null;
  readOnly?: boolean;
}

const MAP_STYLES = {
  dark: process.env.NEXT_PUBLIC_MAP_STYLE_DARK || "https://basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json",
  light: process.env.NEXT_PUBLIC_MAP_STYLE_LIGHT || "https://basemaps.cartocdn.com/gl/positron-gl-style/style.json",
};

const SATELLITE_TILE_URL =
  process.env.NEXT_PUBLIC_SATELLITE_TILE_URL ||
  "https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}";

const HIGH_ZOOM_STREET_LABELS = {
  dark:
    process.env.NEXT_PUBLIC_STREET_LABELS_DARK ||
    "https://a.basemaps.cartocdn.com/rastertiles/dark_only_labels/{z}/{x}/{y}.png",
  light:
    process.env.NEXT_PUBLIC_STREET_LABELS_LIGHT ||
    "https://a.basemaps.cartocdn.com/rastertiles/voyager_only_labels/{z}/{x}/{y}.png",
};

export function MapCanvas({ places = [], trip, readOnly = false }: MapCanvasProps) {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<maplibregl.Map | null>(null);
  const markersRef = useRef<maplibregl.Marker[]>([]);
  const satelliteErrorsRef = useRef<number>(0);
  const { t, locale } = useTranslation();

  const {
    mode,
    theme,
    basemap,
    setBasemap,
    toggleBasemap,
    selectedDayId,
    statusFilters,
    categoryFilters,
    openPlaceDrawer,
    satelliteErrorToast,
    setSatelliteErrorToast,
    quickPinMode,
    setQuickPinMode,
    droppedPin,
    setDroppedPin,
    activeRegionFilter,
  } = useMapUIStore();

  const quickPinModeRef = useRef(quickPinMode);
  quickPinModeRef.current = quickPinMode;
  const localeRef = useRef(locale);
  localeRef.current = locale;

  // Function to attach/re-attach custom GeoJSON and raster sources & layers
  const setupCustomLayers = useCallback((map: maplibregl.Map) => {
    // 1. Add Satellite Raster Source & Layer if not exists
    if (!map.getSource("satellite-tiles")) {
      map.addSource("satellite-tiles", {
        type: "raster",
        tiles: [SATELLITE_TILE_URL],
        tileSize: 256,
        attribution: "Tiles &copy; Esri, Maxar, Earthstar Geographics",
      });
    }

    if (!map.getLayer("satellite-raster-layer")) {
      // Find the first symbol layer to place raster below labels if available
      const layers = map.getStyle().layers;
      let firstSymbolId: string | undefined;
      if (layers) {
        for (const l of layers) {
          if (l.type === "symbol") {
            firstSymbolId = l.id;
            break;
          }
        }
      }

      map.addLayer(
        {
          id: "satellite-raster-layer",
          type: "raster",
          source: "satellite-tiles",
          paint: {
            "raster-opacity": basemap === "satellite" ? 1 : 0,
            "raster-fade-duration": 300,
          },
        },
        firstSymbolId
      );
    }

    // 2. Add Visited Countries GeoJSON for Scratch Map
    if (!map.getSource("visited-countries")) {
      map.addSource("visited-countries", {
        type: "geojson",
        data: VISITED_COUNTRIES_GEOJSON,
      });
    }

    if (!map.getLayer("visited-countries-fill")) {
      map.addLayer({
        id: "visited-countries-fill",
        type: "fill",
        source: "visited-countries",
        paint: {
          "fill-color": "#38bdf8",
          "fill-opacity": basemap === "satellite" ? 0.20 : 0.25,
        },
      });
    }

    if (!map.getLayer("visited-countries-outline")) {
      map.addLayer({
        id: "visited-countries-outline",
        type: "line",
        source: "visited-countries",
        paint: {
          "line-color": "#38bdf8",
          "line-width": 1.5,
          "line-dasharray": [2, 2],
        },
      });
    }

    // 3. Add Route Polyline
    if (!map.getSource("trip-route")) {
      map.addSource("trip-route", {
        type: "geojson",
        data: {
          type: "Feature",
          properties: {},
          geometry: { type: "LineString", coordinates: [] },
        },
      });
    }

    // Outer dark casing halo for contrast over satellite/bright maps
    if (!map.getLayer("trip-route-casing")) {
      map.addLayer({
        id: "trip-route-casing",
        type: "line",
        source: "trip-route",
        layout: {
          "line-join": "round",
          "line-cap": "round",
        },
        paint: {
          "line-color": "#000000",
          "line-width": 6.5,
          "line-opacity": basemap === "satellite" ? 0.85 : 0.35,
        },
      });
    }

    if (!map.getLayer("trip-route-line")) {
      map.addLayer({
        id: "trip-route-line",
        type: "line",
        source: "trip-route",
        layout: {
          "line-join": "round",
          "line-cap": "round",
        },
        paint: {
          "line-color": "#38bdf8",
          "line-width": 3.5,
          "line-opacity": 0.95,
          "line-dasharray": [1.5, 1.5],
        },
      });
    }

    // 4. Add High-Zoom Street & Road Labels Layer (zoom >= 12)
    if (!map.getSource("high-zoom-street-labels")) {
      map.addSource("high-zoom-street-labels", {
        type: "raster",
        tiles: [HIGH_ZOOM_STREET_LABELS[theme]],
        tileSize: 256,
        attribution: "&copy; OpenStreetMap contributors &copy; CARTO",
        minzoom: 12,
        maxzoom: 20,
      });
    }

    if (!map.getLayer("high-zoom-street-labels-layer")) {
      map.addLayer({
        id: "high-zoom-street-labels-layer",
        type: "raster",
        source: "high-zoom-street-labels",
        minzoom: 12,
        paint: {
          "raster-opacity": [
            "interpolate",
            ["linear"],
            ["zoom"],
            12,
            0.6,
            14,
            1.0,
          ],
          "raster-fade-duration": 200,
        },
      });
    }
  }, [basemap, theme]);

  // Initialize Map
  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) return;

    const map = new maplibregl.Map({
      container: mapContainerRef.current,
      style: MAP_STYLES[theme],
      center: [139.75, 35.68],
      zoom: 4,
      attributionControl: false,
    });

    map.addControl(new maplibregl.NavigationControl({ showCompass: false }), "top-right");

    map.on("load", () => {
      setupCustomLayers(map);
    });

    // Map click for Quick Pin tool (FEAT-02)
    map.on("click", async (e: maplibregl.MapMouseEvent) => {
      if (!quickPinModeRef.current) return;
      const { lat, lng } = e.lngLat;
      const currentLoc = localeRef.current;
      const defaultPinName =
        currentLoc === "ru"
          ? `Точка на карте (${lat.toFixed(4)}, ${lng.toFixed(4)})`
          : `Dropped Pin (${lat.toFixed(4)}, ${lng.toFixed(4)})`;
      const defaultAddress =
        currentLoc === "ru" ? "Определение адреса..." : "Reverse geocoding...";

      setDroppedPin({
        lat,
        lng,
        name: defaultPinName,
        address: defaultAddress,
      });
      setQuickPinMode(false);

      try {
        const geo = await apiClient.reverseGeocode(lat, lng, currentLoc);
        setDroppedPin({
          lat,
          lng,
          name: geo.name || defaultPinName,
          address: geo.address,
          city: geo.city,
          country: geo.country,
        });
      } catch {
        // fallback remains
      }
    });

    // Re-attach custom GeoJSON & raster layers whenever style finishes loading (styledata / style.load)
    map.on("style.load", () => {
      setupCustomLayers(map);
    });

    // Satellite tile error detection and graceful fallback
    map.on("error", (e: any) => {
      if (e?.error && (e.error?.message?.includes("satellite") || e.sourceId === "satellite-tiles" || e.source?.id === "satellite-tiles")) {
        satelliteErrorsRef.current += 1;
        if (satelliteErrorsRef.current >= 3) {
          setBasemap("streets");
          setSatelliteErrorToast(true);
          setTimeout(() => setSatelliteErrorToast(false), 5000);
        }
      }
    });

    mapRef.current = map;

    return () => {
      map.remove();
      mapRef.current = null;
    };
  }, []);

  // Update Style on Theme Switch while preserving layers
  useEffect(() => {
    if (!mapRef.current) return;
    mapRef.current.setStyle(MAP_STYLES[theme]);
  }, [theme]);

  // Update cursor on Quick Pin mode toggle (FEAT-02)
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;
    map.getCanvas().style.cursor = quickPinMode ? "crosshair" : "";
  }, [quickPinMode]);

  // Handle UN Regional filter camera flyTo (FEAT-06)
  useEffect(() => {
    const map = mapRef.current;
    if (!map || !activeRegionFilter || mode !== "world") return;

    const REGION_CENTERS: Record<string, { center: [number, number]; zoom: number }> = {
      "Western Europe": { center: [4.5, 48.5], zoom: 4.5 },
      "Southern Europe": { center: [14.0, 40.5], zoom: 4.5 },
      "Northern Europe": { center: [16.0, 60.0], zoom: 3.8 },
      "Eastern Europe": { center: [28.0, 52.0], zoom: 3.8 },
      "South-Eastern Asia": { center: [104.0, 13.0], zoom: 4.2 },
      "Eastern Asia": { center: [120.0, 36.0], zoom: 3.6 },
      "Southern Asia": { center: [78.0, 22.0], zoom: 4.0 },
      "Northern America": { center: [-98.0, 40.0], zoom: 3.2 },
      "South America": { center: [-60.0, -18.0], zoom: 3.0 },
      "Northern Africa": { center: [18.0, 26.0], zoom: 3.5 },
      "Sub-Saharan Africa": { center: [22.0, -1.0], zoom: 3.0 },
      "Australia and New Zealand": { center: [138.0, -26.0], zoom: 3.5 },
    };

    const target = REGION_CENTERS[activeRegionFilter];
    if (target) {
      map.flyTo({
        center: target.center,
        zoom: target.zoom,
        duration: 1600,
        essential: true,
      });
    }
  }, [activeRegionFilter, mode]);

  // Update Satellite visibility & contrast on basemap switch
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;

    if (map.getLayer("satellite-raster-layer")) {
      map.setPaintProperty(
        "satellite-raster-layer",
        "raster-opacity",
        basemap === "satellite" ? 1 : 0
      );
    }
    if (map.getLayer("visited-countries-fill")) {
      map.setPaintProperty(
        "visited-countries-fill",
        "fill-opacity",
        basemap === "satellite" ? 0.20 : 0.25
      );
    }
    if (map.getLayer("trip-route-casing")) {
      map.setPaintProperty(
        "trip-route-casing",
        "line-opacity",
        basemap === "satellite" ? 0.85 : 0.35
      );
    }
  }, [basemap]);

  // Update Markers & Layers by Mode
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;

    // Clear existing markers
    markersRef.current.forEach((m) => m.remove());
    markersRef.current = [];

    if (mode === "world") {
      // Toggle World layers visibility
      if (map.getLayer("visited-countries-fill")) {
        map.setLayoutProperty("visited-countries-fill", "visibility", "visible");
        map.setLayoutProperty("visited-countries-outline", "visibility", "visible");
      }
      if (map.getLayer("trip-route-line")) {
        map.setLayoutProperty("trip-route-line", "visibility", "none");
      }
      if (map.getLayer("trip-route-casing")) {
        map.setLayoutProperty("trip-route-casing", "visibility", "none");
      }

      // Filter and place pins
      const filteredPlaces = places.filter(
        (p) => statusFilters.includes(p.status) && categoryFilters.includes(p.category)
      );

      filteredPlaces.forEach((p) => {
        const el = document.createElement("div");
        el.className = "cursor-pointer transition-transform hover:scale-125 duration-150";

        let pinBg = "bg-amber-400";
        if (p.status === "planned") pinBg = "bg-sky-400";
        if (p.status === "wishlist") pinBg = "bg-pink-400";

        el.innerHTML = `
          <div class="w-7 h-7 rounded-full ${pinBg} text-slate-900 shadow-xl flex items-center justify-center font-bold text-xs ring-2 ring-white/70">
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="10" r="3"/><path d="M12 21.7C17.3 17 20 13 20 10a8 8 0 1 0-16 0c0 3 2.7 7 8 11.7z"/></svg>
          </div>
        `;

        el.addEventListener("click", () => {
          openPlaceDrawer(p.id);
        });

        const marker = new maplibregl.Marker({ element: el })
          .setLngLat([p.lng, p.lat])
          .addTo(map);

        markersRef.current.push(marker);
      });
    } else if (mode === "trip" && trip) {
      // Hide country fills, show route line
      if (map.getLayer("visited-countries-fill")) {
        map.setLayoutProperty("visited-countries-fill", "visibility", "none");
        map.setLayoutProperty("visited-countries-outline", "visibility", "none");
      }
      if (map.getLayer("trip-route-line")) {
        map.setLayoutProperty("trip-route-line", "visibility", "visible");
      }
      if (map.getLayer("trip-route-casing")) {
        map.setLayoutProperty("trip-route-casing", "visibility", "visible");
      }

      const activeDay = trip.days.find((d) => d.id === selectedDayId) || trip.days[0];
      const activeCoords: [number, number][] = [];

      if (activeDay && activeDay.places.length > 0) {
        activeDay.places.forEach((ref, index) => {
          activeCoords.push([ref.lng, ref.lat]);

          const el = document.createElement("div");
          el.className = "cursor-pointer transition-transform hover:scale-125 duration-150";
          el.innerHTML = `
            <div class="relative flex items-center justify-center">
              <div class="w-8 h-8 rounded-full bg-sky-500 text-white font-extrabold text-sm shadow-2xl flex items-center justify-center ring-2 ring-white/90">
                ${index + 1}
              </div>
              <div class="absolute -bottom-1 w-2 h-2 bg-sky-500 rotate-45"></div>
            </div>
          `;

          el.addEventListener("click", () => {
            openPlaceDrawer(ref.placeId);
          });

          const marker = new maplibregl.Marker({ element: el })
            .setLngLat([ref.lng, ref.lat])
            .addTo(map);

          markersRef.current.push(marker);
        });

        // Update Polyline
        const routeSource = map.getSource("trip-route") as maplibregl.GeoJSONSource | undefined;
        if (routeSource) {
          routeSource.setData({
            type: "Feature",
            properties: {},
            geometry: {
              type: "LineString",
              coordinates: activeCoords,
            },
          });
        }

        // Fly to active day bounds
        fitBoundsPlaces(map, activeCoords, { paddingLeft: readOnly ? 100 : 440 });
      }
    }

    // Render Dropped Quick Pin marker (FEAT-02)
    if (droppedPin) {
      const el = document.createElement("div");
      el.className = "cursor-pointer transition-transform hover:scale-125 duration-150";
      el.innerHTML = `
        <div class="relative flex items-center justify-center">
          <div class="w-8 h-8 rounded-full bg-amber-500 text-slate-900 font-extrabold text-sm shadow-2xl flex items-center justify-center ring-4 ring-amber-400/50">
            📍
          </div>
          <div class="absolute -inset-2 rounded-full border-2 border-amber-400 animate-ping pointer-events-none"></div>
        </div>
      `;
      el.addEventListener("click", () => {
        // Re-focus drawer
      });
      const marker = new maplibregl.Marker({ element: el })
        .setLngLat([droppedPin.lng, droppedPin.lat])
        .addTo(map);
      markersRef.current.push(marker);
    }
  }, [mode, places, trip, selectedDayId, statusFilters, categoryFilters, openPlaceDrawer, readOnly, droppedPin]);

  return (
    <div className="absolute inset-0 w-full h-full overflow-hidden">
      <div ref={mapContainerRef} className="w-full h-full" />

      {/* Floating Streets ↔ Satellite Toggle Pill (MH-D) */}
      <div className="absolute top-4 right-16 z-20 flex items-center rounded-xl p-1 bg-[var(--glass)] backdrop-blur-xl border border-[var(--glass-border)] shadow-2xl text-xs font-semibold text-[var(--text)]">
        <button
          onClick={() => setBasemap("streets")}
          className={`px-3 py-1.5 rounded-lg transition-all ${
            basemap === "streets"
              ? "bg-[var(--accent)] text-white shadow-md font-bold"
              : "text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-black/5 dark:hover:bg-white/10"
          }`}
        >
          {t("basemap.streets")}
        </button>
        <button
          onClick={() => setBasemap("satellite")}
          className={`px-3 py-1.5 rounded-lg transition-all ${
            basemap === "satellite"
              ? "bg-[var(--accent)] text-white shadow-md font-bold"
              : "text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-black/5 dark:hover:bg-white/10"
          }`}
        >
          {t("basemap.satellite")}
        </button>
      </div>

      {/* Satellite Error Fallback Toast */}
      {satelliteErrorToast && (
        <div className="absolute bottom-6 left-1/2 -translate-x-1/2 z-50 px-4 py-2.5 rounded-xl bg-amber-500/90 text-slate-950 font-medium text-xs shadow-2xl backdrop-blur-md flex items-center gap-2 animate-bounce">
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          {t("basemap.satelliteFallback")}
        </div>
      )}
    </div>
  );
}
