package user

import (
	"context"

	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error){
	// err := i.accessService.Check(ctx, "/user_v1.UserV1/Delete")
	// if err != nil {
	// 	return nil, err
	// }

	err := i.userservice.Delete(ctx, req.Id)
	if err != nil{
		return nil, errors.Errorf("failed to delete user: %v", err)
	}

	return &emptypb.Empty{}, nil
}