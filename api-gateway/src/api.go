package main

import "github.com/gin-gonic/gin"
// import "net/http"

type Logout struct {
    Username string `json:"username" binding:"required"`
    // Password string `form:"password" json:"password" binding:"required"`
}

type Items struct {
	Item_Id uint `json:"item_id" binding:"required"`
}


func main() {
	router := gin.Default()

	// TODO: ADD Auth with redis cache
	// group: user
	user := router.Group("/user")
	{
		user.POST("/login", login)
		user.POST("/logout", logout)
		user.POST("/signup", signup)
	}

	// group: items
	// authorized := r.Group("/")
	// // per group middleware! in this case we use the custom created
	// // AuthRequired() middleware just in the "authorized" group.
	// authorized.Use(AuthRequired())
	items := router.Group("/items")
	{
		items.GET("/get-items", getItems)
		items.POST("/add-item", addItem)
		items.GET("/list-user-items", listUserItems)
		items.GET("/price", price)
	}

	router.Run(":8080")
}

func login(c *gin.Context) {
	/* 
		Description: Allow user to login and obtain auth session
		Input: Form body - username, password
		Output: JSON Object - AuthToken: Auth session token
	*/
	// TODO: ADD in gRPC request to user service

	username := c.PostForm("username")
	password := c.PostForm("password")
	// response body
	c.JSON(200, gin.H{
		"message": "Login success",
		"username": username,
		"password": password,
	})
}

func signup(c *gin.Context) {
	/* 
		Description: Allow user to signup
		Input: Form body - username, password
		Output: NIL
	*/
	// TODO: ADD in gRPC request to user service
	username := c.PostForm("username")
	password := c.PostForm("password")
	// response body
	c.JSON(200, gin.H{
		"message": "signup success",
		"username": username,
		"password": password,
	})
}

func logout(c *gin.Context) {
	/* 
		Description: Allow user to logout
		Input: JSON Object - username
		Output: NIL
		Header: Auth session token
	*/
	// TODO: ADD in gRPC request to user service
	var json Logout
	c.BindJSON(&json)
	// response body
	c.JSON(200, gin.H{
		"message": "Logout success",
		"username": json.Username,
	})
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

func addItem(c *gin.Context) {
	/* 
		Description: Add an item to user's tracking list
		Input: JSON Object - item_id
		Output: NIL
		Header: Auth session token
	*/
	// TODO: ADD in gRPC request to items service
	var json Items
	c.BindJSON(&json)
	// response body
	if json.Item_Id == 0 {
		c.JSON(404, gin.H{
			"message": "Add items unsuccessful",
		})
	} else{
		c.JSON(200, gin.H{
			"message": "Add items success",
			"item_id": json.Item_Id,
		})
	}

}

func listUserItems(c *gin.Context) {
	/* 
		Description: Retrieve all items from database that user have added to tracker
		Args: NIL
		Output: JSON Object with list of all user tracked itmes
		Header: Auth session token
	*/
	// TODO: ADD in gRPC request to items service
	// response body
	c.JSON(200, gin.H{
		"message": "All user items listed",
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
