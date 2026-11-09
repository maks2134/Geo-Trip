"use client";

import React, { useState, useEffect } from "react";
import { Plus, PieChart, RefreshCw, SlidersHorizontal, ArrowRightLeft } from "lucide-react";
import { Trip, Expense } from "@/entities/trip";
import { ProgressBar } from "@/shared/ui/ProgressBar";
import { formatCurrency } from "@/shared/lib/utils";
import { apiClient, FXRatesResponse } from "@/shared/lib/api";
import { useTranslation } from "@/shared/i18n/i18n-context";

interface BudgetTabProps {
  trip: Trip;
  onTripUpdate: (updated: Trip) => void;
}

const POPULAR_CURRENCIES = [
  "EUR",
  "USD",
  "JPY",
  "THB",
  "GBP",
  "RUB",
  "TRY",
  "GEL",
  "AED",
  "CHF",
  "CNY",
  "CAD",
  "AUD",
  "SGD",
  "IDR",
  "KRW",
];

export function BudgetTab({ trip, onTripUpdate }: BudgetTabProps) {
  const { t } = useTranslation();
  const [isAdding, setIsAdding] = useState(false);
  const [title, setTitle] = useState("");
  const [amount, setAmount] = useState("");
  const [selectedCurrency, setSelectedCurrency] = useState(trip.currency || "EUR");
  const [category, setCategory] = useState<"stay" | "food" | "transport" | "activity" | "other">("food");
  const [fxRates, setFxRates] = useState<FXRatesResponse | null>(null);
  const [isCustomRate, setIsCustomRate] = useState(false);
  const [customRateInput, setCustomRateInput] = useState("");

  // Fetch FX rates based on trip base currency
  useEffect(() => {
    let isMounted = true;
    apiClient.getFXRates(trip.currency).then((res) => {
      if (isMounted) setFxRates(res);
    });
    return () => {
      isMounted = false;
    };
  }, [trip.currency]);

  // Derived effective exchange rate: how many Foreign Currency units per 1 Base Currency unit
  const foreignPerBase = fxRates?.rates?.[selectedCurrency] || 1.0;
  const effectiveForeignPerBase = isCustomRate && parseFloat(customRateInput) > 0
    ? parseFloat(customRateInput)
    : foreignPerBase;

  // Rate multiplier to convert Foreign Amount -> Base Amount: Base = Foreign * (1 / foreignPerBase)
  const baseFxRate = effectiveForeignPerBase > 0 ? 1 / effectiveForeignPerBase : 1.0;

  const rawAmount = parseFloat(amount) || 0;
  const computedBaseAmount = selectedCurrency === trip.currency
    ? rawAmount
    : Math.round(rawAmount * baseFxRate * 100) / 100;

  // Aggregate totals using base amounts
  const totalSpent = trip.expenses.reduce(
    (acc, exp) => acc + (exp.baseAmount ?? exp.amount),
    0
  );
  const percentUsed = Math.min(100, Math.round((totalSpent / trip.budgetLimit) * 100));

  const categoryTotals: Record<string, number> = {};
  trip.expenses.forEach((e) => {
    const val = e.baseAmount ?? e.amount;
    categoryTotals[e.category] = (categoryTotals[e.category] || 0) + val;
  });

  const handleAddExpense = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title || !amount) return;

    await apiClient.addExpense(trip.id, {
      tripId: trip.id,
      title,
      amount: computedBaseAmount,
      currency: trip.currency,
      category,
      paidAt: new Date().toISOString().split("T")[0],
      originalCurrency: selectedCurrency,
      originalAmount: rawAmount,
      fxRate: baseFxRate,
      baseAmount: computedBaseAmount,
      isCustomRate,
    });

    const refreshed = await apiClient.getTrip(trip.id);
    if (refreshed) onTripUpdate(refreshed);

    setTitle("");
    setAmount("");
    setIsAdding(false);
    setIsCustomRate(false);
    setCustomRateInput("");
  };

  const currencyOptions = Array.from(new Set([trip.currency, ...POPULAR_CURRENCIES]));

  return (
    <div className="space-y-4 text-xs">
      {/* Total Budget Card */}
      <div className="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-3">
        <div className="flex items-center justify-between">
          <div>
            <div className="text-[11px] text-[var(--text-muted)] uppercase tracking-wider">
              {t("budget.totalBudget")}
            </div>
            <div className="text-xl font-black text-[var(--text)] mt-0.5">
              {formatCurrency(totalSpent, trip.currency)}{" "}
              <span className="text-xs font-normal text-[var(--text-muted)]">
                / {formatCurrency(trip.budgetLimit, trip.currency)}
              </span>
            </div>
          </div>
          <div className="text-right">
            <span
              className={`font-extrabold text-sm ${
                percentUsed > 90
                  ? "text-red-400"
                  : percentUsed > 70
                  ? "text-amber-400"
                  : "text-emerald-400"
              }`}
            >
              {percentUsed}%
            </span>
            <div className="text-[10px] text-[var(--text-muted)]">{t("budget.allocated")}</div>
          </div>
        </div>

        <ProgressBar
          value={percentUsed}
          barColor={
            percentUsed > 90
              ? "bg-red-500"
              : percentUsed > 70
              ? "bg-amber-500"
              : "bg-emerald-500"
          }
        />
      </div>

      {/* Category Breakdown */}
      <div className="p-3 rounded-2xl bg-white/[0.03] border border-white/5 space-y-2">
        <div className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)] flex items-center gap-1.5">
          <PieChart className="w-3.5 h-3.5" />
          <span>{t("budget.categoryBreakdown")}</span>
        </div>
        <div className="grid grid-cols-2 gap-2 pt-1">
          {Object.entries(categoryTotals).map(([cat, val]) => (
            <div
              key={cat}
              className="p-2 rounded-xl bg-white/[0.02] flex items-center justify-between"
            >
              <span className="capitalize text-[var(--text-muted)]">{cat}</span>
              <span className="font-semibold text-[var(--text)]">
                {formatCurrency(val, trip.currency)}
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* Section Header & Add Expense Toggle */}
      <div className="flex items-center justify-between pt-1">
        <span className="font-semibold text-xs text-[var(--text)]">{t("budget.recentExpenses")}</span>
        <button
          onClick={() => setIsAdding(!isAdding)}
          className="flex items-center gap-1 px-2.5 py-1 rounded-lg bg-[var(--accent)] text-white font-medium hover:opacity-90 transition-opacity"
        >
          <Plus className="w-3.5 h-3.5" />
          <span>{t("budget.add")}</span>
        </button>
      </div>

      {/* Multi-Currency Expense Form (FEAT-03) */}
      {isAdding && (
        <form
          onSubmit={handleAddExpense}
          className="p-3.5 rounded-2xl bg-white/[0.05] border border-white/10 space-y-3"
        >
          <input
            type="text"
            placeholder="Expense title (e.g. Shinkansen ticket, Ramen)"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full p-2.5 rounded-xl bg-black/20 border border-white/10 text-[var(--text)] placeholder:text-[var(--text-muted)] font-medium focus:outline-none focus:ring-1 focus:ring-[var(--accent)]"
            required
          />

          <div className="grid grid-cols-3 gap-2">
            <input
              type="number"
              step="any"
              placeholder="Amount"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="col-span-2 p-2.5 rounded-xl bg-black/20 border border-white/10 text-[var(--text)] font-semibold text-sm focus:outline-none focus:ring-1 focus:ring-[var(--accent)]"
              required
            />
            <select
              value={selectedCurrency}
              onChange={(e) => setSelectedCurrency(e.target.value)}
              className="p-2.5 rounded-xl bg-black/30 border border-white/10 text-[var(--text)] font-bold focus:outline-none focus:ring-1 focus:ring-[var(--accent)]"
            >
              {currencyOptions.map((c) => (
                <option key={c} value={c} className="bg-slate-900 text-white">
                  {c}
                </option>
              ))}
            </select>
          </div>

          <div className="grid grid-cols-2 gap-2">
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value as any)}
              className="p-2.5 rounded-xl bg-black/20 border border-white/10 text-[var(--text)] capitalize font-medium"
            >
              <option value="food">Food</option>
              <option value="stay">Stay</option>
              <option value="transport">Transport</option>
              <option value="activity">Activity</option>
              <option value="other">Other</option>
            </select>

            {/* Converted Base Amount Indicator */}
            <div className="p-2.5 rounded-xl bg-white/[0.03] border border-white/5 flex items-center justify-between text-xs">
              <span className="text-[var(--text-muted)]">Base ({trip.currency}):</span>
              <span className="font-bold text-[var(--text)]">
                {formatCurrency(computedBaseAmount, trip.currency)}
              </span>
            </div>
          </div>

          {/* FX Conversion Notice & Custom Override */}
          {selectedCurrency !== trip.currency && (
            <div className="p-2.5 rounded-xl bg-sky-500/10 border border-sky-500/20 space-y-2">
              <div className="flex items-center justify-between text-[11px]">
                <div className="flex items-center gap-1.5 text-sky-300">
                  <ArrowRightLeft className="w-3.5 h-3.5" />
                  <span>
                    1 {trip.currency} ≈ {effectiveForeignPerBase.toFixed(2)} {selectedCurrency}
                  </span>
                </div>
                <button
                  type="button"
                  onClick={() => setIsCustomRate(!isCustomRate)}
                  className="text-sky-400 hover:underline font-medium text-[10px]"
                >
                  {isCustomRate ? "Use Live Rate" : "Custom Rate"}
                </button>
              </div>

              {isCustomRate && (
                <div className="flex items-center gap-2 pt-1 border-t border-sky-500/20">
                  <span className="text-[10px] text-sky-200 shrink-0">
                    1 {trip.currency} =
                  </span>
                  <input
                    type="number"
                    step="any"
                    placeholder={foreignPerBase.toFixed(2)}
                    value={customRateInput}
                    onChange={(e) => setCustomRateInput(e.target.value)}
                    className="flex-1 p-1.5 rounded-lg bg-black/30 border border-sky-500/30 text-white text-xs font-semibold"
                  />
                  <span className="text-[10px] text-sky-200 shrink-0">
                    {selectedCurrency}
                  </span>
                </div>
              )}
            </div>
          )}

          <div className="flex justify-end gap-2 pt-1">
            <button
              type="button"
              onClick={() => {
                setIsAdding(false);
                setIsCustomRate(false);
              }}
              className="px-3 py-1.5 rounded-xl bg-white/10 text-[var(--text-muted)] hover:bg-white/15 transition-colors"
            >
              {t("common.cancel")}
            </button>
            <button
              type="submit"
              className="px-4 py-1.5 rounded-xl bg-[var(--accent)] hover:opacity-95 text-white font-semibold shadow-md shadow-[var(--accent)]/30 transition-all"
            >
              {t("budget.save")}
            </button>
          </div>
        </form>
      )}

      {/* Expense List with Dual-Currency Display (FEAT-03) */}
      <div className="space-y-1.5 max-h-60 overflow-y-auto pr-1">
        {trip.expenses.map((exp) => {
          const hasForeign =
            exp.originalCurrency &&
            exp.originalCurrency !== trip.currency &&
            exp.originalAmount !== undefined;

          return (
            <div
              key={exp.id}
              className="flex items-center justify-between p-2.5 rounded-xl bg-white/[0.02] border border-white/5 hover:bg-white/[0.04] transition-colors"
            >
              <div className="space-y-0.5">
                <div className="font-semibold text-[var(--text)]">{exp.title}</div>
                <div className="text-[10px] text-[var(--text-muted)] capitalize flex items-center gap-1.5">
                  <span>{exp.category}</span>
                  <span>•</span>
                  <span>{exp.paidAt}</span>
                  {exp.isCustomRate && (
                    <span className="px-1.5 py-0.2 rounded bg-amber-500/20 text-amber-300 font-bold text-[9px]">
                      custom FX
                    </span>
                  )}
                </div>
              </div>

              <div className="text-right">
                <div className="font-bold text-[var(--text)]">
                  {formatCurrency(exp.baseAmount ?? exp.amount, trip.currency)}
                </div>
                {hasForeign && (
                  <div className="text-[10px] text-sky-400 font-medium">
                    {exp.originalAmount?.toLocaleString()} {exp.originalCurrency}
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
