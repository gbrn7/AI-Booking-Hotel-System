"use client";

import { useState, useEffect, useRef } from "react";
import { MessageSquare, X, Send, Sparkles, RefreshCw } from "lucide-react";
import ChatMessage from "./ChatMessage";
import QuickChips from "./QuickChips";

interface MessageItem {
  id: string;
  sender: "user" | "assistant" | "human_agent";
  content: string;
  timestamp: string;
}

const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080";

export default function LiveChatWidget() {
  const [isOpen, setIsOpen] = useState(false);
  const [sessionId, setSessionId] = useState<string>("");
  const [inputMessage, setInputMessage] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const [messages, setMessages] = useState<MessageItem[]>([
    {
      id: "initial-1",
      sender: "assistant",
      content: "Halo Kak! Selamat datang di Grand Oasis Hotel 🌊",
      timestamp: "Baru saja",
    },
    {
      id: "initial-2",
      sender: "assistant",
      content: "Saya Rama, Concierge AI Anda. Ada tanggal atau tipe kamar yang ingin Anda tanyakan?",
      timestamp: "Baru saja",
    },
  ]);

  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Initialize or restore session ID
  useEffect(() => {
    let sid = localStorage.getItem("grand_oasis_chat_sid");
    if (!sid) {
      sid = "web:" + Math.random().toString(36).substring(2, 10);
      localStorage.setItem("grand_oasis_chat_sid", sid);
    }
    setSessionId(sid);

    // Fetch past history if exists
    fetch(`${BACKEND_URL}/api/chat/history?session_id=${encodeURIComponent(sid)}`)
      .then((res) => res.json())
      .then((data) => {
        if (data.messages && data.messages.length > 0) {
          const loaded: MessageItem[] = data.messages.map((m: any) => ({
            id: String(m.id),
            sender: m.sender === "user" ? "user" : m.sender === "human_agent" ? "human_agent" : "assistant",
            content: m.content,
            timestamp: new Date(m.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
          }));
          setMessages(loaded);
        }
      })
      .catch(() => {});

    // Listen for custom open-ai-chat event from room cards
    const handleOpenChat = (e: Event) => {
      const customEvent = e as CustomEvent<{ query: string }>;
      setIsOpen(true);
      if (customEvent.detail?.query) {
        sendMessage(customEvent.detail.query);
      }
    };

    window.addEventListener("open-ai-chat", handleOpenChat);
    return () => window.removeEventListener("open-ai-chat", handleOpenChat);
  }, []);

  useEffect(() => {
    if (isOpen) {
      messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages, isOpen, isTyping]);

  const sendMessage = async (textToSend: string) => {
    const text = textToSend.trim();
    if (!text || isTyping) return;

    setInputMessage("");
    const nowTime = new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });

    // Append user message
    const userMsg: MessageItem = {
      id: "user-" + Date.now(),
      sender: "user",
      content: text,
      timestamp: nowTime,
    };
    setMessages((prev) => [...prev, userMsg]);
    setIsTyping(true);

    try {
      const response = await fetch(`${BACKEND_URL}/api/chat/stream`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          session_id: sessionId || "web:guest",
          message: text,
        }),
      });

      if (!response.ok || !response.body) {
        throw new Error("Gagal terhubung ke concierge");
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder("utf-8");
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n\n");
        buffer = lines.pop() || "";

        for (const line of lines) {
          if (line.startsWith("data: ")) {
            const jsonStr = line.replace("data: ", "").trim();
            if (!jsonStr) continue;
            try {
              const parsed = JSON.parse(jsonStr);
              if (parsed.type === "bubble" || parsed.type === "takeover") {
                setMessages((prev) => [
                  ...prev,
                  {
                    id: "asst-" + Date.now() + "-" + Math.random(),
                    sender: parsed.type === "takeover" ? "human_agent" : "assistant",
                    content: parsed.content,
                    timestamp: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
                  },
                ]);
              }
            } catch {
              // Ignore malformed chunks
            }
          }
        }
      }
    } catch {
      setMessages((prev) => [
        ...prev,
        {
          id: "err-" + Date.now(),
          sender: "assistant",
          content: "Mohon maaf, sistem sedang memproses antrean. Silakan coba kembali sesaat lagi.",
          timestamp: nowTime,
        },
      ]);
    } finally {
      setIsTyping(false);
    }
  };

  const handleResetSession = () => {
    const newSid = "web:" + Math.random().toString(36).substring(2, 10);
    localStorage.setItem("grand_oasis_chat_sid", newSid);
    setSessionId(newSid);
    setMessages([
      {
        id: "init-reset",
        sender: "assistant",
        content: "Halo! Sesi obrolan baru dimulai. Ada yang bisa Rama bantu seputar kamar Grand Oasis Hotel?",
        timestamp: "Baru saja",
      },
    ]);
  };

  return (
    <>
      {/* Floating Trigger Button */}
      {!isOpen && (
        <button
          onClick={() => setIsOpen(true)}
          className="fixed bottom-6 right-6 z-50 flex items-center gap-2.5 px-5 py-3.5 rounded-full bg-slate-900 text-white shadow-2xl hover:bg-brand-600 active:scale-95 transition-all group border border-slate-700/50"
        >
          <div className="relative">
            <Sparkles className="w-5 h-5 text-brand-400 group-hover:rotate-12 transition-transform" />
            <span className="absolute -top-1 -right-1 w-2.5 h-2.5 bg-emerald-500 rounded-full border-2 border-slate-900 animate-pulse" />
          </div>
          <span className="text-sm font-semibold tracking-wide">Tanya Concierge AI</span>
        </button>
      )}

      {/* Chat Window Modal */}
      {isOpen && (
        <div className="fixed bottom-4 right-4 sm:bottom-6 sm:right-6 z-50 w-[calc(100vw-2rem)] sm:w-[420px] h-[580px] max-h-[85vh] bg-white rounded-3xl shadow-2xl border border-slate-200/80 flex flex-col overflow-hidden animate-in fade-in slide-in-from-bottom-6 duration-300">
          {/* Header */}
          <div className="p-4 bg-slate-900 text-white flex items-center justify-between shadow-sm">
            <div className="flex items-center gap-3">
              <div className="relative">
                <div className="w-10 h-10 rounded-full bg-brand-500 text-white flex items-center justify-center font-bold text-sm shadow-md">
                  R
                </div>
                <span className="absolute bottom-0 right-0 w-3 h-3 bg-emerald-500 rounded-full border-2 border-slate-900" />
              </div>
              <div>
                <div className="flex items-center gap-1.5 font-bold text-sm">
                  <span>Rama Concierge</span>
                  <span className="text-[10px] px-1.5 py-0.5 rounded bg-brand-500/20 text-brand-300 font-semibold border border-brand-500/30">
                    AI Online
                  </span>
                </div>
                <div className="text-[11px] text-slate-400 font-light">
                  Grand Oasis Hotel Guest Assistant
                </div>
              </div>
            </div>

            <div className="flex items-center gap-1">
              <button
                onClick={handleResetSession}
                title="Mulai obrolan baru"
                className="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
              >
                <RefreshCw className="w-4 h-4" />
              </button>
              <button
                onClick={() => setIsOpen(false)}
                className="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
          </div>

          {/* Messages body */}
          <div className="flex-1 p-4 overflow-y-auto bg-slate-50/70">
            {messages.map((msg) => (
              <ChatMessage
                key={msg.id}
                sender={msg.sender}
                content={msg.content}
                timestamp={msg.timestamp}
              />
            ))}

            {isTyping && (
              <div className="flex gap-2 items-center text-xs text-slate-400 mb-3 px-2">
                <div className="w-2 h-2 rounded-full bg-brand-500 animate-bounce" />
                <div className="w-2 h-2 rounded-full bg-brand-500 animate-bounce [animation-delay:0.2s]" />
                <div className="w-2 h-2 rounded-full bg-brand-500 animate-bounce [animation-delay:0.4s]" />
                <span className="ml-1 text-[11px]">Rama sedang mengetik balasan...</span>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>

          {/* Bottom input section */}
          <div className="p-3 bg-white border-t border-slate-100">
            <QuickChips onSelect={(prompt) => sendMessage(prompt)} disabled={isTyping} />

            <form
              onSubmit={(e) => {
                e.preventDefault();
                sendMessage(inputMessage);
              }}
              className="flex items-center gap-2 mt-2"
            >
              <input
                type="text"
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                placeholder="Ketik pertanyaan kamar atau booking..."
                disabled={isTyping}
                className="flex-1 px-4 py-2.5 rounded-xl border border-slate-200 text-xs sm:text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 disabled:bg-slate-50"
              />
              <button
                type="submit"
                disabled={!inputMessage.trim() || isTyping}
                className="w-10 h-10 rounded-xl bg-slate-900 hover:bg-brand-600 disabled:opacity-40 text-white flex items-center justify-center transition-all shrink-0"
              >
                <Send className="w-4 h-4" />
              </button>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
