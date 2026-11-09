import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: ["class", '[data-theme="dark"]'],
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        map: {
          base: "var(--map-base)",
        },
        glass: {
          DEFAULT: "var(--glass)",
          border: "var(--glass-border)",
        },
        brand: {
          text: "var(--text)",
          muted: "var(--text-muted)",
          accent: "var(--accent)",
        },
        pin: {
          food: "var(--pin-food)",
          stay: "var(--pin-stay)",
          sight: "var(--pin-sight)",
          nature: "var(--pin-nature)",
          transit: "var(--pin-transit)",
          visited: "var(--pin-visited)",
        },
      },
      backdropBlur: {
        glass: "20px",
      },
      boxShadow: {
        glass: "0 8px 32px 0 rgba(0, 0, 0, 0.25)",
        "glass-light": "0 8px 32px 0 rgba(0, 0, 0, 0.08)",
      },
    },
  },
  plugins: [],
};

export default config;
