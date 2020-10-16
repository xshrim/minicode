package main

import "kubecp/kubectl"

func main() {
	restconfig, err, k8sclient := kubectl.InitRestClient()
	if err != nil {
		return
	}
	cp := kubectl.NewCopyer("cattle-node-agent-gl8s5", "cattle-system", "agent", restconfig, k8sclient)
	_ = cp.FromPod("/etc/sysctl.conf", "/home/xshrim/code/kubecp/")
	// cp.ToPod(...)
}
