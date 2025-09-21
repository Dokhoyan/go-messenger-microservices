package converter

import (
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	userPb "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserToProto - конвертирует модель пользователя в proto
func UserToProto(user *model.User) *userPb.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &userPb.User{
		Id:        user.ID,
		Info:      UserInfoToProto(user.Info),
		UpdatedAt: updatedAt,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

// UserInfoToProto - конвертирует информацию о пользователе в proto
func UserInfoToProto(info model.UserInfo) *userPb.UserInfo {
	return &userPb.UserInfo{
		Username: info.Username,
		Name:     info.Name,
		Role:     userPb.UserRole(info.Role),
		BirthDate: timestamppb.New(info.Birth_date),
		Email:    info.Email,
		AvatarUrl: info.Avatar_url,
	}
}

// ProtoToUserInfo - конвертирует proto в модель информации о пользователе
func ProtoToUserInfo(info *userPb.UserInfo) model.UserInfo {
	return model.UserInfo{
		Username: info.Username,
		Name:     info.Name,
		Role:     model.UserRole(info.Role),
		Email:    info.Email,
		Avatar_url: info.AvatarUrl,
		Birth_date: info.BirthDate.AsTime(),
	}
}

// ProtoToUserInfoUpdate - конвертирует proto в модель информации о пользователе
func ProtoToUserInfoUpdate(info *userPb.UpdateInfo) model.UserInfo {
	return model.UserInfo{
		Username: info.Username.Value,
		Name:     info.Name.Value,
		Role:     model.UserRole(info.Role),
		Email:    info.Email.Value,
		Avatar_url: info.AvatarUrl.Value,
		Birth_date: info.BirthDate.AsTime(),
	}
}