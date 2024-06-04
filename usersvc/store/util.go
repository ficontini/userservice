package store

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/ficontini/user-search/types"
)

const (
	numUsers       = 20
	maxLengthPhone = 9
)

var (
	cities = []string{
		"Phoenix",
		"LA",
		"San Diego",
		"San Antonio",
		"Dallas",
	}
	names = []string{
		"Bob",
		"Steve",
		"Grace",
		"Eve",
		"Frank",
	}
)

func randomName() string {
	return names[rand.Intn(len(names))]
}
func randomCity() string {
	return cities[rand.Intn(len(cities))]
}
func randomPhone() string {
	var sb strings.Builder
	for i := 0; i < maxLengthPhone; i++ {
		fmt.Fprintf(&sb, "%d", rand.Intn(10))
	}
	return sb.String()
}

func GenerateRandomUsers() []*types.User {
	var users []*types.User
	for i := 0; i < numUsers; i++ {
		users = append(users,
			GenerateRandomUser(i))
	}
	return users
}

func GenerateRandomUser(id int) *types.User {
	return types.NewUser(
		int64(id),
		randomName(),
		randomCity(),
		randomPhone(),
		rand.Float64()+5,
		rand.Intn(2) == 0)
}
