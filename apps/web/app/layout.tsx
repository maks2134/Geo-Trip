import type { Metadata } from "next";
import { QueryProvider } from "@/shared/lib/query-client";
import { I18nProvider } from "@/shared/i18n/i18n-context";
import "./globals.css";

export const metadata: Metadata = {
  title: "AtlasFlow — NomadGlass Travel Scratch Map & Planner",
  description: "Personal travel scratch-map and day-by-day trip planner with adaptive glass UI.",
};

const themeInitScript = `
(function() {
  try {
    var storedTheme = localStorage.getItem('atlasflow_theme');
    var isDark = true;
    if (storedTheme === 'light') {
      isDark = false;
    } else if (storedTheme === 'dark') {
      isDark = true;
    } else {
      isDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light');
    var storedLocale = localStorage.getItem('atlasflow_locale');
    if (storedLocale) {
      document.documentElement.setAttribute('lang', storedLocale);
    }
  } catch (e) {}
})();
`;

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" data-theme="dark" suppressHydrationWarning>
      <head>
        <script dangerouslySetInnerHTML={{ __html: themeInitScript }} />
      </head>
      <body className="w-screen h-screen overflow-hidden">
        <QueryProvider>
          <I18nProvider>{children}</I18nProvider>
        </QueryProvider>
      </body>
    </html>
  );
}
