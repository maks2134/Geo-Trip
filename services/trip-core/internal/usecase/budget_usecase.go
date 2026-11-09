package usecase

import (
	"context"
	"fmt"

	"github.com/makskozlov/atlasflow/services/trip-core/internal/domain"
)

type BudgetUsecase struct {
	repo  TripRepository
	clock Clock
	ids   IDGen
}

func NewBudgetUsecase(repo TripRepository, clock Clock, ids IDGen) *BudgetUsecase {
	return &BudgetUsecase{repo: repo, clock: clock, ids: ids}
}

type AddExpenseInput struct {
	TripID      string
	UserID      string
	DayID       *string
	PlaceItemID *string
	Category    domain.ExpenseCategory
	Amount      float64
	Currency    string
	Title       string
	Notes       string
}

func (u *BudgetUsecase) AddExpense(ctx context.Context, in AddExpenseInput) (domain.Expense, domain.BudgetSummary, error) {
	t, err := u.repo.GetTrip(ctx, in.TripID)
	if err != nil {
		return domain.Expense{}, domain.BudgetSummary{}, err
	}
	if err := t.AssertOwner(in.UserID); err != nil {
		return domain.Expense{}, domain.BudgetSummary{}, err
	}

	if in.Amount < 0 {
		return domain.Expense{}, domain.BudgetSummary{}, fmt.Errorf("%w: amount cannot be negative", domain.ErrValidation)
	}

	curr := in.Currency
	if curr == "" {
		curr = t.Currency
	}

	exp := domain.Expense{
		ID:          u.ids.New(),
		TripID:      in.TripID,
		DayID:       in.DayID,
		PlaceItemID: in.PlaceItemID,
		Category:    in.Category,
		Amount:      in.Amount,
		Currency:    curr,
		Title:       in.Title,
		Notes:       in.Notes,
		CreatedAt:   u.clock.Now(),
	}

	if err := u.repo.CreateExpense(ctx, exp); err != nil {
		return domain.Expense{}, domain.BudgetSummary{}, fmt.Errorf("create expense: %w", err)
	}

	summary, err := u.GetBudgetSummary(ctx, in.TripID, in.UserID)
	if err != nil {
		return domain.Expense{}, domain.BudgetSummary{}, err
	}

	return exp, summary, nil
}

func (u *BudgetUsecase) ListExpenses(ctx context.Context, tripID, userID string, dayID *string, cat *domain.ExpenseCategory) ([]domain.Expense, domain.BudgetSummary, error) {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return nil, domain.BudgetSummary{}, err
	}
	if err := t.AssertOwner(userID); err != nil {
		return nil, domain.BudgetSummary{}, err
	}

	expenses, err := u.repo.ListExpenses(ctx, tripID, dayID, cat)
	if err != nil {
		return nil, domain.BudgetSummary{}, err
	}

	allExpenses, err := u.repo.ListExpenses(ctx, tripID, nil, nil)
	if err != nil {
		return nil, domain.BudgetSummary{}, err
	}

	summary := domain.CalculateBudgetSummary(t.BudgetLimit, t.Currency, allExpenses)
	return expenses, summary, nil
}

func (u *BudgetUsecase) DeleteExpense(ctx context.Context, tripID, expenseID, userID string) (domain.BudgetSummary, error) {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return domain.BudgetSummary{}, err
	}
	if err := t.AssertOwner(userID); err != nil {
		return domain.BudgetSummary{}, err
	}

	if err := u.repo.DeleteExpense(ctx, tripID, expenseID); err != nil {
		return domain.BudgetSummary{}, err
	}

	allExpenses, err := u.repo.ListExpenses(ctx, tripID, nil, nil)
	if err != nil {
		return domain.BudgetSummary{}, err
	}

	return domain.CalculateBudgetSummary(t.BudgetLimit, t.Currency, allExpenses), nil
}

func (u *BudgetUsecase) GetBudgetSummary(ctx context.Context, tripID, userID string) (domain.BudgetSummary, error) {
	t, err := u.repo.GetTrip(ctx, tripID)
	if err != nil {
		return domain.BudgetSummary{}, err
	}
	if err := t.AssertOwner(userID); err != nil {
		return domain.BudgetSummary{}, err
	}

	all, err := u.repo.ListExpenses(ctx, tripID, nil, nil)
	if err != nil {
		return domain.BudgetSummary{}, err
	}

	return domain.CalculateBudgetSummary(t.BudgetLimit, t.Currency, all), nil
}