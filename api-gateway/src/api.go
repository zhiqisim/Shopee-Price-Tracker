package main

import (
	// "fmt"
	"api-gateway/proto"
	"net/http"

	// "strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Logout struct {
	Username string `json:"username" binding:"required"`
	// Password string `form:"password" json:"password" binding:"required"`
}

type Items struct {
	Item_Id uint `json:"item_id" binding:"required"`
}

func main() {
	// gRPC connection
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewUserServiceClient(conn)
	router := gin.Default()

	// TODO: when auth change this username to user's username
	username := "zhiqisim"

	// API-gateway Endpoints

	// TODO: ADD Auth with redis cache
	// group: user
	/*
		User Service
	*/
	user := router.Group("/user")
	{
		user.POST("/login", func(c *gin.Context) {
			/*
				Description: Allow user to login and obtain auth session
				Input: Form body - username, password
				Output: JSON Object - AuthToken: Auth session token
			*/

			req := &proto.LoginRequest{User: &proto.User{
				Username: c.PostForm("username"),
				Password: c.PostForm("password"),
			}}
			if response, err := client.Login(c, req); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"message": response.Message,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		})
		user.POST("/logout", func(c *gin.Context) {
			/*
				Description: Allow user to logout
				Input: JSON Object - username
				Output: NIL
				Header: Auth session token
			*/

			req := &proto.LogoutRequest{Username: "zhiqisim"}
			if response, err := client.Logout(c, req); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"message": response.Message,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		})
		user.POST("/signup", func(c *gin.Context) {
			/*
				Description: Allow user to signup
				Input: Form body - username, password
				Output: NIL
			*/

			req := &proto.SignupRequest{User: &proto.User{
				Username: c.PostForm("username"),
				Password: c.PostForm("password"),
			}}
			if response, err := client.Signup(c, req); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"message": response.Message,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		})
		user.POST("/add-item", func(c *gin.Context) {
			/*
				Description: Add an item to user's tracking list
				Input: Form body - item_id
				Output: NIL
				Header: Auth session token
			*/

			// response body
			req := &proto.AddItemRequest{UserItem: &proto.UserItem{
				Username: username,
				ItemId:   c.PostForm("item_id"),
			}}
			if response, err := client.AddItem(c, req); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"message": response.Message,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}

		})
		user.GET("/watchlist", func(c *gin.Context) {
			/*
				Description: Retrieve all items from database that user have added to tracker
				Args: NIL
				Output: JSON Object with list of all user tracked itmes
				Header: Auth session token
			*/

			req := &proto.ListItemsRequest{Username: username}
			if response, err := client.ListItems(c, req); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"message": response.Message,
					"item_id": response.ItemId,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		})
	}

	// group: items
	// authorized := r.Group("/")
	// // per group middleware! in this case we use the custom created
	// // AuthRequired() middleware just in the "authorized" group.
	// authorized.Use(AuthRequired())
	/*
		Item Service
	*/
	items := router.Group("/item")
	{

		items.GET("/get-items", getItems)
		items.GET("/price", price)
	}

	router.Run(":8080")
}

func getItems(c *gin.Context) {
	/*
		Description: Retrieve all items from database for user to add to tracker
		Args: NIL
		Output: JSON Object with list of all itmes
		Header: Auth session token
	*/
	// TODO: ADD in gRPC request to items service
	// response body
	c.JSON(200, gin.H{
		"message": "All items listed",
	})
}

func price(c *gin.Context) {
	/*
		Description: Retrieve price changelog of item
		Args: item_id
		Output: JSON Object with list of all price history of item
		Header: Auth session token
	*/
	// TODO: ADD in gRPC request to items service

	item_id := c.Query("itemid")
	// response body
	c.JSON(200, gin.H{
		"message": "Price history shown",
		"item_id": item_id,
	})
}
