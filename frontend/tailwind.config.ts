import type { Config } from "tailwindcss";

export default {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50: "#fdf8ee",
          100: "#faefd4",
          200: "#f4dea7",
          300: "#ebc673",
          400: "#e3ad44",
          500: "#d99324",
          600: "#bc731b",
          700: "#955219",
          800: "#79421b",
          900: "#65371a",
        },
        navy: {
          800: "#0f172a",
          900: "#090d16",
          950: "#05080e",
        },
      },
    },
  },
  plugins: [],
} satisfies Config;
