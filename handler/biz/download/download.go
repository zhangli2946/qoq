package download

import (
	"qoq/actor/emq"
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/hashicorp/go-getter"
)

type download struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func init() {
	emq.DefaultHandler["download"] = Download
}

// Download handler
func Download(i *emq.Emq, cmd json.RawMessage) (err error) {
	var d download
	var pwd string
	err = json.Unmarshal(cmd, &d)
	if err != nil {
		return
	}
	pwd, err = os.Getwd()
	if err != nil {
		return
	}
	ctx, cancel := context.WithCancel(i.Ctx)
	client := &getter.Client{
		Ctx:     ctx,
		Src:     d.Src,
		Dst:     d.Dst,
		Pwd:     pwd,
		Mode:    getter.ClientModeAny,
		Options: []getter.ClientOption{},
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	errChan := make(chan error, 2)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := client.Get(); err != nil {
			errChan <- err
		}
	}()
	select {
	case <-ctx.Done():
		wg.Wait()
	case err = <-errChan:
		wg.Wait()
		return
	}
	return nil
}
