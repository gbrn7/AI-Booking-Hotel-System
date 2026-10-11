import Link from "next/link";
import { MapPin, Phone, Mail, Instagram, Shield } from "lucide-react";

export default function Footer() {
  return (
    <footer className="bg-slate-900 text-white pt-16 pb-12 border-t border-slate-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-10 mb-12">
          {/* Brand info */}
          <div className="md:col-span-1">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-9 h-9 rounded-full bg-brand-500 text-white flex items-center justify-center font-serif text-lg font-bold">
                G
              </div>
              <div>
                <span className="font-serif text-lg tracking-wider font-bold block text-white">
                  GRAND OASIS
                </span>
                <span className="text-[9px] tracking-widest text-brand-400 uppercase font-semibold">
                  Resort & Spa
                </span>
              </div>
            </div>
            <p className="text-xs text-slate-400 leading-relaxed font-light mb-6">
              Resor bintang lima di pesisir eksotis dengan pengalaman menginap tak tertandingi dan layanan AI Concierge terdepan.
            </p>
          </div>

          {/* Contact */}
          <div>
            <h4 className="text-xs font-semibold uppercase tracking-wider text-brand-400 mb-4">
              Kontak & Reservasi
            </h4>
            <ul className="space-y-3 text-xs text-slate-300">
              <li className="flex items-start gap-2.5">
                <MapPin className="w-4 h-4 text-brand-400 shrink-0 mt-0.5" />
                <span>Jl. Pantai Indah Nirwana No. 88, Badung, Bali 80361</span>
              </li>
              <li className="flex items-center gap-2.5">
                <Phone className="w-4 h-4 text-brand-400 shrink-0" />
                <span>+62 (361) 888-9999</span>
              </li>
              <li className="flex items-center gap-2.5">
                <Mail className="w-4 h-4 text-brand-400 shrink-0" />
                <span>reservation@grandoasis.com</span>
              </li>
            </ul>
          </div>

          {/* Quick links */}
          <div>
            <h4 className="text-xs font-semibold uppercase tracking-wider text-brand-400 mb-4">
              Eksplorasi
            </h4>
            <ul className="space-y-2 text-xs text-slate-400">
              <li><a href="#rooms" className="hover:text-white transition-colors">Tipe Kamar & Tarif</a></li>
              <li><a href="#amenities" className="hover:text-white transition-colors">Fasilitas Kolam & Restoran</a></li>
              <li><a href="#experience" className="hover:text-white transition-colors">Paket Honeymoon & Spa</a></li>
              <li><Link href="/admin" className="hover:text-brand-400 transition-colors flex items-center gap-1.5"><Shield className="w-3.5 h-3.5"/>Login Staf Hotel</Link></li>
            </ul>
          </div>

          {/* Policy & Banking */}
          <div>
            <h4 className="text-xs font-semibold uppercase tracking-wider text-brand-400 mb-4">
              Rekening Resmi Hotel
            </h4>
            <p className="text-xs text-slate-400 leading-relaxed font-light mb-2">
              Pembayaran hold reservasi resmi hotel ditujukan ke:
            </p>
            <div className="p-3 rounded-xl bg-slate-800/80 border border-slate-700 text-xs">
              <div className="font-semibold text-white">Bank Central Asia (BCA)</div>
              <div className="text-brand-300 font-mono">8870-1234-5678-9000</div>
              <div className="text-[11px] text-slate-400">a/n Grand Oasis Hotel PT</div>
            </div>
          </div>
        </div>

        <div className="pt-8 border-t border-slate-800 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-slate-500">
          <div>
            &copy; 2026 Grand Oasis Hotel. Hak Cipta Dilindungi Undang-Undang.
          </div>
          <div className="flex items-center gap-6">
            <span className="text-[11px]">Powered by Go GenAI Agent Engine & Next.js 15</span>
          </div>
        </div>
      </div>
    </footer>
  );
}
