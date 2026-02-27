package test

import (
	"main/internal/domain/test"
	"main/pkg"

	"go.uber.org/fx"
)

type testRepo struct {
	db pkg.PostgresDB
}

func NewUserRepository(db pkg.PostgresDB) test.TestRepo {
	return &testRepo{db: db}
}

func (r *testRepo) GetOne() (int, error) {
	var result int
	err := r.db.QueryRow(GetOneQuery).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

var Module = fx.Options(
	fx.Provide(NewUserRepository),
)
