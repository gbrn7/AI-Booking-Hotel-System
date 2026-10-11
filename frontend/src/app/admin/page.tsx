"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { AdminStats } from "@/types";
import { fetchAdminStats } from "@/lib/api";
import {
  BedDouble,
  Clock,
  CheckCircle2,
  MessageSquare,
  ArrowRight,
  ShieldAlert,
  Sparkles,
} from "lucide-react";

export default function AdminDashboardPage() {
  const [stats, setStats] = useState<AdminStats>({
    total_rooms: 4,
    pending_hold_count: 0,
    confirmed_count: 0,
    active_conversations: 0,
  });

  useEffect(() => {
    fetchAdminStats("")
      .then((data) => setStats(data))
      .catch(() => {});
  }, []);

  const cards = [
    {
      title: "Total Tipe Kamar",
      value: stats.total_rooms,
      subtitle: "Katalog kamar terdaftar",
      icon: BedDouble,
      color: "bg-blue-500",
      link: "/admin/rooms",
    },
    {
      title: "Menunggu Verifikasi (HOLD)",
      value: stats.pending_hold_count,
      subtitle: "Reservasi butuh cek transfer",
      icon: Clock,
      color: "bg-amber-500",
      link: "/admin/bookings?status=HOLD",
    },
    {
      title: "Reservasi Terkonfirmasi",
      value: stats.confirmed_count,
      subtitle: "Pembayaran lunas & valid",
      icon: CheckCircle2,
      color: "bg-emerald-500",
      link: "/admin/bookings?status=CONFIRMED",
    },
    {
      title: "Percakapan Masuk",
      value: stats.active_conversations,
      subtitle: "WhatsApp & Web Widget",
      icon: MessageSquare,
      color: "bg-purple-500",
      link: "/admin/takeover",
    },
  ];

  return (
    <div className="space-y-8">
      {/* Welcome Banner */}
      <div className="bg-gradient-to-r from-slate-900 via-slate-800 to-brand-900 rounded-3xl p-8 text-white shadow-xl relative overflow-hidden">
        <div className="relative z-10 max-w-2xl">
          <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/10 backdrop-blur-md text-brand-300 text-xs font-semibold mb-4 border border-white/10">
            <Sparkles className="w-3.5 h-3.5" />
            <span>AI Concierge & Front Desk Synchronization</span>
          </div>
          <h1 className="text-3xl font-serif font-bold tracking-tight mb-2">
            Selamat Datang di Command Center Hotel
          </h1>
          <p className="text-sm text-slate-300 font-light leading-relaxed">
            AI Sales Agent Rama secara mandiri melayani reservasi tamu di WhatsApp dan Website.
            Gunakan dashboard ini untuk memverifikasi mutasi pembayaran dan mengambil alih obrolan kapan saja dibutuhkan.
          </p>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {cards.map((c, i) => (
          <Link
            key={i}
            href={c.link}
            className="p-6 rounded-2xl bg-white border border-slate-200 shadow-sm hover:shadow-md transition-all group flex flex-col justify-between"
          >
            <div>
              <div className="flex items-center justify-between mb-4">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
                  {c.title}
                </span>
                <div className={`w-9 h-9 rounded-xl ${c.color} text-white flex items-center justify-center shadow-sm`}>
                  <c.icon className="w-4 h-4" />
                </div>
              </div>
              <div className="text-3xl font-serif font-bold text-slate-900 mb-1">
                {c.value}
              </div>
              <p className="text-xs text-slate-400 font-light">{c.subtitle}</p>
            </div>

            <div className="mt-4 pt-3 border-t border-slate-100 flex items-center justify-between text-xs font-semibold text-slate-600 group-hover:text-brand-600 transition-colors">
              <span>Buka Menu</span>
              <ArrowRight className="w-3.5 h-3.5 group-hover:translate-x-1 transition-transform" />
            </div>
          </Link>
        ))}
      </div>

      {/* Quick operational links */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="p-6 rounded-2xl bg-white border border-slate-200 shadow-sm">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 rounded-xl bg-amber-50 border border-amber-200 text-amber-700 flex items-center justify-center font-bold">
              <ShieldAlert className="w-5 h-5" />
            </div>
            <div>
              <h3 className="font-serif font-bold text-slate-900 text-base">
                Antrean Verifikasi Booking HOLD
              </h3>
              <p className="text-xs text-slate-500 font-light">
                Periksa bukti transfer tamu sebelum batas waktu 2 jam kedaluwarsa.
              </p>
            </div>
          </div>
          <Link
            href="/admin/bookings?status=HOLD"
            className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-slate-900 text-white text-xs font-semibold hover:bg-slate-800 transition-colors"
          >
            <span>Buka Reservasi Menunggu Verifikasi</span>
            <ArrowRight className="w-3.5 h-3.5" />
          </Link>
        </div>

        <div className="p-6 rounded-2xl bg-white border border-slate-200 shadow-sm">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 rounded-xl bg-purple-50 border border-purple-200 text-purple-700 flex items-center justify-center font-bold">
              <MessageSquare className="w-5 h-5" />
            </div>
            <div>
              <h3 className="font-serif font-bold text-slate-900 text-base">
                Live Takeover Center
              </h3>
              <p className="text-xs text-slate-500 font-light">
                Ambil alih percakapan tamu WhatsApp atau Web saat ada eskalasi komplain.
              </p>
            </div>
          </div>
          <Link
            href="/admin/takeover"
            className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-purple-600 text-white text-xs font-semibold hover:bg-purple-700 transition-colors"
          >
            <span>Buka Live Chat Staf</span>
            <ArrowRight className="w-3.5 h-3.5" />
          </Link>
        </div>
      </div>
    </div>
  );
}
