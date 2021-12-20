package user

var (
	Queue             = "user"
	UpdateUserSubject = "faceit-user-updateUser"
)

// User type used to define queue messages
type User struct {
	ID        string
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string
	Country   string
	CreatedAt string
	UpdatedAt string
}
