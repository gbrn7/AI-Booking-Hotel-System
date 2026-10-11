import { Room, RoomAvailability, Booking, Conversation, Message, AdminStats } from "@/types";

const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080";

export async function fetchRooms(): Promise<Room[]> {
  const res = await fetch(`${BACKEND_URL}/api/rooms`, { cache: "no-store" });
  if (!res.ok) throw new Error("Gagal mengambil data kamar");
  const data = await res.json();
  return data.rooms || [];
}

export async function checkRoomAvailability(checkIn: string, checkOut: string, roomType?: string): Promise<RoomAvailability[]> {
  const params = new URLSearchParams({ check_in: checkIn, check_out: checkOut });
  if (roomType) params.append("room_type", roomType);

  const res = await fetch(`${BACKEND_URL}/api/availability?${params.toString()}`, { cache: "no-store" });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Gagal memeriksa ketersediaan");
  }
  const data = await res.json();
  return data.availability || [];
}

export async function fetchBookingByCode(code: string): Promise<Booking> {
  const res = await fetch(`${BACKEND_URL}/api/bookings/${code}`, { cache: "no-store" });
  if (!res.ok) throw new Error("Booking tidak ditemukan");
  const data = await res.json();
  return data.booking;
}

export async function fetchAdminStats(token: string): Promise<AdminStats> {
  const res = await fetch(`${BACKEND_URL}/api/admin/stats`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });
  if (!res.ok) throw new Error("Gagal memuat statistik");
  return res.json();
}

export async function fetchAdminConversations(token: string): Promise<Conversation[]> {
  const res = await fetch(`${BACKEND_URL}/api/admin/conversations`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });
  if (!res.ok) throw new Error("Gagal memuat percakapan");
  const data = await res.json();
  return data.conversations || [];
}

export async function fetchConversationMessages(id: string, token: string): Promise<Message[]> {
  const res = await fetch(`${BACKEND_URL}/api/admin/conversations/${encodeURIComponent(id)}/messages`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });
  if (!res.ok) throw new Error("Gagal memuat pesan");
  const data = await res.json();
  return data.messages || [];
}

export async function toggleTakeover(id: string, isTakeover: boolean, token: string): Promise<void> {
  const res = await fetch(`${BACKEND_URL}/api/admin/conversations/${encodeURIComponent(id)}/takeover`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ is_takeover: isTakeover }),
  });
  if (!res.ok) throw new Error("Gagal mengubah status takeover");
}

export async function sendStaffReply(id: string, message: string, token: string): Promise<void> {
  const res = await fetch(`${BACKEND_URL}/api/admin/conversations/${encodeURIComponent(id)}/reply`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ message }),
  });
  if (!res.ok) throw new Error("Gagal mengirim balasan staf");
}

export async function fetchAdminBookings(status?: string, token?: string): Promise<Booking[]> {
  const url = status
    ? `${BACKEND_URL}/api/admin/bookings?status=${status}`
    : `${BACKEND_URL}/api/admin/bookings`;

  const headers: Record<string, string> = {};
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(url, { headers, cache: "no-store" });
  if (!res.ok) throw new Error("Gagal memuat daftar reservasi");
  const data = await res.json();
  return data.bookings || [];
}

export async function confirmBookingPayment(id: number, token?: string): Promise<Booking> {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${BACKEND_URL}/api/admin/bookings/${id}/confirm`, {
    method: "POST",
    headers,
  });
  if (!res.ok) throw new Error("Gagal mengonfirmasi booking");
  const data = await res.json();
  return data.booking;
}
