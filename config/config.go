package config

import (
    "os"
)

func GetMongoURI() string {
    uri := "mongodb+srv://vdcluster0:admin@cluster0.7k2pl.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
    return uri
}

func GetPort() string {
    port := os.Getenv("PORT")
    if port == "" {
        return "8080"
    }
    return port
}