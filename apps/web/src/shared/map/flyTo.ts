import maplibregl from "maplibre-gl";

interface FlyToOptions {
  paddingLeft?: number;
  duration?: number;
}

export function fitBoundsPlaces(
  map: maplibregl.Map,
  coordinates: [number, number][],
  options: FlyToOptions = {}
) {
  if (!map || coordinates.length === 0) return;

  const { paddingLeft = 460, duration = 1200 } = options;

  if (coordinates.length === 1) {
    map.flyTo({
      center: coordinates[0],
      zoom: 14,
      duration,
      padding: { top: 60, bottom: 60, left: paddingLeft, right: 60 },
      essential: true,
    });
    return;
  }

  const bounds = new maplibregl.LngLatBounds(coordinates[0], coordinates[0]);
  for (const coord of coordinates) {
    bounds.extend(coord);
  }

  map.fitBounds(bounds, {
    padding: { top: 80, bottom: 80, left: paddingLeft, right: 80 },
    maxZoom: 15,
    duration,
  });
}