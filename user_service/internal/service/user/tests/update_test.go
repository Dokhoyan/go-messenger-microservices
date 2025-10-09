package tests

import (
	"context"
	"testing"
	"time"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	logsRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/logs"
	userRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user"
	userservice "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/user"

	"github.com/Dokhoyan/common/pkg/client/db"
	dbmocks "github.com/Dokhoyan/common/pkg/client/db/mocks"
	"github.com/Dokhoyan/common/pkg/client/db/transaction"
	"github.com/Dokhoyan/common/pkg/storage"
	storagemocks "github.com/Dokhoyan/common/pkg/storage/mocks"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdate(t *testing.T) {
	type dbClientMock func(mc *gomock.Controller) db.Client
	type txManagerMock func(mc *gomock.Controller) db.TxManager
	type storageMock func(mc *gomock.Controller) storage.Redis

	ctx := context.Background()
	id := int64(1)

	userDTO := &model.UserUpdate{
		Info: model.UserInfo{
			Email:      "email",
			Role:       model.UserRole(1),
			Name:       "name",
			Username:   "username",
			Avatar_url: "avatar_url",
			Birth_date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	tests := []struct {
		name        string
		err         error
		dbClient    dbClientMock
		txManager   txManagerMock
		storageMock storageMock
	}{
		{
			name: "successful test",
			dbClient: func(mc *gomock.Controller) db.Client {
				client := dbmocks.NewMockClient(mc)
				dbb := dbmocks.NewMockDB(mc)
				row := dbmocks.NewMockRow(mc)

				// row.Scan — возвращает id обновлённого пользователя
				row.EXPECT().
					Scan(gomock.Any()).
					DoAndReturn(func(dest ...interface{}) error {
						if res, ok := dest[0].(*int64); ok {
							*res = id
						}
						return nil
					}).
					AnyTimes()

				// QueryRowContext — симулируем UPDATE ... RETURNING id
				dbb.EXPECT().
					QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(row).
					AnyTimes()

				// ExecContext — успешное выполнение update
				dbb.EXPECT().
					ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(pgconn.CommandTag{}, nil).
					AnyTimes()

				client.EXPECT().
					DB().
					Return(dbb).
					AnyTimes()

				return client
			},
			txManager: func(mc *gomock.Controller) db.TxManager {
				tx := dbmocks.NewMockTx(mc)
				transactor := dbmocks.NewMockTransactor(mc)

				txOptions := pgx.TxOptions{
					IsoLevel: pgx.ReadCommitted,
				}

				// Commit проходит успешно
				tx.EXPECT().
					Commit(gomock.Any()).
					Return(nil).
					AnyTimes()

				// Транзакция начинается без ошибок
				transactor.EXPECT().
					BeginTx(gomock.Any(), txOptions).
					Return(tx, nil).
					AnyTimes()

				return transaction.NewTransactionManager(transactor)
			},
			storageMock: func(mc *gomock.Controller) storage.Redis {
				return storagemocks.NewMockRedis(mc)
			},
			err: nil,
		},
		{
			name: "error at update (commit failed)",
			dbClient: func(mc *gomock.Controller) db.Client {
				client := dbmocks.NewMockClient(mc)
				dbb := dbmocks.NewMockDB(mc)
				row := dbmocks.NewMockRow(mc)

				row.EXPECT().
					Scan(gomock.Any()).
					DoAndReturn(func(dest ...interface{}) error {
						if res, ok := dest[0].(*int64); ok {
							*res = id
						}
						return nil
					}).
					AnyTimes()

				dbb.EXPECT().
					QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(row).
					AnyTimes()

				dbb.EXPECT().
					ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(pgconn.CommandTag{}, nil).
					AnyTimes()

				client.EXPECT().
					DB().
					Return(dbb).
					AnyTimes()

				return client
			},
			txManager: func(mc *gomock.Controller) db.TxManager {
				tx := dbmocks.NewMockTx(mc)
				transactor := dbmocks.NewMockTransactor(mc)

				txOptions := pgx.TxOptions{
					IsoLevel: pgx.ReadCommitted,
				}

				// Commit возвращает ошибку
				tx.EXPECT().
					Commit(gomock.Any()).
					Return(errors.New("commit error")).
					AnyTimes()

				transactor.EXPECT().
					BeginTx(gomock.Any(), txOptions).
					Return(tx, nil).
					AnyTimes()

				return transaction.NewTransactionManager(transactor)
			},
			storageMock: func(mc *gomock.Controller) storage.Redis {
				return storagemocks.NewMockRedis(mc)
			},
			err: errors.New("tx commit failed: commit error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			dbMockClient := test.dbClient(mc)
			txManager := test.txManager(mc)
			redis := test.storageMock(mc)

			userRepo := userRepository.NewRepository(dbMockClient)
			logRepo := logsRepository.NewRepository(dbMockClient)

			userServ := userservice.NewService(userRepo, txManager, logRepo, redis)

			err := userServ.Update(ctx, userDTO)

			if err != nil && test.err != nil {
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}
