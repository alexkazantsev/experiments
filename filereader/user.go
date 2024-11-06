package main

import "strconv"

type User struct {
	first  string
	last   string
	age    uint8
	domain string
	email  string
}

func (u User) filter() bool {
	return u.age <= 20
}

func (u User) toRecord() []string {
	return []string{
		u.first,
		u.last,
		strconv.FormatUint(uint64(u.age), 10),
		u.domain,
		u.email,
	}
}
