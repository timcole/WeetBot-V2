package main

type User struct {
	ID          int64
	Login       string
	DisplayName string
	Followers   int64
	Views       int64
	Type        bool
	Join        bool
}
