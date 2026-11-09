"use client";

import React, { createContext, useContext, useEffect, useState, useMemo } from "react";
import enMessages from "./locales/en.json";
import ruMessages from "./locales/ru.json";

export type Locale = "en" | "ru";

type TranslationDict = typeof enMessages;

const dictionaries: Record<Locale, TranslationDict> = {
  en: enMessages,
  ru: ruMessages,
};

export type PluralForms = { one: string; few: string; many: string };

interface I18nContextValue {
  locale: Locale;
  setLocale: (l: Locale) => void;
  t: (keyPath: string, variables?: Record<string, string | number>) => string;
  plural: (
    count: number,
    singularOrForms: PluralForms | string,
    few?: string,
    many?: string
  ) => string;
  formatDate: (date: Date | string, options?: Intl.DateTimeFormatOptions) => string;
  formatNumber: (n: number) => string;
}

const I18nContext = createContext<I18nContextValue | null>(null);

function resolveKey(obj: any, path: string): string | undefined {
  const parts = path.split(".");
  let cur = obj;
  for (const p of parts) {
    if (cur === undefined || cur === null) return undefined;
    cur = cur[p];
  }
  return typeof cur === "string" ? cur : undefined;
}

export function I18nProvider({ children }: { children: React.ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>("en");
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
    const saved = localStorage.getItem("atlasflow_locale") as Locale | null;
    if (saved === "en" || saved === "ru") {
      setLocaleState(saved);
      document.documentElement.lang = saved;
      return;
    }

    // Auto-detect browser language
    const browserLangs = navigator.languages || [navigator.language];
    let detected: Locale = "en";
    for (const l of browserLangs) {
      if (l.toLowerCase().startsWith("ru")) {
        detected = "ru";
        break;
      }
    }
    setLocaleState(detected);
    document.documentElement.lang = detected;
    localStorage.setItem("atlasflow_locale", detected);
  }, []);

  const setLocale = (newLocale: Locale) => {
    setLocaleState(newLocale);
    if (typeof window !== "undefined") {
      localStorage.setItem("atlasflow_locale", newLocale);
      document.documentElement.lang = newLocale;
    }
  };

  const pluralRules = useMemo(() => {
    try {
      return new Intl.PluralRules(locale);
    } catch {
      return new Intl.PluralRules("en");
    }
  }, [locale]);

  const t = (keyPath: string, variables?: Record<string, string | number>): string => {
    const dict = dictionaries[locale] || dictionaries.en;
    let str = resolveKey(dict, keyPath);
    if (str === undefined) {
      // Fallback to English
      str = resolveKey(dictionaries.en, keyPath) || keyPath;
    }

    if (variables) {
      for (const [k, v] of Object.entries(variables)) {
        str = str.replace(new RegExp(`\\{${k}\\}`, "g"), String(v));
      }
    }
    return str;
  };

  const plural = (
    count: number,
    singularOrForms: PluralForms | string,
    few?: string,
    many?: string
  ): string => {
    const rule = pluralRules.select(count);
    let oneStr = "";
    let fewStr = "";
    let manyStr = "";

    if (typeof singularOrForms === "object") {
      oneStr = singularOrForms.one;
      fewStr = singularOrForms.few;
      manyStr = singularOrForms.many;
    } else {
      oneStr = singularOrForms;
      fewStr = few || singularOrForms;
      manyStr = many || few || singularOrForms;
    }

    if (locale === "ru") {
      if (rule === "one") return `${count} ${oneStr}`;
      if (rule === "few") return `${count} ${fewStr}`;
      return `${count} ${manyStr}`;
    }
    // English rules
    if (rule === "one") return `${count} ${oneStr}`;
    return `${count} ${manyStr}`;
  };

  const formatDate = (date: Date | string, options?: Intl.DateTimeFormatOptions): string => {
    const d = typeof date === "string" ? new Date(date) : date;
    if (isNaN(d.getTime())) return "";
    try {
      return new Intl.DateTimeFormat(locale, options || { month: "short", day: "numeric", year: "numeric" }).format(d);
    } catch {
      return d.toLocaleDateString();
    }
  };

  const formatNumber = (n: number): string => {
    try {
      return new Intl.NumberFormat(locale).format(n);
    } catch {
      return String(n);
    }
  };

  return (
    <I18nContext.Provider value={{ locale, setLocale, t, plural, formatDate, formatNumber }}>
      {children}
    </I18nContext.Provider>
  );
}

export function useTranslation() {
  const ctx = useContext(I18nContext);
  if (!ctx) {
    // Fallback safe dummy translation if rendered outside provider
    return {
      locale: "en" as Locale,
      setLocale: () => {},
      t: (k: string) => k,
      plural: (
        c: number,
        s: PluralForms | string,
        _f?: string,
        m?: string
      ) => {
        const one = typeof s === "object" ? s.one : s;
        const many = typeof s === "object" ? s.many : m || s;
        return `${c} ${c === 1 ? one : many}`;
      },
      formatDate: (d: any) => String(d),
      formatNumber: (n: number) => String(n),
    };
  }
  return ctx;
}
