package runner

import (
	"testing"
)

func TestRunner_Watch(t *testing.T) {
	//scripts := make([]*script.Script, 0)
	//
	//scriptsMgr := &scriptManagerMock{}
	//scriptsMgr.On("Get").Return(scripts, nil)
	//
	//ctx, ctxCancel := context.WithCancel(context.Background())
	//var wg *sync.WaitGroup
	//
	//rnr := &Runner{
	//	pool:           make(map[string]*Job),
	//	scriptsManager: scriptsMgr,
	//	logger:         zap.NewNop(),
	//}
	//
	//time.AfterFunc(time.Millisecond*200, func() {
	//	ctxCancel()
	//})
	//
	//var canceled bool
	//
	//go func() {
	//	rnr.Watch(ctx, ctxCancel, wg, false)
	//	canceled = true
	//}()
	//
	//<-time.After(time.Millisecond * 500)
	//assert.True(t, canceled)
	//
	//scriptsMgr.AssertCalled(t, "Get")
	//scriptsMgr.AssertExpectations(t)
}

func TestRunner_Watch_Error(t *testing.T) {
	//scripts := make([]*script.Script, 0)
	//
	//e := fmt.Errorf("error1")
	//
	//scriptsMgr := &scriptManagerMock{}
	//scriptsMgr.On("Get").Return(scripts, e)
	//
	//ctx, ctxCancel := context.WithCancel(context.Background())
	//var wg *sync.WaitGroup
	//
	//core, logs := observer.New(zap.DebugLevel)
	//logger := zap.New(core)
	//
	//rnr := &Runner{
	//	updateInterval: time.Second,
	//	pool:           make(map[string]*Job),
	//	scriptsManager: scriptsMgr,
	//	logger:         logger,
	//}
	//
	//time.AfterFunc(time.Millisecond*200, func() {
	//	ctxCancel()
	//})
	//
	//var canceled bool
	//
	//go func() {
	//	rnr.Watch(ctx, ctxCancel, wg, false)
	//	canceled = true
	//}()
	//
	//<-time.After(time.Millisecond * 500)
	//assert.True(t, canceled)
	//
	//scriptsMgr.AssertCalled(t, "Get")
	//scriptsMgr.AssertExpectations(t)
	//
	//assert.Equal(t, 1, logs.FilterMessage("error get scripts").FilterField(zap.Error(e)).Len())
}
