package auth

import (
	"time"

	"gojet/internal/config"
	"gojet/internal/infra/apperror"
	"gojet/internal/infra/jwt"
	"gojet/internal/user"
)

// Service 认证业务逻辑。
type Service struct {
	repo *user.UserRepository
	cfg  *config.Config
}

// New 创建认证服务。
func New(repo *user.UserRepository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

// Login 校验账号密码并签发双层 Token。
func (s *Service) Login(req LoginReq) (*LoginResp, error) {
	user, err := s.repo.GetUserByUserName(req.Username)
	if err != nil {
		return nil, apperror.Wrap(err, 404, apperror.UserNotFound)
	}

	if !user.CompareSimple(req.Password) {
		return nil, apperror.New(401, apperror.AuthFailed)
	}

	accessDuration := time.Duration(s.cfg.JWT.AccessExpireHours) * time.Hour
	accessToken, err := jwt.SignAccess(jwt.Context{ID: user.ID, Username: user.Username}, s.cfg.JWT.Secret, accessDuration)
	if err != nil {
		return nil, apperror.Wrap(err, 500, "生成 Access Token 失败")
	}

	refreshDuration := time.Duration(s.cfg.JWT.RefreshExpireDays) * 24 * time.Hour
	refreshToken, err := jwt.SignRefresh(jwt.Context{ID: user.ID, Username: user.Username}, s.cfg.JWT.Secret, refreshDuration)
	if err != nil {
		return nil, apperror.Wrap(err, 500, "生成 Refresh Token 失败")
	}

	return &LoginResp{
		Userid:       user.ID,
		Username:     user.Username,
		NickName:     user.NickName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    accessDuration.Seconds(),
	}, nil
}
