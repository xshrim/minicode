package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jonluo94/baasmanager/baas-core/common/log"
	"github.com/spf13/viper"
)

var Config *viper.Viper

var logger = log.GetLogger("ido.config", log.INFO)

func init() {
	//监听改变动态跟新配置
	go watchConfig()
	//加载配置
	loadConfig()
}

//监听配置改变
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Info("Config file changed:", e.Name)
		//改变重新加载
		loadConfig()
	})
}

//加载配置
// TODO 配置当前项目环境变量
func loadConfig() {
	prefix := "IDO" // 环境变量以IDO为前缀, 如IDO_ROOTPATH
	envs := []string{"FABRICENGINEPORT", "ROOTPATH", "NFSSHARED", "NFSREALPATH", "NFSSERVER", "NFSCAPACITY", "KUBEENGINE", "TEMPLATE", "VMENDPOINT", "CHAINCODEGITHUB", "MSPTYPE", "BATCHTIMEOUT", "MAXMESSAGECOUNT", "ABSOLUTEMAXBYTES", "PREFERREDMAXBYTES", "CLUSTERID", "PROJECTID", "NODESELECTOR", "CPULIMIT", "MEMORYLIMIT"}
	viper.SetConfigName("feconfig")              // name of kubeconfig file
	viper.AddConfigPath(".")                     // optionally look for kubeconfig in the working directory
	viper.AddConfigPath("/etc/baas")             // path to look for the kubeconfig file in
	if err := viper.ReadInConfig(); err != nil { // Find and read the feconfig.yaml file
		// Handle errors reading the kubeconfig file
		logger.Errorf("Error reading config file: %s \n", err)
	}

	viper.SetEnvPrefix(prefix)
	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			logger.Errorf("Error reading environment: %s \n", err)
		}
	}
	viper.AutomaticEnv()
	//全局配置
	Config = viper.GetViper()
	logger.Infof("%v", Config.AllSettings())
}
