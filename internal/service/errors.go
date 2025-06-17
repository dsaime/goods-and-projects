package service

import "errors"

var (
	ErrIDMustBeGTZero        = errors.New("id должен быть больше нуля")
	ErrProjectIDMustBeGTZero = errors.New("projectID должен быть больше нуля")
	ErrPriorityMustBeGTZero  = errors.New("priority должен быть больше нуля")
	ErrNameMustNotBeEmpty    = errors.New("name не может быть пустым")
	ErrLimitMustBePositive   = errors.New("limit должен быть равен или больше нуля")
	ErrOffsetMustBePositive  = errors.New("offset должен быть равен или больше нуля")
)
