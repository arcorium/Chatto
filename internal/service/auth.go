package service

import (
	"chatto/internal/constant"
	"chatto/internal/dto"
	"chatto/internal/model/common"
	"chatto/internal/repository"
	"errors"
	"log"
	"net/http"
	"strings"

	"chatto/internal/config"
	"chatto/internal/model"
	"chatto/internal/util"

	"github.com/golang-jwt/jwt"
)

type IAuthService interface {
	SignIn(input *dto.SignInInput, sysInfo *common.SystemInfo) (dto.SignInOutput, common.Error)
	SignUp(input *dto.SignUpInput) common.Error
	GetLoginDevices(userId string) ([]model.Device, common.Error)
	Logout(userId string, tokenId string) common.Error
	LogoutAllDevice(userId string) common.Error
	RefreshToken(input *dto.RefreshTokenInput) (dto.RefreshTokenOutput, common.Error)
}

func NewAuthService(conf *config.AppConfig, authRepos repository.IAuthRepository, userRepos repository.IUserRepository) IAuthService {
	return &authService{config: conf, authRepos: authRepos, userRepos: userRepos}
}

type authService struct {
	config    *config.AppConfig
	authRepos repository.IAuthRepository
	userRepos repository.IUserRepository
}

func (a *authService) SignIn(input *dto.SignInInput, sysInfo *common.SystemInfo) (dto.SignInOutput, common.Error) {
	// Check Credential
	user, err := a.userRepos.FindUserByName(input.Username)
	if err != nil {
		return dto.SignInOutput{}, common.NewError(common.INTERNAL_REPOSITORY_ERROR, constant.MSG_FAILED_USER_LOGIN)
	}

	err = util.ValidatePassword(user.Password, input.Password)
	if err != nil {
		return dto.SignInOutput{}, common.NewError(common.AUTH_TOKEN_NOT_VALIDATED_ERROR, constant.MSG_FAILED_USER_LOGIN)
	}

	refreshToken, err := a.generateRefreshToken(user.Id, sysInfo.Name+" "+sysInfo.Os)
	if err != nil {
		return dto.SignInOutput{}, common.NewError(common.CREATE_TOKEN_ERROR, constant.MSG_FAILED_USER_LOGIN)
	}

	err = a.authRepos.CreateToken(&refreshToken)
	if err != nil {
		return dto.SignInOutput{}, common.NewError(common.CREATE_TOKEN_ERROR, constant.MSG_FAILED_USER_LOGIN)
	}

	// Access Token
	accessToken, err := a.generateAccessToken(user, refreshToken.Id)
	return dto.NewSignInOutput(accessToken), common.NewConditionalError(err, common.CREATE_TOKEN_ERROR, constant.MSG_FAILED_USER_LOGIN)
}

func (a *authService) SignUp(input *dto.SignUpInput) common.Error {
	password, err := util.HashPassword(input.Password)
	if err != nil {
		return common.NewError(common.HASH_PASSWORD_ERROR, constant.MSG_FAILED_SIGNUP)
	}
	input.Password = password

	user := dto.NewUserFromSignUpInput(input)
	err = a.userRepos.CreateUser(&user)
	return common.NewConditionalError(err, common.USER_SIGNUP_ERROR, constant.MSG_FAILED_SIGNUP)
}
func (a *authService) GetLoginDevices(userId string) ([]model.Device, common.Error) {
	devices, err := a.authRepos.FindDevicesByUserId(userId)
	return devices, common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}
func (a *authService) Logout(userId string, refreshId string) common.Error {
	refreshToken, err := a.authRepos.FindTokenById(refreshId)
	if err != nil {
		return common.NewError(common.AUTH_TOKEN_NOT_FOUND_ERROR, constant.MSG_LOGOUT_FAILED)
	}
	// Check user id on token db
	// NOTE: It can be moved into repository to find token which have refresh id and user id
	if refreshToken.UserId != userId {
		return common.NewError(common.AUTH_TOKEN_BAD_OWNERSHIP_ERROR, constant.MSG_LOGOUT_FAILED)
	}
	err = a.authRepos.RemoveTokenById(refreshId)
	return common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_LOGOUT_FAILED)
}

func (a *authService) LogoutAllDevice(userId string) common.Error {
	err := a.authRepos.RemoveTokensByUserId(userId)
	return common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_LOGOUT_FAILED)
}

func (a *authService) RefreshToken(input *dto.RefreshTokenInput) (dto.RefreshTokenOutput, common.Error) {

	// Parse token
	jwtAccessToken, err := util.ParseToken(input.AccessToken, a.config.JWTKeyFunc)
	if err != nil {
		return dto.RefreshTokenOutput{}, common.NewError(http.StatusBadRequest, constant.MSG_BAD_FORMAT_TOKEN)
	}

	// Get access claims
	accessClaims, err := a.extractAccessTokenClaims(jwtAccessToken.Claims)
	if err != nil {
		return dto.RefreshTokenOutput{}, common.NewError(http.StatusBadRequest, constant.MSG_BAD_FORMAT_TOKEN)
	}

	// Check if the user id is exists
	user, err := a.userRepos.FindUserById(accessClaims.UserId)
	if err != nil {
		// Delete all refresh id belong to it because it can be malicious action because there is no such user
		err = a.authRepos.RemoveTokensByUserId(accessClaims.UserId)
		if err != nil {
			log.Println("Error Remove Tokens: ", err)
		}
		return dto.RefreshTokenOutput{}, common.NewError(common.AUTH_TOKEN_BAD_OWNERSHIP_ERROR, constant.MSG_NO_OWNER_TOKEN)
	}

	// Get Token by refreshId
	refreshToken, err := a.authRepos.FindTokenById(accessClaims.RefreshId)
	if err != nil {
		return dto.RefreshTokenOutput{}, common.NewError(common.AUTH_TOKEN_NOT_FOUND_ERROR, constant.MSG_TOKEN_NOT_FOUND)
	}

	// Check user id on token db
	// NOTE: It can be moved into repository to find token which have refresh id and user id
	if refreshToken.UserId != accessClaims.UserId {
		return dto.RefreshTokenOutput{}, common.NewError(common.AUTH_TOKEN_BAD_OWNERSHIP_ERROR, constant.MSG_NO_OWNER_TOKEN)
	}

	// Validate refresh token
	err = util.ValidateToken(refreshToken.Token, a.config.JWTKeyFunc)
	if err != nil {
		// Create new refresh token
		refreshToken, err = a.generateRefreshToken(accessClaims.UserId, refreshToken.DeviceName)
		if err != nil {
			return dto.RefreshTokenOutput{}, common.NewError(common.CREATE_TOKEN_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
		}

		// Replace token
		err = a.authRepos.UpdateToken(accessClaims.RefreshId, &refreshToken)
		if err != nil {
			return dto.RefreshTokenOutput{}, common.NewError(common.AUTH_TOKEN_UPDATE_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
		}
	}

	// Generate new access token
	accessToken, err := a.generateAccessToken(user, refreshToken.Id)
	return dto.NewRefreshTokenOutput(accessToken), common.NewConditionalError(err, common.CREATE_TOKEN_ERROR, constant.MSG_TOKEN_REFRESH_FAILED)
}

func (a *authService) generateAccessToken(user *model.User, refreshId string) (string, error) {
	accessClaims := make(jwt.MapClaims)
	accessClaims["user_id"] = user.Id
	accessClaims["name"] = user.Name
	accessClaims["role"] = user.Role
	accessClaims["refresh_id"] = refreshId // For exact refresh_id for multiple login
	return util.CreateAccessToken(accessClaims, a.config)
}
func (a *authService) extractAccessTokenClaims(claims jwt.Claims) (model.AccessTokenClaims, error) {
	// Extract
	accessTokenClaims := model.AccessTokenClaims{}
	rawAccessClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return accessTokenClaims, errors.New("broken claims")
	}
	refreshId, exist := rawAccessClaims["refresh_id"]
	if !exist {
		return accessTokenClaims, errors.New("refresh id not found")
	}
	userId, exist := rawAccessClaims["user_id"]
	if !exist {
		return accessTokenClaims, errors.New("user id not found")
	}

	// Cast
	refreshIdString, ok := refreshId.(string)
	if !ok {
		return accessTokenClaims, errors.New("refresh id value is not string")
	}
	userIdString, ok := userId.(string)
	if !ok {
		return accessTokenClaims, errors.New("user id value is not string")
	}

	accessTokenClaims.UserId = userIdString
	accessTokenClaims.RefreshId = refreshIdString
	return accessTokenClaims, nil
}
func (a *authService) generateRefreshToken(userId string, deviceName string) (model.Credential, error) {
	refreshClaims := make(jwt.MapClaims)
	//refreshClaims["user_id"] = userId	// NOTE: UserId already set on the database next to token, maybe not necessary?
	rawRefreshToken, err := util.CreateRefreshToken(refreshClaims, a.config)
	return model.NewCredential(userId, strings.TrimSpace(deviceName), rawRefreshToken), err
}
