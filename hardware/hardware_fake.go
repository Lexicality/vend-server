// +build !linux,!arm

package hardware

import (
	"context"
	"sync"
	"time"

	logging "github.com/op/go-logging"

	"github.com/lexicality/vending/vend"
)

type physicalHardware struct {
	sync.Mutex
	log *logging.Logger
}

func (hw *physicalHardware) Setup(ctx context.Context) error {
	if hw.log != nil {
		hw.log.Info("Hello I'm not ARM!")
	}
	return nil
}

func (hw *physicalHardware) Vend(ctx context.Context, location uint8) vend.Result {
	hw.Lock()
	defer hw.Unlock()
	err := ctx.Err()
	if err != nil {
		if hw.log != nil {
			hw.log.Warningf("Canceling vend attempt: %s", err)
		}
		return vend.ResultAborted
	}

	if hw.log != nil {
		hw.log.Infof("Starting simulated vend of item %d", location)
	}

	<-time.After(time.Second * 15)

	if hw.log != nil {
		hw.log.Infof("Completed simulated vend of item %d", location)
	}

	return vend.ResultSuccess
}
