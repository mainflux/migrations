package export

type metadata map[string]interface{}

type pageMetadata struct {
	Total  uint64
	Offset uint64 `json:"offset,omitempty"`
	Limit  uint64 `json:"limit,omitempty"`
}

// user is a user entity.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type user struct {
	ID       string   `json:"id,omitempty"` // since version 0.10.0 and 0.11.0 do not have ID
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Metadata metadata `json:"metadata"`
}

// usersPage is a page of users.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type usersPage struct {
	pageMetadata
	Users []user
}

// dbUser is a user entity in the database.
// this is for version 0.10.0, 0.11.0, 0.12.0 and 0.13.0.
type dbUser struct {
	ID       string `db:"id,omitempty"` // since version 0.10.0 and 0.11.0 do not have ID
	Email    string `db:"email"`
	Password string `db:"password"`
	Metadata []byte `db:"metadata"`
}

