package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type Product struct {
    ID          int
    Name        string
    Price       float64
    CategoryID  int
    CategoryName string
}

func main() {
    db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/demo")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer tx.Rollback() 
    _, err = tx.ExecContext(ctx, "INSERT INTO Products (Name, Price, CategoryID) VALUES (?, ?, ?)", "New Product", 99.99, 1)
    if err != nil {
        log.Fatal(err)
    }
    var product Product
    err = tx.QueryRowContext(ctx, "SELECT p.ID, p.Name, p.Price, c.ID, c.Name FROM Products p JOIN Categories c ON p.CategoryID = c.ID WHERE p.Name = ?", "New Product").Scan(&product.ID, &product.Name, &product.Price, &product.CategoryID, &product.CategoryName)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Yangi mahsulot: %+v\n", product)
    _, err = tx.ExecContext(ctx, "UPDATE Products SET Price = ? WHERE ID = ?", 149.99, product.ID)
    if err != nil {
        log.Fatal(err)
    }
    err = tx.QueryRowContext(ctx, "SELECT p.Price, c.Name FROM Products p JOIN Categories c ON p.CategoryID = c.ID WHERE p.ID = ?", product.ID).Scan(&product.Price, &product.CategoryName)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Yangilangan mahsulot: Narx: %.2f, Categoriya: %s\n", product.Price, product.CategoryName)
    _, err = tx.ExecContext(ctx, "DELETE FROM Products WHERE ID = ?", product.ID)
    if err != nil {
        log.Fatal(err)
    }
    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }
}
