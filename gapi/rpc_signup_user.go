package gapi

import (
	"context"
	"strings"

	"github.com/thanhpk/randstr"
	"github.com/wpcodevo/golang-mongodb/models"
	"github.com/wpcodevo/golang-mongodb/pb"
	"github.com/wpcodevo/golang-mongodb/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (authServer *AuthServer) SignUpUser(ctx context.Context, req *pb.SignUpUserInput) (*pb.GenericResponse, error) {
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	user := models.SignUpInput{
		Name:            req.GetName(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}

	newUser, err := authServer.authService.SignUpUser(&user)

	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())

		}
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	// Generate Verification Code
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	updateData := &models.UpdateInput{
		VerificationCode: verificationCode,
	}

	// Update User in Database
	authServer.userService.UpdateUserById(newUser.ID.Hex(), updateData)

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[0]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       authServer.config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(newUser, &emailData, "verificationCode.html")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "There was an error sending email: %s", err.Error())

	}

	message := "We sent an email with a verification code to " + newUser.Email

	res := &pb.GenericResponse{
		Status:  "success",
		Message: message,
	}
	return res, nil
}
