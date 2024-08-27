package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Loan struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Description string             `json:"description"`
	Amount      float64            `json:"amount"`
	Status      string             `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
type LoanFilter struct {
	Status string
	Order  string
}

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan Loan) (primitive.ObjectID, error)
	GetLoanByID(ctx context.Context, id primitive.ObjectID) (Loan, error)
	GetAllLoans(ctx context.Context) ([]Loan, error)
	UpdateLoanStatus(ctx context.Context, id primitive.ObjectID, status string) error
	DeleteLoan(ctx context.Context, id primitive.ObjectID) error
}

type LoanUsecase interface {
	ApplyForLoan(ctx context.Context, loan Loan) (primitive.ObjectID, error)
	ViewLoanStatus(ctx context.Context, id primitive.ObjectID) (Loan, error)
	ViewAllLoans(ctx context.Context) ([]Loan, error)
	ApproveOrRejectLoan(ctx context.Context, id primitive.ObjectID, status string) error
	DeleteLoan(ctx context.Context, id primitive.ObjectID) error
}
