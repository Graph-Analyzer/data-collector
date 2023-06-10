package repository

import (
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"
)

func InitDBConnection(host string, port string, username string, password string, realm string) neo4j.Driver {
	driver, err := neo4j.NewDriver(
		fmt.Sprintf("%s:%s", host, port),
		neo4j.BasicAuth(
			username,
			password,
			realm,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Try to verify database connection 10 times
	for i := 1; i <= 10; i++ {
		log.Debugf("Trying to verify database connection, try #%d of %d\n", i, 10)
		err = driver.VerifyConnectivity()
		if err == nil {
			break
		}

		if i == 10 {
			log.Fatalln("Connection to database could not be established")
		}

		log.Warningf("Failed to establish database connection, retrying in %d seconds\n", 10)

		time.Sleep(10 * time.Second)
	}

	log.Infoln("Database connection established")

	return driver
}
