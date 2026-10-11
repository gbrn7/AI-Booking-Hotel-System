"use client";

import { Room } from "@/types";
import { Users, Check, MessageSquareCode } from "lucide-react";

interface RoomCardProps {
  room: Room;
  onAskAI?: (roomName: string) => void;
}

export default function RoomCard({ room, onAskAI }: RoomCardProps) {
  let images: string[] = [];
  try {
    images = JSON.parse(room.image_urls || "[]");
  } catch {
    images = ["https://images.unsplash.com/photo-1590490360182-c33d57733427?auto=format&fit=crop&w=800&q=80"];
  }
  const coverImage = images[0] || "https://images.unsplash.com/photo-1590490360182-c33d57733427?auto=format&fit=crop&w=800&q=80";

  let amenitiesList: string[] = [];
  try {
    amenitiesList = JSON.parse(room.amenities || "[]");
  } catch {
    amenitiesList = ["Free WiFi", "AC", "Breakfast"];
  }

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      maximumFractionDigits: 0,
    }).format(price);
  };

  const handleAskAI = () => {
    if (onAskAI) {
      onAskAI(room.name);
    } else {
      window.dispatchEvent(
        new CustomEvent("open-ai-chat", {
          detail: { query: `Halo Rama, bisakah jelaskan detail dan promo untuk kamar ${room.name}?` },
        })
      );
    }
  };

  return (
    <div className="bg-white rounded-2xl overflow-hidden border border-slate-100 shadow-sm hover:shadow-md transition-all duration-300 flex flex-col group">
      {/* Image container */}
      <div className="relative h-64 overflow-hidden bg-slate-100">
        <img
          src={coverImage}
          alt={room.name}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
        />
        <div className="absolute top-3.5 left-3.5 px-3 py-1 rounded-full bg-slate-900/80 backdrop-blur-md text-white text-[11px] font-semibold uppercase tracking-wider">
          {room.type}
        </div>
        <div className="absolute top-3.5 right-3.5 flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-white/90 backdrop-blur-md text-slate-800 text-xs font-medium shadow-sm">
          <Users className="w-3.5 h-3.5 text-brand-600" />
          <span>Maks {room.capacity_guests} Tamu</span>
        </div>
      </div>

      {/* Content */}
      <div className="p-6 flex-1 flex flex-col justify-between">
        <div>
          <h3 className="text-xl font-bold text-slate-900 mb-2 group-hover:text-brand-600 transition-colors">
            {room.name}
          </h3>
          <p className="text-sm text-slate-500 line-clamp-2 mb-4 leading-relaxed font-light">
            {room.description}
          </p>

          {/* Amenities tags */}
          <div className="flex flex-wrap gap-1.5 mb-6">
            {amenitiesList.slice(0, 4).map((amenity, idx) => (
              <span
                key={idx}
                className="inline-flex items-center gap-1 px-2.5 py-1 rounded-md bg-slate-50 border border-slate-100 text-slate-600 text-xs"
              >
                <Check className="w-3 h-3 text-brand-500" />
                <span>{amenity}</span>
              </span>
            ))}
          </div>
        </div>

        {/* Pricing and Action */}
        <div className="pt-4 border-t border-slate-100 flex items-center justify-between gap-3">
          <div>
            <span className="text-xs text-slate-400 block font-light">Mulai dari</span>
            <div className="text-lg font-bold text-slate-900">
              {formatPrice(room.base_price_per_night)}
              <span className="text-xs font-normal text-slate-500"> /malam</span>
            </div>
          </div>

          <button
            onClick={handleAskAI}
            className="flex items-center gap-1.5 px-3.5 py-2 rounded-xl bg-brand-50 hover:bg-brand-100 text-brand-700 text-xs font-semibold transition-colors border border-brand-200"
          >
            <MessageSquareCode className="w-3.5 h-3.5" />
            <span>Tanya Rama</span>
          </button>
        </div>
      </div>
    </div>
  );
}
