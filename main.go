package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

type Product struct {
    ID          int
    Name        string
    Price       float64
    CategoryID  int
    CategoryName string
}

func main() {
    db, err := sql.Open("postgres", "user=postgres password=yourpassword dbname=demo sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }
    defer tx.Rollback()
    _, err = tx.Exec("INSERT INTO Products (Name, Price, CategoryID) VALUES ($1, $2, (SELECT ID FROM Categories WHERE Name = $3))", "New Product", 99.99, "Electronics")
    if err != nil {
        log.Fatal(err)
    }
    var product Product
    err = tx.QueryRow("SELECT p.ID, p.Name, p.Price, p.CategoryID, c.Name FROM Products p JOIN Categories c ON p.CategoryID = c.ID WHERE p.Name = $1", "New Product").Scan(&product.ID, &product.Name, &product.Price, &product.CategoryID, &product.CategoryName)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Yangi mahsulot: %+v\n", product)
    _, err = tx.Exec("UPDATE Products SET Price = $1 WHERE ID = $2", 149.99, product.ID)
    if err != nil {
        log.Fatal(err)
    }
    err = tx.QueryRow("SELECT p.Price, c.Name FROM Products p JOIN Categories c ON p.CategoryID = c.ID WHERE p.ID = $1", product.ID).Scan(&product.Price, &product.CategoryName)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Yangilangan mahsulot: Narx: %.2f, Kategoriya: %s\n", product.Price, product.CategoryName)
    _, err = tx.Exec("DELETE FROM Products WHERE ID = $1", product.ID)
    if err != nil {
        log.Fatal(err)
    }
    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }
}
