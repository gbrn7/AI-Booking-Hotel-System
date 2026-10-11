import localFont from "next/font/local";

export const plusJakartaSans = localFont({
  src: [
    {
      path: "../assets/font/Plus_Jakarta_Sans/PlusJakartaSans-VariableFont_wght.ttf",
      style: "normal",
    },
    {
      path: "../assets/font/Plus_Jakarta_Sans/PlusJakartaSans-Italic-VariableFont_wght.ttf",
      style: "italic",
    },
  ],
  variable: "--font-plus-jakarta-sans",
  display: "swap",
});
