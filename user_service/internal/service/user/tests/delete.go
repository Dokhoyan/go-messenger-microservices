package tests

import (
	"context"
	"testing"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/kafka/producer"
	logsRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/logs"
	userRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user"
	userservice "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/user"
	"go.uber.org/mock/gomock"

	"github.com/Dokhoyan/common/pkg/client/db"
	dbmocks "github.com/Dokhoyan/common/pkg/client/db/mocks"
	"github.com/Dokhoyan/common/pkg/client/db/transaction"
	"github.com/Dokhoyan/common/pkg/storage"
	storagemocks "github.com/Dokhoyan/common/pkg/storage/mocks"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	type dbClientMock func(mc *gomock.Controller) db.Client
	type txManagerMock func(mc *gomock.Controller) db.TxManager
	type storageMock func(mc *gomock.Controller) storage.Redis

	ctx := context.Background()
	id := int64(1)

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

				row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
					res, ok := dest[0].(*int64)
					if ok {
						*res = id
					}
					return nil
				}).AnyTimes()

				dbb.EXPECT().QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(row).AnyTimes()
				dbb.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil).AnyTimes()
				client.EXPECT().DB().Return(dbb).AnyTimes()

				return client
			},
			txManager: func(mc *gomock.Controller) db.TxManager {
				tx := dbmocks.NewMockTx(mc)
				transactor := dbmocks.NewMockTransactor(mc)

				txOptions := pgx.TxOptions{
					IsoLevel: pgx.ReadCommitted,
				}

				tx.EXPECT().Commit(gomock.Any()).Return(nil).AnyTimes()
				transactor.EXPECT().BeginTx(gomock.Any(), txOptions).Return(tx, nil).AnyTimes()

				return transaction.NewTransactionManager(transactor)
			},
			storageMock: func(mc *gomock.Controller) storage.Redis {
				return storagemocks.NewMockRedis(mc)
			},
			err: nil,
		},
		{
			name: "error at delete (tx begin failed)",
			dbClient: func(mc *gomock.Controller) db.Client {
				client := dbmocks.NewMockClient(mc)
				return client
			},
			txManager: func(mc *gomock.Controller) db.TxManager {
				transactor := dbmocks.NewMockTransactor(mc)
				txOptions := pgx.TxOptions{
					IsoLevel: pgx.ReadCommitted,
				}

				transactor.EXPECT().BeginTx(gomock.Any(), txOptions).Return(nil, errors.New("tx error")).AnyTimes()

				return transaction.NewTransactionManager(transactor)
			},
			storageMock: func(mc *gomock.Controller) storage.Redis {
				return storagemocks.NewMockRedis(mc)
			},
			err: errors.New("can`t begin transaction: tx error"),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			dbMockClient := test.dbClient(mc)
			txManager := test.txManager(mc)
			redis := test.storageMock(mc)

			userRepo := userRepository.NewRepository(dbMockClient)
			logRepo := logsRepository.NewRepository(dbMockClient)

			br := []string{"localhost:9092","localhost:9093","localhost:9094"}
			producer , err := producer.NewProducer(br)
			userServ := userservice.NewService(userRepo, txManager, logRepo, redis, producer)

			err = userServ.Delete(ctx, id)

			if err != nil && test.err != nil {
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}
