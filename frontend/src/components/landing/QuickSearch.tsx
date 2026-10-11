"use client";

import { useState } from "react";
import { Calendar, Users, Search, AlertCircle, CheckCircle2 } from "lucide-react";
import { checkRoomAvailability } from "@/lib/api";
import { RoomAvailability } from "@/types";

interface QuickSearchProps {
  onResults?: (results: RoomAvailability[]) => void;
}

export default function QuickSearch({ onResults }: QuickSearchProps) {
  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  const dayAfter = new Date();
  dayAfter.setDate(dayAfter.getDate() + 3);

  const [checkIn, setCheckIn] = useState(tomorrow.toISOString().split("T")[0]);
  const [checkOut, setCheckOut] = useState(dayAfter.toISOString().split("T")[0]);
  const [guests, setGuests] = useState(2);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successCount, setSuccessCount] = useState<number | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccessCount(null);

    try {
      const avails = await checkRoomAvailability(checkIn, checkOut);
      setSuccessCount(avails.length);
      if (onResults) onResults(avails);
      // Smooth scroll to rooms section
      const el = document.getElementById("rooms");
      if (el) el.scrollIntoView({ behavior: "smooth" });
    } catch (err: any) {
      setError(err.message || "Gagal memeriksa ketersediaan");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full max-w-4xl mx-auto bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-slate-100 p-4 sm:p-6 transition-all">
      <form onSubmit={handleSearch} className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 items-end">
        <div>
          <label className="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1.5">
            Check-In
          </label>
          <div className="relative">
            <input
              type="date"
              value={checkIn}
              onChange={(e) => setCheckIn(e.target.value)}
              className="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
              required
            />
            <Calendar className="w-4 h-4 text-slate-400 absolute right-3 top-3 pointer-events-none" />
          </div>
        </div>

        <div>
          <label className="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1.5">
            Check-Out
          </label>
          <div className="relative">
            <input
              type="date"
              value={checkOut}
              onChange={(e) => setCheckOut(e.target.value)}
              className="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
              required
            />
            <Calendar className="w-4 h-4 text-slate-400 absolute right-3 top-3 pointer-events-none" />
          </div>
        </div>

        <div>
          <label className="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1.5">
            Tamu
          </label>
          <div className="relative">
            <select
              value={guests}
              onChange={(e) => setGuests(Number(e.target.value))}
              className="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 bg-white"
            >
              <option value={1}>1 Dewasa</option>
              <option value={2}>2 Dewasa</option>
              <option value={3}>3 Dewasa</option>
              <option value={4}>4 Dewasa / Keluarga</option>
              <option value={5}>5+ Orang</option>
            </select>
            <Users className="w-4 h-4 text-slate-400 absolute right-3 top-3 pointer-events-none" />
          </div>
        </div>

        <div>
          <button
            type="submit"
            disabled={loading}
            className="w-full h-[42px] rounded-xl bg-brand-500 hover:bg-brand-600 active:scale-[0.98] text-white font-medium text-sm flex items-center justify-center gap-2 transition-all shadow-md shadow-brand-500/25 disabled:opacity-70"
          >
            {loading ? (
              <span className="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
            ) : (
              <>
                <Search className="w-4 h-4" />
                <span>Cek Ketersediaan</span>
              </>
            )}
          </button>
        </div>
      </form>

      {error && (
        <div className="mt-3.5 flex items-center gap-2 p-2.5 rounded-lg bg-rose-50 border border-rose-200 text-rose-700 text-xs">
          <AlertCircle className="w-4 h-4 shrink-0" />
          <span>{error}</span>
        </div>
      )}

      {successCount !== null && (
        <div className="mt-3.5 flex items-center gap-2 p-2.5 rounded-lg bg-emerald-50 border border-emerald-200 text-emerald-800 text-xs">
          <CheckCircle2 className="w-4 h-4 shrink-0" />
          <span>Ditemukan {successCount} tipe kamar dengan status ketersediaan pada tanggal tersebut.</span>
        </div>
      )}
    </div>
  );
}
