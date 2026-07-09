package version

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Current 当前版本号（构建时可通过 ldflags 覆盖）。
var Current = "0.1.0"

// Repo 仓库名（owner/repo），构建期通过 ldflags 注入并打包进 exe，
// 例如：-X mflow-go/internal/version.Repo=niuhuan/mflow-go
var Repo = ""

// repo 返回仓库名：优先使用编译期注入的 Repo，其次回退到环境变量（便于开发调试）。
func repo() string {
	if Repo != "" {
		return Repo
	}
	return os.Getenv("GITHUB_REPOSITORY")
}

// ReleaseURL 返回发布页地址。
func ReleaseURL() string {
	r := repo()
	if r == "" {
		return ""
	}
	return "https://github.com/" + r + "/releases"
}

func latestReleaseTag() (string, error) {
	r := repo()
	if r == "" {
		return "", nil
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/"+r+"/releases/latest", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "mflow-go")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	return data.TagName, nil
}

var semverRe = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

func parse(v string) (int, int, int, bool) {
	m := semverRe.FindStringSubmatch(v)
	if m == nil {
		return 0, 0, 0, false
	}
	a, _ := strconv.Atoi(m[1])
	b, _ := strconv.Atoi(m[2])
	c, _ := strconv.Atoi(m[3])
	return a, b, c, true
}

// NewVersion 若有更新则返回新版本号，否则返回空串。
func NewVersion() string {
	tag, err := latestReleaseTag()
	if err != nil || tag == "" {
		return ""
	}
	latest := strings.TrimPrefix(tag, "v")
	la, lb, lc, ok1 := parse(latest)
	ca, cb, cc, ok2 := parse(strings.TrimPrefix(Current, "v"))
	if !ok1 || !ok2 {
		return ""
	}
	if la > ca || (la == ca && lb > cb) || (la == ca && lb == cb && lc > cc) {
		return latest
	}
	return ""
}
