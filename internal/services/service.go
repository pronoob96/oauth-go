package services

import (
	"context"
	"log"
	"oauth/internal/repository"
	"oauth/pkg/dto"
	m "oauth/pkg/models"
	"os"
	"time"

	r "oauth/redis"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	AddUserService(context.Context, *dto.User) (map[string]string, error)
	LoginService(context.Context, *dto.User) (map[string]string, error)
}
type userService struct {
	userRepo repository.UserRepository
}

func NewUsersService(user repository.UserRepository) UserService {
	return &userService{
		userRepo: user,
	}
}

func (s *userService) AddUserService(ctx context.Context, userRequest *dto.User) (map[string]string, error) {
	_, err := s.userRepo.GetUserByUsername(ctx, userRequest.Username)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			var user m.User
			err = copier.Copy(&user, userRequest)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			user.Password = string(pass)

			userResponse, err := s.userRepo.AddUser(ctx, &user)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			token, err := createToken(&user)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			saveErr := createAuth(userResponse.ID.Hex(), token)
			if saveErr != nil {
				log.Println(saveErr)
				return nil, saveErr
			}

			tokens := map[string]string{
				"access_token":  token.AccessToken,
				"refresh_token": token.RefreshToken,
			}

			return tokens, nil

		}
		return nil, err
	}
	err = errors.Errorf("User Already Present")
	return nil, err

}

func (s *userService) LoginService(ctx context.Context, userRequest *dto.User) (map[string]string, error) {
	var user m.User
	err := copier.Copy(&user, userRequest)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	gotUser, err := s.userRepo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(gotUser.Password), []byte(user.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println(err)
		return nil, err
	}

	token, err := createToken(gotUser)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	saveErr := createAuth(gotUser.ID.Hex(), token)
	if saveErr != nil {
		log.Println(saveErr)
		return nil, saveErr
	}

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	return tokens, nil

}

func createToken(gotUser *m.User) (*dto.TokenDetails, error) {
	td := &dto.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = gotUser.ID.Hex()
	atClaims["username"] = gotUser.Username
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf")
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	atClaims["user_id"] = gotUser.ID.Hex()
	atClaims["username"] = gotUser.Username
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func createAuth(userid string, td *dto.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := r.Client.Set(td.AccessUuid, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := r.Client.Set(td.RefreshUuid, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
