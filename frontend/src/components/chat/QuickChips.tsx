interface QuickChipsProps {
  onSelect: (prompt: string) => void;
  disabled?: boolean;
}

export default function QuickChips({ onSelect, disabled }: QuickChipsProps) {
  const chips = [
    "Cek kamar kosong akhir pekan ini",
    "Rekomendasi kamar untuk honeymoon",
    "Tarif Deluxe Ocean View berapa?",
    "Cara pesan dan batas hold booking",
  ];

  return (
    <div className="flex gap-1.5 overflow-x-auto pb-2 scrollbar-none text-xs">
      {chips.map((chip, idx) => (
        <button
          key={idx}
          disabled={disabled}
          onClick={() => onSelect(chip)}
          className="whitespace-nowrap px-3 py-1.5 rounded-full bg-slate-100 hover:bg-brand-50 hover:text-brand-700 hover:border-brand-200 border border-slate-200 text-slate-600 transition-colors shrink-0 disabled:opacity-50"
        >
          {chip}
        </button>
      ))}
    </div>
  );
}
