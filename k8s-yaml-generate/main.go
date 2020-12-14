package main

import (
	"encoding/json"
	"fmt"
	"github.com/jiashiran/k8s-yaml-generate/util"
	"io/ioutil"
	av1 "k8s.io/api/apps/v1"
	bv1 "k8s.io/api/batch/v1"
	cv1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"log"
	"net/http"
	"sigs.k8s.io/yaml"
	"strings"
)

var (
	yaml_types = make(map[string]interface{})
)

func init() {
	configMap := cv1.ConfigMap{}
	yaml_types["ConfigMap"] = configMap
	ingress := v1beta1.Ingress{}
	yaml_types["Ingress"] = ingress
	service := cv1.Service{}
	yaml_types["Service"] = service
	replicaSet := av1.ReplicaSet{}
	yaml_types["ReplicaSet"] = replicaSet
	deployment := av1.Deployment{}
	yaml_types["Deployment"] = deployment
	secret := cv1.Secret{}
	yaml_types["Secret"] = secret
	namespace := cv1.Namespace{}
	yaml_types["Namespace"] = namespace
	pod := cv1.Pod{}
	yaml_types["Pod"] = pod
	volume := cv1.Volume{}
	yaml_types["Volume"] = volume
	job := bv1.Job{}
	yaml_types["Job"] = job
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {
		methodType := request.Header.Get("type")
		switch methodType {
		case "getK8sTypes":
			types := make([]string, 0)
			for k, _ := range yaml_types {
				types = append(types, k)
			}
			bs, _ := json.Marshal(types)
			writer.Write(bs)
			break
		case "generate":
			yaml_type := request.Header.Get("yaml_type")
			in := yaml_types[yaml_type]
			fmt.Println("in:", yaml_type, in)
			yamlStr := toYamlString(in)
			if _, ok := yaml_types[yaml_type]; ok {
				writer.Write([]byte(yamlStr))
			} else {
				writer.Write([]byte("no type:" + yaml_type))
			}
			break
		case "tidy":
			yamlStr, _ := ioutil.ReadAll(request.Body)
			fmt.Println("yamlStr:", string(yamlStr))
			yamlString := tidy(string(yamlStr))
			writer.Write([]byte(yamlString))
			break
		}
	})
	http.ListenAndServe(":8080", nil)

}

func tidy(yamlStr string) string {
	defer func() {
		if err := recover(); err != nil {
			log.Println("OrganizeData err :", err)
		}
	}()
	lines := strings.Split(yamlStr, "\n")
	newYaml := ""
	for _, l := range lines {
		if strings.Trim(l, " ") != ":" && !strings.Contains(l, ": 44") && !strings.Contains(l, "{}") && !strings.Contains(l, "string-value") {
			newYaml = newYaml + l + "\n"
		}

	}
	//fmt.Println(newYaml)
	bs, err := util.ToJSON([]byte(newYaml))
	if err != nil {
		fmt.Println(err, newYaml)
	}
	lines = strings.Split(string(bs), "\n")
	//fmt.Println(lines)
	newYaml = ""
	for _, l := range lines {
		if strings.Trim(l, " ") != ":" && !strings.Contains(l, ": null") {
			newYaml = newYaml + l + "\n"
		}

	}
	//fmt.Println(newYaml)
	bs, err = yaml.JSONToYAML([]byte(newYaml))
	if err != nil {
		fmt.Println(err, newYaml)
	}
	//fmt.Println(string(bs))
	yamlStr = string(bs)
	lines = strings.Split(yamlStr, "\n")
	newYaml = ""
	for _, l := range lines {
		if strings.Trim(l, " ") != ":" && !strings.Contains(l, ": 44") && !strings.Contains(l, "{}") && !strings.Contains(l, "string-value") {
			newYaml = newYaml + l + "\n"
		}

	}
	//fmt.Println(newYaml)
	return newYaml
}

func toYamlString(object interface{}) string {
	value := util.InitializeStruct(object)
	//fmt.Println(value)
	bytes, err := json.Marshal(value.Interface())
	if err != nil {
		log.Println(err)
	}
	//fmt.Println(string(bytes))
	bytes, err = yaml.JSONToYAML(bytes)
	if err != nil {
		log.Println(err)
	}
	return string(bytes)
}
