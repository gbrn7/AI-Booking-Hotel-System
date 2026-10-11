import Navbar from "@/components/landing/Navbar";
import HeroSection from "@/components/landing/HeroSection";
import RoomGrid from "@/components/landing/RoomGrid";
import Amenities from "@/components/landing/Amenities";
import Footer from "@/components/landing/Footer";

export default function HomePage() {
  return (
    <div className="min-h-screen flex flex-col bg-slate-50">
      <Navbar />
      <main className="flex-1">
        <HeroSection />
        <RoomGrid />
        <Amenities />
      </main>
      <Footer />
    </div>
  );
}
