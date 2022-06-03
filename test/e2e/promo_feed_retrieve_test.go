package e2e

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/goprom/internal/domain"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestFlow(t *testing.T) {
	wd, _ := os.Getwd()
	td := path.Join(wd, "_testdata")
	promotions := path.Join(wd, "_testdata", "promotions.csv")

	f, err := os.Create(promotions)
	if err != nil {
		log.Fatalf("failed to create promotions file %s", promotions)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	_, _ = w.WriteString("d018ef0b-dbd9-48f1-ac1a-eb4d90e57118,60.683466,2018-08-04 05:32:31 +0200 CEST\n")
	_, _ = w.WriteString("e2649ca5-7e05-4d53-a8ff-919917a4922e,66.640497,2018-08-22 18:34:11 +0200 CEST")
	_ = w.Flush()

	id := "d018ef0b-dbd9-48f1-ac1a-eb4d90e57118"
	ticker := time.NewTicker(10 * time.Second)
	timeout := time.After(2 * time.Minute)
	ch := make(chan bool)

	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("checking if promotions were imported ...")
				files, err := ioutil.ReadDir(td)
				if err != nil {
					log.Fatal(err)
				}

				if len(files) > 2 {
					// assume that the process has ended
					ch <- true
					return
				}

			case <-timeout:
				ch <- false
				return
			}

			time.Sleep(10 * time.Second)
		}
	}()

	c := <-ch
	if !c {
		log.Fatal("could ensure that the promotions were imported")
	}

	log.Println("promotions were imported. continuing with tests")

	res, err := http.Get(fmt.Sprintf("http://localhost:8080/promotions/%s", id))
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	promo := &domain.Promotion{}
	err = json.NewDecoder(res.Body).Decode(promo)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, promo.ID, id)
	assert.Equal(t, promo.Price, 60.68)
	assert.Equal(t, promo.ExpirationDate, "2018-08-04 05:32:31")

	files, err := ioutil.ReadDir(td)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), ".txt") || strings.Contains(file.Name(), "--imported") {
			err = os.Remove(path.Join(td, file.Name()))
			if err != nil {
				log.Fatalf("error removing test data file %s. reason: %s", file.Name(), err.Error())
			}
		}
	}

	log.Println("finished testing main workflow")
}
