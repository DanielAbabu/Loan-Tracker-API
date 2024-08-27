package usecase

import (
	"context"
	"errors"
	"loan-tracker/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type loanUsecase struct {
	loanRepo domain.LoanRepository
}

// NewLoanUsecase creates a new instance of LoanUsecase
func NewLoanUsecase(loanRepo domain.LoanRepository) domain.LoanUsecase {
	return &loanUsecase{
		loanRepo: loanRepo,
	}
}

// ApplyForLoan handles the business logic for applying for a loan
func (uc *loanUsecase) ApplyForLoan(ctx context.Context, loan domain.Loan) (primitive.ObjectID, error) {
	loan.ID = primitive.NewObjectID()
	loan.Status = "pending"
	loan.CreatedAt = time.Now()
	loan.UpdatedAt = time.Now()

	loanID, err := uc.loanRepo.CreateLoan(ctx, loan)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return loanID, nil
}

// ViewLoanStatus retrieves the status of a specific loan
func (uc *loanUsecase) ViewLoanStatus(ctx context.Context, id primitive.ObjectID) (domain.Loan, error) {
	loan, err := uc.loanRepo.GetLoanByID(ctx, id)
	if err != nil {
		return domain.Loan{}, err
	}
	return loan, nil
}

// ViewAllLoans retrieves all loan applications based on the provided filter
func (uc *loanUsecase) ViewAllLoans(ctx context.Context) ([]domain.Loan, error) {

	loans, err := uc.loanRepo.GetAllLoans(ctx)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

// ApproveOrRejectLoan handles the business logic for approving or rejecting a loan application
func (uc *loanUsecase) ApproveOrRejectLoan(ctx context.Context, id primitive.ObjectID, status string) error {
	if status != "approved" && status != "rejected" {
		return errors.New("invalid status value, you can only enter approved or rejected")
	}

	err := uc.loanRepo.UpdateLoanStatus(ctx, id, status)
	if err != nil {
		return err
	}
	return nil
}

// DeleteLoan handles the business logic for deleting a loan application
func (uc *loanUsecase) DeleteLoan(ctx context.Context, id primitive.ObjectID) error {
	err := uc.loanRepo.DeleteLoan(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
