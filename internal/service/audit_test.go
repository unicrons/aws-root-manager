package service

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/unicrons/aws-root-manager/internal/aws"

// 	awssdk "github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // MockAWSProvider is a mock implementation of AWSProvider for testing
// type MockAWSProvider struct {
// 	mock.Mock
// }

// func (m *MockAWSProvider) GetIAMService() aws.IAMService {
// 	args := m.Called()
// 	return args.Get(0).(aws.IAMService)
// }

// func (m *MockAWSProvider) GetSTSService() aws.STSService {
// 	args := m.Called()
// 	return args.Get(0).(aws.STSService)
// }

// func (m *MockAWSProvider) GetIAMServiceForConfig(config awssdk.Config) aws.IAMService {
// 	args := m.Called(config)
// 	return args.Get(0).(aws.IAMService)
// }

// // MockIAMService is a mock implementation of IAMService for testing
// type MockIAMService struct {
// 	mock.Mock
// }

// func (m *MockIAMService) CheckOrganizationRootAccess(ctx context.Context, rootSessionsRequired bool) error {
// 	args := m.Called(ctx, rootSessionsRequired)
// 	return args.Error(0)
// }

// func (m *MockIAMService) GetLoginProfile(ctx context.Context, accountId string) (bool, error) {
// 	args := m.Called(ctx, accountId)
// 	return args.Bool(0), args.Error(1)
// }

// func (m *MockIAMService) DeleteLoginProfile(ctx context.Context, accountId string) error {
// 	args := m.Called(ctx, accountId)
// 	return args.Error(0)
// }

// func (m *MockIAMService) ListAccessKeys(ctx context.Context, accountId string) ([]string, error) {
// 	args := m.Called(ctx, accountId)
// 	return args.Get(0).([]string), args.Error(1)
// }

// func (m *MockIAMService) DeleteAccessKeys(ctx context.Context, accountId string, accessKeyIds []string) error {
// 	args := m.Called(ctx, accountId, accessKeyIds)
// 	return args.Error(0)
// }

// func (m *MockIAMService) ListMFADevices(ctx context.Context, accountId string) ([]string, error) {
// 	args := m.Called(ctx, accountId)
// 	return args.Get(0).([]string), args.Error(1)
// }

// func (m *MockIAMService) DeactivateMFADevices(ctx context.Context, accountId string, mfaSerialNumbers []string) error {
// 	args := m.Called(ctx, accountId, mfaSerialNumbers)
// 	return args.Error(0)
// }

// func (m *MockIAMService) ListSigningCertificates(ctx context.Context, accountId string) ([]string, error) {
// 	args := m.Called(ctx, accountId)
// 	return args.Get(0).([]string), args.Error(1)
// }

// func (m *MockIAMService) DeleteSigningCertificates(ctx context.Context, accountId string, certificates []string) error {
// 	args := m.Called(ctx, accountId, certificates)
// 	return args.Error(0)
// }

// func (m *MockIAMService) EnableOrganizationsRootCredentialsManagement(ctx context.Context) error {
// 	args := m.Called(ctx)
// 	return args.Error(0)
// }

// func (m *MockIAMService) EnableOrganizationsRootSessions(ctx context.Context) error {
// 	args := m.Called(ctx)
// 	return args.Error(0)
// }

// func (m *MockIAMService) CreateLoginProfile(ctx context.Context) error {
// 	args := m.Called(ctx)
// 	return args.Error(0)
// }

// // MockSTSService is a mock implementation of STSService for testing
// type MockSTSService struct {
// 	mock.Mock
// }

// func (m *MockSTSService) GetAssumeRootConfig(ctx context.Context, accountId, taskPolicyName string) (awssdk.Config, error) {
// 	args := m.Called(ctx, accountId, taskPolicyName)
// 	return args.Get(0).(awssdk.Config), args.Error(1)
// }

// func (m *MockSTSService) assumeRoot(ctx context.Context, accountId, taskPolicyName string) (awssdk.Credentials, error) {
// 	args := m.Called(ctx, accountId, taskPolicyName)
// 	return args.Get(0).(awssdk.Credentials), args.Error(1)
// }

// // TestAuditAccounts demonstrates how to test the AuditAccounts function
// func TestAuditAccounts(t *testing.T) {
// 	// Set up test case
// 	ctx := context.Background()
// 	accountIDs := []string{"123456789012"}

// 	// Create mock services
// 	mockIAM := new(MockIAMService)
// 	mockSTS := new(MockSTSService)
// 	mockIAMForConfig := new(MockIAMService)

// 	// Set up mock expectations
// 	mockIAM.On("CheckOrganizationRootAccess", ctx, false).Return(nil)

// 	// Mock STS.GetAssumeRootConfig for each account
// 	mockSTS.On("GetAssumeRootConfig", ctx, "123456789012", "IAMAuditRootUserCredentials").Return(awssdk.Config{}, nil)

// 	// Mock IAM for config responses
// 	mockIAMForConfig.On("GetLoginProfile", ctx, "123456789012").Return(true, nil)
// 	mockIAMForConfig.On("ListAccessKeys", ctx, "123456789012").Return([]string{"AKIA123456789"}, nil)
// 	mockIAMForConfig.On("ListMFADevices", ctx, "123456789012").Return([]string{"arn:aws:iam::123456789012:mfa/root"}, nil)
// 	mockIAMForConfig.On("ListSigningCertificates", ctx, "123456789012").Return([]string{"ABCDEF123456"}, nil)

// 	// Configure iam for config provider function
// 	iamForConfigProvider := func(cfg awssdk.Config) aws.IAMService {
// 		return mockIAMForConfig
// 	}

// 	// Create the service with the specialized constructor
// 	service := NewCompleteService(mockIAM, mockSTS, nil, iamForConfigProvider)

// 	// Call the method being tested
// 	result, err := service.AuditAccounts(ctx, accountIDs)

// 	// Assertions
// 	assert.NoError(t, err)
// 	assert.Len(t, result, 1)
// 	assert.Equal(t, "123456789012", result[0].AccountId)
// 	assert.True(t, result[0].LoginProfile)
// 	assert.Equal(t, []string{"AKIA123456789"}, result[0].AccessKeys)
// 	assert.Equal(t, []string{"arn:aws:iam::123456789012:mfa/root"}, result[0].MfaDevices)
// 	assert.Equal(t, []string{"ABCDEF123456"}, result[0].SigningCertificates)

// 	// Verify that all expected calls were made
// 	mockIAM.AssertExpectations(t)
// 	mockSTS.AssertExpectations(t)
// 	mockIAMForConfig.AssertExpectations(t)
// }

// // TestAuditAccountsError demonstrates testing error handling
// func TestAuditAccountsError(t *testing.T) {
// 	// Set up test case
// 	ctx := context.Background()
// 	accountIDs := []string{"123456789012"}
// 	expectedError := errors.New("organization access not enabled")

// 	// Create mock service with error behavior
// 	mockIAM := new(MockIAMService)
// 	mockIAM.On("CheckOrganizationRootAccess", ctx, false).Return(expectedError)

// 	// Create the service with only the IAM service
// 	service := NewIAMService(mockIAM)

// 	// Call the method being tested
// 	_, err := service.AuditAccounts(ctx, accountIDs)

// 	// Assertions
// 	assert.Error(t, err)
// 	assert.Equal(t, expectedError, err)

// 	// Verify that all expected calls were made
// 	mockIAM.AssertExpectations(t)
// }
