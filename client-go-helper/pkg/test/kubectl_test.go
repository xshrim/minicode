package test

import (
	"log"
	"testing"

	"github.com/ica10888/client-go-helper/pkg/kubectl"
)

var yaml = `
apiVersion: apps/v1 # for versions before 1.9.0 use apps/v1beta2
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2 # tells deployment to run 2 pods matching the template
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
`

func TestApply(t *testing.T) {
	e := kubectl.Apply(yaml, "default")
	if e != nil {
		log.Print(e)
	}

}

func TestCreate(t *testing.T) {
	e := kubectl.Create(yaml, "default")
	if e != nil {
		log.Print(e)
	}
}

func TestExec(t *testing.T) {
	e := kubectl.Exec(yaml, "local")
	if e != nil {
		log.Print(e)
	}
}
