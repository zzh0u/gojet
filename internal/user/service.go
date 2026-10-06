package user

import (
	"log/slog"

	"gojet/internal/infra/apperror"
)

// Service 用户业务逻辑。
type Service struct {
	repo *UserRepository
}

// New 创建用户服务。
func New(repo *UserRepository) *Service {
	return &Service{repo: repo}
}

// CreateUser 使用请求信息创建用户。
func (s *Service) CreateUser(req CreateUserRequest) (*UserResponse, error) {
	existingUser, err := s.repo.GetUserByUserName(req.Username)
	if err != nil && err.Error() != apperror.RecordNotFound {
		return nil, apperror.Wrap(err, 500, apperror.DBQueryError)
	}
	if existingUser != nil {
		return nil, apperror.New(409, apperror.UserNameExists)
	}

	existingUserByEmail, err := s.repo.GetUserByEmail(req.Email)
	if err != nil && err.Error() != apperror.RecordNotFound {
		return nil, apperror.Wrap(err, 500, apperror.DBQueryError)
	}
	if existingUserByEmail != nil {
		return nil, apperror.New(409, apperror.EmailExists)
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		slog.Error("密码哈希失败", "username", req.Username, "error", err)
		return nil, apperror.Wrap(err, 500, "密码加密失败")
	}

	user := &User{
		Username: req.Username,
		NickName: req.NickName,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err = s.repo.Create(user); err != nil {
		slog.Error("创建用户失败", "用户", user.Username, "error", err)
		return nil, apperror.Wrap(err, 500, apperror.UserCreateFailed)
	}

	slog.Info("创建用户成功", "id", user.ID, "username", user.Username)

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		NickName: user.NickName,
		Email:    user.Email,
	}, nil
}

// CreateInitialData 创建初始用户数据。
func (s *Service) CreateInitialData() error {
	existingUsers, err := s.repo.GetAll()
	if err != nil {
		return apperror.Wrap(err, 500, "检查现有数据失败")
	}
	if len(existingUsers) > 0 {
		slog.Info("初始数据已存在，跳过插入")
		return nil
	}

	users := []*User{
		{Username: "包子", NickName: "包子", Password: "123456", Email: "baozi@example.com"},
		{Username: "玉米", NickName: "玉米", Password: "123456", Email: "corn@example.com"},
		{Username: "花卷", NickName: "花卷", Password: "123456", Email: "flower@example.com"},
		{Username: "吐司", NickName: "吐司", Password: "123456", Email: "toast@example.com"},
	}

	for _, user := range users {
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			slog.Error("密码哈希失败", "username", user.Username, "error", err)
			return apperror.Wrap(err, 500, "密码哈希失败")
		}
		user.Password = hashedPassword
	}

	if err = s.repo.CreateBatch(users); err != nil {
		slog.Error("创建初始数据失败", "error", err)
		return apperror.Wrap(err, 500, apperror.DBInsertError)
	}

	slog.Info("初始数据创建成功", "count", len(users))
	return nil
}

// GetAllUsers 获取所有用户。
func (s *Service) GetAllUsers() ([]*UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, apperror.Wrap(err, 500, "获取用户列表失败")
	}

	userResList := make([]*UserResponse, len(users))
	for i, user := range users {
		userResList[i] = &UserResponse{
			ID:       user.ID,
			Username: user.Username,
			NickName: user.NickName,
			Email:    user.Email,
		}
	}
	return userResList, nil
}

// GetUserByID 根据 ID 获取用户。
func (s *Service) GetUserByID(id int) (*UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		NickName: user.NickName,
		Email:    user.Email,
	}, nil
}

// UpdateUser 更新用户信息。
func (s *Service) UpdateUser(id int, req UpdateUserRequest) (*UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user.Username != req.Username {
		existingUser, err := s.repo.GetUserByUserName(req.Username)
		if err != nil && err.Error() != apperror.RecordNotFound {
			return nil, apperror.Wrap(err, 500, apperror.DBQueryError)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, apperror.New(409, apperror.UserNameExists)
		}
	}

	if user.Email != req.Email {
		existingUserByEmail, err := s.repo.GetUserByEmail(req.Email)
		if err != nil && err.Error() != apperror.RecordNotFound {
			return nil, apperror.Wrap(err, 500, apperror.DBQueryError)
		}
		if existingUserByEmail != nil && existingUserByEmail.ID != id {
			return nil, apperror.New(409, apperror.EmailExists)
		}
	}

	user.Username = req.Username
	user.NickName = req.NickName
	user.Email = req.Email

	if err = s.repo.Update(user); err != nil {
		slog.Error("更新用户失败", "id", id, "error", err)
		return nil, apperror.Wrap(err, 500, apperror.UserUpdateFailed)
	}

	slog.Info("更新用户成功", "id", id, "name", req.Username)

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		NickName: user.NickName,
		Email:    user.Email,
	}, nil
}

// DeleteUser 删除用户。
func (s *Service) DeleteUser(id int) error {
	if err := s.repo.Delete(id); err != nil {
		slog.Error("删除用户失败", "id", id, "error", err)
		return apperror.Wrap(err, 500, apperror.UserDeleteFailed)
	}
	slog.Info("删除用户成功", "id", id)
	return nil
}
