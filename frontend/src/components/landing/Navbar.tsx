import Link from "next/link";
import { Sparkles, Shield } from "lucide-react";

export default function Navbar() {
  return (
    <header className="sticky top-0 z-40 bg-white/90 backdrop-blur-md border-b border-slate-100 shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-20 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-brand-500 text-white flex items-center justify-center text-xl font-bold shadow-md">
            G
          </div>
          <div>
            <span className="text-xl tracking-wider font-bold text-slate-900 block leading-tight">
              GRAND OASIS
            </span>
            <span className="text-[10px] tracking-widest text-brand-600 uppercase font-semibold">
              Resort & Coastal Spa
            </span>
          </div>
        </Link>

        <nav className="hidden md:flex items-center gap-8 text-sm font-medium text-slate-600">
          <a href="#rooms" className="hover:text-brand-600 transition-colors">Tipe Kamar</a>
          <a href="#amenities" className="hover:text-brand-600 transition-colors">Fasilitas</a>
          <a href="#experience" className="hover:text-brand-600 transition-colors">Pengalaman</a>
          <a href="#faq" className="hover:text-brand-600 transition-colors">FAQ</a>
        </nav>

        <div className="flex items-center gap-3">
          <div className="hidden sm:flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-brand-50 text-brand-800 text-xs font-medium border border-brand-200">
            <Sparkles className="w-3.5 h-3.5 text-brand-600 animate-pulse" />
            <span>AI Concierge 24/7</span>
          </div>
          <Link
            href="/admin"
            className="flex items-center gap-1.5 px-4 py-2 rounded-lg bg-slate-900 hover:bg-slate-800 text-white text-xs font-semibold tracking-wide transition-all shadow-sm"
          >
            <Shield className="w-3.5 h-3.5" />
            <span>Portal Staf</span>
          </Link>
        </div>
      </div>
    </header>
  );
}
