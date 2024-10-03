package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	_ "github.com/lib/pq"
)

type DbStore struct {
	Id        uuid.UUID `json:"id" db:"Id"`
	Hostname  string    `json:"hostname" db:"Hostname"`
	AuthToken string    `json:"authToken" db:"AuthToken"`
	UserUid   uuid.UUID `json:"userUid" db:"UserUid"`
}

type Store map[string]*sqlx.DB

var (
	DB      *sqlx.DB
	dbStore Store = make(Store)
	mu      sync.Mutex
)

func (s *DbStore) insertRecord() {

	query := `
		INSERT INTO DbStore (
			Id,
			Hostname,
			AuthToken,
			UserUid
		)
		VALUES (
			:Id,
			:Hostname,
			:AuthToken,
			:UserUid
		)
	`

	_, err := DB.NamedExec(query, &s)
	if err != nil {
		log.Fatal(err)
	}

}

func GetDb(userId uuid.UUID) (*sqlx.DB, error) {
	db := dbStore[userId.String()]
	if db == nil {
		//
		query := `
			SELECT * FROM DbStore WHERE UserUid = ?
		`

		var storeRec DbStore

		if err := DB.Get(&storeRec, query, userId); err != nil {
			if err == sql.ErrNoRows {
				mu.Lock()
				defer mu.Unlock()
				return createDatabase(userId)
			}

			return nil, errors.New("an error occured")
		}

		decodedJt := decodeJWT(storeRec.AuthToken)

		db, ndbErr := sqlx.Connect("libsql", "libsql://"+storeRec.Hostname+"?authToken="+decodedJt)
		if ndbErr != nil {
			return nil, ndbErr
		}
		mu.Lock()
		defer mu.Unlock()
		dbStore[userId.String()] = db
		return db, nil

	} else {
		return db, nil
	}

}

func InitDB(connString string) {
	var err error
	DB, err = sqlx.Connect("libsql", connString)

	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func createDatabase(token uuid.UUID) (*sqlx.DB, error) {
	if dbStore[token.String()] != nil {
		return dbStore[token.String()], nil
	}

	// DB, _ := GetDb(token)
	// if DB != nil {
	// 	return DB, nil
	// }

	body := []byte(fmt.Sprintf(`{
		"name": "%s",
		"group": "sillysocksapp",
		"seed": {
			"type": "database",
			"name": "masterdb"
		}

	}`, token.String()))

	req, _ := http.NewRequest("POST", "https://api.turso.tech/v1/organizations/jalenarms1/databases", bytes.NewReader(body))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TURSO_API_TOKEN")))
	client := &http.Client{}
	res, resErr := client.Do(req)
	if resErr != nil {
		return nil, resErr
	}

	data, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	var jd map[string]interface{}

	err := json.Unmarshal(data, &jd)
	if err != nil || jd["error"] != nil {
		return nil, errors.New("please verify the database does not exists already")
	}
	hostName := jd["database"].(map[string]interface{})["Hostname"]
	dbName := jd["database"].(map[string]interface{})["Name"].(string)
	fmt.Println("Getting token")
	dbToken, tErr := getDatabaseToken(dbName)
	if tErr != nil {
		return nil, tErr
	}
	fmt.Println("Got token")
	jt := tokenToJWT(dbToken)
	newId, _ := uuid.NewV4()
	newStore := &DbStore{
		Id:        newId,
		Hostname:  hostName.(string),
		AuthToken: jt,
		UserUid:   token,
	}

	go newStore.insertRecord()

	newDbUrl := "libsql://" + hostName.(string) + "?authToken=" + dbToken

	nDb, nErr := sqlx.Connect("libsql", newDbUrl)
	if nErr != nil {
		return nil, nErr
	}
	fmt.Println("Got new db")
	mu.Lock()
	defer mu.Unlock()
	dbStore[token.String()] = nDb
	return nDb, nil
}

func getDatabaseToken(dbName string) (string, error) {

	req, _ := http.NewRequest("POST", fmt.Sprintf("https://api.turso.tech/v1/organizations/jalenarms1/databases/%s/auth/tokens?expiration=2w&authorization=full-access", dbName), nil)

	req.Header.Add("Authorization", "Bearer "+os.Getenv("TURSO_API_TOKEN"))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	data, _ := io.ReadAll(resp.Body)

	var jd map[string]interface{}

	jErr := json.Unmarshal(data, &jd)
	if jErr != nil {
		return "", jErr
	}

	token := jd["jwt"].(string)

	return token, nil

}

func tokenToJWT(token string) string {
	jToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"name": token,
		},
	)

	tokenStr, err := jToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	return tokenStr
}

func decodeJWT(jt string) string {
	token, err := jwt.Parse(jt, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["name"].(string)
	}

	return ""

}
