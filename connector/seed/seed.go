// +build ignore

package main

import (
	"bufio"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"os"
)

// This file generates 002_customer_database_seed.sql

func seed() error {
	file, err := os.Create("./002_customer_database_seed.sql")
	if err != nil {
		return err
	}

	defer file.Close()
	w := bufio.NewWriter(file)

	fmt.Fprintln(w, "-- Do not edit, this is a generated file. See seed.go")
	fmt.Fprintln(w, "BEGIN;")

	schema := "(processor, card_id, card_number, value, ssn, vendor)"
	fmt.Fprintf(w, "insert into transactions %s values\n", schema)
	iterations := 10000

	for i := 0; i < iterations; i++ {

		cc := gofakeit.CreditCard()
		ssn := gofakeit.SSN()
		job := gofakeit.Job()
		price := gofakeit.Price(0.99, 2000)

		fmt.Fprintf(w, "\t('%s', %d, %d, %f, %s, '%s')", cc.Type, i, cc.Number, price, ssn, job.Company)
		if i < iterations-1 {
			fmt.Fprintf(w, ",")
		}

		fmt.Fprintf(w, "\n")
	}

	fmt.Fprintln(w, ";")

	fmt.Fprintln(w, "COMMIT;\n")
	fmt.Fprintln(w, "---- create above / drop below ----\n")
	fmt.Fprintln(w, "BEGIN;")
	fmt.Fprintln(w, "delete from transactions;")
	fmt.Fprintln(w, "COMMIT;")
	return w.Flush()
}

func main() {
	err := seed()
	if err != nil {
		panic(err)
	}
}
