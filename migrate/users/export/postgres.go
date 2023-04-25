package export

import (
	"context"
	"encoding/json"

	mf13errors "github.com/mainflux/mainflux/pkg/errors"
	"github.com/mainflux/mainflux/users/postgres" // for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0
)

// dbRetrieveUsers retrieves users from the database with the given page navigation parameters.
func dbRetrieveUsers(ctx context.Context, query string, db postgres.Database, pm pageMetadata) (usersPage, error) {
	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return usersPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []user
	for rows.Next() {
		dbu := dbUser{}
		if err := rows.StructScan(&dbu); err != nil {
			return usersPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		user, err := toUser(dbu)
		if err != nil {
			return usersPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
		}

		items = append(items, user)
	}

	cq := `SELECT COUNT(*) FROM users;`

	total, err := total(ctx, db, cq, params)
	if err != nil {
		return usersPage{}, mf13errors.Wrap(mf13errors.ErrViewEntity, err)
	}

	page := usersPage{
		Users: items,
		pageMetadata: pageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

func total(ctx context.Context, db postgres.Database, query string, params interface{}) (uint64, error) {
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

// toUser converts dbUser to user.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
func toUser(dbUser dbUser) (user, error) {
	var metadata map[string]interface{}
	if dbUser.Metadata != nil {
		if err := json.Unmarshal(dbUser.Metadata, &metadata); err != nil {
			return user{}, mf13errors.Wrap(mf13errors.ErrMalformedEntity, err)
		}
	}

	return user{
		ID:       dbUser.ID,
		Email:    dbUser.Email,
		Password: dbUser.Password,
		Metadata: metadata,
	}, nil
}
