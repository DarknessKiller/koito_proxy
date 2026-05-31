package koito_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"koito_proxy/internal/config"
	"koito_proxy/internal/model"
	"koito_proxy/internal/proxy/koito"
	"koito_proxy/internal/rules"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

func TestKoito(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Koito Suite")
}

var _ = Describe("InterceptMerge", func() {
	It("adds a merge rule and proxies the merge to upstream", func() {
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// handle GET source/target and POST merge
			switch r.Method {
			case http.MethodGet:
				if r.URL.Path == "/apis/web/v1/track/123" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"id":123,"title":"Old Track","artists":[{"id":1,"name":"Old Artist"}],"album_id":1}`))
					return
				}
				if r.URL.Path == "/apis/web/v1/track/456" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"id":456,"title":"New Track","artists":[{"id":2,"name":"New Artist"}],"album_id":2}`))
					return
				}
				// album endpoints used by fetchAlbum
				if r.URL.Path == "/apis/web/v1/album/1" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"id":1,"title":"Old Album"}`))
					return
				}
				if r.URL.Path == "/apis/web/v1/album/2" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"id":2,"title":"New Album"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
				return
			case http.MethodPost:
				if r.URL.Path == "/apis/web/v1/track/456/merge" {
					body, err := io.ReadAll(r.Body)
					Expect(err).NotTo(HaveOccurred())
					var req map[string]interface{}
					Expect(json.Unmarshal(body, &req)).To(Succeed())
					Expect(req).To(HaveKey("merge_from_id"))

					w.Header().Set("X-Test-Proxy", "true")
					w.WriteHeader(http.StatusAccepted)
					_, _ = w.Write([]byte(`{"status":"merged"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
		}))
		defer upstream.Close()

		cfg := &config.Config{UpstreamURL: upstream.URL}

		// prepare in-memory sqlite DB and store
		db, err := sql.Open("sqlite", ":memory:")
		Expect(err).NotTo(HaveOccurred())
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS overwrite_rules (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,

		    match_track_name TEXT NULL,
		    match_artist_name TEXT NULL,
		    match_release_name TEXT NULL,

		    replace_track_name TEXT NULL,
		    replace_artist_name TEXT NULL,
		    replace_release_name TEXT NULL,

		    enabled INTEGER NOT NULL DEFAULT 1,

		    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`)
		Expect(err).NotTo(HaveOccurred())

		store := rules.NewStore(db)
		engine := rules.NewRuleEngine([]model.Rule{})
		h := koito.NewHandler(engine, store, cfg)

		reqBody, err := json.Marshal(map[string]interface{}{"merge_from_id": 123})
		Expect(err).NotTo(HaveOccurred())

		req := httptest.NewRequest(http.MethodPost, "/apis/web/v1/track/456/merge", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		respRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(respRecorder)
		c.Request = req
		c.Params = gin.Params{{Key: "entity", Value: "track"}, {Key: "id", Value: "456"}}

		h.InterceptMerge(c)

		Expect(respRecorder.Code).To(Equal(http.StatusAccepted))
		Expect(respRecorder.Body.String()).To(Equal(`{"status":"merged"}`))
		Expect(respRecorder.Header().Get("X-Test-Proxy")).To(Equal("true"))

		// verify rule was inserted into DB
		row := db.QueryRow("SELECT replace_track_name, replace_artist_name FROM overwrite_rules WHERE match_track_name = ?", "Old Track")
		var replaceTrack, replaceArtist sql.NullString
		Expect(row.Scan(&replaceTrack, &replaceArtist)).To(Succeed())
		Expect(replaceTrack.Valid).To(BeTrue())
		Expect(replaceTrack.String).To(Equal("New Track"))
		Expect(replaceArtist.Valid).To(BeTrue())
		Expect(replaceArtist.String).To(Equal("New Artist"))
	})
})
