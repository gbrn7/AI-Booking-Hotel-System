import { Waves, Utensils, Sparkles, Compass, ShieldCheck, HeartHandshake } from "lucide-react";

export default function Amenities() {
  const features = [
    {
      icon: Waves,
      title: "Infinity Beachfront Pool",
      description: "Berenang di kolam tanpa tepi berlatar panorama laut lepas dengan cabana privat eksklusif.",
    },
    {
      icon: Utensils,
      title: "Seaside Fine Dining",
      description: "Restoran bintang lima dengan sajian hidangan laut segar dan racikan kuliner fusion chef internasional.",
    },
    {
      icon: Sparkles,
      title: "Holistic Wellness Spa",
      description: "Perawatan tubuh dan terapi aromaterapi tradisional untuk meremajakan jiwa dan raga.",
    },
    {
      icon: Compass,
      title: "Island Excursions",
      description: "Tur perahu pesisir, snorkeling, dan kunjungan privat ke pantai tersembunyi sekitar teluk.",
    },
    {
      icon: ShieldCheck,
      title: "Higienitas & Privasi",
      description: "Standar protokol sanitasi tinggi dan jaminan keamanan resor 24 jam dengan staf berdedikasi.",
    },
    {
      icon: HeartHandshake,
      title: "Layanan Concierge Rama",
      description: "Asisten cerdas terintegrasi siap melayani kebutuhan jadwal, booking, hingga special request kapan saja.",
    },
  ];

  return (
    <section id="amenities" className="py-24 bg-white border-t border-slate-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-2xl mx-auto mb-16">
          <span className="text-xs font-semibold tracking-widest text-brand-600 uppercase mb-2 block">
            FASILITAS & PENGALAMAN RESOR
          </span>
          <h2 className="text-3xl sm:text-4xl font-serif font-bold text-slate-900 mb-4">
            Kemewahan Lengkap Untuk Setiap Momen
          </h2>
          <p className="text-slate-600 text-sm sm:text-base font-light">
            Dirancang secara teliti untuk memanjakan liburan romantis, liburan keluarga, maupun retret santai Anda.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((item, idx) => (
            <div
              key={idx}
              className="p-8 rounded-2xl bg-slate-50 border border-slate-100 hover:border-brand-300 hover:bg-brand-50/20 transition-all group"
            >
              <div className="w-12 h-12 rounded-xl bg-white border border-slate-200 flex items-center justify-center text-brand-600 mb-6 shadow-sm group-hover:scale-110 transition-transform">
                <item.icon className="w-6 h-6" />
              </div>
              <h3 className="text-lg font-serif font-bold text-slate-900 mb-2">
                {item.title}
              </h3>
              <p className="text-sm text-slate-500 leading-relaxed font-light">
                {item.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
