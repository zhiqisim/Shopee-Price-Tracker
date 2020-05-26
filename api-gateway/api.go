package main

import (
	"api-gateway/proto"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	// "github.com/processout/grpc-go-pool"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"google.golang.org/grpc"
)

var (
	//RedisHost ...
	RedisHost = os.Getenv("REDISHOST")
	//RedisPort ...
	RedisPort = os.Getenv("REDISPORT")
	// Server Port
	serverPort = os.Getenv("APPID")
	// user service host
	userHost = os.Getenv("USERHOST")
	// item service host
	itemHost = os.Getenv("ITEMHOST")
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
			c.JSON(http.StatusOK, gin.H{
				"message": "error",
				"error":   "auth",
			})
			c.Abort()
		} else {
			session := sessions.Default(c)
			sessionID := session.Get(token)
			if sessionID == nil {
				c.JSON(http.StatusOK, gin.H{
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
		// http.StatusForbidden
		c.JSON(http.StatusOK, gin.H{
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

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	// router := gin.Default()
	router := gin.New()

	conn, err := grpc.Dial(userHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewUserServiceClient(conn)

	// var factory grpcpool.Factory
	// factory = func() (*grpc.ClientConn, error) {
	// 	itemC, err := grpc.Dial(itemHost+":50051", grpc.WithInsecure())
	// 	if err != nil {
	// 		log.Fatalf("Failed to start gRPC connection: %v", err)
	// 	}
	// 	return itemC, err
	// }

	// pool, err := grpcpool.New(factory, 5, 5, time.Second)
	// if err != nil {
	// 	log.Fatalf("Failed to create gRPC pool: %v", err)
	// }
	itemConn, err := grpc.Dial(itemHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	itemClient := proto.NewItemServiceClient(itemConn)

	p := ginprometheus.NewPrometheus("gin")
	// Preserving a low cardinality for the request counter - deleting query params
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		return c.Request.URL.Path
	}
	p.Use(router)
	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3001"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// connect to REDIS session store
	store, _ := redis.NewStore(10, "tcp", RedisHost+":"+RedisPort, "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// max age of each session in seconds
	sessionTimeout := 60 * 60 // 1 hour

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/metrics"},
	}))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	/*
		Group: public (no auth required)
		API-gateway Endpoints
	*/

	router.POST("/api/login", func(c *gin.Context) {
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

	router.POST("/api/signup", func(c *gin.Context) {
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
	user := router.Group("/api/user")
	user.Use(authRequired())
	{

		user.GET("/is-auth", func(c *gin.Context) {
			/*
				Description: Check if user is auth-ed to server
				Args: NIL
				Output: message that user is auth-ed
			*/
			c.JSON(http.StatusOK, gin.H{
				"message": "success",
			})
		})

		user.POST("/logout", func(c *gin.Context) {
			/*
				Description: Allow user to logout
				Input: NIL
				Output: NIL
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
				Input: Form body - item_id, item_name
				Output: NIL
			*/

			// obtain username from session to track user
			username := getUsername(c)

			req := &proto.AddItemRequest{UserItem: &proto.UserItem{
				Username: username,
				ItemId:   c.PostForm("item_id"),
				ItemName: c.PostForm("item_name"),
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
			*/
			// obtain username from session to track user
			username := getUsername(c)
			limit := c.Query("limit")
			offset := c.Query("offset")
			o, err := strconv.Atoi(offset)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "error",
					"error":   "invalid offset param",
				})
				return
			}
			l, err := strconv.Atoi(limit)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "error",
					"error":   "invalid limit param",
				})
				return
			}
			req := &proto.ListItemsRequest{Username: username, Limit: int64(l), Offset: int64(o)}
			if response, err := client.ListItems(c, req); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"message":      response.Message,
					"item_details": response.Item,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		})

		user.POST("/add-new-item", func(c *gin.Context) {
			/*
				Description: Add an item to user's tracking list
				Input: Form body - item_id, item_name, shop_id
				Output: NIL
			*/

			// // obtain connection from pool
			// itemConn, err := pool.Get(c)
			// defer itemConn.Close()
			// if err != nil {
			// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			// 	return
			// }
			// itemClient := proto.NewItemServiceClient(itemConn)

			// obtain username from session to track user
			username := getUsername(c)

			req := &proto.AddItemRequest{UserItem: &proto.UserItem{
				Username: username,
				ItemId:   c.PostForm("item_id"),
				ItemName: c.PostForm("item_name"),
			}}
			req2 := &proto.AddNewItemRequest{Item: &proto.ItemWithoutPrice{
				ItemId:   c.PostForm("item_id"),
				ItemName: c.PostForm("item_name"),
				ShopId:   c.PostForm("shop_id"),
			}}
			if response2, err2 := itemClient.AddNewItem(c, req2); err2 == nil {
				if response, err := client.AddItem(c, req); err == nil {
					if response.Message == "success" {
						c.JSON(http.StatusOK, gin.H{
							"message": response2.Message,
						})
					} else {
						c.JSON(http.StatusOK, gin.H{
							"message": response.Message,
						})
					}
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			}

		})
	}

	/*
		Group: items (no auth required)
		Item Service
	*/
	items := router.Group("/api/item")
	{

		items.GET("/get-items", func(c *gin.Context) {
			/*
				Description: Retrieve all items from database for user to add to tracker
				Args: NIL
				Output: JSON Object with list of all itmes
			*/

			// // obtain connection from pool
			// itemConn, err := pool.Get(c)
			// defer itemConn.Close()
			// if err != nil {
			// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			// 	return
			// }
			// itemClient := proto.NewItemServiceClient(itemConn)

			// obtain offset and limit from query params
			offset := c.Query("offset")
			limit := c.Query("limit")
			o, err := strconv.Atoi(offset)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "error",
				})
				return
			}
			l, err := strconv.Atoi(limit)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "error",
				})
				return
			}
			req := &proto.ListAllItemsRequest{Offset: int64(o), Limit: int64(l)}
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
			*/

			// // obtain connection from pool
			// itemConn, err := pool.Get(c)
			// defer itemConn.Close()
			// if err != nil {
			// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			// 	return
			// }
			// itemClient := proto.NewItemServiceClient(itemConn)

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

	// router.Run(":8000")
	router.Run(":" + serverPort)
}
