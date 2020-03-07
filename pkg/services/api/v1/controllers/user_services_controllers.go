package controllers

import (
	"context"
	"database/sql"
	"os"
	v1 "protoChatServices/pkg/api/v1"
	"protoChatServices/pkg/services/api/v1/global"
	"protoChatServices/pkg/services/api/v1/models"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServices implemented on version 1 proto interface
type chatServices struct {
	db *sql.DB
}

// New sending otp services create sending otp service
func NewChatServicesService(db *sql.DB) v1.ChatServicesServer {
	return &chatServices{db: db}
}

// Checking Api version
func (s *chatServices) CheckApi(api string) error {
	if len(api) > 0 {
		if os.Getenv("API_VERSION") != api {
			return status.Errorf(codes.Unimplemented, "Unsupported API Version: Service API implement using '%s', but asked for '%s'", os.Getenv("API_VERSION"), api)
		}
	}
	return nil
}

// Function to check token valid or not
func Auth(data string) (string, bool) {
	tokenize := &global.Tokenization{}

	tkn, err := jwt.ParseWithClaims(data, tokenize, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN")), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "Invalid signature", false
		}
		return "Bas Request", false
	}

	if !tkn.Valid {
		return "Invalid Token", false
	} else {
		return tokenize.Username, true
	}
}

// Function for registering username
func (s *chatServices) RegisterUsername(ctx context.Context, req *v1.RegisterUsernameRequest) (*v1.RegisterUsernameResponse, error) {
	var code, message, status, token, registeredDate string

	// Checking API Version
	if err := s.CheckApi(req.Api); err != nil {
		return nil, err
	} else {
		code, status, message, token, registeredDate = models.RegisterAccount(req.Username)
	}

	return &v1.RegisterUsernameResponse{
		Code:    code,
		Message: message,
		Status:  status,
		Data: &v1.DataResponseRegistration{
			Token:      token,
			LoggedTime: registeredDate,
		},
	}, nil
}

// Func register account
func (s *chatServices) SendingMessage(ctx context.Context, req *v1.RequestChat) (*v1.ResponseChat, error) {
	// Get data IP Address
	// timeRequest := time.Now().Format("2006-01-02 15:04:05")
	// var code, status, message, statusMessage string

	// headers, _ := metadata.FromIncomingContext(ctx)
	// data := headers["authorization"]

	// decodedData, _ := models.Decrypt(data[0])
	// fmt.Println(decodedData)

	// data, status := Auth(decodedData)

	// if status == false {
	// 	code = "403"
	// 	status = "Unauthorized"
	// 	message = "Token invalid"
	// } else {
	// 	code, status, message, statusMessage = models.
	// }

	// return &v1.ResponseChat{
	// 	Code: code,
	// 	Status: status,
	// 	Message: message,
	// 	Data: &v1.DataResponseChat{
	// 		Username: data,
	// 		MessageStatus: statusMessage,
	// 		MessageTime: timeRequest
	// 	},
	// }, nil

	return nil, nil
}
