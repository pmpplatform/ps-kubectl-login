package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var logger = log.New(os.Stdout, "", log.LUTC)

const (
	configFile       = ".kubectl-login.json"
	remoteConfigFile = "https://raw.githubusercontent.com/pmpplatform/ps-kubectl-login/master/.kubectl-login.json"
	timeout          = time.Second * 120
)

type configuration struct {
	DexURL  string   `json:"dex-url"`
	Aliases []string `json:"aliases"`
}

type app struct {
	cluster   string
	namespace string
}

var update, list bool

var cmd = &cobra.Command{
	Use:   "kubectl login [namespace]",
	Short: "Authenticates users against OIDC and writes the required kubeconfig.",
	Long:  "",
	// Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if update {
			fmt.Println("Updating config...")
			err := downloadFile(configFile, remoteConfigFile)
			if err != nil {
				logger.Printf("error: failed to download config file: %v", err)
				return err
			}
		}
		if list {
			rawConfig := getRawConfig()
			fmt.Printf("%-16s  %s\n", "Cluster", "Alias")
			for k, v := range rawConfig {
				for _, alias := range v.Aliases {
					fmt.Printf("%-16s: %s\n", k, alias)
				}
			}
			os.Exit(1)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
		var a app
		rawConfig := getRawConfig()
		alias := getAlias(os.Args[1:])
		config, cluster := getConfigByAlias(alias, rawConfig)
		a.namespace = alias
		a.cluster = cluster
		timer := time.AfterFunc(timeout, func() {
			log.Printf("\nLogin timeout... exiting")
			os.Exit(0)
		})
		defer timer.Stop()
		a.switchContext()
		if isLoggedIn() {
			log.Printf("Logged in: %v", cluster)
			os.Exit(0)
		}
		return login(cluster, config.DexURL)
	},
}

func init() {
	cmd.Flags().BoolVarP(&update, "update", "u", false, "update config file from github")
	cmd.Flags().BoolVarP(&list, "list", "l", false, "list aliases")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
}

func downloadFile(filename string, url string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath.Join(home, filename))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func getRawConfig() map[string]*configuration {
	configPath := os.Getenv("HOME") + string(os.PathSeparator) + configFile

	file, err := os.Open(configPath)
	if err != nil {
		logger.Fatalf("error: cannot open config file at %s: %v", configPath, err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		closeFile(file)
		logger.Fatalf("error: cannot read config file at %s: %v", configPath, err)
	}

	var cfg map[string]*configuration
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		closeFile(file)
		logger.Fatalf("error: cannot unmarshal contents of config file at %s: %v", configPath, err)
	}

	closeFile(file)
	return cfg
}

func getAlias(args []string) string {
	if len(args) == 0 {
		logger.Fatalf("Alias is mandatory i.e %s. try '%s' to get this value.",
			Bold(Cyan("kubectl-login <ALIAS>")), Bold(Cyan("cat $HOME/"+configFile)))
	}
	return args[0]
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		logger.Printf("warning: couldn't close config file: %v", err)
	}
}

func getConfigByAlias(alias string, rawConfig map[string]*configuration) (*configuration, string) {
	for k, v := range rawConfig {
		if containsAlias(v, alias) {
			return v, k
		}
	}
	logger.Fatalf("Alias \"%s\" not found. Try '%s' to get this value.",
		Bold(Cyan(alias)), Bold(Cyan("cat $HOME/"+configFile)))
	return nil, ""
}

func containsAlias(c *configuration, s string) bool {
	for _, val := range c.Aliases {
		if val == s {
			return true
		}
	}
	return false
}

func RunWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}

func login(cluster string, url string) error {
	for {
		fmt.Printf("Logging in to cluster %s\n", cluster)

		dir, err := os.MkdirTemp("", "chromedp-example")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(dir)
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.DisableGPU,
			chromedp.UserDataDir(dir),
			chromedp.Flag("headless", false),
		)

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		browserCtx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		// bail out if the browser is closed
		chromedp.ListenBrowser(browserCtx, func(ev interface{}) {
			if ev, ok := ev.(*target.EventTargetDestroyed); ok {
				if c := chromedp.FromContext(browserCtx); c != nil {
					if c.Target.TargetID == ev.TargetID {
						cancel()
					}
				}
			}
		})

		var res string

		err = chromedp.Run(
			browserCtx,
			RunWithTimeOut(&browserCtx, 300, chromedp.Tasks{
				chromedp.Navigate(url),
				chromedp.WaitVisible(`//*[@value="Login To Cluster"]`),
				chromedp.Click(`//*[@value="Login To Cluster"]`, chromedp.BySearch),
				chromedp.Text(`#idMergeConfig`, &res, chromedp.NodeVisible),
			}),
		)
		if err != nil {
			cancel()
			log.Fatal(err)
		}

		if res != "" {
			err := kubeMerge(strings.TrimSpace(res))
			if err != nil {
				return err
			}
			log.Printf("Logged in: %v", cluster)
			return nil
		}
	}
}

func kubeMerge(config string) error {
	cmd := exec.Command("/bin/sh", "-c", config)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func isLoggedIn() bool {
	err := exec.Command("kubectl", "get", "namespace").Run()
	return err == nil
}

func (a *app) switchContext() {
	clusterArg := fmt.Sprintf("--cluster=%s", a.cluster)
	user := fmt.Sprintf("--user=%s", a.cluster)
	cmd := exec.Command("kubectl", "config", "set-context", a.cluster, user, clusterArg, "--namespace="+a.namespace)
	err := cmd.Run()
	if err != nil {
		logger.Fatalf("error: cannot set kubectl login context: %v", err)
	}

	cmd = exec.Command("kubectl", "config", "use-context", a.cluster)
	err = cmd.Run()
	if err != nil {
		logger.Fatalf("error: cannot switch to kubectl login context: %v", err)
	}
}
