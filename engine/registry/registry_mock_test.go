package registry // import "github.com/docker/docker/registry"

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/containerd/log"
	"github.com/docker/docker/api/types/registry"
	"gotest.tools/v3/assert"
)

var (
	testHTTPServer  *httptest.Server
	testHTTPSServer *httptest.Server
)

func init() {
	r := http.NewServeMux()

	// /v1/
	r.HandleFunc("/v1/_ping", handlerGetPing)
	r.HandleFunc("/v1/search", handlerSearch)

	// /v2/
	r.HandleFunc("/v2/version", handlerGetPing)

	testHTTPServer = httptest.NewServer(handlerAccessLog(r))
	testHTTPSServer = httptest.NewTLSServer(handlerAccessLog(r))

	// override net.LookupIP
	lookupIP = func(host string) ([]net.IP, error) {
		if host == "127.0.0.1" {
			// I believe in future Go versions this will fail, so let's fix it later
			return net.LookupIP(host)
		}
		mockHosts := map[string][]net.IP{
			"":            {net.ParseIP("0.0.0.0")},
			"localhost":   {net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
			"example.com": {net.ParseIP("42.42.42.42")},
			"other.com":   {net.ParseIP("43.43.43.43")},
		}
		for h, addrs := range mockHosts {
			if host == h {
				return addrs, nil
			}
			for _, addr := range addrs {
				if addr.String() == host {
					return []net.IP{addr}, nil
				}
			}
		}
		return nil, errors.New("lookup: no such host")
	}
}

func handlerAccessLog(handler http.Handler) http.Handler {
	logHandler := func(w http.ResponseWriter, r *http.Request) {
		log.G(context.TODO()).Debugf(`%s "%s %s"`, r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(logHandler)
}

func makeURL(req string) string {
	return testHTTPServer.URL + req
}

func makeHTTPSURL(req string) string {
	return testHTTPSServer.URL + req
}

func makeIndex(req string) *registry.IndexInfo {
	return &registry.IndexInfo{
		Name: makeURL(req),
	}
}

func makeHTTPSIndex(req string) *registry.IndexInfo {
	return &registry.IndexInfo{
		Name: makeHTTPSURL(req),
	}
}

func makePublicIndex() *registry.IndexInfo {
	return &registry.IndexInfo{
		Name:     IndexServer,
		Secure:   true,
		Official: true,
	}
}

func makeServiceConfig(mirrors []string, insecureRegistries []string) (*serviceConfig, error) {
	return newServiceConfig(ServiceOptions{
		Mirrors:            mirrors,
		InsecureRegistries: insecureRegistries,
	})
}

func writeHeaders(w http.ResponseWriter) {
	h := w.Header()
	h.Add("Server", "docker-tests/mock")
	h.Add("Expires", "-1")
	h.Add("Content-Type", "application/json")
	h.Add("Pragma", "no-cache")
	h.Add("Cache-Control", "no-cache")
}

func writeResponse(w http.ResponseWriter, message interface{}, code int) {
	writeHeaders(w)
	w.WriteHeader(code)
	body, err := json.Marshal(message)
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
		return
	}
<<<<<<< HEAD
	_, _ = w.Write(body)
=======
	w.Write(body)
}

func readJSON(r *http.Request, dest interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dest)
}

func apiError(w http.ResponseWriter, message string, code int) {
	body := map[string]string{
		"error": message,
	}
	writeResponse(w, body, code)
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v == %v", a, b)
	}
	t.Fatal(message)
}

// Similar to assertEqual, but does not stop test
func checkEqual(t *testing.T, a interface{}, b interface{}, messagePrefix string) {
	if a == b {
		return
	}
	message := fmt.Sprintf("%v != %v", a, b)
	if len(messagePrefix) != 0 {
		message = messagePrefix + ": " + message
	}
	t.Error(message)
}

// Similar to assertNotEqual, but does not stop test
func checkNotEqual(t *testing.T, a interface{}, b interface{}, messagePrefix string) {
	if a != b {
		return
	}
	message := fmt.Sprintf("%v == %v", a, b)
	if len(messagePrefix) != 0 {
		message = messagePrefix + ": " + message
	}
	t.Error(message)
}

func requiresAuth(w http.ResponseWriter, r *http.Request) bool {
	writeCookie := func() {
		value := fmt.Sprintf("FAKE-SESSION-%d", time.Now().UnixNano())
		cookie := &http.Cookie{Name: "session", Value: value, MaxAge: 3600}
		http.SetCookie(w, cookie)
		// FIXME(sam): this should be sent only on Index routes
		value = fmt.Sprintf("FAKE-TOKEN-%d", time.Now().UnixNano())
		w.Header().Add("X-Docker-Token", value)
	}
	if len(r.Cookies()) > 0 {
		writeCookie()
		return true
	}
	if len(r.Header.Get("Authorization")) > 0 {
		writeCookie()
		return true
	}
	w.Header().Add("WWW-Authenticate", "token")
	apiError(w, "Wrong auth", http.StatusUnauthorized)
	return false
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
}

func handlerGetPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeResponse(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	writeResponse(w, true, http.StatusOK)
}

<<<<<<< HEAD
=======
func handlerGetImage(w http.ResponseWriter, r *http.Request) {
	if !requiresAuth(w, r) {
		return
	}
	vars := mux.Vars(r)
	layer, exists := testLayers[vars["image_id"]]
	if !exists {
		http.NotFound(w, r)
		return
	}
	writeHeaders(w)
	layerSize := len(layer["layer"])
	w.Header().Add("X-Docker-Size", strconv.Itoa(layerSize))
	io.WriteString(w, layer[vars["action"]])
}

func handlerPutImage(w http.ResponseWriter, r *http.Request) {
	if !requiresAuth(w, r) {
		return
	}
	vars := mux.Vars(r)
	imageID := vars["image_id"]
	action := vars["action"]
	layer, exists := testLayers[imageID]
	if !exists {
		if action != "json" {
			http.NotFound(w, r)
			return
		}
		layer = make(map[string]string)
		testLayers[imageID] = layer
	}
	if checksum := r.Header.Get("X-Docker-Checksum"); checksum != "" {
		if checksum != layer["checksum_simple"] && checksum != layer["checksum_tarsum"] {
			apiError(w, "Wrong checksum", http.StatusBadRequest)
			return
		}
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apiError(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}
	layer[action] = string(body)
	writeResponse(w, true, http.StatusOK)
}

func handlerGetDeleteTags(w http.ResponseWriter, r *http.Request) {
	if !requiresAuth(w, r) {
		return
	}
	repositoryName, err := reference.WithName(mux.Vars(r)["repository"])
	if err != nil {
		apiError(w, "Could not parse repository", http.StatusBadRequest)
		return
	}
	tags, exists := testRepositories[repositoryName.String()]
	if !exists {
		apiError(w, "Repository not found", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodDelete {
		delete(testRepositories, repositoryName.String())
		writeResponse(w, true, http.StatusOK)
		return
	}
	writeResponse(w, tags, http.StatusOK)
}

func handlerGetTag(w http.ResponseWriter, r *http.Request) {
	if !requiresAuth(w, r) {
		return
	}
	vars := mux.Vars(r)
	repositoryName, err := reference.WithName(vars["repository"])
	if err != nil {
		apiError(w, "Could not parse repository", http.StatusBadRequest)
		return
	}
	tagName := vars["tag"]
	tags, exists := testRepositories[repositoryName.String()]
	if !exists {
		apiError(w, "Repository not found", http.StatusNotFound)
		return
	}
	tag, exists := tags[tagName]
	if !exists {
		apiError(w, "Tag not found", http.StatusNotFound)
		return
	}
	writeResponse(w, tag, http.StatusOK)
}

func handlerPutTag(w http.ResponseWriter, r *http.Request) {
	if !requiresAuth(w, r) {
		return
	}
	vars := mux.Vars(r)
	repositoryName, err := reference.WithName(vars["repository"])
	if err != nil {
		apiError(w, "Could not parse repository", http.StatusBadRequest)
		return
	}
	tagName := vars["tag"]
	tags, exists := testRepositories[repositoryName.String()]
	if !exists {
		tags = make(map[string]string)
		testRepositories[repositoryName.String()] = tags
	}
	tagValue := ""
	readJSON(r, tagValue)
	tags[tagName] = tagValue
	writeResponse(w, true, http.StatusOK)
}

func handlerUsers(w http.ResponseWriter, r *http.Request) {
	code := http.StatusOK
	if r.Method == http.MethodPost {
		code = http.StatusCreated
	} else if r.Method == http.MethodPut {
		code = http.StatusNoContent
	}
	writeResponse(w, "", code)
}

func handlerImages(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(testHTTPServer.URL)
	w.Header().Add("X-Docker-Endpoints", fmt.Sprintf("%s 	,  %s ", u.Host, "test.example.com"))
	w.Header().Add("X-Docker-Token", fmt.Sprintf("FAKE-SESSION-%d", time.Now().UnixNano()))
	if r.Method == http.MethodPut {
		if strings.HasSuffix(r.URL.Path, "images") {
			writeResponse(w, "", http.StatusNoContent)
			return
		}
		writeResponse(w, "", http.StatusOK)
		return
	}
	if r.Method == http.MethodDelete {
		writeResponse(w, "", http.StatusNoContent)
		return
	}
	var images []map[string]string
	for imageID, layer := range testLayers {
		image := make(map[string]string)
		image["id"] = imageID
		image["checksum"] = layer["checksum_tarsum"]
		image["Tag"] = "latest"
		images = append(images, image)
	}
	writeResponse(w, images, http.StatusOK)
}

func handlerAuth(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, "OK", http.StatusOK)
}

>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
func handlerSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeResponse(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	result := &registry.SearchResults{
		Query:      "fakequery",
		NumResults: 1,
		Results:    []registry.SearchResult{{Name: "fakeimage", StarCount: 42}},
	}
	writeResponse(w, result, http.StatusOK)
}

func TestPing(t *testing.T) {
	res, err := http.Get(makeURL("/v1/_ping"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK, "")
	assert.Equal(t, res.Header.Get("Server"), "docker-tests/mock")
}
