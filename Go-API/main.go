package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	Valor    string `json:"valor"`
	Quantity int    `json:"quantity"`
}

type WishlistItem struct {
	UserID string `json:"user_id"`
	ItemID string `json:"item_id"`
}

var inventory = []Item{
	{ID: "1", Nome: "In Search of Lost Time", Valor: "30.00", Quantity: 2},
	{ID: "2", Nome: "The Great Gatsby", Valor: "50.00", Quantity: 5},
	{ID: "3", Nome: "War and Peace", Valor: "25.50", Quantity: 6},
}

var wishlist = make(map[string][]string)

func getItem(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, inventory)
}

func getItemByID(c *gin.Context) {
	id := c.Param("id")
	item, err := getItemByIDFromInventory(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, item)
}

func checkoutItem(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameter"})
		return
	}

	item, err := getItemByIDFromInventory(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
		return
	}

	if item.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "item not available"})
		return
	}

	item.Quantity--
	c.IndentedJSON(http.StatusOK, item)
}

func returnItem(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameter"})
		return
	}

	item, err := getItemByIDFromInventory(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
		return
	}

	item.Quantity++
	c.IndentedJSON(http.StatusOK, item)
}

func getItemByIDFromInventory(id string) (*Item, error) {
	for i, item := range inventory {
		if item.ID == id {
			return &inventory[i], nil
		}
	}
	return nil, errors.New("item not found")
}

//Criando item
func createItem(c *gin.Context) {
	var newItem Item

	if err := c.BindJSON(&newItem); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	inventory = append(inventory, newItem)
	c.IndentedJSON(http.StatusCreated, newItem)
}

//adicionando a wishlist
func addToWishlist(c *gin.Context) {
	userID := c.Query("user_id")
	itemID := c.Query("item_id")

	if userID == "" || itemID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing user_id or item_id query parameter"})
		return
	}

	wishlistItems, ok := wishlist[userID]
	if !ok {
		wishlistItems = []string{}
	}

	wishlistItems = append(wishlistItems, itemID)
	wishlist[userID] = wishlistItems

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "item added to wishlist"})
}

//remover da wishlist
func removeFromWishlist(c *gin.Context) {
	userID := c.Query("user_id")
	itemID := c.Query("item_id")

	if userID == "" || itemID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing user_id or item_id query parameter"})
		return
	}

	wishlistItems, ok := wishlist[userID]
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "wishlist not found for user"})
		return
	}

	for i, id := range wishlistItems {
		if id == itemID {
			wishlistItems = append(wishlistItems[:i], wishlistItems[i+1:]...)
			wishlist[userID] = wishlistItems
			c.IndentedJSON(http.StatusOK, gin.H{"message": "item removed from wishlist"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found in wishlist"})
}

func main() {
	router := gin.Default()
	router.GET("/item", getItem)
	router.GET("/item/:id", getItemByID)
	router.POST("/item", createItem)
	router.PATCH("/checkout", checkoutItem)
	router.PATCH("/return", returnItem)
	router.POST("/wishlist/add", addToWishlist)
	router.DELETE("/wishlist/remove", removeFromWishlist)
	router.Run("localhost:8080")
}

