import { Attachment, PinCategory, Place } from "./place";

export interface PlaceRef {
  id: string;
  placeId: string;
  name: string;
  category: PinCategory;
  lat: number;
  lng: number;
  position: string;
  dayId?: string | null;
  notes?: string;
  durationMinutes?: number;
  cost?: number;
  transitToNext?: {
    mode: "walk" | "transit" | "drive";
    durationMinutes: number;
    distanceKm: number;
  };
}

export interface Day {
  id: string;
  tripId: string;
  dayNumber: number;
  date: string;
  title: string;
  places: PlaceRef[];
}

export interface Expense {
  id: string;
  tripId: string;
  dayId?: string | null;
  placeId?: string | null;
  title: string;
  amount: number;
  currency: string;
  category: "stay" | "food" | "transport" | "activity" | "other";
  paidAt: string;
}

export interface ChecklistItem {
  id: string;
  tripId: string;
  section: "pre-flight" | "packing" | "buy-on-site";
  title: string;
  completed: boolean;
}

export interface Trip {
  id: string;
  title: string;
  slug: string;
  description: string;
  status: "draft" | "upcoming" | "active" | "completed";
  startDate: string;
  endDate: string;
  coverImageUrl: string;
  budgetLimit: number;
  currency: string;
  countries: string[];
  cities: string[];
  days: Day[];
  backlog: PlaceRef[];
  expenses: Expense[];
  checklists: ChecklistItem[];
}

export interface PublicPlaceView {
  id: string;
  name: string;
  category: PinCategory;
  lat: number;
  lng: number;
  notes?: string;
  durationMinutes?: number;
  orderNumber: number;
  attachments?: Array<Pick<Attachment, "id" | "name" | "fileType" | "url">>;
}

export interface PublicDayView {
  id: string;
  dayNumber: number;
  date: string;
  title: string;
  places: PublicPlaceView[];
}

export interface PublicTripView {
  slug: string;
  title: string;
  description: string;
  startDate: string;
  endDate: string;
  coverImageUrl: string;
  countries: string[];
  cities: string[];
  days: PublicDayView[];
}
