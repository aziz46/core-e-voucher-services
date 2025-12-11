package connector

import (
	"context"
	"fmt"
	"math/rand"
)

// Provider defines the interface for provider connectors
type Provider interface {
	Inquiry(ctx context.Context, request InquiryRequest) (InquiryResponse, error)
	Pay(ctx context.Context, request PayRequest) (PayResponse, error)
	Cancel(ctx context.Context, request CancelRequest) (CancelResponse, error)
}

// InquiryRequest is the request for inquiry
type InquiryRequest struct {
	ProductCode string `json:"product_code"`
	CustomerNo  string `json:"customer_no"`
}

// InquiryResponse is the response from inquiry
type InquiryResponse struct {
	CustomerNo   string `json:"customer_no"`
	CustomerName string `json:"customer_name"`
	Amount       int64  `json:"amount"`
	AdminFee     int64  `json:"admin_fee"`
	Status       string `json:"status"`
}

// PayRequest is the request for payment
type PayRequest struct {
	ProductCode string `json:"product_code"`
	CustomerNo  string `json:"customer_no"`
	Amount      int64  `json:"amount"`
	RefNo       string `json:"ref_no"`
}

// PayResponse is the response from payment
type PayResponse struct {
	RefNo         string `json:"ref_no"`
	ProviderRefNo string `json:"provider_ref_no"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

// CancelRequest is the request for cancellation
type CancelRequest struct {
	RefNo         string `json:"ref_no"`
	ProviderRefNo string `json:"provider_ref_no"`
}

// CancelResponse is the response from cancellation
type CancelResponse struct {
	RefNo   string `json:"ref_no"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// MockProvider is a mock provider for testing and local development
type MockProvider struct {
	failureRate int // 0-100
}

// NewMockProvider creates a new mock provider
func NewMockProvider(failureRate int) *MockProvider {
	return &MockProvider{
		failureRate: failureRate,
	}
}

// Inquiry simulates an inquiry call
func (m *MockProvider) Inquiry(ctx context.Context, request InquiryRequest) (InquiryResponse, error) {
	if request.CustomerNo == "" {
		return InquiryResponse{}, fmt.Errorf("invalid customer number")
	}

	return InquiryResponse{
		CustomerNo:   request.CustomerNo,
		CustomerName: fmt.Sprintf("Customer %s", request.CustomerNo),
		Amount:       50000,
		AdminFee:     2500,
		Status:       "success",
	}, nil
}

// Pay simulates a payment call
func (m *MockProvider) Pay(ctx context.Context, request PayRequest) (PayResponse, error) {
	if request.CustomerNo == "" {
		return PayResponse{}, fmt.Errorf("invalid customer number")
	}

	// Simulate random failure based on failureRate
	if rand.Intn(100) < m.failureRate {
		return PayResponse{
			RefNo:   request.RefNo,
			Status:  "failed",
			Message: "Provider temporarily unavailable",
		}, fmt.Errorf("provider payment failed")
	}

	return PayResponse{
		RefNo:         request.RefNo,
		ProviderRefNo: fmt.Sprintf("MOCK-%d", rand.Int63()),
		Status:        "success",
		Message:       "Payment successful",
	}, nil
}

// Cancel simulates a cancellation call
func (m *MockProvider) Cancel(ctx context.Context, request CancelRequest) (CancelResponse, error) {
	if request.RefNo == "" {
		return CancelResponse{}, fmt.Errorf("invalid reference number")
	}

	return CancelResponse{
		RefNo:   request.RefNo,
		Status:  "success",
		Message: "Cancellation successful",
	}, nil
}
