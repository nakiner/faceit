package user

import (
	"fmt"
	"time"
)

type Conditions map[string]interface{}

type User struct {
	ID        string    `gorm:"primaryKey,size:64"`
	FirstName string    `gorm:"size:64"`
	LastName  string    `gorm:"size:64"`
	Nickname  string    `gorm:"size:64"`
	Password  string    `gorm:"size:64"`
	Email     string    `gorm:"size:64"`
	Country   string    `gorm:"size:64"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
}

func (User) TableName() string {
	return "users"
}

func (User) TimeToString(timeVal time.Time) string {
	return timeVal.Format("2006-01-02T15:04:05.999999999")
}

// GetPreparedStatement performs parse operation to get conditions in string from slice of GET params
// argument comes validated to this function
func (c Conditions) GetPreparedStatement() (res string) {
	var sliced []string
	for key, _ := range c {
		sliced = append(sliced, fmt.Sprintf("%s = @%s", key, key))
	}

	if len(sliced) == 1 {
		res = sliced[0]
		return
	}

	for k, cond := range sliced {
		if k == len(sliced)-1 {
			res = fmt.Sprintf("%s %s", res, cond)
		} else if k == 0 {
			res = fmt.Sprintf("%s AND", cond)
		} else {
			res = fmt.Sprintf("%s %s AND", res, cond)
		}
	}

	return
}
