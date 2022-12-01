package postgres

import (
	"database/sql"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

type Conn interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Rebind(sql string) string
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
}

type Repo struct {
	DB  *sqlx.DB
	Ctx context.Context

	tx *sqlx.Tx
}

func AllStatements(script string) (ret []string) {
	for _, s := range strings.Split(removeComments(script), ";") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		ret = append(ret, s)
	}
	return
}

func Open(url string) (*Repo, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, errors.Wrap(err, 1)
	}

	return &Repo{DB: db}, nil
}

// ensure various types conform to Conn interface
var _ Conn = (*sqlx.Tx)(nil)
var _ Conn = (*sqlx.DB)(nil)

// SQLBlockComments is a regex that matches against SQL block comments
var sqlBlockComments = regexp.MustCompile(`/\*.*?\*/`)

// SQLLineComments is a regex that matches against SQL line comments
var sqlLineComments = regexp.MustCompile("--.*?\n")

func removeComments(script string) string {
	withoutBlocks := sqlBlockComments.ReplaceAllString(script, "")
	return sqlLineComments.ReplaceAllString(withoutBlocks, "")
}
