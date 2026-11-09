import React from "react";
import { notFound } from "next/navigation";
import { apiClient } from "@/shared/lib/api";
import { PublicShareView } from "@/features/share/PublicShareView";
import { PrintDossier } from "@/features/trip-workspace/PrintDossier";
import { MapCanvas } from "@/shared/map/MapCanvas";
import { Trip } from "@/entities/trip";

interface SharePageProps {
  params: Promise<{ slug: string }>;
}

export default async function SharePage({ params }: SharePageProps) {
  const { slug } = await params;
  const publicTrip = await apiClient.getPublicShare(slug);

  if (!publicTrip) {
    notFound();
  }

  const readOnlyTripAdapter: Trip = {
    id: publicTrip.slug,
    title: publicTrip.title,
    slug: publicTrip.slug,
    description: publicTrip.description,
    status: "active",
    startDate: publicTrip.startDate,
    endDate: publicTrip.endDate,
    coverImageUrl: publicTrip.coverImageUrl,
    budgetLimit: 0,
    currency: "USD",
    countries: publicTrip.countries,
    cities: publicTrip.cities,
    days: publicTrip.days.map((d) => ({
      id: d.id,
      tripId: publicTrip.slug,
      dayNumber: d.dayNumber,
      date: d.date,
      title: d.title,
      places: d.places.map((p) => ({
        id: p.id,
        placeId: p.id,
        name: p.name,
        category: p.category,
        lat: p.lat,
        lng: p.lng,
        position: `pos-${p.orderNumber}`,
        dayId: d.id,
        durationMinutes: p.durationMinutes,
        notes: p.notes,
      })),
    })),
    backlog: [],
    expenses: [],
    checklists: [],
  };

  return (
    <main className="relative w-screen h-screen overflow-hidden">
      <MapCanvas places={[]} trip={readOnlyTripAdapter} readOnly={true} />
      <PublicShareView trip={publicTrip} />
      <PrintDossier trip={readOnlyTripAdapter} isPublicShare={true} />
    </main>
  );
}
