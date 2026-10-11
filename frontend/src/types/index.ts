export interface Room {
  id: number;
  name: string;
  type: string;
  base_price_per_night: number;
  capacity_guests: number;
  total_units: number;
  description: string;
  amenities: string; // JSON array string
  image_urls: string; // JSON array string
  is_active: boolean;
}

export interface RoomAvailability {
  room_id: number;
  name: string;
  type: string;
  base_price_per_night: number;
  capacity_guests: number;
  total_units: number;
  available_units: number;
  total_nights: number;
  estimated_total_price: number;
  description: string;
  amenities: string;
  image_urls: string;
}

export interface Guest {
  id: number;
  name: string;
  phone_number: string;
  email?: string;
}

export type BookingStatus = "HOLD" | "CONFIRMED" | "CANCELLED" | "CHECKED_IN" | "CHECKED_OUT";

export interface Booking {
  id: number;
  booking_code: string;
  guest_id: number;
  guest?: Guest;
  room_id: number;
  room?: Room;
  check_in_date: string;
  check_out_date: string;
  total_nights: number;
  total_rooms: number;
  total_price: number;
  status: BookingStatus;
  hold_expires_at?: string;
  payment_proof_url?: string;
  special_requests?: string;
  created_at: string;
}

export interface Conversation {
  id: string;
  guest_id?: number;
  guest?: Guest;
  channel: "whatsapp" | "web";
  is_human_takeover: boolean;
  last_message_at: string;
  created_at: string;
}

export interface Message {
  id: number;
  conversation_id: string;
  sender: "user" | "assistant" | "agent_tool" | "human_agent";
  content: string;
  thought_process?: string;
  created_at: string;
}

export interface AdminStats {
  total_rooms: number;
  pending_hold_count: number;
  confirmed_count: number;
  active_conversations: number;
}
