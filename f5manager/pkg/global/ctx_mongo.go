package global

// import (
// 	"context"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/casbin/casbin/v2"
// 	"github.com/xshrim/f5m/pkg/database"
// 	"github.com/xshrim/gol"
// 	"gorm.io/gorm"
// )

// var Ctx *GlobalContext

// type GlobalContext struct {
// 	DB       *gorm.DB
// 	Enforcer *casbin.Enforcer
// 	Context  context.Context
// 	Cancel   context.CancelFunc
// }

// // func init() {
// // 	Ctx = New()
// // }

// func New() *GlobalContext {

// 	ctx, cancel := context.WithCancel(context.Background())

// 	db, err := database.NewDB("mongodb://localhost:27017", "lbm", 60, 5)
// 	if err != nil {
// 		gol.Error(err)
// 		return nil
// 	}
// 	efc, err := NewEnforcer(db)
// 	if err != nil {
// 		gol.Error(err)
// 		return nil
// 	}

// 	// Load the policy from DB.
// 	efc.LoadPolicy()
// 	efc.AddPolicy("platform-admin", "*", ".*/.*", "*")
// 	efc.AddPolicy("platform-auditor", "*", ".*/.*", "read")
// 	efc.AddPolicy("project-admin", "cicd", ".*/.*", "*")
// 	efc.AddPolicy("namespace-admin", "baas", "region/dev", "*")
// 	efc.AddPolicy("namespace-developer", "echat", "region/dev", "read")
// 	efc.AddGroupingPolicy("tom", "project-admin", "*")
// 	// efc.SavePolicy()

// 	return &GlobalContext{
// 		DB:       db,
// 		Enforcer: efc,
// 		Context:  ctx,
// 		Cancel:   cancel,
// 	}
// }

// func (ctx *GlobalContext) GetDB() *gorm.DB {
// 	return ctx.DB
// }

// func (ctx *GlobalContext) Run(exitFunc func()) {
// 	c := make(chan os.Signal)
// 	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
// 	defer func() {
// 		signal.Stop(c)
// 		ctx.Cancel()
// 	}()

// 	for {
// 		select {
// 		case sig := <-c:
// 			switch sig {
// 			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
// 				gol.Info("Exit:", sig)
// 				if exitFunc != nil {
// 					exitFunc()
// 				}
// 			case syscall.SIGUSR1:
// 				gol.Info("Usr1:", sig)
// 			case syscall.SIGUSR2:
// 				gol.Info("Usr2:", sig)
// 			default:
// 				gol.Info("Other:", sig)
// 			}
// 			ctx.Cancel()
// 			return
// 		case <-ctx.Context.Done():
// 			gol.Info("Done")
// 			return
// 		}
// 	}
// }
