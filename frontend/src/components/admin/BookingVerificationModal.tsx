"use client";

import { useState } from "react";
import { Booking } from "@/types";
import { confirmBookingPayment } from "@/lib/api";
import { X, CheckCircle, ShieldCheck, Clock, CreditCard } from "lucide-react";

interface BookingVerificationModalProps {
  booking: Booking;
  onClose: () => void;
  onSuccess: (updated: Booking) => void;
}

export default function BookingVerificationModal({
  booking,
  onClose,
  onSuccess,
}: BookingVerificationModalProps) {
  const [loading, setLoading] = useState(false);

  const handleConfirm = async () => {
    setLoading(true);
    try {
      const updated = await confirmBookingPayment(booking.id);
      onSuccess(updated);
      onClose();
    } catch (err: any) {
      alert("Gagal mengonfirmasi: " + err.message);
    } finally {
      setLoading(false);
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      maximumFractionDigits: 0,
    }).format(price);
  };

  return (
    <div className="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4 animate-in fade-in duration-200">
      <div className="bg-white rounded-3xl max-w-lg w-full overflow-hidden shadow-2xl border border-slate-100">
        {/* Header */}
        <div className="p-6 bg-slate-900 text-white flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <ShieldCheck className="w-5 h-5 text-brand-400" />
            <h3 className="text-base font-bold">Verifikasi Pembayaran Reservasi</h3>
          </div>
          <button
            onClick={onClose}
            className="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-5">
          <div className="p-4 rounded-2xl bg-amber-50 border border-amber-200 text-amber-900 text-xs flex items-start gap-2.5">
            <Clock className="w-4 h-4 text-amber-700 shrink-0 mt-0.5" />
            <div>
              <span className="font-bold">Status Saat Ini: HOLD (Menunggu Pembayaran)</span>
              <p className="mt-0.5 text-amber-800 font-light">
                Periksa mutasi rekening bank hotel sebelum melakukan konfirmasi agar status kamar terkunci secara permanen.
              </p>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4 text-xs">
            <div className="p-3.5 rounded-xl bg-slate-50 border border-slate-100">
              <span className="text-slate-400 block mb-0.5">Kode Reservasi</span>
              <span className="font-mono font-bold text-slate-900 text-sm">
                {booking.booking_code}
              </span>
            </div>
            <div className="p-3.5 rounded-xl bg-slate-50 border border-slate-100">
              <span className="text-slate-400 block mb-0.5">Nama Tamu</span>
              <span className="font-bold text-slate-900 text-sm">
                {booking.guest?.name || "Tamu Hotel"}
              </span>
            </div>
            <div className="p-3.5 rounded-xl bg-slate-50 border border-slate-100">
              <span className="text-slate-400 block mb-0.5">Tipe Kamar</span>
              <span className="font-bold text-slate-900">
                {booking.room?.name || "Kamar"}
              </span>
            </div>
            <div className="p-3.5 rounded-xl bg-slate-50 border border-slate-100">
              <span className="text-slate-400 block mb-0.5">Durasi Menginap</span>
              <span className="font-bold text-slate-900">
                {booking.check_in_date} s.d. {booking.check_out_date} ({booking.total_nights} Malam)
              </span>
            </div>
          </div>

          {/* Payment summary */}
          <div className="p-4 rounded-2xl bg-slate-900 text-white flex items-center justify-between">
            <div>
              <span className="text-xs text-slate-400 flex items-center gap-1.5">
                <CreditCard className="w-3.5 h-3.5 text-brand-400" />
                <span>Total Tagihan Booking</span>
              </span>
              <span className="text-xl font-bold text-brand-400">
                {formatPrice(booking.total_price)}
              </span>
            </div>
            <div className="text-right text-[11px] text-slate-400">
              <span>BCA 8870-1234-5678-9000</span>
              <div className="text-white font-medium">Grand Oasis Hotel</div>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="p-6 bg-slate-50 border-t border-slate-100 flex items-center justify-end gap-3">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2.5 rounded-xl border border-slate-200 text-xs font-semibold text-slate-600 hover:bg-slate-100 transition-colors"
          >
            Tutup
          </button>
          <button
            type="button"
            onClick={handleConfirm}
            disabled={loading}
            className="px-5 py-2.5 rounded-xl bg-emerald-600 hover:bg-emerald-700 active:scale-[0.98] text-white text-xs font-semibold flex items-center gap-2 transition-all shadow-md shadow-emerald-600/20 disabled:opacity-50"
          >
            {loading ? (
              <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
            ) : (
              <>
                <CheckCircle className="w-4 h-4" />
                <span>Verifikasi & Set CONFIRMED</span>
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
