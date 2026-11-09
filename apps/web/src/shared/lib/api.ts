import { GeoStats, Place } from "@/entities/place";
import { ChecklistItem, Expense, PlaceRef, PublicTripView, Trip } from "@/entities/trip";

const MOCK_GEO_STATS: GeoStats = {
  unlockedCountries: 24,
  unlockedCities: 68,
  totalKm: 48250,
  worldPercentage: 12.3,
  visitedCountryCodes: ["JPN", "ITA", "FRA", "USA", "THA", "ESP", "DEU", "GBR", "CHE", "NOR"],
};

const MOCK_PLACES: Place[] = [
  {
    id: "pl-tokyo-skytree",
    name: "Tokyo Skytree & Solamachi",
    category: "sight",
    status: "planned",
    lat: 35.710063,
    lng: 139.8107,
    city: "Tokyo",
    country: "Japan",
    countryCode: "JPN",
    address: "1-1-2 Oshiage, Sumida-ku, Tokyo",
    rating: 4.8,
    imageUrl: "https://images.unsplash.com/photo-1542051841857-5f90071e7989?w=800&q=80",
    notes: "Pre-booked fast pass ticket for sunset window at 17:30.",
    externalLink: "https://www.tokyo-skytree.jp/en/",
    estimatedCost: 32,
    durationMinutes: 90,
    attachments: [
      {
        id: "att-1",
        name: "Skytree_Voucher_FastPass.pdf",
        fileType: "pdf",
        url: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
        sizeBytes: 245000,
        isPrivate: true,
        uploadedAt: "2026-10-15T09:30:00Z",
      },
    ],
  },
  {
    id: "pl-sensoji",
    name: "Senso-ji Ancient Temple",
    category: "sight",
    status: "planned",
    lat: 35.714765,
    lng: 139.796655,
    city: "Tokyo",
    country: "Japan",
    countryCode: "JPN",
    rating: 4.9,
    imageUrl: "https://images.unsplash.com/photo-1503899036084-c55cdd92da26?w=800&q=80",
    notes: "Stroll down Nakamise-dori shopping street; try ningyo-yaki cakes.",
    durationMinutes: 75,
    estimatedCost: 10,
  },
  {
    id: "pl-ichiran-asakusa",
    name: "Ichiran Ramen Asakusa",
    category: "food",
    status: "planned",
    lat: 35.7121,
    lng: 139.7984,
    city: "Tokyo",
    country: "Japan",
    countryCode: "JPN",
    rating: 4.7,
    imageUrl: "https://images.unsplash.com/photo-1569718212165-3a8278d5f624?w=800&q=80",
    notes: "Classic tonkotsu broth with private dining booths.",
    durationMinutes: 45,
    estimatedCost: 18,
  },
  {
    id: "pl-teamlab-planets",
    name: "teamLab Planets TOKYO",
    category: "sight",
    status: "planned",
    lat: 35.6491,
    lng: 139.7898,
    city: "Tokyo",
    country: "Japan",
    countryCode: "JPN",
    rating: 4.9,
    imageUrl: "https://images.unsplash.com/photo-1509198397868-475647b2a1e5?w=800&q=80",
    notes: "Water installation: roll pants up above knees. Lockers available.",
    durationMinutes: 120,
    estimatedCost: 40,
    attachments: [
      {
        id: "att-2",
        name: "TeamLab_QR_Pass.png",
        fileType: "image",
        url: "https://images.unsplash.com/photo-1509198397868-475647b2a1e5?w=800&q=80",
        sizeBytes: 850000,
        isPrivate: true,
        uploadedAt: "2026-10-18T14:10:00Z",
      },
    ],
  },
  {
    id: "pl-shibuya-sky",
    name: "SHIBUYA SKY Observation Deck",
    category: "sight",
    status: "visited",
    lat: 35.6586,
    lng: 139.7013,
    city: "Tokyo",
    country: "Japan",
    countryCode: "JPN",
    rating: 4.9,
    imageUrl: "https://images.unsplash.com/photo-1536098561742-ca998e48cbcc?w=800&q=80",
    notes: "Unobstructed open rooftop sky edge view over Shibuya scramble.",
    durationMinutes: 90,
    estimatedCost: 25,
  },
  {
    id: "pl-fushimi-inari",
    name: "Fushimi Inari Taisha",
    category: "nature",
    status: "wishlist",
    lat: 34.9671,
    lng: 135.7727,
    city: "Kyoto",
    country: "Japan",
    countryCode: "JPN",
    rating: 4.9,
    imageUrl: "https://images.unsplash.com/photo-1493976040374-85c8e12f0c0e?w=800&q=80",
    notes: "Thousands of vermilion torii gates winding up the sacred mountain.",
    durationMinutes: 150,
    estimatedCost: 0,
  },
];

let MOCK_TRIPS: Trip[] = [
  {
    id: "trip-japan-2026",
    title: "Tokyo Highlights & Culinary Adventure",
    slug: "tokyo-adventure-2026",
    description: "Cherry blossom season, neon lights, ancient temples, and ramen crawls across Tokyo.",
    status: "active",
    startDate: "2026-04-02",
    endDate: "2026-04-10",
    coverImageUrl: "https://images.unsplash.com/photo-1503899036084-c55cdd92da26?w=1200&q=80",
    budgetLimit: 2800,
    currency: "USD",
    countries: ["Japan"],
    cities: ["Tokyo", "Kyoto"],
    days: [
      {
        id: "day-1",
        tripId: "trip-japan-2026",
        dayNumber: 1,
        date: "2026-04-02",
        title: "Asakusa Temples & Skytree Sunset",
        places: [
          {
            id: "it-1",
            placeId: "pl-sensoji",
            name: "Senso-ji Temple",
            category: "sight",
            lat: 35.714765,
            lng: 139.796655,
            dayId: "day-1",
            position: "a0",
            durationMinutes: 75,
            cost: 10,
            notes: "Wear slip-on shoes for temple hall entrance.",
          },
          {
            id: "it-2",
            placeId: "pl-ichiran-asakusa",
            name: "Ichiran Ramen Asakusa",
            category: "food",
            lat: 35.7121,
            lng: 139.7984,
            dayId: "day-1",
            position: "a1",
            durationMinutes: 45,
            cost: 18,
            notes: "Order extra chashu pork and rich garlic broth level.",
          },
          {
            id: "it-3",
            placeId: "pl-tokyo-skytree",
            name: "Tokyo Skytree",
            category: "sight",
            lat: 35.710063,
            lng: 139.8107,
            dayId: "day-1",
            position: "a2",
            durationMinutes: 90,
            cost: 32,
            notes: "Observation deck fast-pass slot booked for 17:30 sunset.",
          },
        ],
      },
      {
        id: "day-2",
        tripId: "trip-japan-2026",
        dayNumber: 2,
        date: "2026-04-03",
        title: "Digital Art & Waterfront Scenery",
        places: [
          {
            id: "it-4",
            placeId: "pl-teamlab-planets",
            name: "teamLab Planets TOKYO",
            category: "sight",
            lat: 35.6491,
            lng: 139.7898,
            dayId: "day-2",
            position: "b0",
            durationMinutes: 120,
            cost: 40,
            notes: "Water mirror room: bare feet required.",
          },
        ],
      },
    ],
    backlog: [
      {
        id: "it-bk-1",
        placeId: "pl-shibuya-sky",
        name: "SHIBUYA SKY",
        category: "sight",
        lat: 35.6586,
        lng: 139.7013,
        dayId: null,
        position: "bk0",
        durationMinutes: 90,
        cost: 25,
        notes: "Best night view over the scramble crossing.",
      },
      {
        id: "it-bk-2",
        placeId: "pl-fushimi-inari",
        name: "Fushimi Inari Taisha",
        category: "nature",
        lat: 34.9671,
        lng: 135.7727,
        dayId: null,
        position: "bk1",
        durationMinutes: 150,
        cost: 0,
        notes: "Early morning hike before crowds arrive.",
      },
    ],
    expenses: [
      { id: "exp-1", tripId: "trip-japan-2026", category: "stay", title: "Asakusa View Hotel (4 nights)", amount: 580, currency: "USD", paidAt: "2026-04-02" },
      { id: "exp-2", tripId: "trip-japan-2026", category: "activity", title: "Tokyo Skytree Combo Tickets", amount: 64, currency: "USD", paidAt: "2026-04-02" },
      { id: "exp-3", tripId: "trip-japan-2026", category: "food", title: "Ichiran Ramen dinner", amount: 36, currency: "USD", paidAt: "2026-04-02" },
      { id: "exp-4", tripId: "trip-japan-2026", category: "activity", title: "teamLab Planets tickets (2p)", amount: 80, currency: "USD", paidAt: "2026-04-03" },
      { id: "exp-5", tripId: "trip-japan-2026", category: "transport", title: "Suica IC Card top-up", amount: 100, currency: "USD", paidAt: "2026-04-02" },
    ],
    checklists: [
      { id: "chk-1", tripId: "trip-japan-2026", section: "pre-flight", title: "Japan Visit Web immigration & customs QR", completed: true },
      { id: "chk-2", tripId: "trip-japan-2026", section: "pre-flight", title: "Passport valid 6+ months", completed: true },
      { id: "chk-3", tripId: "trip-japan-2026", section: "pre-flight", title: "eSIM activation confirmed", completed: false },
      { id: "chk-4", tripId: "trip-japan-2026", section: "packing", title: "Universal plug adapter Type A/B", completed: true },
      { id: "chk-5", tripId: "trip-japan-2026", section: "packing", title: "Comfortable walking sneakers (20k steps/day)", completed: true },
      { id: "chk-6", tripId: "trip-japan-2026", section: "packing", title: "Coin purse for 100/500 yen coins", completed: false },
      { id: "chk-7", tripId: "trip-japan-2026", section: "buy-on-site", title: "Load Suica/Pasmo card on Apple Wallet", completed: false },
      { id: "chk-8", tripId: "trip-japan-2026", section: "buy-on-site", title: "Buy 7-Eleven egg salad sandwich on arrival", completed: false },
    ],
  },
];

export interface PresignUploadPayload {
  filename: string;
  contentType: string;
  fileSize: number;
  visibility: "public" | "private";
  placeId?: string;
  dayId?: string;
  ownerType?: "trip_cover" | "place" | "day";
}

export const apiClient = {
  async getGeoStats(): Promise<GeoStats> {
    try {
      const res = await fetch("/api/v1/geo/stats");
      if (res.ok) return await res.json();
    } catch {}
    return MOCK_GEO_STATS;
  },

  async getPlaces(): Promise<Place[]> {
    try {
      const res = await fetch("/api/v1/geo/places");
      if (res.ok) return await res.json();
    } catch {}
    return MOCK_PLACES;
  },

  async getTrips(): Promise<Trip[]> {
    try {
      const res = await fetch("/api/v1/trips");
      if (res.ok) return await res.json();
    } catch {}
    return MOCK_TRIPS;
  },

  async getTrip(tripId: string): Promise<Trip | null> {
    try {
      const res = await fetch(`/api/v1/trips/${tripId}`);
      if (res.ok) return await res.json();
    } catch {}
    return MOCK_TRIPS.find((t) => t.id === tripId) || MOCK_TRIPS[0];
  },

  async updateTripCover(tripId: string, coverUrl: string): Promise<void> {
    const trip = MOCK_TRIPS.find((t) => t.id === tripId);
    if (trip) {
      trip.coverImageUrl = coverUrl;
    }
  },

  async getPublicShare(slug: string): Promise<PublicTripView | null> {
    try {
      const res = await fetch(`/api/v1/share/${slug}`);
      if (res.ok) {
        const data = await res.json();
        if (data.trip) return data.trip;
        return data;
      }
    } catch {}
    const trip = MOCK_TRIPS.find((t) => t.slug === slug) || MOCK_TRIPS[0];
    if (!trip) return null;

    return {
      slug: trip.slug,
      title: trip.title,
      description: trip.description,
      startDate: trip.startDate,
      endDate: trip.endDate,
      coverImageUrl: trip.coverImageUrl,
      countries: trip.countries,
      cities: trip.cities,
      days: trip.days.map((day) => ({
        id: day.id,
        dayNumber: day.dayNumber,
        date: day.date,
        title: day.title,
        places: day.places.map((p, idx) => ({
          id: p.id,
          name: p.name,
          category: p.category,
          lat: p.lat,
          lng: p.lng,
          notes: p.notes,
          durationMinutes: p.durationMinutes,
          orderNumber: idx + 1,
        })),
      })),
    };
  },

  async clonePublicTrip(slug: string): Promise<Trip> {
    try {
      const res = await fetch(`/api/v1/share/${slug}/clone`, {
        method: "POST",
      });
      if (res.ok) {
        const data = await res.json();
        return data.trip;
      }
    } catch {}

    // Local fallback clone
    const orig = MOCK_TRIPS.find((t) => t.slug === slug) || MOCK_TRIPS[0];
    const newId = `trip-cloned-${Date.now()}`;
    const cloned: Trip = {
      ...JSON.parse(JSON.stringify(orig)),
      id: newId,
      title: `${orig.title} (Copy)`,
      budgetLimit: 0,
      expenses: [],
      slug: `copy-${slug}-${Date.now()}`,
    };
    MOCK_TRIPS.push(cloned);
    return cloned;
  },

  async createPresignedUpload(tripId: string, payload: PresignUploadPayload): Promise<{
    uploadUrl: string;
    downloadUrl: string;
    attachmentId: string;
  }> {
    try {
      const res = await fetch(`/api/v1/trips/${tripId}/attachments/presign`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (res.ok) return await res.json();
    } catch {}

    // Mock presign fallback
    const mockId = `att-${Date.now()}`;
    return {
      uploadUrl: `http://localhost:9000/atlasflow-media/trips/${tripId}/${mockId}?upload=true`,
      downloadUrl: URL.createObjectURL(new Blob(["mock content"], { type: payload.contentType })),
      attachmentId: mockId,
    };
  },

  async uploadFileToMinIO(uploadUrl: string, file: File): Promise<void> {
    try {
      await fetch(uploadUrl, {
        method: "PUT",
        headers: {
          "Content-Type": file.type,
        },
        body: file,
      });
    } catch {
      // In offline/mock mode, simulate network upload delay
      await new Promise((r) => setTimeout(r, 600));
    }
  },

  async confirmUpload(tripId: string, attachmentId: string): Promise<any> {
    try {
      const res = await fetch(`/api/v1/trips/${tripId}/attachments/${attachmentId}/confirm`, {
        method: "POST",
      });
      if (res.ok) return await res.json();
    } catch {}
    return { id: attachmentId, isReady: true };
  },

  async applyChecklistTemplate(tripId: string, templateKey: string): Promise<ChecklistItem[]> {
    try {
      const res = await fetch(`/api/v1/trips/${tripId}/checklist/templates/${templateKey}`, {
        method: "POST",
      });
      if (res.ok) {
        const data = await res.json();
        return data.items;
      }
    } catch {}

    // Local fallback for starter templates
    const templates: Record<string, { section: "pre-flight" | "packing" | "buy-on-site"; title: string }[]> = {
      digital_nomad: [
        { section: "pre-flight", title: "Passport & Visa" },
        { section: "pre-flight", title: "Travel insurance policy" },
        { section: "pre-flight", title: "Offline maps & boarding pass downloaded" },
        { section: "packing", title: "Laptop & charger" },
        { section: "packing", title: "Universal power adapter" },
        { section: "packing", title: "Power bank 20000mAh" },
        { section: "packing", title: "Noise-canceling headphones" },
        { section: "packing", title: "Backup credit / debit cards" },
        { section: "buy-on-site", title: "Local SIM / eSIM data plan" },
        { section: "buy-on-site", title: "Local cash / ATM emergency" },
      ],
      beach_sun: [
        { section: "pre-flight", title: "Hotel / resort confirmation" },
        { section: "packing", title: "Swimwear & beach towel" },
        { section: "packing", title: "Sunscreen SPF 50+ & lip balm" },
        { section: "packing", title: "Sunglasses & sun hat" },
        { section: "packing", title: "Waterproof phone pouch" },
        { section: "buy-on-site", title: "After-sun aloe vera gel" },
      ],
      mountain_trekking: [
        { section: "pre-flight", title: "Trail permits & offline topo maps" },
        { section: "pre-flight", title: "Rescue insurance coverage" },
        { section: "packing", title: "Trail hiking boots (broken in)" },
        { section: "packing", title: "Waterproof shell jacket" },
        { section: "packing", title: "First aid kit & blister pads" },
        { section: "packing", title: "Headlamp with extra batteries" },
        { section: "buy-on-site", title: "Trail energy bars & electrolytes" },
      ],
      winter_city: [
        { section: "pre-flight", title: "Museum reservation slots" },
        { section: "packing", title: "Thermal base layers & socks" },
        { section: "packing", title: "Insulated winter jacket" },
        { section: "packing", title: "Gloves, scarf & beanie" },
        { section: "packing", title: "Compact storm umbrella" },
        { section: "buy-on-site", title: "City transit card" },
      ],
    };

    const trip = MOCK_TRIPS.find((t) => t.id === tripId);
    if (!trip) return [];

    const itemsToAdd = templates[templateKey] || templates.digital_nomad;
    const existingTitles = new Set(trip.checklists.map((c) => c.title.toLowerCase()));

    for (const item of itemsToAdd) {
      if (!existingTitles.has(item.title.toLowerCase())) {
        trip.checklists.push({
          id: `chk-${Date.now()}-${Math.random().toString(36).substring(7)}`,
          tripId,
          section: item.section,
          title: item.title,
          completed: false,
        });
      }
    }
    return [...trip.checklists];
  },

  async reorderPlace(tripId: string, itemId: string, targetDayId: string | null, newIndex: number): Promise<Trip> {
    const trip = MOCK_TRIPS.find((t) => t.id === tripId) || MOCK_TRIPS[0];
    let movedItem: PlaceRef | undefined;

    for (const day of trip.days) {
      const idx = day.places.findIndex((p) => p.id === itemId);
      if (idx !== -1) {
        movedItem = day.places.splice(idx, 1)[0];
        break;
      }
    }
    if (!movedItem) {
      const bIdx = trip.backlog.findIndex((p) => p.id === itemId);
      if (bIdx !== -1) {
        movedItem = trip.backlog.splice(bIdx, 1)[0];
      }
    }

    if (!movedItem) return trip;

    movedItem.dayId = targetDayId;
    movedItem.position = `pos-${Date.now()}`;

    if (targetDayId === null) {
      trip.backlog.splice(newIndex, 0, movedItem);
    } else {
      const targetDay = trip.days.find((d) => d.id === targetDayId);
      if (targetDay) {
        targetDay.places.splice(newIndex, 0, movedItem);
      }
    }

    MOCK_TRIPS = [...MOCK_TRIPS];
    return { ...trip };
  },

  async toggleChecklistItem(tripId: string, itemId: string): Promise<ChecklistItem[]> {
    const trip = MOCK_TRIPS.find((t) => t.id === tripId);
    if (!trip) return [];
    const item = trip.checklists.find((c) => c.id === itemId);
    if (item) {
      item.completed = !item.completed;
    }
    return [...trip.checklists];
  },

  async addExpense(tripId: string, expense: Omit<Expense, "id">): Promise<Expense> {
    const trip = MOCK_TRIPS.find((t) => t.id === tripId);
    const newExp: Expense = {
      ...expense,
      id: `exp-${Date.now()}`,
    };
    if (trip) {
      trip.expenses.unshift(newExp);
    }
    return newExp;
  },
};
