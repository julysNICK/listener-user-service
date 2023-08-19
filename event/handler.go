package event

import (
	"context"
	userGRPC "listener-user-service/users"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ResponseGetAllUsers struct {
	Users []userGRPC.User
}

func GetAllUsersViaGRPC() error {

	conn, err := grpc.Dial("user-service:5003", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		return err
	}

	defer conn.Close()

	u := userGRPC.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := u.GetAllUsers(ctx, &emptypb.Empty{})

	if err != nil {
		return err
	}

	for _, v := range res.Users {
		println(v.FirstName)
	}

	return nil

}

type RequestGetOneUser struct {
	Id string `json:"id"`
}

func GetUserViaGRPC(payload RequestGetOneUser) error {
	conn, err := grpc.Dial("user-service:5003", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		return err
	}
	defer conn.Close()
	u := userGRPC.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	convID, err := strconv.Atoi(payload.Id)

	if err != nil {

		return err
	}
	_, err = u.GetOneUser(ctx, &userGRPC.UserRequestGetOne{
		Id: int32(convID),
	})
	if err != nil {

		return err
	}
	return nil
}

type RequestCreateUser struct {
	Email string `json:"email"`
}

func UserDeleteViaGRPC(payload RequestCreateUser) error {
	conn, err := grpc.Dial("user-service:5003", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		return err
	}

	defer conn.Close()

	u := userGRPC.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	_, err = u.DeleteUser(ctx, &userGRPC.UserRequestDelete{
		Email: payload.Email,
	})

	if err != nil {
		return err
	}

	return nil
}

type UserUpdateViaGRPCPayload struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func UserUpdateViaGRPC(payload UserUpdateViaGRPCPayload) error {
	conn, err := grpc.Dial("user-service:5003", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {

		return err
	}

	defer conn.Close()

	u := userGRPC.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	_, err = u.UpdateUser(ctx, &userGRPC.UserRequestUpdate{
		Email: payload.Email,
	})

	if err != nil {
		return err
	}

	return nil
}
