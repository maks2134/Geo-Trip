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
  {
    id: "pl-eiffel-tower",
    name: "Eiffel Tower",
    category: "sight",
    status: "visited",
    lat: 48.8584,
    lng: 2.2945,
    city: "Paris",
    country: "France",
    countryCode: "FRA",
    rating: 4.7,
    notes: "Visited at night during champagne sparkle illumination.",
  },
  {
    id: "pl-colosseum",
    name: "Colosseum & Roman Forum",
    category: "sight",
    status: "visited",
    lat: 41.8902,
    lng: 12.4922,
    city: "Rome",
    country: "Italy",
    countryCode: "ITA",
    rating: 4.8,
  },
  {
    id: "pl-golden-gate",
    name: "Golden Gate Bridge Overlook",
    category: "nature",
    status: "visited",
    lat: 37.8199,
    lng: -122.4783,
    city: "San Francisco",
    country: "USA",
    countryCode: "USA",
    rating: 4.9,
  },
  {
    id: "pl-wat-arun",
    name: "Wat Arun (Temple of Dawn)",
    category: "sight",
    status: "visited",
    lat: 13.7437,
    lng: 100.4889,
    city: "Bangkok",
    country: "Thailand",
    countryCode: "THA",
    rating: 4.7,
  },
];

let MOCK_TRIPS: Trip[] = [
  {
    id: "trip-japan-2026",
    title: "Japan Autumn Odyssey",
    slug: "japan-autumn-2026",
    description: "Tokyo neon alleys, historic Asakusa, modern digital art, and Kyoto mountain shrines.",
    status: "active",
    startDate: "2026-11-12",
    endDate: "2026-11-19",
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
        date: "2026-11-12",
        title: "Historical East Tokyo & Skytree",
        places: [
          {
            id: "ref-1-1",
            placeId: "pl-sensoji",
            name: "Senso-ji Ancient Temple",
            category: "sight",
            lat: 35.714765,
            lng: 139.796655,
            position: "a0",
            dayId: "day-1",
            durationMinutes: 75,
            cost: 10,
            notes: "Morning stroll before crowds arrive.",
            transitToNext: { mode: "walk", durationMinutes: 6, distanceKm: 0.4 },
          },
          {
            id: "ref-1-2",
            placeId: "pl-ichiran-asakusa",
            name: "Ichiran Ramen Asakusa",
            category: "food",
            lat: 35.7121,
            lng: 139.7984,
            position: "a1",
            dayId: "day-1",
            durationMinutes: 45,
            cost: 18,
            notes: "Lunch stop for tonkotsu refill.",
            transitToNext: { mode: "transit", durationMinutes: 14, distanceKm: 1.8 },
          },
          {
            id: "ref-1-3",
            placeId: "pl-tokyo-skytree",
            name: "Tokyo Skytree & Solamachi",
            category: "sight",
            lat: 35.710063,
            lng: 139.8107,
            position: "a2",
            dayId: "day-1",
            durationMinutes: 90,
            cost: 32,
            notes: "Sunset panoramic deck at 350m.",
          },
        ],
      },
      {
        id: "day-2",
        tripId: "trip-japan-2026",
        dayNumber: 2,
        date: "2026-11-13",
        title: "Digital Art & Shibuya Heights",
        places: [
          {
            id: "ref-2-1",
            placeId: "pl-teamlab-planets",
            name: "teamLab Planets TOKYO",
            category: "sight",
            lat: 35.6491,
            lng: 139.7898,
            position: "a0",
            dayId: "day-2",
            durationMinutes: 120,
            cost: 40,
            notes: "Barefoot immersive museum.",
            transitToNext: { mode: "transit", durationMinutes: 28, distanceKm: 9.2 },
          },
          {
            id: "ref-2-2",
            placeId: "pl-shibuya-sky",
            name: "SHIBUYA SKY Observation Deck",
            category: "sight",
            lat: 35.6586,
            lng: 139.7013,
            position: "a1",
            dayId: "day-2",
            durationMinutes: 90,
            cost: 25,
            notes: "Scramble crossing bird-eye viewpoint.",
          },
        ],
      },
      {
        id: "day-3",
        tripId: "trip-japan-2026",
        dayNumber: 3,
        date: "2026-11-14",
        title: "Kyoto Shinkansen & Torii Gates",
        places: [
          {
            id: "ref-3-1",
            placeId: "pl-fushimi-inari",
            name: "Fushimi Inari Taisha",
            category: "nature",
            lat: 34.9671,
            lng: 135.7727,
            position: "a0",
            dayId: "day-3",
            durationMinutes: 150,
            cost: 0,
            notes: "Hike up to Yotsutsuji intersection for Kyoto city view.",
          },
        ],
      },
    ],
    backlog: [
      {
        id: "ref-bl-1",
        placeId: "pl-harajuku",
        name: "Takeshita Street & Meiji Jingu",
        category: "sight",
        lat: 35.6702,
        lng: 139.7027,
        position: "b0",
        dayId: null,
        durationMinutes: 80,
        cost: 15,
        notes: "Crepes + serene forest shrine.",
      },
      {
        id: "ref-bl-2",
        placeId: "pl-tsukiji",
        name: "Tsukiji Outer Market Food Crawl",
        category: "food",
        lat: 35.6655,
        lng: 139.7708,
        position: "b1",
        dayId: null,
        durationMinutes: 90,
        cost: 35,
        notes: "Tamagoyaki, fresh uni, and wagyu skewers.",
      },
    ],
    expenses: [
      {
        id: "exp-1",
        tripId: "trip-japan-2026",
        dayId: "day-1",
        placeId: "pl-tokyo-skytree",
        title: "Skytree FastPass Admission (x2)",
        amount: 64,
        currency: "USD",
        category: "activity",
        paidAt: "2026-11-12",
      },
      {
        id: "exp-2",
        tripId: "trip-japan-2026",
        dayId: "day-1",
        placeId: "pl-ichiran-asakusa",
        title: "Ichiran Ramen bowls + matcha beer",
        amount: 36,
        currency: "USD",
        category: "food",
        paidAt: "2026-11-12",
      },
      {
        id: "exp-3",
        tripId: "trip-japan-2026",
        dayId: null,
        placeId: null,
        title: "Tokyo Bay Hotel 4 nights",
        amount: 720,
        currency: "USD",
        category: "stay",
        paidAt: "2026-10-20",
      },
      {
        id: "exp-4",
        tripId: "trip-japan-2026",
        dayId: null,
        placeId: null,
        title: "7-Day Japan Rail Pass",
        amount: 380,
        currency: "USD",
        category: "transport",
        paidAt: "2026-10-25",
      },
      {
        id: "exp-5",
        tripId: "trip-japan-2026",
        dayId: "day-2",
        placeId: "pl-teamlab-planets",
        title: "teamLab Planets entry passes",
        amount: 80,
        currency: "USD",
        category: "activity",
        paidAt: "2026-11-13",
      },
    ],
    checklists: [
      { id: "chk-1", tripId: "trip-japan-2026", section: "pre-flight", title: "Fill Visit Japan Web customs QR code", completed: true },
      { id: "chk-2", tripId: "trip-japan-2026", section: "pre-flight", title: "Order e-SIM data package (10GB)", completed: true },
      { id: "chk-3", tripId: "trip-japan-2026", section: "pre-flight", title: "Confirm travel insurance certificate", completed: true },
      { id: "chk-4", tripId: "trip-japan-2026", section: "packing", title: "Universal plug adapter Type A/B", completed: true },
      { id: "chk-5", tripId: "trip-japan-2026", section: "packing", title: "Comfortable walking sneakers (20k steps/day)", completed: true },
      { id: "chk-6", tripId: "trip-japan-2026", section: "packing", title: "Coin purse for 100/500 yen coins", completed: false },
      { id: "chk-7", tripId: "trip-japan-2026", section: "buy-on-site", title: "Load Suica/Pasmo card on Apple Wallet", completed: false },
      { id: "chk-8", tripId: "trip-japan-2026", section: "buy-on-site", title: "Buy 7-Eleven egg salad sandwich on arrival", completed: false },
    ],
  },
  {
    id: "trip-amalfi-2026",
    title: "Amalfi Coast & Rome Escapade",
    slug: "amalfi-roadtrip",
    description: "Cliffside coastal roads, limoncello tasting, and ancient Roman history.",
    status: "upcoming",
    startDate: "2026-09-04",
    endDate: "2026-09-12",
    coverImageUrl: "https://images.unsplash.com/photo-1533105079780-92b9be482077?w=1200&q=80",
    budgetLimit: 3400,
    currency: "EUR",
    countries: ["Italy"],
    cities: ["Rome", "Positano", "Capri"],
    days: [],
    backlog: [],
    expenses: [],
    checklists: [],
  },
];

export const apiClient = {
  async getGeoStats(): Promise<GeoStats> {
    try {
      const res = await fetch("/api/geo/stats");
      if (res.ok) return await res.json();
    } catch {}
    return MOCK_GEO_STATS;
  },

  async getPlaces(): Promise<Place[]> {
    try {
      const res = await fetch("/api/places");
      if (res.ok) return await res.json();
    } catch {}
    return MOCK_PLACES;
  },

  async getTrips(): Promise<Trip[]> {
    try {
      const res = await fetch("/api/trips");
      if (res.ok) return await res.json();
    } catch {}
    return MOCK_TRIPS;
  },

  async getTrip(tripId: string): Promise<Trip | null> {
    try {
      const res = await fetch(`/api/trips/${tripId}`);
      if (res.ok) return await res.json();
    } catch {}
    return MOCK_TRIPS.find((t) => t.id === tripId) || MOCK_TRIPS[0];
  },

  async getPublicShare(slug: string): Promise<PublicTripView | null> {
    try {
      const res = await fetch(`/api/share/${slug}`);
      if (res.ok) return await res.json();
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
