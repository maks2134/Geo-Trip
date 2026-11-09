// GeoJSON bounding polygons highlighting visited countries (Japan, Italy, France, USA, Thailand)
export const VISITED_COUNTRIES_GEOJSON: GeoJSON.FeatureCollection = {
  type: "FeatureCollection",
  features: [
    {
      type: "Feature",
      properties: { name: "Japan", code: "JPN" },
      geometry: {
        type: "Polygon",
        coordinates: [
          [
            [129.5, 31.0],
            [130.0, 33.5],
            [133.0, 34.5],
            [137.0, 37.0],
            [141.0, 41.5],
            [144.5, 44.0],
            [146.0, 43.5],
            [141.5, 38.5],
            [140.0, 35.0],
            [136.0, 33.5],
            [131.0, 31.0],
            [129.5, 31.0],
          ],
        ],
      },
    },
    {
      type: "Feature",
      properties: { name: "Italy", code: "ITA" },
      geometry: {
        type: "Polygon",
        coordinates: [
          [
            [7.0, 44.0],
            [8.5, 46.0],
            [12.5, 46.5],
            [13.8, 45.8],
            [13.0, 42.0],
            [16.0, 41.0],
            [18.5, 40.0],
            [15.5, 38.0],
            [15.0, 40.5],
            [11.0, 42.5],
            [9.0, 44.0],
            [7.0, 44.0],
          ],
        ],
      },
    },
    {
      type: "Feature",
      properties: { name: "France", code: "FRA" },
      geometry: {
        type: "Polygon",
        coordinates: [
          [
            [-4.5, 48.5],
            [2.5, 51.0],
            [7.5, 49.0],
            [7.0, 44.0],
            [3.0, 42.5],
            [-1.5, 43.5],
            [-4.5, 48.5],
          ],
        ],
      },
    },
  ],
};