"use client";

import React, { useState } from "react";
import { Plus, PieChart } from "lucide-react";
import { Trip, Expense } from "@/entities/trip";
import { ProgressBar } from "@/shared/ui/ProgressBar";
import { formatCurrency } from "@/shared/lib/utils";
import { apiClient } from "@/shared/lib/api";

interface BudgetTabProps {
  trip: Trip;
  onTripUpdate: (updated: Trip) => void;
}

export function BudgetTab({ trip, onTripUpdate }: BudgetTabProps) {
  const [isAdding, setIsAdding] = useState(false);
  const [title, setTitle] = useState("");
  const [amount, setAmount] = useState("");
  const [category, setCategory] = useState<"stay" | "food" | "transport" | "activity" | "other">("food");

  const totalSpent = trip.expenses.reduce((acc, exp) => acc + exp.amount, 0);
  const percentUsed = Math.min(100, Math.round((totalSpent / trip.budgetLimit) * 100));

  const categoryTotals: Record<string, number> = {};
  trip.expenses.forEach((e) => {
    categoryTotals[e.category] = (categoryTotals[e.category] || 0) + e.amount;
  });

  const handleAddExpense = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title || !amount) return;

    await apiClient.addExpense(trip.id, {
      tripId: trip.id,
      title,
      amount: parseFloat(amount),
      currency: trip.currency,
      category,
      paidAt: new Date().toISOString().split("T")[0],
    });

    const refreshed = await apiClient.getTrip(trip.id);
    if (refreshed) onTripUpdate(refreshed);

    setTitle("");
    setAmount("");
    setIsAdding(false);
  };

  return (
    <div className="space-y-4 text-xs">
      <div className="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-3">
        <div className="flex items-center justify-between">
          <div>
            <div className="text-[11px] text-[var(--text-muted)] uppercase tracking-wider">Total Budget</div>
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
                percentUsed > 90 ? "text-red-400" : percentUsed > 70 ? "text-amber-400" : "text-emerald-400"
              }`}
            >
              {percentUsed}%
            </span>
            <div className="text-[10px] text-[var(--text-muted)]">Allocated</div>
          </div>
        </div>

        <ProgressBar
          value={percentUsed}
          barColor={percentUsed > 90 ? "bg-red-500" : percentUsed > 70 ? "bg-amber-500" : "bg-emerald-500"}
        />
      </div>

      <div className="p-3 rounded-2xl bg-white/[0.03] border border-white/5 space-y-2">
        <div className="text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)] flex items-center gap-1.5">
          <PieChart className="w-3.5 h-3.5" />
          <span>Category Breakdown</span>
        </div>
        <div className="grid grid-cols-2 gap-2 pt-1">
          {Object.entries(categoryTotals).map(([cat, val]) => (
            <div key={cat} className="p-2 rounded-xl bg-white/[0.02] flex items-center justify-between">
              <span className="capitalize text-[var(--text-muted)]">{cat}</span>
              <span className="font-semibold text-[var(--text)]">{formatCurrency(val, trip.currency)}</span>
            </div>
          ))}
        </div>
      </div>

      <div className="flex items-center justify-between pt-1">
        <span className="font-semibold text-xs text-[var(--text)]">Recent Expenses</span>
        <button
          onClick={() => setIsAdding(!isAdding)}
          className="flex items-center gap-1 px-2.5 py-1 rounded-lg bg-[var(--accent)] text-white font-medium hover:opacity-90 transition-opacity"
        >
          <Plus className="w-3.5 h-3.5" />
          <span>Add</span>
        </button>
      </div>

      {isAdding && (
        <form onSubmit={handleAddExpense} className="p-3 rounded-xl bg-white/[0.06] border border-white/10 space-y-2">
          <input
            type="text"
            placeholder="Expense description (e.g. Shinkansen ticket)"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full p-2 rounded-lg bg-black/20 border border-white/10 text-[var(--text)]"
            required
          />
          <div className="flex gap-2">
            <input
              type="number"
              placeholder="Amount"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="w-1/2 p-2 rounded-lg bg-black/20 border border-white/10 text-[var(--text)]"
              required
            />
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value as any)}
              className="w-1/2 p-2 rounded-lg bg-black/20 border border-white/10 text-[var(--text)] capitalize"
            >
              <option value="stay">Stay</option>
              <option value="food">Food</option>
              <option value="transport">Transport</option>
              <option value="activity">Activity</option>
              <option value="other">Other</option>
            </select>
          </div>
          <div className="flex justify-end gap-2 pt-1">
            <button
              type="button"
              onClick={() => setIsAdding(false)}
              className="px-3 py-1 rounded-lg bg-white/10 text-[var(--text-muted)]"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-3 py-1 rounded-lg bg-[var(--accent)] text-white font-semibold"
            >
              Save
            </button>
          </div>
        </form>
      )}

      <div className="space-y-1.5 max-h-56 overflow-y-auto pr-1">
        {trip.expenses.map((exp) => (
          <div
            key={exp.id}
            className="flex items-center justify-between p-2.5 rounded-xl bg-white/[0.02] border border-white/5"
          >
            <div>
              <div className="font-semibold text-[var(--text)]">{exp.title}</div>
              <div className="text-[10px] text-[var(--text-muted)] capitalize">
                {exp.category} • {exp.paidAt}
              </div>
            </div>
            <div className="font-bold text-[var(--text)]">{formatCurrency(exp.amount, exp.currency)}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
