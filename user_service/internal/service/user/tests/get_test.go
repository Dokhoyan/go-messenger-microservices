package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/kafka/producer"
	logsRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/logs"
	userRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user"
	userservice "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/user"
	"go.uber.org/mock/gomock"

	// "github.com/a1exCross/common/pkg/client/db"
	dbmocks "github.com/Dokhoyan/common/pkg/client/db/mocks"
	"github.com/Dokhoyan/common/pkg/client/db/transaction"

	// "github.com/a1exCross/common/pkg/client/db/transaction"
	// "github.com/a1exCross/common/pkg/storage"
	storagemocks "github.com/Dokhoyan/common/pkg/storage/mocks"

	// "github.com/gojuno/minimock/v3"
	"github.com/Dokhoyan/common/pkg/client/db"
	"github.com/Dokhoyan/common/pkg/storage"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	type dbClientMock func(mc *gomock.Controller) db.Client
	type txManagerMock func(mc *gomock.Controller) db.TxManager
	type storageMock func(mc *gomock.Controller) storage.Redis

	ctx := context.Background()
	mc := gomock.NewController(t)
	id := int64(1)
	timeNow := time.Now()

	user := &model.User{
		ID: id,
		Info: model.UserInfo{
			Email:      "email",
			Role:       model.UserRole(1),
			Name:       "name",
			Username:   "username",
			Avatar_url: "avatar_url",
			Birth_date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
		},
		UpdatedAt: sql.NullTime{
			Valid: true,
			Time:  timeNow,
		},
		CreatedAt: timeNow,
	}

	tests := []struct {
		name      string
		err       error
		expected  *model.User
		dbClient  dbClientMock
		txManager txManagerMock
		storageMock
	}{
		{
			name: "successful test",
	dbClient: func(mc *gomock.Controller) db.Client {
		client := dbmocks.NewMockClient(mc)
		dbb := dbmocks.NewMockDB(mc)
		row := dbmocks.NewMockRow(mc)

		// row.Scan — возвращает ID пользователя
		row.EXPECT().
			Scan(gomock.Any()).
			DoAndReturn(func(dest ...interface{}) error {
				if res, ok := dest[0].(*int64); ok {
					*res = id
				}
				return nil
			}).
			AnyTimes()

		// QueryRowContext — возвращает мок строки
		dbb.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(row).
			AnyTimes()

		// ScanOneContext — заполняет структуру пользователя
		dbb.EXPECT().
			ScanOneContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
				if res, ok := dest.(*model.User); ok {
					res.ID = user.ID
					res.Info = user.Info
					res.CreatedAt = user.CreatedAt
					res.UpdatedAt = user.UpdatedAt
				}
				return nil
			}).
			AnyTimes()

		// DB() возвращает мокнутый dbb
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
			
				// Ожидаем успешный Commit транзакции
				tx.EXPECT().
					Commit(gomock.Any()).
					Return(nil).
					AnyTimes()
			
				// Ожидаем вызов BeginTx — возвращаем мок транзакции без ошибок
				transactor.EXPECT().
					BeginTx(gomock.Any(), txOptions).
					Return(tx, nil).
					AnyTimes()
			
				// Создаём менеджер транзакций на основе мокнутого транзактора
				txManager := transaction.NewTransactionManager(transactor)
			
				return txManager
			},
			storageMock: func(mc *gomock.Controller) storage.Redis {
				mock := storagemocks.NewMockRedis(mc)

				return mock
			},
			err: nil,
			expected: &model.User{
				ID: id,
				Info: model.UserInfo{
					Email:      "email",
			Role:       model.UserRole(1),
			Name:       "name",
			Username:   "username",
			Avatar_url: "avatar_url",
			Birth_date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
				},
				UpdatedAt: sql.NullTime{
					Valid: true,
					Time:  timeNow,
				},
				CreatedAt: timeNow,
			},
		},
		{
			name: "error at get",
			dbClient: func(mc *gomock.Controller) db.Client {
				client := dbmocks.NewMockClient(mc)

				return client
			},
			txManager: func(mc *gomock.Controller) db.TxManager {
				transactor := dbmocks.NewMockTransactor(mc)
			
				txOptions := pgx.TxOptions{
					IsoLevel: pgx.ReadCommitted,
				}
			
				// Ожидаем, что BeginTx вернёт ошибку при попытке начать транзакцию
				transactor.EXPECT().
					BeginTx(gomock.Any(), txOptions).
					Return(nil, errors.New("tx error")).
					AnyTimes()
			
				// Создаём менеджер транзакций на основе мокнутого транзактора
				txManager := transaction.NewTransactionManager(transactor)
			
				return txManager
			},
			storageMock: func(mc *gomock.Controller) storage.Redis {
				mock := storagemocks.NewMockRedis(mc)

				return mock
			},
			err:      errors.New("can`t begin transaction: tx error"),
			expected: nil,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			dbMockClient := test.dbClient(mc)
			txManager := test.txManager(mc)
			redis := test.storageMock(mc)

			userRepo := userRepository.NewRepository(dbMockClient)
			logRepo := logsRepository.NewRepository(dbMockClient)

			br := []string{"localhost:9092","localhost:9093","localhost:9094"}
			producer , err := producer.NewProducer(br)
			userServ := userservice.NewService(userRepo, txManager, logRepo, redis, producer)

			res, err := userServ.Get(ctx, id)

			require.Equal(t, test.expected, res)

			if err != nil && test.err != nil {
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}