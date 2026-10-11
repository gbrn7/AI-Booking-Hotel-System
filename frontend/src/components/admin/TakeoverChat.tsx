"use client";

import { useState, useEffect } from "react";
import { Conversation, Message } from "@/types";
import {
  fetchAdminConversations,
  fetchConversationMessages,
  toggleTakeover,
  sendStaffReply,
} from "@/lib/api";
import { formatChatText } from "@/lib/formatChat";
import {
  MessageSquare,
  Send,
  ToggleLeft,
  ToggleRight,
  RefreshCw,
  Phone,
  Globe,
  User,
  Sparkles,
  ShieldAlert,
} from "lucide-react";

export default function TakeoverChat() {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [replyText, setReplyText] = useState("");
  const [loading, setLoading] = useState(false);
  const [sending, setSending] = useState(false);

  const loadConversations = async () => {
    try {
      const data = await fetchAdminConversations("");
      setConversations(data);
      if (!selectedId && data.length > 0) {
        setSelectedId(data[0].id);
      }
    } catch {
      // Fallback empty
    }
  };

  const loadMessages = async (id: string) => {
    setLoading(true);
    try {
      const msgs = await fetchConversationMessages(id, "");
      setMessages(msgs);
    } catch {
      setMessages([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadConversations();
    const interval = setInterval(loadConversations, 5000);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    if (selectedId) {
      loadMessages(selectedId);
    }
  }, [selectedId]);

  const selectedConv = conversations.find((c) => c.id === selectedId);

  const handleToggleTakeover = async () => {
    if (!selectedConv) return;
    const nextState = !selectedConv.is_human_takeover;
    try {
      await toggleTakeover(selectedConv.id, nextState, "");
      setConversations((prev) =>
        prev.map((c) => (c.id === selectedConv.id ? { ...c, is_human_takeover: nextState } : c))
      );
    } catch (err: any) {
      alert("Gagal mengubah takeover: " + err.message);
    }
  };

  const handleSendReply = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedId || !replyText.trim() || sending) return;

    setSending(true);
    const text = replyText.trim();
    setReplyText("");

    try {
      await sendStaffReply(selectedId, text, "");
      setMessages((prev) => [
        ...prev,
        {
          id: Date.now(),
          conversation_id: selectedId,
          sender: "human_agent",
          content: text,
          created_at: new Date().toISOString(),
        },
      ]);
    } catch (err: any) {
      alert("Gagal mengirim pesan: " + err.message);
    } finally {
      setSending(false);
    }
  };

  return (
    <div className="bg-white rounded-2xl shadow-sm border border-slate-200 h-[calc(100vh-10rem)] flex overflow-hidden">
      {/* Left Column: Conversations List */}
      <div className="w-80 border-r border-slate-200 flex flex-col bg-slate-50/50">
        <div className="p-4 border-b border-slate-200 flex items-center justify-between">
          <div className="flex items-center gap-2 font-bold text-sm text-slate-800">
            <MessageSquare className="w-4 h-4 text-brand-600" />
            <span>Chat Masuk ({conversations.length})</span>
          </div>
          <button
            onClick={loadConversations}
            title="Refresh"
            className="p-1 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-200 transition-colors"
          >
            <RefreshCw className="w-3.5 h-3.5" />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto divide-y divide-slate-100">
          {conversations.length === 0 ? (
            <div className="p-8 text-center text-xs text-slate-400">
              Belum ada percakapan aktif
            </div>
          ) : (
            conversations.map((conv) => {
              const isSelected = conv.id === selectedId;
              const isWA = conv.channel === "whatsapp" || conv.id.startsWith("wa:");
              return (
                <button
                  key={conv.id}
                  onClick={() => setSelectedId(conv.id)}
                  className={`w-full text-left p-3.5 transition-colors flex items-start gap-3 ${
                    isSelected ? "bg-white border-l-4 border-l-brand-500 shadow-sm" : "hover:bg-slate-100/70"
                  }`}
                >
                  <div className="w-9 h-9 rounded-full bg-slate-200 flex items-center justify-center shrink-0">
                    {isWA ? (
                      <Phone className="w-4 h-4 text-emerald-600" />
                    ) : (
                      <Globe className="w-4 h-4 text-blue-600" />
                    )}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-xs font-bold text-slate-900 truncate">
                        {conv.guest?.name || conv.id}
                      </span>
                      <span className="text-[10px] text-slate-400">
                        {new Date(conv.last_message_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                      </span>
                    </div>

                    <div className="flex items-center gap-1.5">
                      <span
                        className={`text-[9px] px-1.5 py-0.5 rounded font-semibold uppercase ${
                          isWA ? "bg-emerald-100 text-emerald-800" : "bg-blue-100 text-blue-800"
                        }`}
                      >
                        {isWA ? "WhatsApp" : "Web Chat"}
                      </span>
                      {conv.is_human_takeover && (
                        <span className="text-[9px] px-1.5 py-0.5 rounded bg-rose-100 text-rose-800 font-semibold uppercase">
                          Takeover
                        </span>
                      )}
                    </div>
                  </div>
                </button>
              );
            })
          )}
        </div>
      </div>

      {/* Right Column: Active Conversation */}
      <div className="flex-1 flex flex-col bg-white">
        {selectedConv ? (
          <>
            {/* Header with Takeover Toggle */}
            <div className="p-4 border-b border-slate-200 flex items-center justify-between bg-slate-50/70">
              <div>
                <h3 className="text-sm font-bold text-slate-900 flex items-center gap-2">
                  <span>{selectedConv.guest?.name || selectedConv.id}</span>
                  <span className="text-xs font-normal text-slate-500 font-mono">
                    ({selectedConv.id})
                  </span>
                </h3>
                <p className="text-[11px] text-slate-500">
                  Channel: {selectedConv.channel} • Terakhir aktif:{" "}
                  {new Date(selectedConv.last_message_at).toLocaleTimeString()}
                </p>
              </div>

              <div className="flex items-center gap-3">
                <div className="text-right">
                  <div className="text-xs font-semibold text-slate-800">
                    {selectedConv.is_human_takeover ? "Mode Resepsionis (Takeover)" : "Mode AI Aktif (Rama)"}
                  </div>
                  <div className="text-[10px] text-slate-400">
                    {selectedConv.is_human_takeover ? "AI tidak akan membalas" : "AI membalas otomatis"}
                  </div>
                </div>

                <button
                  onClick={handleToggleTakeover}
                  className={`flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-semibold transition-all ${
                    selectedConv.is_human_takeover
                      ? "bg-rose-500 hover:bg-rose-600 text-white shadow-sm"
                      : "bg-slate-200 hover:bg-slate-300 text-slate-700"
                  }`}
                >
                  {selectedConv.is_human_takeover ? (
                    <>
                      <ToggleRight className="w-4 h-4" />
                      <span>Takeover ON</span>
                    </>
                  ) : (
                    <>
                      <ToggleLeft className="w-4 h-4" />
                      <span>Ambil Alih Chat</span>
                    </>
                  )}
                </button>
              </div>
            </div>

            {/* Messages Body */}
            <div className="flex-1 p-6 overflow-y-auto space-y-4 bg-slate-50/40">
              {loading ? (
                <div className="p-8 text-center text-xs text-slate-400">
                  Memuat riwayat percakapan...
                </div>
              ) : messages.length === 0 ? (
                <div className="p-8 text-center text-xs text-slate-400">
                  Belum ada pesan tercatat di sesi ini.
                </div>
              ) : (
                messages.map((m) => {
                  const isUser = m.sender === "user";
                  const isStaff = m.sender === "human_agent";
                  const bubbles = (m.content || "")
                    .split(/\[next\]/i)
                    .map((b) => b.trim())
                    .filter(Boolean);
                  const displayBubbles = bubbles.length > 0 ? bubbles : [m.content];

                  return (
                    <div
                      key={m.id}
                      className={`flex gap-3 ${isUser ? "justify-start" : "justify-end"}`}
                    >
                      <div className={`flex flex-col gap-1.5 max-w-[70%] ${isUser ? "items-start" : "items-end"}`}>
                        {displayBubbles.map((bubbleText, bIdx) => (
                          <div
                            key={bIdx}
                            className={`rounded-2xl px-4 py-3 text-xs leading-relaxed shadow-sm ${
                              isUser
                                ? "bg-white border border-slate-200 text-slate-800"
                                : isStaff
                                ? "bg-emerald-600 text-white"
                                : "bg-slate-900 text-white"
                            }`}
                          >
                            {bIdx === 0 && (
                              <div className="flex items-center gap-1.5 font-bold mb-1 opacity-80 text-[10px]">
                                {isUser ? (
                                  <>
                                    <User className="w-3 h-3" />
                                    <span>Tamu</span>
                                  </>
                                ) : isStaff ? (
                                  <>
                                    <ShieldAlert className="w-3 h-3" />
                                    <span>Staf Resepsionis (Anda)</span>
                                  </>
                                ) : (
                                  <>
                                    <Sparkles className="w-3 h-3 text-brand-400" />
                                    <span>Rama (Concierge AI)</span>
                                  </>
                                )}
                              </div>
                            )}
                            <div className="whitespace-pre-wrap">{formatChatText(bubbleText)}</div>
                          </div>
                        ))}
                      </div>
                    </div>
                  );
                })
              )}
            </div>

            {/* Bottom Manual Reply Box */}
            <div className="p-4 bg-white border-t border-slate-200">
              <form onSubmit={handleSendReply} className="flex gap-2">
                <input
                  type="text"
                  value={replyText}
                  onChange={(e) => setReplyText(e.target.value)}
                  placeholder={
                    selectedConv.is_human_takeover
                      ? "Ketik pesan balasan manual staf di sini..."
                      : "Aktifkan Takeover terlebih dahulu untuk membalas tanpa intervensi AI..."
                  }
                  className="flex-1 px-4 py-2.5 rounded-xl border border-slate-200 text-xs focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
                <button
                  type="submit"
                  disabled={!replyText.trim() || sending}
                  className="px-4 py-2.5 rounded-xl bg-slate-900 hover:bg-brand-600 disabled:opacity-40 text-white text-xs font-semibold flex items-center gap-2 transition-all"
                >
                  <Send className="w-3.5 h-3.5" />
                  <span>Kirim</span>
                </button>
              </form>
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-slate-400 text-sm">
            Pilih percakapan di kolom kiri untuk melihat riwayat atau mengambil alih chat.
          </div>
        )}
      </div>
    </div>
  );
}
