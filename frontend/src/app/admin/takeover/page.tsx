import TakeoverChat from "@/components/admin/TakeoverChat";

export default function TakeoverPage() {
  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-slate-900">
          Live Chat & Receptionist Takeover
        </h1>
        <p className="text-xs text-slate-500 font-light">
          Pantau seluruh interaksi tamu (WhatsApp & Web Widget), ambil alih percakapan (Takeover Mode), dan kirim balasan staf secara langsung.
        </p>
      </div>
      <TakeoverChat />
    </div>
  );
}
