package grpcuser

import (
	"sort"
	"sync"
	
	"demo-user/util"
	"demo-user/models"
	transactionpb "demo-user/proto/models/transaction"
)

func convertToTransactionDetailList(data []*transactionpb.TransactionDetail) []models.TransactionDetail {
	var (
		result = make([]models.TransactionDetail, 0)
		wg     sync.WaitGroup
	)

	total := len(data)
	// Add process
	wg.Add(total)

	for index := range data {
		go func(index int) {
			defer wg.Done()

			// Convert to TransactionDetail
			transaction := convertToTransactionDetail(data[index])

			// Append
			result = append(result, transaction)
		}(index)
	}

	// Wait process
	wg.Wait()

	// Sort
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result
}

func convertToTransactionDetail(transaction *transactionpb.TransactionDetail) models.TransactionDetail {
	var (
		id ,_       = util.HelperParseStringToObjectID(transaction.Id)
		companyID,_ = util.HelperParseStringToObjectID(transaction.CompanyID)
		branchID,_  = util.HelperParseStringToObjectID(transaction.BranchID)
		userID,_    = util.HelperParseStringToObjectID(transaction.UserID)
		createdAt = util.HelperConvertTimestampProtoToTime(transaction.CreatedAt)
	)

	// TransactionDetail
	result := models.TransactionDetail{
		ID:                     id,
		CompanyID:              companyID,
		BranchID:               branchID,
		UserID:                 userID,
		Amount:                 transaction.Amount,
		Commission:             transaction.Commission,
		CompanyCashbackPercent: transaction.CompanyCashbackPercent,
		CreatedAt:              createdAt,
	}
	return result
}
