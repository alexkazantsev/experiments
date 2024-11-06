package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type User struct {
	first  string
	last   string
	age    uint8
	domain string
	email  string
}

func (u User) filter() bool {
	return strings.Contains(u.email, "@k")
}

func (u User) toRecord() []string {
	return []string{
		u.first,
		u.last,
		strconv.FormatUint(uint64(u.age), 10),
		u.domain,
		u.email,
	}

}

// 1. read
// 2. transform
// 3. filter
// 4. write
func main() {
	var (
		files       = []string{"./filereader/files/1.csv", "./filereader/files/2.csv", "./filereader/files/3.csv"}
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	)

	defer cancel()

	write(
		filter(
			transform(
				read(ctx, files),
			),
		),
	)
}

func read(c context.Context, input []string) (context.Context, <-chan []string) {
	var (
		out    = make(chan []string)
		g, ctx = errgroup.WithContext(c)
	)

	for _, path := range input {
		g.Go(func() error {
			var (
				err  error
				file *os.File
			)

			if file, err = os.Open(path); err != nil {
				return err
			}
			defer file.Close()

			reader := csv.NewReader(file)

			// skip the header
			if _, err = reader.Read(); err != nil {
				return err
			}

			for {
				record, err := reader.Read()
				if err != nil {
					if err == io.EOF {
						break
					}

					return err
				}

				out <- record

				// simulate latency
				// time.Sleep(500 * time.Millisecond)
			}

			return nil
		})
	}

	go func() {
		if err := g.Wait(); err != nil {
			fmt.Printf("errgroup err: %+v\n", err)
		}

		close(out)
	}()

	return ctx, out
}

func transform(ctx context.Context, in <-chan []string) (context.Context, <-chan User) {
	var out = make(chan User)

	go func() {
		for {
			select {
			case line, done := <-in:

				if !done {
					close(out)

					return
				}

				var user = User{
					first:  line[0],
					last:   line[1],
					domain: line[3],
					email:  line[4],
				}

				if age, err := strconv.ParseUint(line[2], 10, 0); err == nil {
					user.age = uint8(age)
				}

				out <- user

			case <-ctx.Done():
				close(out)

				return
			}
		}
	}()

	return ctx, out
}

func filter(ctx context.Context, in <-chan User) (context.Context, <-chan User) {
	var out = make(chan User)

	go func() {
		for {
			select {
			case u, ok := <-in:

				if !ok {
					close(out)

					return
				}

				if u.filter() {
					out <- u
				}
			case <-ctx.Done():
				close(out)

				return
			}
		}
	}()

	return ctx, out
}

func printF(in <-chan User) {
	for {
		user, ok := <-in
		if !ok {
			return
		}

		fmt.Printf("user: %+v\n", user)
	}
}

func write(ctx context.Context, in <-chan User) {
	var (
		path = "./filereader/files/out.csv"
	)

	if isExist(path) {
		_ = os.Remove(path)
	}

	file, _ := os.Create(path)
	writer := csv.NewWriter(file)

	for {
		select {
		case a, ok := <-in:
			if !ok {
				writer.Flush()

				return
			}

			if err := writer.Write(a.toRecord()); err != nil {
				fmt.Printf("can not write a record: %+v", a)
			}
		case <-ctx.Done():
			writer.Flush()

			return
		}
	}
}

func isExist(path string) bool {
	var (
		file *os.File
		err  error
	)

	if file, err = os.Open(path); err != nil {
		return true
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	return false
}
