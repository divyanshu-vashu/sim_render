package config

import (
    "fmt"
    "os"
)

func GetDBConfig() string {
    host := "ep-summer-dust-a1syu99u.ap-southeast-1.pg.koyeb.app"
    user := "vashu-admin"
    password := "npg_sjFa1wcy5WQp"
    dbname := "koyebdb"

    return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", 
        host, user, password, dbname)
}

func GetEmailPassword() string {
    return "xfmz rlod pixm mjvi"
}

func GetPort() string {
    port := os.Getenv("PORT")
    if port == "" {
        return "8080"
    }
    return port
}