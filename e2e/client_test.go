package e2e

import (
	"testing"
	"github.com/gt-tallinn/gte-backend/client"
	"strconv"
	"time"
	"fmt"
	"sync"
)

func TestAdd(t *testing.T) {
	requestID := "reqID" + strconv.Itoa(int(time.Now().UnixNano()))
	fmt.Println(requestID)
	serviceName := "service" + strconv.Itoa(int(time.Now().UnixNano()))
	fmt.Println(serviceName)
	cl := client.New("http://localhost:9999/add", serviceName)
	ctx := cl.GetCtx(requestID, "receivedCtx.receivedCtx", "component1")
	ctx2 := ctx.GetCtx("component2")
	ctx2.Finish()
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		serviceName := "service1" + strconv.Itoa(int(time.Now().UnixNano()))
		fmt.Println(serviceName)
		cl1 := client.New("http://localhost:9999/add", serviceName)
		ctx1 := cl1.GetCtx(ctx2.GetReqID(), ctx2.GetCtxString(), "component11")
		ctx12 := ctx1.GetCtx("component12")
		ctx12.Finish()
		ctx1.Finish()
	}()
	ctx.Finish()
	wg.Wait()
	time.Sleep(5 * time.Second)
}
