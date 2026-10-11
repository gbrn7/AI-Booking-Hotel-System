import QuickSearch from "./QuickSearch";
import { Star, ShieldCheck, Clock } from "lucide-react";

export default function HeroSection() {
  return (
    <section className="relative pt-12 pb-24 md:pt-20 md:pb-32 bg-gradient-to-b from-brand-50/50 via-white to-slate-50 overflow-hidden">
      {/* Decorative background blur circles */}
      <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[300px] bg-brand-200/30 rounded-full blur-3xl pointer-events-none" />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center relative z-10">
        <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-white border border-slate-200 shadow-sm text-xs font-semibold text-slate-800 mb-6">
          <div className="flex text-amber-400">
            {[...Array(5)].map((_, i) => (
              <Star key={i} className="w-3.5 h-3.5 fill-current" />
            ))}
          </div>
          <span>Bintang 5 Beachfront Sanctuary</span>
        </div>

        <h1 className="text-4xl sm:text-5xl md:text-6xl font-serif text-slate-900 tracking-tight font-bold max-w-4xl mx-auto leading-tight mb-6">
          Menemukan Kedamaian Elegan di Pesisir Samudra
        </h1>

        <p className="text-base sm:text-lg text-slate-600 max-w-2xl mx-auto mb-10 leading-relaxed font-light">
          Nikmati kenyamanan kamar berkelas dunia dengan pemandangan lautan tak berbatas.
          Didukung layanan pemesanan cerdas dan asisten concierge AI 24 jam nonstop.
        </p>

        {/* Quick Availability Search Bar */}
        <div className="mb-12">
          <QuickSearch />
        </div>

        {/* Value props badges */}
        <div className="flex flex-wrap items-center justify-center gap-6 sm:gap-10 text-xs font-medium text-slate-500">
          <div className="flex items-center gap-2">
            <ShieldCheck className="w-4 h-4 text-brand-600" />
            <span>Garansi Harga Terbaik & Transparan</span>
          </div>
          <div className="flex items-center gap-2">
            <Clock className="w-4 h-4 text-brand-600" />
            <span>Hold Reservasi Instan 2 Jam</span>
          </div>
          <div className="flex items-center gap-2">
            <Star className="w-4 h-4 text-brand-600" />
            <span>Resepsionis Manusia & AI Terpadu</span>
          </div>
        </div>
      </div>
    </section>
  );
}
