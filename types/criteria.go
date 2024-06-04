package types

type SearchCriteria interface {
	Meets(*User) bool
	Name() string
	City() string
}

type SearchField string

const (
	Name SearchField = "name"
	City SearchField = "city"
)

type searchCriteria struct {
	criteria map[SearchField]string
}

func NewSearchCriteria(name, city string) SearchCriteria {
	return &searchCriteria{
		criteria: map[SearchField]string{
			Name: name,
			City: city,
		},
	}
}
func (c *searchCriteria) Name() string {
	return c.criteria[Name]
}
func (c *searchCriteria) City() string {
	return c.criteria[City]
}
func (c *searchCriteria) Meets(u *User) bool {
	if !isEmpty(c.Name()) && !isEquals(u.FName, c.Name()) {
		return false
	}
	if !isEmpty(c.City()) && !isEquals(u.City, c.City()) {
		return false
	}
	return true
}

func isEmpty(value string) bool {
	return value == ""
}
func isEquals(actual, expected string) bool {
	return actual == expected
}
