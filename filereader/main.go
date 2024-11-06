package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

// 1. read
// 2. transform
// 3. filter
// 4. write
func main() {
	var (
		files       = []string{"./filereader/files/1.csv", "./filereader/files/2.csv", "./filereader/files/3.csv"}
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		done        = make(chan bool, 1)
	)

	defer close(done)
	defer cancel()

	chR := read(ctx, files)
	chT := transform(ctx, chR)
	chF := filter(ctx, chT)
	chP := printF(ctx, chF)

	write(ctx, chP, done)

	select {
	case <-ctx.Done():
		fmt.Println("ctx done")
	case <-done:
		fmt.Println("all jobs done")
	}
}

func read(c context.Context, input []string) <-chan []string {
	var (
		out  = make(chan []string)
		g, _ = errgroup.WithContext(c)
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

	return out
}

func transform(ctx context.Context, in <-chan []string) <-chan User {
	var out = make(chan User)

	go func() {
		defer close(out)

		for {
			select {
			case line, done := <-in:

				if !done {
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
				return
			}
		}
	}()

	return out
}

func filter(ctx context.Context, in <-chan User) <-chan User {
	var out = make(chan User)

	go func() {
		defer close(out)

		for {
			select {
			case u, ok := <-in:

				if !ok {
					return
				}

				if u.filter() {
					out <- u
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func printF(ctx context.Context, in <-chan User) <-chan User {
	var out = make(chan User)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case user, ok := <-in:
				if !ok {
					return
				}

				fmt.Printf("user: %+v\n", user)

				out <- user
			}
		}
	}()

	return out
}

func write(ctx context.Context, in <-chan User, done chan<- bool) {
	var (
		path = "./filereader/files/out.csv"
	)

	if isExist(path) {
		_ = os.Remove(path)
	}

	file, _ := os.Create(path)
	writer := csv.NewWriter(file)

loop:
	for {
		select {
		case a, ok := <-in:
			if !ok {
				break loop
			}

			if err := writer.Write(a.toRecord()); err != nil {
				fmt.Printf("can not write a record: %+v", a)
			}

		case <-ctx.Done():
			_ = os.Remove(path)
			return
		}
	}

	writer.Flush()
	done <- true
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

type User struct {
	first  string
	last   string
	age    uint8
	domain string
	email  string
}

func (u User) filter() bool {
	return u.age <= 20
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
