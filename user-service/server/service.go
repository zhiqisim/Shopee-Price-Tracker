package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"
	"user-service/proto"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

var db *sql.DB

func InitializeMySQL() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		"root",
		"root",
		"localhost:32000",
		"userdb",
		"parseTime=true")
	dBConnection, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Connection Failed!!")
	}
	err = dBConnection.Ping()
	if err != nil {
		fmt.Println("Ping Failed!!")
	}
	db = dBConnection
	dBConnection.SetMaxOpenConns(10)
	dBConnection.SetMaxIdleConns(5)
	dBConnection.SetConnMaxLifetime(time.Second * 10)
}

func GetMySQLConnection() *sql.DB {
	return db
}

func main() {
	InitializeMySQL()
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterUserServiceServer(srv, &server{})
	reflection.Register(srv)
	fmt.Println("Serving gRPC server on port 4040")

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}

}

// Add item for user
func (s *server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	sqlQuery := "SELECT username, user_password FROM user WHERE username=?"
	stmt, err := db.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return &proto.LoginResponse{
			Message: "error",
		}, nil
	}
	user := req.User
	username := user.Username
	password := user.Password
	var fromdb_username string
	var fromdb_password string
	errr := stmt.QueryRow(username).Scan(&fromdb_username, &fromdb_password)
	if errr != nil {
		fmt.Println(errr)
		return &proto.LoginResponse{
			Message: "error",
		}, nil
	}
	// check password hash
	errf := bcrypt.CompareHashAndPassword([]byte(fromdb_password), []byte(password))
	if errf != nil { //Password does not match!
		fmt.Println(errf)
		return &proto.LoginResponse{
			Message: "error",
		}, nil
	}
	fmt.Println("success")
	return &proto.LoginResponse{
		Message: "success",
	}, nil
}

// Add item for user
func (s *server) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	return &proto.LogoutResponse{
		Message: "success",
	}, nil
}

// Add item for user
func (s *server) Signup(ctx context.Context, req *proto.SignupRequest) (*proto.SignupResponse, error) {

	sqlQuery := "INSERT user SET username = ?, user_password = ?"
	stmt, err := db.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return &proto.SignupResponse{
			Message: "error",
		}, nil
	}
	username := req.User.Username
	password, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return &proto.SignupResponse{
			Message: "error",
		}, nil
	}
	res, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println(err)
		return &proto.SignupResponse{
			Message: "error",
		}, nil
	}
	_, err = res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return &proto.SignupResponse{
			Message: "error",
		}, nil
	}
	return &proto.SignupResponse{
		Message: "success",
	}, nil
}

// Add item for user
func (s *server) AddItem(ctx context.Context, req *proto.AddItemRequest) (*proto.AddItemResponse, error) {
	sqlQuery := "INSERT user_item SET username = ?, item_id = ?"
	stmt, err := db.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return &proto.AddItemResponse{
			Message: "error",
		}, nil
	}
	username := req.UserItem.Username
	item_id := req.UserItem.ItemId
	res, err := stmt.Exec(username, item_id)
	if err != nil {
		fmt.Println(err)
		return &proto.AddItemResponse{
			Message: "error",
		}, nil
	}
	_, err = res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return &proto.AddItemResponse{
			Message: "error",
		}, nil
	}
	fmt.Println("AddItems success!")
	return &proto.AddItemResponse{
		Message: "success",
	}, nil
}

// Read all items from user
func (s *server) ListItems(ctx context.Context, req *proto.ListItemsRequest) (*proto.ListItemsResponse, error) {
	sqlQuery := "SELECT `item_id` FROM user_item WHERE `username`=?"
	stmt, err := db.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return &proto.ListItemsResponse{
			Message: "error",
		}, nil
	}
	list := []string{}
	res, err := stmt.Query(req.Username)
	defer res.Close()
	if err != nil {
		fmt.Println(err)
		return &proto.ListItemsResponse{
			Message: "error",
		}, nil
	}
	for res.Next() {
		var td string
		res.Scan(&td)
		list = append(list, td)
	}
	fmt.Println("ListItems success!")
	return &proto.ListItemsResponse{
		Message: "success",
		ItemId:  list,
	}, nil
}
