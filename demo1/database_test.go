package demo1

import (
	"context"
	"demo1/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"strconv"
	"testing"
)

var (
	_repo *userRepo
	//_conn *sqlx.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	postgres, postgresUrl := util.LaunchPostgres(ctx)

	_conn, err := sqlx.Connect("postgres", postgresUrl)
	if err != nil {
		log.Fatal("connect:", err)
	}

	if err := runMigrations(_conn); err != nil {
		log.Fatal("runMigrations:", err)
	}

	_repo = NewRepo(_conn)
	exitCode := m.Run()
	log.Println("exit code from tests:", exitCode)

	err = postgres.Terminate(ctx)
	if err != nil {
		log.Fatal("error shutting down Postgres container:", err)
	}
	log.Println("Postgres container shut down")

	//os.Exit(exitCode)
}

func TestRepoImp(t *testing.T) {
	t.Run("create and get single user", func(t *testing.T) {
		user, err := _repo.CreateUser("username")
		require.NoError(t, err)

		getUser, err := _repo.GetUserByID(user.ID)
		require.NoError(t, err)
		assert.Equal(t, user, getUser)
	})

	t.Run("get all users", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			_, err := _repo.CreateUser(strconv.Itoa(i))
			require.NoError(t, err)
		}
		users, err := _repo.GetAllUsers()
		require.NoError(t, err)
		assert.Len(t, users, 11) // 10 + 1 previously
	})
}
