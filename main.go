package main

import (
	"assignment-2/database"
	"assignment-2/routes"
)

func main()  {
	var PORT = ":8080"
	
	database.StartDB()
	routes.StartServer().Run(PORT)
}
