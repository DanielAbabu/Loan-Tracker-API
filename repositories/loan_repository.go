package repositories

import (
	"context"
	"loan-tracker/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type loanRepository struct {
	db *mongo.Collection
}

func NewLoanRepository(db *mongo.Client) domain.LoanRepository {
	return &loanRepository{
		db: db.Database("loan-tracker").Collection("loans"),
	}
}

func (r *loanRepository) CreateLoan(ctx context.Context, loan domain.Loan) (primitive.ObjectID, error) {
	result, err := r.db.InsertOne(ctx, loan)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *loanRepository) GetLoanByID(ctx context.Context, id primitive.ObjectID) (domain.Loan, error) {
	var loan domain.Loan
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&loan)
	return loan, err
}

func (r *loanRepository) GetAllLoans(ctx context.Context) ([]domain.Loan, error) {
	var loans []domain.Loan
	cursor, err := r.db.Find(ctx, bson.M{}, options.Find())

	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &loans)
	return loans, err
}

func (r *loanRepository) UpdateLoanStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (r *loanRepository) DeleteLoan(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
