package listenbrainz_test

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
	"koito_proxy/internal/proxy/listenbrainz"
	"koito_proxy/internal/rules"

	"github.com/gin-gonic/gin"
)

func TestListenBrainz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ListenBrainz Suite")
}

var _ = Describe("InterceptSubmitListen", func() {
	It("applies rule substitutions before proxying to upstream", func() {
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.Method).To(Equal(http.MethodPost))
			Expect(r.URL.Path).To(Equal("/apis/listenbrainz/1/submit-listens"))

			body, err := io.ReadAll(r.Body)
			Expect(err).NotTo(HaveOccurred())

			var req model.ListenBrainzSubmitRequest
			Expect(json.Unmarshal(body, &req)).To(Succeed())
			Expect(req.Payload).To(HaveLen(1))
			Expect(req.Payload[0].TrackMetaData.TrackName).To(Equal("New Track"))
			Expect(req.Payload[0].TrackMetaData.ArtistName).To(Equal("New Artist"))

			w.Header().Set("X-Test-Proxy", "true")
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"status":"proxied"}`))
		}))
		defer upstream.Close()

		cfg := &config.Config{UpstreamURL: upstream.URL}

		rule := rules.Rule{
			MatchTrackName:    sql.NullString{String: "Old Track", Valid: true},
			MatchArtistName:   sql.NullString{String: "Old Artist", Valid: true},
			ReplaceTrackName:  sql.NullString{String: "New Track", Valid: true},
			ReplaceArtistName: sql.NullString{String: "New Artist", Valid: true},
		}
		engine := rules.NewRuleEngine([]rules.Rule{rule})
		h := listenbrainz.NewHandler(engine, cfg)

		reqBody, err := json.Marshal(model.ListenBrainzSubmitRequest{
			ListenType: "single",
			Payload: []model.ListenBrainzSubmitPayload{{
				ListenedAt: 1234567890,
				TrackMetaData: model.ListenBrainzTrackMetaData{
					TrackName:  "Old Track",
					ArtistName: "Old Artist",
				},
			}},
		})
		Expect(err).NotTo(HaveOccurred())

		req := httptest.NewRequest(http.MethodPost, "/apis/listenbrainz/1/submit-listens", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		respRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(respRecorder)
		c.Request = req

		h.InterceptSubmitListen(c)

		Expect(respRecorder.Code).To(Equal(http.StatusAccepted))
		Expect(respRecorder.Body.String()).To(Equal(`{"status":"proxied"}`))
		Expect(respRecorder.Header().Get("X-Test-Proxy")).To(Equal("true"))
	})

	It("returns bad request for invalid JSON", func() {
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer upstream.Close()

		cfg := &config.Config{UpstreamURL: upstream.URL}
		h := listenbrainz.NewHandler(nil, cfg)

		req := httptest.NewRequest(http.MethodPost, "/apis/listenbrainz/1/submit-listens", bytes.NewReader([]byte("not json")))
		req.Header.Set("Content-Type", "application/json")

		respRecorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(respRecorder)
		c.Request = req

		h.InterceptSubmitListen(c)

		Expect(respRecorder.Code).To(Equal(http.StatusBadRequest))
		Expect(respRecorder.Body.String()).To(ContainSubstring("error"))
	})
})
