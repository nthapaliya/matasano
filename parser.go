package matasano

import (
	"errors"
	"strconv"
	"strings"
)

type pair struct {
	key, value string
}

var UserProfiles = make(map[string]string)

func eatMeta(r rune) rune {
	if r == '&' || r == '=' {
		return -1
	}
	return r
}

func Parse(input string) ([]pair, error) {
	var m []pair

	for _, v := range strings.Split(input, "&") {
		keyvalpair := strings.Split(v, "=")
		if len(keyvalpair) != 2 {
			return nil, errors.New("improper input to parser")
		}
		key := strings.Map(eatMeta, keyvalpair[0])
		val := strings.Map(eatMeta, keyvalpair[1])
		m = append(m, pair{key: key, value: val})
	}
	if m == nil {
		return nil, errors.New("input string returning empty data, check again")
	}
	return m, nil
}

func Encode(pairs []pair) string {
	out := make([]string, len(pairs))
	for i, p := range pairs {
		out[i] = strings.Map(eatMeta, p.key) + "=" +
			strings.Map(eatMeta, p.value)
	}
	return strings.Join(out, "&")
}

func ProfileFor(email string) string {
	email = strings.Map(eatMeta, email)
	v, ok := UserProfiles[email]
	if ok {
		return v
	}
	newProfile := []pair{
		{"email", email},
		{"user", strconv.Itoa(RND.Int())},
		{"role", "user"},
	}
	UserProfiles[email] = Encode(newProfile)
	// fmt.Println(newProfile)
	return UserProfiles[email]
}
