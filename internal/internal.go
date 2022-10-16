package internal

var id int64 = 0

func GetId() int64 {
	id += 1
	return id
}
