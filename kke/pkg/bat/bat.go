package bat

import (
	"fmt"
	"kke/pkg/utils"
	"strings"
	"sync"
)

type Host struct {
	Addr string
	Port int
	User string
	Cred string
}

type Response struct {
	Code int    `json:"code"`
	Out  string `json:"out"`
	Err  string `json:"err"`
}

type Task struct {
	Hosts  []Host
	Mode   string
	Args   []string
	Result map[string]Response
}

func run(idx int, host Host, mode string, args []string, resp *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	r := Response{}

	if mode == "ping" {
		if len(args) < 1 {
			if ok := utils.Ping(host.Addr, args...); ok {
				r.Code = 0
				r.Out = "Y"
			} else {
				r.Code = -1
				r.Out = "N"
			}
		} else {
			out, ok := utils.PortScan(host.Addr, args...)
			if ok {
				r.Code = 0
			} else {
				r.Code = -1
			}
			r.Out = out
		}

		resp.Store(idx, r)
		return
	}

	client, err := New(host)
	if err != nil {
		r.Code = -1
		r.Err = err.Error()
	} else {
		// var stdoutBuf bytes.Buffer
		// var stderrBuf bytes.Buffer
		// session.Stdout = &stdoutBuf
		// session.Stderr = &stderrBuf

		var out []byte
		switch mode {
		case "execute":
			out, err = client.Execute(strings.Join(args, " "))
		case "script":
			out, err = client.Script(strings.Join(args, " "))
		case "template":
			if len(args) < 3 {
				err = fmt.Errorf("no enough arguments for file template")
			} else {
				out, err = client.Template(args[0], args[1], args[2])
			}
		case "shell":
			_ = client.Shell(strings.Join(args, " "))
		case "push":
			if len(args) == 0 {
				err = fmt.Errorf("no local/remote path for file transfer")
			} else {
				localPath := args[0]
				remotePath := ""
				if len(args) > 1 {
					remotePath = args[1]
				}
				err = client.Push(localPath, remotePath)
			}
		case "pull":
			if len(args) == 0 {
				err = fmt.Errorf("no local/remote path for file transfer")
			} else {
				localPath := ""
				remotePath := args[0]
				if len(args) > 1 {
					localPath = args[1]
				}
				err = client.Pull(remotePath, localPath)
			}
		default:
			err = fmt.Errorf("unknown batch mode")
		}
		defer client.Close()

		if err != nil {
			r.Code = 1
			r.Err = err.Error()
		} else {
			r.Code = 0
		}
		r.Out = strings.TrimSpace(string(out))
	}

	resp.Store(idx, r)
}

func (task *Task) Do() *Task {
	var wg sync.WaitGroup
	resp := &sync.Map{}

	for idx, host := range task.Hosts {
		wg.Add(1)
		go run(idx, host, task.Mode, task.Args, resp, &wg)
	}

	wg.Wait()

	if task.Result == nil {
		task.Result = make(map[string]Response)
	}

	for idx, host := range task.Hosts {
		if out, ok := resp.Load(idx); ok {
			task.Result[fmt.Sprintf("%s@%s:%d", host.User, host.Addr, host.Port)] = out.(Response)
		}
	}

	return task
}

func (task *Task) Display() error {
	return nil
}
