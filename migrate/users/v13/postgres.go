package users13

import (
	"context"
	"encoding/json"

	// required for SQL access
	mf13errors "github.com/mainflux/mainflux/pkg/errors"
	mf13users "github.com/mainflux/mainflux/users"
	mf13postgres "github.com/mainflux/mainflux/users/postgres"
)

type db13User struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Metadata []byte `db:"metadata"`
}

// dbRetrieveUsers retrieves users from the database with the given page navigation parameters
func dbRetrieveUsers(ctx context.Context, db mf13postgres.Database, pm mf13users.PageMetadata) (mf13users.UserPage, error) {
	q := `SELECT id, email, password, metadata FROM users LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return mf13users.UserPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []mf13users.User
	for rows.Next() {
		dbu := db13User{}
		if err := rows.StructScan(&dbu); err != nil {
			return mf13users.UserPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		user, err := to13User(dbu)
		if err != nil {
			return mf13users.UserPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		items = append(items, user)
	}

	cq := `SELECT COUNT(*) FROM users;`

	total, err := total(ctx, db, cq, params)
	if err != nil {
		return mf13users.UserPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}

	page := mf13users.UserPage{
		Users: items,
		PageMetadata: mf13users.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

func total(ctx context.Context, db mf13postgres.Database, query string, params interface{}) (uint64, error) {
	rows, err := db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}
	return total, nil
}

func to13User(u db13User) (mf13users.User, error) {
	var metadata map[string]interface{}
	if u.Metadata != nil {
		if err := json.Unmarshal([]byte(u.Metadata), &metadata); err != nil {
			return mf13users.User{}, mf13errors.Wrap(mf13errors.ErrMalformedEntity, err)
		}
	}
	return mf13users.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Metadata: metadata,
	}, nil
}
