package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
	"os"
	"time"
	"user-service/proto"

	log "github.com/sirupsen/logrus"

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
		"user-db:3306",
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
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	// log.SetLevel(log.InfoLevel)
	InitializeMySQL()
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterUserServiceServer(srv, &server{})
	reflection.Register(srv)
	log.Info("Serving gRPC server on port 4040")

	if e := srv.Serve(listener); e != nil {
		log.Panic(e)
	}

}

// Login for user
func (s *server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	sqlQuery := "SELECT username, user_password FROM user WHERE username=?"
	stmt, err := db.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		log.Error(err)
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
		log.Error(errr)
		return &proto.LoginResponse{
			Message: "error",
		}, nil
	}
	// check password hash
	errf := bcrypt.CompareHashAndPassword([]byte(fromdb_password), []byte(password))
	if errf != nil { //Password does not match!
		log.Error(errf)
		return &proto.LoginResponse{
			Message: "error",
		}, nil
	}
	log.Info("Login gRPC endpoint ran successfully.")
	return &proto.LoginResponse{
		Message: "success",
	}, nil
}

// Logout for user
func (s *server) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	log.Info("Logout gRPC endpoint ran successfully.")
	return &proto.LogoutResponse{
		Message: "success",
	}, nil
}

// Signup for user
func (s *server) Signup(ctx context.Context, req *proto.SignupRequest) (*proto.SignupResponse, error) {

	sqlQuery := "INSERT user SET username = ?, user_password = ?"
	stmt, err := db.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		log.Error(err)
		return &proto.SignupResponse{
			Message: "error",
		}, nil
	}
	username := req.User.Username
	password, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return &proto.SignupResponse{
			Message: "error",
		}, nil
	}
	res, err := stmt.Exec(username, password)
	if err != nil {
		log.Error(err)
		return &proto.SignupResponse{
			Message: "error",
		}, nil
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Error(err)
		return &proto.SignupResponse{
			Message: "error",
		}, nil
	}
	log.Info("Signup gRPC endpoint ran successfully.")
	return &proto.SignupResponse{
		Message: "success",
	}, nil
}

// Add item for user
func (s *server) AddItem(ctx context.Context, req *proto.AddItemRequest) (*proto.AddItemResponse, error) {
	sqlQuery := "INSERT user_item SET username = ?, item_id = ?, item_name = ?"
	stmt, err := db.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		log.Error(err)
		return &proto.AddItemResponse{
			Message: "error",
		}, nil
	}
	username := req.UserItem.Username
	itemID := req.UserItem.ItemId
	itemName := req.UserItem.ItemName
	res, err := stmt.Exec(username, itemID, itemName)
	if err != nil {
		log.Error(err)
		return &proto.AddItemResponse{
			Message: "error",
		}, nil
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Error(err)
		return &proto.AddItemResponse{
			Message: "error",
		}, nil
	}
	log.Info("AddItem gRPC endpoint ran successfully.")
	return &proto.AddItemResponse{
		Message: "success",
	}, nil
}

// Read all items in watchlist of user
func (s *server) ListItems(ctx context.Context, req *proto.ListItemsRequest) (*proto.ListItemsResponse, error) {
	sqlQuery := "SELECT `item_id`, `item_name` FROM user_item WHERE `username`=? LIMIT ? OFFSET ?"
	stmt, err := db.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		log.Error(err)
		return &proto.ListItemsResponse{
			Message: "error",
		}, nil
	}
	list := []*proto.ItemDetails{}
	res, err := stmt.Query(req.Username, req.Limit, req.Offset)
	defer res.Close()
	if err != nil {
		log.Error(err)
		return &proto.ListItemsResponse{
			Message: "error",
		}, nil
	}
	for res.Next() {
		td := new(proto.ItemDetails)
		res.Scan(&td.ItemId, &td.ItemName)
		list = append(list, td)
	}
	log.Info("ListItems gRPC endpoint ran successfully.")
	return &proto.ListItemsResponse{
		Message: "success",
		Item:    list,
	}, nil
}
