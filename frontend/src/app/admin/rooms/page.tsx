"use client";

import { useState, useEffect } from "react";
import { Room } from "@/types";
import { fetchRooms } from "@/lib/api";
import { BedDouble, Users, CheckCircle, RefreshCw, Layers } from "lucide-react";

export default function RoomsAdminPage() {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);

  const loadData = async () => {
    setLoading(true);
    try {
      const data = await fetchRooms();
      setRooms(data);
    } catch {
      setRooms([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      maximumFractionDigits: 0,
    }).format(price);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-serif font-bold text-slate-900">
            Inventori Tipe Kamar Hotel
          </h1>
          <p className="text-xs text-slate-500 font-light mt-0.5">
            Daftar tipe kamar, total kapasitas unit fisik yang dikunci oleh engine inventori atomik, dan tarif per malam.
          </p>
        </div>

        <button
          onClick={loadData}
          className="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-white border border-slate-200 text-xs font-semibold text-slate-700 hover:bg-slate-50 transition-all shadow-sm"
        >
          <RefreshCw className="w-3.5 h-3.5" />
          <span>Segarkan Data</span>
        </button>
      </div>

      {/* Grid of rooms */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {loading ? (
          <div className="col-span-full p-12 text-center text-xs text-slate-400">
            Memuat data kamar...
          </div>
        ) : (
          rooms.map((room) => {
            let images: string[] = [];
            try {
              images = JSON.parse(room.image_urls || "[]");
            } catch {
              images = [];
            }
            const cover = images[0] || "https://images.unsplash.com/photo-1590490360182-c33d57733427?auto=format&fit=crop&w=800&q=80";

            return (
              <div
                key={room.id}
                className="bg-white rounded-2xl border border-slate-200 overflow-hidden shadow-sm flex flex-col justify-between"
              >
                <div>
                  <div className="relative h-44 bg-slate-100">
                    <img
                      src={cover}
                      alt={room.name}
                      className="w-full h-full object-cover"
                    />
                    <div className="absolute top-3 left-3 px-2.5 py-1 rounded-full bg-slate-900/80 backdrop-blur-md text-white text-[10px] font-semibold uppercase">
                      {room.type}
                    </div>
                  </div>

                  <div className="p-5">
                    <h3 className="font-serif font-bold text-slate-900 text-base mb-1">
                      {room.name}
                    </h3>
                    <p className="text-xs text-slate-500 line-clamp-2 mb-4 font-light">
                      {room.description}
                    </p>

                    <div className="grid grid-cols-2 gap-3 text-xs mb-4">
                      <div className="p-2.5 rounded-xl bg-slate-50 border border-slate-100 flex items-center gap-2">
                        <Users className="w-4 h-4 text-brand-600 shrink-0" />
                        <div>
                          <span className="text-[10px] text-slate-400 block">Kapasitas</span>
                          <span className="font-semibold text-slate-800">
                            {room.capacity_guests} Tamu
                          </span>
                        </div>
                      </div>

                      <div className="p-2.5 rounded-xl bg-slate-50 border border-slate-100 flex items-center gap-2">
                        <Layers className="w-4 h-4 text-brand-600 shrink-0" />
                        <div>
                          <span className="text-[10px] text-slate-400 block">Total Unit Fisik</span>
                          <span className="font-semibold text-slate-800">
                            {room.total_units} Kamar
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="p-5 pt-0 border-t border-slate-100 mt-2 flex items-center justify-between">
                  <div>
                    <span className="text-[10px] text-slate-400 block">Tarif Dasar</span>
                    <span className="font-bold text-slate-900 font-serif text-sm">
                      {formatPrice(room.base_price_per_night)}
                    </span>
                  </div>

                  <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200 text-[10px] font-semibold">
                    <CheckCircle className="w-3 h-3" />
                    <span>Aktif Tersedia</span>
                  </span>
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}
