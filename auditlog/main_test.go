package main

import (
	"auditlog/server"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSend(t *testing.T) {
	payload := `
  {
    "auditID": "ef1d249e-bfac-4fd0-a61f-cbdcad53b9bb",
    "requestURI": "/v3/project/c-bcz5t:p-fdr4s/workloads/deployment:default:nginx",
    "sourceIPs": [
        "::1"
    ],
    "user": {
        "name": "user-f4tt2",
        "group": [
            "system:authenticated"
        ]
    },
    "verb": "PUT",
    "stage": "RequestReceived",
    "stageTimestamp": "2018-07-20 10:28:08 +0800"
}
  `

	payload = `
  {
    "auditID": "ef1d249e-bfac-4fd0-a61f-cbdcad53b9bd",
    "requestURI": "/v3/project/c-bcz5t:p-fdr4s/workloads/deployment:default:nginx",
    "sourceIPs": [
        "::1"
    ],
    "user": {
        "name": "user-f4tt2",
        "group": [
            "system:authenticated"
        ]
    },
    "verb": "PUT",
    "stage": "RequestReceived",
    "stageTimestamp": "2018-07-20 10:28:08 +0800",
    "requestBody": {
        "hostIPC": false,
        "hostNetwork": false,
        "hostPID": false,
        "paused": false,
        "annotations": {},
        "baseType": "workload",
        "containers": [
            {
                "allowPrivilegeEscalation": false,
                "image": "nginx",
                "imagePullPolicy": "Always",
                "initContainer": false,
                "name": "nginx",
                "ports": [
                    {
                        "containerPort": 80,
                        "dnsName": "nginx-nodeport",
                        "kind": "NodePort",
                        "name": "80tcp01",
                        "protocol": "TCP",
                        "sourcePort": 0,
                        "type": "/v3/project/schemas/containerPort"
                    }
                ],
                "privileged": false,
                "readOnly": false,
                "resources": {
                    "type": "/v3/project/schemas/resourceRequirements",
                    "requests": {},
                    "limits": {}
                },
                "restartCount": 0,
                "runAsNonRoot": false,
                "stdin": true,
                "stdinOnce": false,
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "tty": true,
                "type": "/v3/project/schemas/container",
                "environmentFrom": [],
                "capAdd": [],
                "capDrop": [],
                "livenessProbe": null,
                "volumeMounts": []
            }
        ],
        "created": "2018-07-18T07:34:16Z",
        "createdTS": 1531899256000,
        "creatorId": null,
        "deploymentConfig": {
            "maxSurge": 1,
            "maxUnavailable": 0,
            "minReadySeconds": 0,
            "progressDeadlineSeconds": 600,
            "revisionHistoryLimit": 10,
            "strategy": "RollingUpdate"
        },
        "deploymentStatus": {
            "availableReplicas": 1,
            "conditions": [
                {
                    "lastTransitionTime": "2018-07-18T07:34:38Z",
                    "lastTransitionTimeTS": 1531899278000,
                    "lastUpdateTime": "2018-07-18T07:34:38Z",
                    "lastUpdateTimeTS": 1531899278000,
                    "message": "Deployment has minimum availability.",
                    "reason": "MinimumReplicasAvailable",
                    "status": "True",
                    "type": "Available"
                },
                {
                    "lastTransitionTime": "2018-07-18T07:34:16Z",
                    "lastTransitionTimeTS": 1531899256000,
                    "lastUpdateTime": "2018-07-18T07:34:38Z",
                    "lastUpdateTimeTS": 1531899278000,
                    "message": "ReplicaSet \"nginx-64d85666f9\" has successfully progressed.",
                    "reason": "NewReplicaSetAvailable",
                    "status": "True",
                    "type": "Progressing"
                }
            ],
            "observedGeneration": 2,
            "readyReplicas": 1,
            "replicas": 1,
            "type": "/v3/project/schemas/deploymentStatus",
            "unavailableReplicas": 0,
            "updatedReplicas": 1
        },
        "dnsPolicy": "ClusterFirst",
        "id": "deployment:default:nginx",
        "labels": {
            "workload.user.cattle.io/workloadselector": "deployment-default-nginx"
        },
        "name": "nginx",
        "namespaceId": "default",
        "projectId": "c-bcz5t:p-fdr4s",
        "publicEndpoints": [
            {
                "addresses": [
                    "10.64.3.58"
                ],
                "allNodes": true,
                "ingressId": null,
                "nodeId": null,
                "podId": null,
                "port": 30917,
                "protocol": "TCP",
                "serviceId": "default:nginx-nodeport",
                "type": "publicEndpoint"
            }
        ],
        "restartPolicy": "Always",
        "scale": 1,
        "schedulerName": "default-scheduler",
        "selector": {
            "matchLabels": {
                "workload.user.cattle.io/workloadselector": "deployment-default-nginx"
            },
            "type": "/v3/project/schemas/labelSelector"
        },
        "state": "active",
        "terminationGracePeriodSeconds": 30,
        "transitioning": "no",
        "transitioningMessage": "",
        "type": "deployment",
        "uuid": "f998037d-8a5c-11e8-a4cf-0245a7ebb0fd",
        "workloadAnnotations": {
            "deployment.kubernetes.io/revision": "1",
            "field.cattle.io/creatorId": "user-f4tt2"
        },
        "workloadLabels": {
            "workload.user.cattle.io/workloadselector": "deployment-default-nginx"
        },
        "scheduling": {
            "node": {}
        },
        "description": "my description",
        "volumes": []
    }
}
  `
	req, err := http.NewRequest("POST", "http://127.0.0.1:9090/set", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(resp, err)
}

func TestGet(t *testing.T) {
	rjson := `
  {
    "user": "tom",
    "scope": 1,
    "timestamp": 1603362532083
  }
  `

	rjson = `
  {
    "scope": 1
  }
  `
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:9090/get", bytes.NewReader([]byte(rjson)))
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	fmt.Println(err)
	defer response.Body.Close()
	var resp []server.Response
	respBodyBytes, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(respBodyBytes, &resp)
	fmt.Println(resp)
}
