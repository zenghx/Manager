package db

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbname  string = "dbName"
	dbUname string = "dnUSer"
	dbPwd   string = "dbPwd"
	socket  string = "127.0.0.1:3306"
)

var dbConn *gorm.DB

func InitDBConn() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUname, dbPwd, socket, dbname)
	dbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

}

func GetDB() *gorm.DB {
	return dbConn
}

func CreateUser(u *User) {

	dbConn.Create(u)
}

func IsNewEmail(email string) bool {
	res := []User{}
	dbConn.Find(&res, User{Email: email})
	return len(res) == 0
}

func QueryUUIDById(uid int) string {
	var user = User{UID: uid}
	dbConn.First(&user)
	id, _ := uuid.FromBytes(user.UUID)
	return id.String()
}

func GetUserByEmail(email string) User {
	res := []User{}
	dbConn.Where("email = ?", email).Find(&res)
	//fmt.Println(user)
	if len(res) == 0 {
		return User{}
	}
	return res[0]
}

func GetUserById(id int) User {
	user := User{UID: id}
	dbConn.First(&user)
	//fmt.Println(user)
	return user
}

func UpdateVerif(id int) {
	user := User{UID: id}
	dbConn.First(&user)
	user.Verified = true
	user.UUID = []byte(uuid.New().String())
	//fmt.Println(user)
	dbConn.Save(&user)
}
