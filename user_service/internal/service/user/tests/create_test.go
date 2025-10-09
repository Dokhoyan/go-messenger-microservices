package tests

import (
	"context"
	"testing"
	"time"

	logsrepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/logs"
	userRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user"
	userService "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/user"
	"github.com/Dokhoyan/common/pkg/client/db"
	dbmocks "github.com/Dokhoyan/common/pkg/client/db/mocks"
	"github.com/Dokhoyan/common/pkg/client/db/transaction"
	"github.com/Dokhoyan/common/pkg/storage"
	storagemocks "github.com/Dokhoyan/common/pkg/storage/mocks"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	type dbClientMock func(mc *gomock.Controller) db.Client
	type txManagerMock func(mc *gomock.Controller) db.TxManager
	type storageMock func(mc *gomock.Controller) storage.Redis

	ctx := context.Background()
	id := int64(1)

	userDTO := &model.UserCreate{
		Info: model.UserInfo{
			Email:      "email",
			Role:       model.UserRole(1),
			Name:       "name",
			Username:   "username",
			Avatar_url: "avatar_url",
			Birth_date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
		},
		Password: "pass",
	}

	tests := []struct {
		name        string
		err         error
		expected    int64
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

				// row.Scan
				row.EXPECT().
					Scan(gomock.Any()).
					DoAndReturn(func(dest ...interface{}) error {
						if res, ok := dest[0].(*int64); ok {
							*res = id
						}
						return nil
					}).
					AnyTimes()

				// QueryRowContext
				dbb.EXPECT().
					QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(row).
					AnyTimes()

				// ScanOneContext
				dbb.EXPECT().
					ScanOneContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(pgx.ErrNoRows).
					AnyTimes()

				// DB
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

				tx.EXPECT().
					Commit(gomock.Any()).
					Return(nil).
					AnyTimes()

				transactor.EXPECT().
					BeginTx(gomock.Any(), txOptions).
					Return(tx, nil).
					AnyTimes()

				txManager := transaction.NewTransactionManager(transactor)

				return txManager
			},
			storageMock: func(mc *gomock.Controller) storage.Redis {
				mock := storagemocks.NewMockRedis(mc)
				return mock
			},
			err:      nil,
			expected: id,
		},
		{
			name: "error at create",
			dbClient: func(mc *gomock.Controller) db.Client {
				client := dbmocks.NewMockClient(mc)
				dbb := dbmocks.NewMockDB(mc)
				row := dbmocks.NewMockRow(mc)

				// row.Scan возвращает ошибку
				row.EXPECT().
					Scan(gomock.Any()).
					DoAndReturn(func(dest ...interface{}) error {
						return errors.New("failed to scan")
					}).
					AnyTimes()

				dbb.EXPECT().
					QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(row).
					AnyTimes()

				dbb.EXPECT().
					ScanOneContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(pgx.ErrNoRows).
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

				tx.EXPECT().
					Rollback(gomock.Any()).
					Return(nil).
					AnyTimes()

				transactor.EXPECT().
					BeginTx(gomock.Any(), txOptions).
					Return(tx, nil).
					AnyTimes()

				txManager := transaction.NewTransactionManager(transactor)

				return txManager
			},
			storageMock: func(mc *gomock.Controller) storage.Redis {
				mock := storagemocks.NewMockRedis(mc)
				return mock
			},
			err: errors.New("failed executing code inside transaction: failed to scan"),
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			dbMockClient := test.dbClient(mc)
			txManager := test.txManager(mc)
			redisMock := test.storageMock(mc)

			userRepo := userRepository.NewRepository(dbMockClient)
			logRepo := logsrepository.NewRepository(dbMockClient)

			userServ := userService.NewService(userRepo, txManager, logRepo, redisMock)

			res, err := userServ.Create(ctx, userDTO)

			require.Equal(t, test.expected, res)

			if err != nil && test.err != nil {
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}
