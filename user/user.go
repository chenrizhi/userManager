package user

import (
	"awesomeProject/userManager/db"
	"crypto/md5"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
)

type User struct {
	Id       int
	Name     string
	Birthday string
	Tel      string
	Address  string
	Remarks  string
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
		name, birthday, address, tel, remarks, password := v[1], v[2], v[3], v[4], v[5], v[6]

		users[id] = &User{
			Id:       id,
			Name:     name,
			Birthday: birthday,
			Address:  address,
			Tel:      tel,
			Remarks:  remarks,
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
	records = append(records, []string{"id", "name", "birthday", "address", "tel", "remarks", "password"})
	for _, u := range users {
		record := []string{
			strconv.Itoa(u.Id),
			u.Name,
			u.Birthday,
			u.Address,
			u.Tel,
			u.Remarks,
			u.Password,
		}
		records = append(records, record)
	}
	writer := csv.NewWriter(f)
	err = writer.WriteAll(records)
	return err
}

func GetUser(userId int) *User {
	row := db.Db.QueryRow("SELECT id,name,birthday,address,tel,remarks,password FROM users WHERE id = ?", userId)
	if err := row.Err(); err != nil {
		fmt.Println("get user failed, err:", err)
		return nil
	}
	user := &User{}
	if err := row.Scan(&user.Id, &user.Name, &user.Birthday, &user.Address, &user.Tel, &user.Remarks, &user.Password); err != nil {
		fmt.Println("get user scan failed, err:", err)
	}
	return user
}

func GetUsers(args ...string) map[int]*User {
	var (
		rows *sql.Rows
		err  error
	)
	if args == nil {
		rows, err = db.Db.Query("SELECT id,name,birthday,address,tel,remarks,password FROM users")
	} else {
		search := args[0]
		rows, err = db.Db.Query("SELECT id,name,birthday,address,tel,remarks,password FROM users WHERE id = ? OR name LIKE ?", search, "%"+search+"%")
	}

	if err != nil {
		fmt.Println("get users failed, err:", err)
		return nil
	}
	defer rows.Close()

	users := map[int]*User{}
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.Id, &user.Name, &user.Birthday, &user.Address, &user.Tel, &user.Remarks, &user.Password); err != nil {
			fmt.Println("get user scan failed, err:", err)
			continue
		}
		users[user.Id] = user
	}
	return users
}

func InsertUser(user User) {
	_, err := db.Db.Exec("INSERT INTO users(name,birthday,address,tel,remarks,password) VALUES(?,?,?,?,?,?)",
		user.Name, user.Birthday, user.Address, user.Tel, user.Remarks, md5Pass(user.Password))
	if err != nil {
		fmt.Println("insert user failed, err:", err)
		return
	}
}

func UpdateUser(user User) {
	if user.Password != "" {
		user.Password = md5Pass(user.Password)
	}
	_, err := db.Db.Exec("REPLACE INTO users(id,name,birthday,address,tel,remarks,password) VALUES(?,?,?,?,?,?,?)",
		user.Id, user.Name, user.Birthday, user.Address, user.Tel, user.Remarks, user.Password)
	if err != nil {
		fmt.Println("update user failed, err:", err)
		return
	}
}

func DeleteUser(userId int) error {
	_, err := db.Db.Exec("DELETE FROM users WHERE id = ?", userId)
	return err
}
