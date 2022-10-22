package hw10programoptimization

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

const SeparateLine byte = 10

//go:generate easyjson -all stats.go
type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
}

func (u *User) Reset() {
	u.ID = 0
	u.Name = ""
	u.Username = ""
	u.Email = ""
	u.Phone = ""
	u.Password = ""
	u.Address = ""
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	reg, errR := regexp.Compile("\\." + domain)
	if errR != nil {
		return nil, errR
	}
	return countDomains(u, reg)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content := make([]byte, 1024)
	var user User
	numU := 0
	line := make([]byte, 1024)
	line = line[:0]
	update := func() {
		if len(line) < 1 {
			return
		}
		if err = user.UnmarshalJSON(line); err != nil {
			return
		}
		result[numU] = user
		user.Reset()
		line = line[:0]
		numU++
	}
	for {
		n, errR := r.Read(content)
		if errR != nil && errR != io.EOF {
			err = errR
			return
		}
		if n < 1 {
			update()
			return
		}
		for _, item := range content[0:n] {
			if item == SeparateLine {
				update()
				continue
			}
			if err == io.EOF {
				update()
				return
			}
			line = append(line, item)
		}
	}

	return
}

func countDomains(u users, reg *regexp.Regexp) (DomainStat, error) {
	result := make(DomainStat, len(u))
	var i string
	for _, user := range u {
		if user.Email == "" {
			continue
		}
		matched := reg.Match([]byte(user.Email))
		if matched {
			i = strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[i]
			num++
			result[i] = num
		}
	}
	return result, nil
}
