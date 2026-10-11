import { Sparkles, User, ShieldAlert } from "lucide-react";
import { formatChatText } from "@/lib/formatChat";

interface ChatMessageProps {
  sender: "user" | "assistant" | "agent_tool" | "human_agent";
  content: string;
  timestamp?: string;
}

export default function ChatMessage({ sender, content, timestamp }: ChatMessageProps) {
  const isUser = sender === "user";
  const isStaff = sender === "human_agent";

  // Split by [NEXT] indicator (case-insensitive) to distinguish chat bubbles
  const bubbles = (content || "")
    .split(/\[next\]/i)
    .map((b) => b.trim())
    .filter(Boolean);

  const displayBubbles = bubbles.length > 0 ? bubbles : [content];

  return (
    <div className={`flex gap-3 mb-3.5 ${isUser ? "flex-row-reverse" : "flex-row"}`}>
      {/* Avatar */}
      <div
        className={`w-8 h-8 rounded-full flex items-center justify-center shrink-0 text-xs font-bold shadow-sm ${
          isUser
            ? "bg-slate-800 text-white"
            : isStaff
            ? "bg-emerald-600 text-white"
            : "bg-brand-500 text-white"
        }`}
      >
        {isUser ? <User className="w-4 h-4" /> : isStaff ? <ShieldAlert className="w-4 h-4" /> : <Sparkles className="w-4 h-4" />}
      </div>

      {/* Bubble Container */}
      <div className={`flex flex-col gap-1.5 max-w-[82%] ${isUser ? "items-end" : "items-start"}`}>
        {displayBubbles.map((bubbleText, idx) => {
          const isFirstBubble = idx === 0;
          const isLastBubble = idx === displayBubbles.length - 1;
          const hasBookingCode = bubbleText.includes("HTL-");

          return (
            <div
              key={idx}
              className={`rounded-2xl px-4 py-3 text-sm leading-relaxed shadow-sm ${
                isUser
                  ? "bg-slate-900 text-white rounded-tr-none"
                  : isStaff
                  ? "bg-emerald-50 border border-emerald-200 text-emerald-950 rounded-tl-none"
                  : "bg-white border border-slate-200 text-slate-800 rounded-tl-none"
              }`}
            >
              {!isUser && isFirstBubble && (
                <div className="flex items-center gap-1.5 mb-1">
                  <span className="text-[11px] font-bold tracking-wide">
                    {isStaff ? "Front Desk Staff (Manusia)" : "Rama — Concierge AI"}
                  </span>
                  {isStaff && (
                    <span className="text-[9px] px-1.5 py-0.5 rounded bg-emerald-200 text-emerald-800 font-semibold uppercase">
                      Takeover
                    </span>
                  )}
                </div>
              )}

              <div className="whitespace-pre-wrap">{formatChatText(bubbleText)}</div>

              {hasBookingCode && (
                <div className="mt-2.5 p-2.5 rounded-xl bg-amber-50/80 border border-amber-200/80 text-[11px] text-amber-900">
                  <div className="font-semibold mb-0.5">💡 Info Status Reservasi:</div>
                  <div>Batas waktu hold 2 jam. Bukti pembayaran dapat dikonfirmasi via staf resepsionis hotel.</div>
                </div>
              )}

              {timestamp && isLastBubble && (
                <div
                  className={`text-[10px] mt-1.5 text-right ${
                    isUser ? "text-slate-400" : "text-slate-400"
                  }`}
                >
                  {timestamp}
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
