export type PinCategory = "food" | "stay" | "sight" | "nature" | "transport";
export type PinStatus = "visited" | "planned" | "wishlist";

export interface Attachment {
  id: string;
  name: string;
  fileType: "pdf" | "image";
  url: string;
  sizeBytes: number;
  isPrivate: boolean;
  uploadedAt: string;
}

export interface Place {
  id: string;
  name: string;
  category: PinCategory;
  status: PinStatus;
  lat: number;
  lng: number;
  city: string;
  country: string;
  countryCode: string;
  address?: string;
  rating?: number;
  imageUrl?: string;
  notes?: string;
  externalLink?: string;
  estimatedCost?: number;
  durationMinutes?: number;
  attachments?: Attachment[];
}

export interface GeoStats {
  unlockedCountries: number;
  unlockedCities: number;
  totalKm: number;
  worldPercentage: number;
  visitedCountryCodes: string[];
}
