"use client";

import React, { useEffect, useRef } from "react";
import maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";
import { useMapUIStore } from "@/stores/map-ui-store";
import { Place } from "@/entities/place";
import { Trip } from "@/entities/trip";
import { VISITED_COUNTRIES_GEOJSON } from "@/features/world-map/world-geojson";
import { fitBoundsPlaces } from "./flyTo";

interface MapCanvasProps {
  places?: Place[];
  trip?: Trip | null;
  readOnly?: boolean;
}

const MAP_STYLES = {
  dark: "https://basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json",
  light: "https://basemaps.cartocdn.com/gl/positron-gl-style/style.json",
};

export function MapCanvas({ places = [], trip, readOnly = false }: MapCanvasProps) {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<maplibregl.Map | null>(null);
  const markersRef = useRef<maplibregl.Marker[]>([]);

  const {
    mode,
    theme,
    selectedDayId,
    statusFilters,
    categoryFilters,
    openPlaceDrawer,
  } = useMapUIStore();

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
      // Add Country Fill Layer for World Scratch Map
      map.addSource("visited-countries", {
        type: "geojson",
        data: VISITED_COUNTRIES_GEOJSON,
      });

      map.addLayer({
        id: "visited-countries-fill",
        type: "fill",
        source: "visited-countries",
        paint: {
          "fill-color": "#38bdf8",
          "fill-opacity": 0.22,
        },
      });

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

      // Route polyline source
      map.addSource("trip-route", {
        type: "geojson",
        data: {
          type: "Feature",
          properties: {},
          geometry: { type: "LineString", coordinates: [] },
        },
      });

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
          "line-opacity": 0.9,
          "line-dasharray": [1.5, 1.5],
        },
      });
    });

    mapRef.current = map;

    return () => {
      map.remove();
      mapRef.current = null;
    };
  }, []);

  // Update Style on Theme Switch
  useEffect(() => {
    if (!mapRef.current) return;
    mapRef.current.setStyle(MAP_STYLES[theme]);
  }, [theme]);

  // Update Markers & Layers
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
          <div class="w-7 h-7 rounded-full ${pinBg} text-slate-900 shadow-lg flex items-center justify-center font-bold text-xs ring-2 ring-white/60">
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

      const activeDay = trip.days.find((d) => d.id === selectedDayId) || trip.days[0];
      const activeCoords: [number, number][] = [];

      if (activeDay && activeDay.places.length > 0) {
        activeDay.places.forEach((ref, index) => {
          activeCoords.push([ref.lng, ref.lat]);

          const el = document.createElement("div");
          el.className = "cursor-pointer transition-transform hover:scale-125 duration-150";
          el.innerHTML = `
            <div class="relative flex items-center justify-center">
              <div class="w-8 h-8 rounded-full bg-sky-500 text-white font-extrabold text-sm shadow-xl flex items-center justify-center ring-2 ring-white/80">
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
  }, [mode, places, trip, selectedDayId, statusFilters, categoryFilters, openPlaceDrawer, readOnly]);

  return (
    <div className="absolute inset-0 w-full h-full overflow-hidden">
      <div ref={mapContainerRef} className="w-full h-full" />
    </div>
  );
}