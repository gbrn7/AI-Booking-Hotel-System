"use client";

import { useEffect, useState } from "react";
import { Room } from "@/types";
import RoomCard from "./RoomCard";
import { fetchRooms } from "@/lib/api";

const defaultRooms: Room[] = [
  {
    id: 1,
    name: "Superior King Room",
    type: "standard",
    base_price_per_night: 750000,
    capacity_guests: 2,
    total_units: 10,
    description: "Kamar nyaman dengan ranjang King Bed, smart TV 43 inci, pemandangan kota, dan kamar mandi modern.",
    amenities: JSON.stringify(["King Bed", "Smart TV 43\"", "Free WiFi", "Rain Shower"]),
    image_urls: JSON.stringify(["https://images.unsplash.com/photo-1590490360182-c33d57733427?auto=format&fit=crop&w=800&q=80"]),
    is_active: true,
  },
  {
    id: 2,
    name: "Deluxe Ocean View",
    type: "deluxe",
    base_price_per_night: 1250000,
    capacity_guests: 2,
    total_units: 8,
    description: "Kamar elegan dengan balkon pribadi yang menghadap langsung ke lautan biru, bathtub marmer, dan sarapan buffet.",
    amenities: JSON.stringify(["Private Balcony", "Ocean View", "Marble Bathtub", "Buffet Breakfast"]),
    image_urls: JSON.stringify(["https://images.unsplash.com/photo-1582719478250-c89cae4dc85b?auto=format&fit=crop&w=800&q=80"]),
    is_active: true,
  },
  {
    id: 3,
    name: "Executive Ocean Suite",
    type: "suite",
    base_price_per_night: 2200000,
    capacity_guests: 3,
    total_units: 4,
    description: "Suite mewah dengan ruang tamu terpisah, jacuzzi pribadi di teras, akses exclusive lounge, dan pemandangan sunset spektakuler.",
    amenities: JSON.stringify(["Private Jacuzzi", "Separate Living Room", "Executive Lounge", "Sunset View"]),
    image_urls: JSON.stringify(["https://images.unsplash.com/photo-1631049307264-da0ec9d70304?auto=format&fit=crop&w=800&q=80"]),
    is_active: true,
  },
  {
    id: 4,
    name: "Grand Family Villa",
    type: "family_suite",
    base_price_per_night: 3500000,
    capacity_guests: 5,
    total_units: 2,
    description: "Villa keluarga dua lantai dengan kolam renang pribadi, dapur mini, 2 kamar tidur terpisah, dan ruang santai keluarga luas.",
    amenities: JSON.stringify(["Private Plunge Pool", "2 Bedrooms", "Kitchenette", "Garden View"]),
    image_urls: JSON.stringify(["https://images.unsplash.com/photo-1566073771259-6a8506099945?auto=format&fit=crop&w=800&q=80"]),
    is_active: true,
  },
];

export default function RoomGrid() {
  const [rooms, setRooms] = useState<Room[]>(defaultRooms);
  const [filterType, setFilterType] = useState<string>("all");

  useEffect(() => {
    fetchRooms()
      .then((data) => {
        if (data && data.length > 0) setRooms(data);
      })
      .catch(() => {
        // Fallback to default rooms if backend is starting up
      });
  }, []);

  const filteredRooms = filterType === "all"
    ? rooms
    : rooms.filter((r) => r.type.toLowerCase() === filterType.toLowerCase());

  return (
    <section id="rooms" className="py-20 bg-slate-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-12">
          <span className="text-xs font-semibold tracking-widest text-brand-600 uppercase mb-2 block">
            PILIHAN AKOMODASI
          </span>
          <h2 className="text-3xl sm:text-4xl font-bold text-slate-900 mb-4">
            Kamar & Villa Mewah Dirancang Untuk Relaksasi
          </h2>
          <p className="text-slate-600 text-sm sm:text-base font-light">
            Setiap unit dilengkapi amenitas premium, pemandangan memukau, serta koneksi instan ke asisten Rama untuk pertanyaan ketersediaan dan booking.
          </p>

          {/* Filter tabs */}
          <div className="flex flex-wrap items-center justify-center gap-2 mt-8">
            {[
              { id: "all", label: "Semua Kamar" },
              { id: "standard", label: "Superior" },
              { id: "deluxe", label: "Deluxe Ocean" },
              { id: "suite", label: "Executive Suite" },
              { id: "family_suite", label: "Family Villa" },
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setFilterType(tab.id)}
                className={`px-4 py-2 rounded-xl text-xs font-medium transition-all ${
                  filterType === tab.id
                    ? "bg-slate-900 text-white shadow-sm"
                    : "bg-white text-slate-600 hover:bg-slate-100 border border-slate-200"
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>
        </div>

        {/* Room cards grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {filteredRooms.map((room) => (
            <RoomCard key={room.id} room={room} />
          ))}
        </div>
      </div>
    </section>
  );
}
