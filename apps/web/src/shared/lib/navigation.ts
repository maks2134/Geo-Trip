import { NavAppPreference } from "@/stores/map-ui-store";

export interface GeoCoordinate {
  latitude: number;
  longitude: number;
  name?: string;
}

export function detectUserPlatform(): "ios" | "android" | "macos" | "windows" | "other" {
  if (typeof navigator === "undefined") return "other";
  const ua = navigator.userAgent || "";
  if (/iPad|iPhone|iPod/.test(ua)) return "ios";
  if (/Android/.test(ua)) return "android";
  if (/Macintosh|Mac OS X/.test(ua)) return "macos";
  if (/Windows/.test(ua)) return "windows";
  return "other";
}

export function resolveNavApp(preference: NavAppPreference): "google" | "apple" | "yandex" {
  if (preference !== "auto") return preference;
  const platform = detectUserPlatform();
  if (platform === "ios" || platform === "macos") {
    return "apple";
  }
  return "google";
}

/**
 * Builds navigation deep-link for a single place coordinate
 */
export function buildSinglePlaceUrl(
  coord: { latitude: number; longitude: number; name?: string },
  preference: NavAppPreference = "auto"
): string {
  const app = resolveNavApp(preference);
  const { latitude, longitude, name } = coord;
  const encodedName = encodeURIComponent(name || `${latitude},${longitude}`);

  switch (app) {
    case "apple":
      return `http://maps.apple.com/?q=${encodedName}&ll=${latitude},${longitude}`;
    case "yandex":
      return `https://yandex.com/maps/?pt=${longitude},${latitude}&z=16&l=map`;
    case "google":
    default:
      return `https://www.google.com/maps/search/?api=1&query=${latitude},${longitude}`;
  }
}

/**
 * Builds transit route deep-link between two points (Origin -> Destination)
 */
export function buildNavigationRouteUrl(
  origin: GeoCoordinate,
  destination: GeoCoordinate,
  preference: NavAppPreference = "auto",
  mode: "driving" | "walking" = "driving"
): string {
  const app = resolveNavApp(preference);
  const oLat = origin.latitude;
  const oLng = origin.longitude;
  const dLat = destination.latitude;
  const dLng = destination.longitude;

  switch (app) {
    case "apple": {
      const dirFlg = mode === "walking" ? "w" : "d";
      return `http://maps.apple.com/?saddr=${oLat},${oLng}&daddr=${dLat},${dLng}&dirflg=${dirFlg}`;
    }
    case "yandex": {
      const rtt = mode === "walking" ? "pd" : "auto";
      return `https://yandex.com/maps/?rtext=${oLat},${oLng}~${dLat},${dLng}&rtt=${rtt}`;
    }
    case "google":
    default: {
      const travelmode = mode === "walking" ? "walking" : "driving";
      return `https://www.google.com/maps/dir/?api=1&origin=${oLat},${oLng}&destination=${dLat},${dLng}&travelmode=${travelmode}`;
    }
  }
}

/**
 * Calculates straight-line distance in kilometers using the Haversine formula
 */
export function calculateHaversineDistanceKm(
  lat1: number,
  lon1: number,
  lat2: number,
  lon2: number
): number {
  const R = 6371; // Earth's radius in km
  const dLat = ((lat2 - lat1) * Math.PI) / 180;
  const dLon = ((lon2 - lon1) * Math.PI) / 180;
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos((lat1 * Math.PI) / 180) *
      Math.cos((lat2 * Math.PI) / 180) *
      Math.sin(dLon / 2) *
      Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return Math.round(R * c * 10) / 10;
}

export interface TransitEstimate {
  distanceKm: number;
  drivingMinutes: number;
  walkingMinutes: number;
}

export function estimateTransit(
  origin: GeoCoordinate,
  destination: GeoCoordinate
): TransitEstimate {
  const distanceKm = calculateHaversineDistanceKm(
    origin.latitude,
    origin.longitude,
    destination.latitude,
    destination.longitude
  );

  // Urban estimations: driving ~25 km/h avg with stops, walking ~4.5 km/h
  const drivingMinutes = Math.max(3, Math.round((distanceKm / 25) * 60));
  const walkingMinutes = Math.max(2, Math.round((distanceKm / 4.5) * 60));

  return {
    distanceKm,
    drivingMinutes,
    walkingMinutes,
  };
}
