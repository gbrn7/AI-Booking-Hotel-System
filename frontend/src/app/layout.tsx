import type { Metadata } from "next";
import "./globals.css";
import LiveChatWidget from "@/components/chat/LiveChatWidget";

export const metadata: Metadata = {
  title: "Grand Oasis Hotel — Luxury Beachfront Resort & Spa",
  description: "Rasakan kemewahan menginap di Grand Oasis Hotel dengan pemandangan pesisir eksotis dan asisten reservasi pintar 24/7.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="id">
      <body className="min-h-screen bg-slate-50 text-slate-900 antialiased relative">
        {children}
        <LiveChatWidget />
      </body>
    </html>
  );
}
