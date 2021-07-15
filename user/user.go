package user

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	Id       int
	Name     string
	Birthday string
	Tel      string
	Address  string
	Desc     string
	Password string
}

func md5Pass(pass string) string {
	ctx := md5.New()
	ctx.Write([]byte(pass))
	return hex.EncodeToString(ctx.Sum(nil))
}

func getId(users map[int]*User) int {
	var max = 0
	for k := range users {
		if k > max {
			max = k
		}
	}
	return max + 1
}

func (u User) ValidatePassword(p string) bool {
	return u.Password == md5Pass(p)
}

func loadUsers(file string) (map[int]*User, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	reader := csv.NewReader(f)
	if _, err := reader.Read(); err != nil {
		return nil, err
	}
	records, err := reader.ReadAll()
	users := map[int]*User{}
	for _, v := range records {
		id, err := strconv.Atoi(v[0])
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		name, birthday, address, tel, desc, password := v[1], v[2], v[3], v[4], v[5], v[6]

		users[id] = &User{
			Id:       id,
			Name:     name,
			Birthday: birthday,
			Address:  address,
			Tel:      tel,
			Desc:     desc,
			Password: password,
		}
	}
	return users, err
}

func saveUsers(users map[int]*User, file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	var records [][]string
	records = append(records, []string{"id", "name", "birthday", "address", "tel", "desc", "password"})
	for _, u := range users {
		record := []string{
			strconv.Itoa(u.Id),
			u.Name,
			u.Birthday,
			u.Address,
			u.Tel,
			u.Desc,
			u.Password,
		}
		records = append(records, record)
	}
	writer := csv.NewWriter(f)
	err = writer.WriteAll(records)
	return err
}

func GetUser(userId int) *User {
	users, err := loadUsers("user.csv")
	if err != nil {
		panic(err)
	}

	if user, ok := users[userId]; ok {
		return user
	} else {
		return nil
	}
}

func GetUsers(args ...string) map[int]*User {
	users, err := loadUsers("user.csv")
	if err != nil {
		panic(err)
	}

	if args != nil {
		searchUsers := map[int]*User{}
		for _, user := range users {
			for _, arg := range args {
				if strings.Contains(user.Name, arg) ||
					strings.Contains(user.Tel, arg) ||
					strings.Contains(user.Address, arg) ||
					strings.Contains(user.Desc, arg) {
					searchUsers[user.Id] = user
				}
			}
		}
		return searchUsers
	}

	return users
}

func UpdateUser(user User) {
	users, err := loadUsers("user.csv")
	if err != nil {
		panic(err)
	}
	if user.Id == 0 {
		user.Id = getId(users)
	}
	user.Password = md5Pass(user.Password)
	users[user.Id] = &user
	err = saveUsers(users, "user.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func DeleteUser(userId int) error {
	users, err := loadUsers("user.csv")
	if err != nil {
		return err
	}
	delete(users, userId)
	err = saveUsers(users, "user.csv")
	return err
}
