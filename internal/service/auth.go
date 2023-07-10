package service

import (
	"chatto/internal/constant"
	"chatto/internal/model/common"

	"log"
	"net/http"

	"chatto/internal/config"
	"chatto/internal/model"
	"chatto/internal/repository"
	"chatto/internal/util"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt"
)

func NewAuthService(conf *config.AppConfig, authRepos repository.IAuthRepository, userRepos repository.IUserRepository) IAuthService {
	return &authService{serverConfig: conf, authRepos: authRepos, userRepos: userRepos}
}

type authService struct {
	serverConfig *config.AppConfig
	authRepos    repository.IAuthRepository
	userRepos    repository.IUserRepository
}

func (a *authService) SignIn(input *model.SignInInput, sysInfo *common.SystemInfo) (string, common.Error) {
	// Check Credential
	user, err := a.userRepos.FindUserByName(input.Username)
	if err != nil {
		log.Println("Error Login: ", err)
		return "", common.NewError(http.StatusBadRequest, constant.ERR_LOGIN_USERNAME_PASSWORD)
	}

	err = util.ValidatePassword(user.Password, input.Password)
	if err != nil {
		log.Println("Error Login: ", err)
		return "", common.NewError(http.StatusBadRequest, constant.ERR_LOGIN_USERNAME_PASSWORD)
	}

	// Generate JWT Token
	refreshId := uuid.NewString()
	// Refresh Token
	refreshClaims := make(jwt.MapClaims)
	refreshClaims["user_id"] = user.Id

	refreshToken, err := util.CreateToken(refreshClaims, constant.REFRESH_TOKEN_EXP_TIME, a.serverConfig.JWTSecretKey)
	if err != nil {
		log.Println("Error Login: ", err)
		return "", common.NewError(http.StatusBadRequest, constant.ERR_TOKEN_CREATION)
	}

	savedToken := model.TokenDetails{
		Id:         refreshId,
		UserId:     user.Id,
		Token:      refreshToken,
		DeviceName: sysInfo.Name + " " + sysInfo.Os,
	}
	err = a.authRepos.SaveToken(&savedToken)
	if err != nil {
		log.Println("Error Save Token: ", err)
		return "", common.NewError(http.StatusBadRequest, constant.ERR_TOKEN_CREATION)
	}

	// Access Token
	accessToken, err := a.generateAccessToken(user.Id, refreshId)
	if err != nil {
		log.Println("Error Save Token: ", err)
		return "", common.NewError(http.StatusBadRequest, constant.ERR_TOKEN_CREATION)
	}

	return accessToken, common.NoError()
}

func (a *authService) SignUp(input *model.SignInInput) common.Error {
	password, err := util.HashPassword(input.Password)
	if err != nil {
		return common.NewError(http.StatusBadRequest, constant.ERR_SIGNUP)
	}

	user := model.User{
		Name:     input.Username,
		Password: password,
	}
	err = a.userRepos.CreateUser(&user)
	if err != nil {
		return common.NewError(http.StatusBadRequest, constant.ERR_SIGNUP)
	}
	return common.NoError()
}

func (a *authService) Logout(userId string, refreshId string) common.Error {
	token, err := a.authRepos.FindTokenById(refreshId)
	if err != nil {
		log.Println(err)
		return common.NewError(http.StatusBadRequest, constant.ERR_LOGOUT)
	}
	if token.UserId != userId {
		return common.NewError(http.StatusBadRequest, constant.ERR_LOGOUT)
	}
	if err := a.authRepos.RemoveTokenById(refreshId); err != nil {
		log.Println("Logout Error: ", err)
		return common.NewError(http.StatusBadRequest, constant.ERR_LOGOUT)
	}
	return common.NoError()
}

func (a *authService) LogoutAllDevice(userId string) common.Error {
	if err := a.authRepos.RemoveTokensByUserId(userId); err != nil {
		log.Println("Logout Error: ", err)
		return common.NewError(http.StatusBadRequest, constant.ERR_LOGOUT)
	}
	return common.NoError()
}

func (a *authService) RefreshToken(accessToken string) (string, common.Error) {

	// Parse token
	accessJwtToken, err := util.ParseToken(accessToken, false, a.serverConfig.JWTKeyFunc)
	if err != nil {
		log.Println("Error Parse Access Token: ", err)
		return "", common.NewError(http.StatusBadRequest, constant.ERR_TOKEN_FORMAT)
	}

	// Get access claims
	accessClaims := accessJwtToken.Claims.(jwt.MapClaims)
	refreshId := accessClaims["refresh_id"].(string)
	userId := accessClaims["user_id"].(string)

	// Check if the user id is exists
	user, err := a.userRepos.FindUserById(userId)
	if err != nil {
		// Delete all refresh id belong to it
		log.Println("Error Find User: ", err)
		err = a.authRepos.RemoveTokensByUserId(userId)
		if err != nil {
			log.Println("Error Remove Tokens: ", err)
		}
		return "", common.NewError(http.StatusBadRequest, constant.ERR_TOKEN_NO_OWNER)
	}

	// Get Token by refreshId
	refreshToken, err := a.authRepos.FindTokenById(refreshId)
	if err != nil {
		log.Println("Error Find Token: ", err)
		return "", common.NewError(http.StatusBadRequest, constant.ERR_TOKEN_REFRESH)
	}

	// Check relation access token to refresh token
	if refreshToken.Id != refreshId {
		return "", common.NewError(http.StatusBadRequest, constant.ERR_TOKEN_NO_OWNER)
	}

	// Validate refresh token
	err = util.ValidateToken(refreshToken.Token, a.serverConfig.JWTKeyFunc)
	if err != nil {
		log.Println("Error Parse Refresh Token: ", err)
		return "", common.NewError(http.StatusBadRequest, constant.ERR_TOKEN_FORMAT)
	}

	// Generate new access token
	accessToken, err = a.generateAccessToken(user.Id, refreshId)
	if err != nil {
		log.Println("Error Generate Token: ", err)
		return "", common.NewError(http.StatusBadRequest, constant.ERR_TOKEN_CREATION)
	}

	// TODO: Maybe need to rotate refresh token

	return accessToken, common.NoError()
}

func (a *authService) generateAccessToken(userId string, refreshId string) (string, error) {
	accessClaims := make(jwt.MapClaims)
	accessClaims["user_id"] = userId
	accessClaims["refresh_id"] = refreshId // For exact refresh_id for multiple login
	return util.CreateToken(accessClaims, constant.ACCESS_TOKEN_EXP_TIME, a.serverConfig.JWTSecretKey)
}
