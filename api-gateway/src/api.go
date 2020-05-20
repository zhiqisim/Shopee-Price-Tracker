package main

import (
	"api-gateway/proto"
	"fmt"
	"net/http"
	"os"

	// "strconv"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var (
	//RedisHost ...
	RedisHost = os.Getenv("REDISHOST")
	//RedisPort ...
	RedisPort = os.Getenv("REDISPORT")
)

func authRequired() gin.HandlerFunc {
	/*
		Description: Authentication middleware
		Check if user has the cookie that is stored in our session cache
	*/
	return func(c *gin.Context) {
		// check if have the cookie if no return error
		token, err := c.Cookie("auth-token")
		if err == http.ErrNoCookie {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "error",
				"error":   "auth",
			})
			c.Abort()
		} else {
			session := sessions.Default(c)
			sessionID := session.Get(token)
			if sessionID == nil {
				c.JSON(http.StatusForbidden, gin.H{
					"message": "error",
					"error":   "auth",
				})
				c.Abort()
			}
		}

	}
}

func getUsername(c *gin.Context) string {
	/*
		Description: Helper function to get username from session
	*/
	// check if have the cookie if no return error
	token, err := c.Cookie("auth-token")
	if err == http.ErrNoCookie {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "error",
			"error":   "auth",
		})
		c.Abort()
		return "error"
	}
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: 60})
	interfaceUser := session.Get(token)
	username := fmt.Sprintf("%v", interfaceUser)
	return username

}

func main() {
	// gRPC connection
	// Port 4040 for User Service (Go)
	// Port 50051 for Item Service (Python)
	userHost := "user-service"
	conn, err := grpc.Dial(userHost+":4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewUserServiceClient(conn)

	itemHost := "item-service"
	itemConn, err := grpc.Dial(itemHost+":50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	itemClient := proto.NewItemServiceClient(itemConn)
	router := gin.Default()
	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// connect to REDIS session store
	store, _ := redis.NewStore(10, "tcp", RedisHost+":"+RedisPort, "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// max age of each session in seconds
	sessionTimeout := 60 * 60 // 1 hour

	/*
		Group: public (no auth required)
		API-gateway Endpoints
	*/

	router.POST("/login", func(c *gin.Context) {
		/*
			Description: Allow user to login and obtain auth session
			Input: Form body - username, password
			Output: JSON Object - AuthToken: Auth session token
		*/
		sessionToken := uuid.New().String()
		username := c.PostForm("username")
		password := c.PostForm("password")
		req := &proto.LoginRequest{User: &proto.User{
			Username: username,
			Password: password,
		}}
		if response, err := client.Login(c, req); err == nil {
			session := sessions.Default(c)
			session.Options(sessions.Options{MaxAge: sessionTimeout})
			session.Set(sessionToken, username)
			session.Save()
			if response.Message == "error" {
				c.JSON(http.StatusOK, gin.H{
					"message": response.Message,
				})
			} else {
				// set cookie for client
				c.SetCookie("auth-token", sessionToken, sessionTimeout, "/", "localhost", false, true)
				c.JSON(http.StatusOK, gin.H{
					"message": response.Message,
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	router.POST("/signup", func(c *gin.Context) {
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

	/*
		Group: user (Auth required)
		Service: User
	*/
	user := router.Group("/user")
	user.Use(authRequired())
	{

		user.POST("/logout", func(c *gin.Context) {
			/*
				Description: Allow user to logout
				Input: JSON Object - username
				Output: NIL
				Header: Auth session token
			*/

			// check if have the cookie
			token, err := c.Cookie("auth-token")
			if err == http.ErrNoCookie {
				c.JSON(http.StatusForbidden, gin.H{
					"message": "error",
					"error":   "auth",
				})
			} else {
				// delete session token for user
				session := sessions.Default(c)
				session.Options(sessions.Options{MaxAge: 60})
				session.Delete(token)
				session.Save()
				// obtain username from session to track user
				username := getUsername(c)
				req := &proto.LogoutRequest{Username: username}
				if response, err := client.Logout(c, req); err == nil {
					// delete cookie from client
					c.SetCookie("auth-token", "placeholder", -1, "/", "localhost", false, true)
					c.JSON(http.StatusOK, gin.H{
						"message": response.Message,
					})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
			}

		})

		user.POST("/add-item", func(c *gin.Context) {
			/*
				Description: Add an item to user's tracking list
				Input: Form body - item_id
				Output: NIL
				Header: Auth session token
			*/

			// obtain username from session to track user
			username := getUsername(c)
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
			// obtain username from session to track user
			username := getUsername(c)
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

	/*
		Group: items (no auth required)
		Item Service
	*/
	items := router.Group("/item")
	{

		items.GET("/get-items", func(c *gin.Context) {
			/*
				Description: Retrieve all items from database for user to add to tracker
				Args: NIL
				Output: JSON Object with list of all itmes
				Header: Auth session token
			*/
			req := &proto.ListAllItemsRequest{}
			if response, err := itemClient.ListAllItems(c, req); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"message": response.Message,
					"items":   response.Items,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		})
		items.GET("/price", func(c *gin.Context) {
			/*
				Description: Retrieve price changelog of item
				Args: itemid
				Output: JSON Object with list of all price history of item
				Header: Auth session token
			*/

			itemID := c.Query("itemid")
			req := &proto.ItemPriceRequest{ItemId: itemID}
			if response, err := itemClient.ItemPrice(c, req); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"message": response.Message,
					"items":   response.ItemPrice,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		})
	}

	router.Run(":8080")
}
