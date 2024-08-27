package controllers

import (
	"loan-tracker/domain"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoanController struct {
	LoanUsecase domain.LoanUsecase
	LogUsecase  domain.LogUsecase
}

func NewLoanController(LoanUsecase domain.LoanUsecase, logUsecase domain.LogUsecase) *LoanController {
	return &LoanController{
		LoanUsecase: LoanUsecase,
		LogUsecase:  logUsecase,
	}
}
func (c *LoanController) ApplyForLoan(ctx *gin.Context) {
	var loan domain.Loan
	if err := ctx.ShouldBindJSON(&loan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.GetString("userid")
	log.Println("User ID: ", id)
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	loan.UserID = ID
	loanID, err := c.LoanUsecase.ApplyForLoan(ctx, loan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log loan application
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "loan_application",
		Details:   "Loan applied with ID: " + loanID.Hex(),
	}
	if logErr := c.LogUsecase.LogEvent(ctx, logEntry); logErr != nil {
		log.Println("Error logging loan application:", logErr)
	}

	ctx.JSON(http.StatusOK, gin.H{"loan_id": loanID})
}

func (c *LoanController) ViewLoanStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	loan, err := c.LoanUsecase.ViewLoanStatus(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log view loan status
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "view_loan_status",
		Details:   "Loan status retrieved for ID: " + id,
	}
	if logErr := c.LogUsecase.LogEvent(ctx, logEntry); logErr != nil {
		log.Println("Error logging view loan status:", logErr)
	}
	ctx.JSON(http.StatusOK, loan)
}

func (c *LoanController) ViewAllLoans(ctx *gin.Context) {

	loans, err := c.LoanUsecase.ViewAllLoans(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log view all loans
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "view_all_loans",
		Details:   "All loans retrieved",
	}
	if logErr := c.LogUsecase.LogEvent(ctx, logEntry); logErr != nil {
		log.Println("Error logging view all loans:", logErr)
	}

	ctx.JSON(http.StatusOK, loans)
}

func (c *LoanController) ApproveOrRejectLoan(ctx *gin.Context) {
	id := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var status domain.Loan
	if err := ctx.ShouldBindJSON(&status); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.LoanUsecase.ApproveOrRejectLoan(ctx, objID, status.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log loan approval/rejection
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "loan_approval_rejection",
		Details:   "Loan status updated for ID: " + id,
	}
	if logErr := c.LogUsecase.LogEvent(ctx, logEntry); logErr != nil {
		log.Println("Error logging loan approval/rejection:", logErr)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "loan status updated"})
}

func (c *LoanController) DeleteLoan(ctx *gin.Context) {
	id := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.LoanUsecase.DeleteLoan(ctx, objID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log loan deletion
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "loan_deletion",
		Details:   "Loan deleted with ID: " + id,
	}
	if logErr := c.LogUsecase.LogEvent(ctx, logEntry); logErr != nil {
		log.Println("Error logging loan deletion:", logErr)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "loan deleted"})
}
