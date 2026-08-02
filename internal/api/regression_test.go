package api_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"
)

var update = flag.Bool("update", false, "実際のテスト結果でtxtarファイルのgoldenファイルを更新する")

// TestRegression はtestdataディレクトリ配下のtxtarファイルをリグレッションテストとして実行する
// 各テストケースはtestdata/<OperationName>/<testcase_name>.txtarに配置する
// `go test ./internal/api -update`コマンドで実際のテスト結果からgoldenファイルを一括で更新できる
func TestRegression(t *testing.T) {
	operationDirs, err := os.ReadDir("testdata")
	require.NoError(t, err)

	for _, operationDir := range operationDirs {
		if !operationDir.IsDir() {
			continue
		}

		files, err := os.ReadDir(filepath.Join("testdata", operationDir.Name()))
		require.NoError(t, err)

		for _, file := range files {
			name, ok := strings.CutSuffix(file.Name(), ".txtar")
			if !ok {
				continue
			}
			t.Run(operationDir.Name()+"/"+name, func(t *testing.T) {
				runTestCase(t, filepath.Join("testdata", operationDir.Name(), file.Name()))
			})
		}
	}
}

func runTestCase(t *testing.T, path string) {
	t.Helper()

	archive, err := txtar.ParseFile(path)
	require.NoError(t, err)

	files := make(map[string]string, len(archive.Files))
	for _, f := range archive.Files {
		files[f.Name] = string(f.Data)
	}

	if file, ok := files["setup.sql"]; ok {
		require.NoError(t, tdb.TruncateAll(t.Context()))
		require.NoError(t, tdb.ExecScript(t.Context(), file))
	}

	requestFile, ok := files["request"]
	require.True(t, ok, "%s must contain %q file", path, "request")
	resp, err := ts.Client().Do(parseRequestFile(t, requestFile))
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	gotResponse := generateResponseGolden(t, resp, body)
	var gotDB string
	if file, ok := files["db.golden"]; ok {
		gotDB = generateDBGolden(t, file)
	}

	if *update {
		setArchiveFile(archive, "response.golden", gotResponse)
		if gotDB != "" {
			setArchiveFile(archive, "db.golden", gotDB)
		}
		writeTestCase(t, path, archive)
		return
	}

	require.Contains(t, files, "response.golden", "%s must contain %q file", path, "response.golden")
	assertGolden(t, "Response", gotResponse, files["response.golden"])
	if gotDB != "" {
		assertGolden(t, "Database state", gotDB, files["db.golden"])
	}
}

func parseRequestFile(t *testing.T, text string) *http.Request {
	t.Helper()

	vars := map[string]string{"TOKEN": token}
	text = os.Expand(text, func(k string) string { return vars[k] })

	head, body, _ := strings.Cut(text, "\n\n")
	lines := strings.Split(strings.TrimSpace(head), "\n")
	method, path, ok := strings.Cut(lines[0], " ")
	require.True(t, ok, `request file must start with a "<method> <path>" line`)

	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)
	for _, line := range lines[1:] {
		key, value, ok := strings.Cut(line, ":")
		require.True(t, ok, "invalid header line: %q", line)
		req.Header.Set(strings.TrimSpace(key), strings.TrimSpace(value))
	}
	return req
}

func generateDBGolden(t *testing.T, dbGolden string) string {
	t.Helper()

	var b strings.Builder
	for line := range strings.Lines(dbGolden) {
		query, ok := strings.CutPrefix(strings.TrimSpace(line), "> ")
		if !ok {
			continue
		}
		b.WriteString("> " + query + "\n")

		result, err := tdb.DumpQueryResult(t.Context(), query)
		require.NoError(t, err)
		b.WriteString(result)
	}
	return b.String()
}

// generateResponseGolden はレスポンスの実測値からresponse.goldenの内容を生成する
// ヘッダはすべてダンプして検証する
// ただし、実行ごとに値が変わるDateとボディの検証と重複するContent-Lengthは除外する
func generateResponseGolden(t *testing.T, resp *http.Response, body []byte) string {
	t.Helper()

	var b strings.Builder
	b.WriteString(strconv.Itoa(resp.StatusCode) + "\n")
	for _, name := range slices.Sorted(maps.Keys(resp.Header)) {
		if name == "Date" || name == "Content-Length" {
			continue
		}
		for _, value := range resp.Header[name] {
			b.WriteString(name + ": " + value + "\n")
		}
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return b.String()
	}
	var buf bytes.Buffer
	require.NoError(t, json.Indent(&buf, body, "", "  "), "response body: %s", body)
	return b.String() + "\n" + buf.String() + "\n"
}

func setArchiveFile(archive *txtar.Archive, name, data string) {
	for i := range archive.Files {
		if archive.Files[i].Name == name {
			archive.Files[i].Data = []byte(data)
			return
		}
	}
	archive.Files = append(archive.Files, txtar.File{Name: name, Data: []byte(data)})
}

func writeTestCase(t *testing.T, path string, archive *txtar.Archive) {
	t.Helper()

	if comment := strings.TrimRight(string(archive.Comment), "\n"); comment != "" {
		archive.Comment = []byte(comment + "\n\n")
	}
	for i, f := range archive.Files {
		data := strings.TrimRight(string(f.Data), "\n") + "\n"
		if i < len(archive.Files)-1 {
			data += "\n"
		}
		archive.Files[i].Data = []byte(data)
	}
	require.NoError(t, os.WriteFile(path, txtar.Format(archive), 0o644))
}

func assertGolden(t *testing.T, name, got, want string) {
	t.Helper()

	if diff := cmp.Diff(strings.TrimRight(got, "\n"), strings.TrimRight(want, "\n")); diff != "" {
		t.Errorf("%s mismatch (-got +want):\n%s", name, diff)
	}
}
