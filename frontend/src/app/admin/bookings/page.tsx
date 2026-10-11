"use client";

import { useState, useEffect } from "react";
import { Booking } from "@/types";
import { fetchAdminBookings } from "@/lib/api";
import BookingVerificationModal from "@/components/admin/BookingVerificationModal";
import { Search, RefreshCw, CheckCircle2, Clock, XCircle, ShieldCheck } from "lucide-react";

export default function BookingsAdminPage() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [statusFilter, setStatusFilter] = useState("ALL");
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedBooking, setSelectedBooking] = useState<Booking | null>(null);

  const loadData = async () => {
    setLoading(true);
    try {
      const data = await fetchAdminBookings(statusFilter === "ALL" ? undefined : statusFilter);
      setBookings(data);
    } catch {
      setBookings([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [statusFilter]);

  const filtered = bookings.filter((b) => {
    const q = searchQuery.toLowerCase();
    return (
      b.booking_code.toLowerCase().includes(q) ||
      (b.guest?.name && b.guest.name.toLowerCase().includes(q)) ||
      (b.guest?.phone_number && b.guest.phone_number.includes(q))
    );
  });

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
          <h1 className="text-2xl font-bold text-slate-900">
            Manajemen Reservasi Tamu
          </h1>
          <p className="text-xs text-slate-500 font-light mt-0.5">
            Daftar seluruh booking yang dibuat melalui AI Sales Agent (WhatsApp & Web Widget) dan status pembayaran.
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

      {/* Filter and Search Bar */}
      <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm flex flex-col md:flex-row gap-4 justify-between items-center">
        {/* Status Tabs */}
        <div className="flex gap-1.5 overflow-x-auto w-full md:w-auto">
          {[
            { id: "ALL", label: "Semua" },
            { id: "HOLD", label: "HOLD (Menunggu Pembayaran)" },
            { id: "CONFIRMED", label: "CONFIRMED (Lunas)" },
            { id: "CANCELLED", label: "CANCELLED" },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setStatusFilter(tab.id)}
              className={`px-3.5 py-2 rounded-xl text-xs font-semibold whitespace-nowrap transition-all ${
                statusFilter === tab.id
                  ? "bg-slate-900 text-white shadow-sm"
                  : "bg-slate-100 text-slate-600 hover:bg-slate-200"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Search Input */}
        <div className="relative w-full md:w-72">
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Cari kode booking / nama tamu..."
            className="w-full pl-9 pr-4 py-2 rounded-xl border border-slate-200 text-xs focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
          />
          <Search className="w-4 h-4 text-slate-400 absolute left-3 top-2.5 pointer-events-none" />
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-50/80 border-b border-slate-200 text-slate-500 uppercase tracking-wider font-semibold">
              <tr>
                <th className="p-4">Kode Booking</th>
                <th className="p-4">Tamu</th>
                <th className="p-4">Tipe Kamar</th>
                <th className="p-4">Tanggal Menginap</th>
                <th className="p-4">Total Biaya</th>
                <th className="p-4">Status</th>
                <th className="p-4 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {loading ? (
                <tr>
                  <td colSpan={7} className="p-8 text-center text-slate-400">
                    Memuat data reservasi...
                  </td>
                </tr>
              ) : filtered.length === 0 ? (
                <tr>
                  <td colSpan={7} className="p-8 text-center text-slate-400">
                    Tidak ada data reservasi ditemukan.
                  </td>
                </tr>
              ) : (
                filtered.map((item) => (
                  <tr key={item.id} className="hover:bg-slate-50/50 transition-colors">
                    <td className="p-4 font-mono font-bold text-slate-900">
                      {item.booking_code}
                    </td>
                    <td className="p-4">
                      <div className="font-semibold text-slate-900">
                        {item.guest?.name || "Tamu"}
                      </div>
                      <div className="text-[11px] text-slate-400">
                        {item.guest?.phone_number || "-"}
                      </div>
                    </td>
                    <td className="p-4 text-slate-700 font-medium">
                      {item.room?.name || "Kamar Hotel"}
                    </td>
                    <td className="p-4 text-slate-600">
                      <div>{item.check_in_date} s.d. {item.check_out_date}</div>
                      <div className="text-[10px] text-slate-400">({item.total_nights} Malam)</div>
                    </td>
                    <td className="p-4 font-bold text-slate-900">
                      {formatPrice(item.total_price)}
                    </td>
                    <td className="p-4">
                      {item.status === "CONFIRMED" ? (
                        <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200 font-semibold text-[11px]">
                          <CheckCircle2 className="w-3 h-3" />
                          <span>CONFIRMED</span>
                        </span>
                      ) : item.status === "HOLD" ? (
                        <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-amber-50 text-amber-700 border border-amber-200 font-semibold text-[11px]">
                          <Clock className="w-3 h-3" />
                          <span>HOLD</span>
                        </span>
                      ) : (
                        <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-slate-100 text-slate-500 font-medium text-[11px]">
                          <XCircle className="w-3 h-3" />
                          <span>{item.status}</span>
                        </span>
                      )}
                    </td>
                    <td className="p-4 text-right">
                      {item.status === "HOLD" ? (
                        <button
                          onClick={() => setSelectedBooking(item)}
                          className="px-3 py-1.5 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs shadow-sm transition-all"
                        >
                          Verifikasi
                        </button>
                      ) : (
                        <span className="text-[11px] text-slate-400">Selesai</span>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Verification Modal */}
      {selectedBooking && (
        <BookingVerificationModal
          booking={selectedBooking}
          onClose={() => setSelectedBooking(null)}
          onSuccess={(updated) => {
            setBookings((prev) => prev.map((b) => (b.id === updated.id ? updated : b)));
          }}
        />
      )}
    </div>
  );
}
