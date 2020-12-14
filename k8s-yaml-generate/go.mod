module github.com/jiashiran/k8s-yaml-generate

go 1.13

replace k8s.io/api v0.0.0-20191121175643-4ed536977f46 => github.com/kubernetes/api v0.0.0-20191121175643-4ed536977f46

require (
	gopkg.in/yaml.v2 v2.2.7 // indirect
	k8s.io/api v0.0.0-20191121175643-4ed536977f46
	sigs.k8s.io/yaml v1.1.0
)
