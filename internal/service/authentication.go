package service

import (
	"log"
	"net/http"

	"server_client_chat/internal/config"
	"server_client_chat/internal/model"
	"server_client_chat/internal/repository"
	"server_client_chat/internal/util"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt"
)

func NewAuthService(conf *config.AppConfig, authRepos repository.IAuthRepository, userRepos repository.IUserRepository) AuthService {
	return AuthService{serverConfig: conf, authRepos: authRepos, userRepos: userRepos}
}

type AuthService struct {
	serverConfig *config.AppConfig
	authRepos    repository.IAuthRepository
	userRepos    repository.IUserRepository
}

func (a *AuthService) SignIn(input *model.SignInInput, sysInfo *model.SystemInfo) (string, CustomError) {
	// Check Credential
	user, err := a.userRepos.FindUserByName(input.Username)
	if err != nil {
		log.Println("Error Login: ", err)
		return "", NewError(http.StatusBadRequest, util.ERR_LOGIN_USERNAME_PASSWORD)
	}

	err = util.ValidatePassword(user.Password, input.Password)
	if err != nil {
		log.Println("Error Login: ", err)
		return "", NewError(http.StatusBadRequest, util.ERR_LOGIN_USERNAME_PASSWORD)
	}

	// Generate JWT Token
	refreshId := uuid.NewString()
	// Refresh Token
	refreshClaims := make(jwt.MapClaims)
	refreshClaims["user_id"] = user.Id

	refreshToken, err := util.CreateToken(refreshClaims, util.REFRESH_TOKEN_EXP_TIME, a.serverConfig.JWTSecretKey)
	if err != nil {
		log.Println("Error Login: ", err)
		return "", NewError(http.StatusBadRequest, util.ERR_TOKEN_CREATION)
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
		return "", NewError(http.StatusBadRequest, util.ERR_TOKEN_CREATION)
	}

	// Access Token
	accessToken, err := a.generateAccessToken(user.Id, refreshId)
	if err != nil {
		log.Println("Error Save Token: ", err)
		return "", NewError(http.StatusBadRequest, util.ERR_TOKEN_CREATION)
	}

	return accessToken, NoError()
}

func (a *AuthService) SignUp(input *model.SignInInput) CustomError {
	password, err := util.HashPassword(input.Password)
	if err != nil {
		return NewError(http.StatusBadRequest, util.ERR_SIGNUP)
	}

	user := model.User{
		Name:     input.Username,
		Password: password,
	}
	err = a.userRepos.CreateUser(&user)
	if err != nil {
		return NewError(http.StatusBadRequest, util.ERR_SIGNUP)
	}
	return NoError()
}

func (a *AuthService) Logout(userId string, refreshId string) CustomError {
	token, err := a.authRepos.FindTokenById(refreshId)
	if err != nil {
		log.Println(err)
		return NewError(http.StatusBadRequest, util.ERR_LOGOUT)
	}
	if token.UserId != userId {
		return NewError(http.StatusBadRequest, util.ERR_LOGOUT)
	}
	if err := a.authRepos.RemoveTokenById(refreshId); err != nil {
		log.Println("Logout Error: ", err)
		return NewError(http.StatusBadRequest, util.ERR_LOGOUT)
	}
	return NoError()
}

func (a *AuthService) LogoutAllDevice(userId string) CustomError {
	if err := a.authRepos.RemoveTokensByUserId(userId); err != nil {
		log.Println("Logout Error: ", err)
		return NewError(http.StatusBadRequest, util.ERR_LOGOUT)
	}
	return NoError()
}

func (a *AuthService) RefreshToken(accessToken string) (string, CustomError) {

	// Parse token
	accessJwtToken, err := util.ParseToken(accessToken, false, a.serverConfig.JWTKeyFunc)
	if err != nil {
		log.Println("Error Parse Access Token: ", err)
		return "", NewError(http.StatusBadRequest, util.ERR_TOKEN_FORMAT)
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
		return "", NewError(http.StatusBadRequest, util.ERR_TOKEN_NO_OWNER)
	}

	// Get Token by refreshId
	refreshToken, err := a.authRepos.FindTokenById(refreshId)
	if err != nil {
		log.Println("Error Find Token: ", err)
		return "", NewError(http.StatusBadRequest, util.ERR_TOKEN_REFRESH)
	}

	// Check relation access token to refresh token
	if refreshToken.Id != refreshId {
		return "", NewError(http.StatusBadRequest, util.ERR_TOKEN_NO_OWNER)
	}

	// Validate refresh token
	err = util.ValidateToken(refreshToken.Token, a.serverConfig.JWTKeyFunc)
	if err != nil {
		log.Println("Error Parse Refresh Token: ", err)
		return "", NewError(http.StatusBadRequest, util.ERR_TOKEN_FORMAT)
	}

	// Generate new access token
	accessToken, err = a.generateAccessToken(user.Id, refreshId)
	if err != nil {
		log.Println("Error Generate Token: ", err)
		return "", NewError(http.StatusBadRequest, util.ERR_TOKEN_CREATION)
	}

	// TODO: Maybe need to rotate refresh token

	return accessToken, NoError()
}

func (a *AuthService) generateAccessToken(userId string, refreshId string) (string, error) {
	accessClaims := make(jwt.MapClaims)
	accessClaims["user_id"] = userId
	accessClaims["refresh_id"] = refreshId // For exact refresh_id for multiple login
	return util.CreateToken(accessClaims, util.ACCESS_TOKEN_EXP_TIME, a.serverConfig.JWTSecretKey)
}
