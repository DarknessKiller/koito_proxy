package koito_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gin-gonic/gin"

	"koito_proxy/internal/config"
	"koito_proxy/internal/model"
	"koito_proxy/internal/proxy/koito"
	"koito_proxy/internal/rules"
)

// MockRuleRepository
type MockRuleRepository struct {
	CreateFunc  func(ctx context.Context, rule *model.Rule) error
	GetByIDFunc func(ctx context.Context, id string) (*model.Rule, error)
	GetAllFunc  func(ctx context.Context) ([]model.Rule, error)
	UpdateFunc  func(ctx context.Context, id string, rule *model.Rule) error
	DeleteFunc  func(ctx context.Context, id string) error
}

func (m *MockRuleRepository) Create(ctx context.Context, rule *model.Rule) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, rule)
	}
	return nil
}

func (m *MockRuleRepository) GetByID(ctx context.Context, id string) (*model.Rule, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("not GetByIDFunc implemented")
}

func (m *MockRuleRepository) GetAll(ctx context.Context) ([]model.Rule, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, errors.New("not GetAllFunc implemented")
}

func (m *MockRuleRepository) Update(ctx context.Context, id string, rule *model.Rule) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, rule)
	}
	return errors.New("not UpdateFunc implemented")
}

func (m *MockRuleRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return errors.New("not DeleteFunc implemented")
}

var _ = Describe("Intercept.Merge", func() {

	var (
		upstreamHandler func(w http.ResponseWriter, r *http.Request)

		upstream *httptest.Server
		cfg      *config.Config
		mockRepo *MockRuleRepository
		engine   *rules.RuleEngine
		h        *koito.Handler
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)

		mockRepo = &MockRuleRepository{}
		engine = rules.NewRuleEngine()

		upstreamHandler = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}

		upstream = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upstreamHandler(w, r)
		}))

		cfg = &config.Config{UpstreamURL: upstream.URL}
		h = koito.NewHandler(engine, mockRepo, cfg)
	})

	AfterEach(func() {
		upstream.Close()
	})

	type UpstreamFn func(w http.ResponseWriter, r *http.Request)

	run := func(
		entity, id string,
		body any,
		upstreamFn UpstreamFn,
	) *httptest.ResponseRecorder {

		if upstreamFn != nil {
			upstreamHandler = upstreamFn
		}

		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)

		b, _ := json.Marshal(body)

		req := httptest.NewRequest(
			http.MethodPost,
			"/apis/web/v1/"+entity+"/"+id+"/merge",
			bytes.NewReader(b),
		)

		req.Header.Set("Content-Type", "application/json")

		ctx.Request = req
		ctx.Params = gin.Params{
			{Key: "entity", Value: entity},
			{Key: "id", Value: id},
		}

		h.InterceptMerge(ctx)
		return rec
	}

	It("handles track merge", func() {

		mockRepo.CreateFunc = func(ctx context.Context, rule *model.Rule) error {
			Expect(rule.MatchTrackName.String).To(Equal("Old Track"))
			Expect(rule.ReplaceTrackName.String).To(Equal("New Track"))
			return nil
		}

		rec := run("track", "456",
			map[string]any{"merge_from_id": 123},
			func(w http.ResponseWriter, r *http.Request) {

				if r.Method == http.MethodGet && r.URL.Path == "/apis/web/v1/track/123" {
					w.Write([]byte(`{"id":123,"title":"Old Track","artists":[{"id":1,"name":"Old Artist"}],"album_id":1}`))
					return
				}

				if r.Method == http.MethodGet && r.URL.Path == "/apis/web/v1/track/456" {
					w.Write([]byte(`{"id":456,"title":"New Track","artists":[{"id":2,"name":"New Artist"}],"album_id":2}`))
					return
				}

				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte(`{"status":"merged"}`))
					return
				}
			},
		)

		Expect(rec.Code).To(Equal(http.StatusAccepted))
	})

	It("handles artist merge", func() {

		mockRepo.CreateFunc = func(ctx context.Context, rule *model.Rule) error {
			Expect(rule.MatchArtistName.String).To(Equal("Old Artist Name"))
			Expect(rule.ReplaceArtistName.String).To(Equal("New Artist Name"))
			return nil
		}

		rec := run("artist", "999",
			map[string]any{"merge_from_id": 789},
			func(w http.ResponseWriter, r *http.Request) {

				if r.Method == http.MethodGet && r.URL.Path == "/apis/web/v1/artist/789" {
					w.Write([]byte(`{"id":789,"name":"Old Artist Name"}`))
					return
				}

				if r.Method == http.MethodGet && r.URL.Path == "/apis/web/v1/artist/999" {
					w.Write([]byte(`{"id":999,"name":"New Artist Name"}`))
					return
				}

				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte(`{"status":"merged"}`))
					return
				}
			},
		)

		Expect(rec.Code).To(Equal(http.StatusAccepted))
	})

	It("handles album merge", func() {

		mockRepo.CreateFunc = func(ctx context.Context, rule *model.Rule) error {
			Expect(rule.MatchReleaseName.String).To(Equal("Old Album"))
			Expect(rule.ReplaceReleaseName.String).To(Equal("New Album"))
			return nil
		}

		rec := run("album", "2",
			map[string]any{"merge_from_id": 1},
			func(w http.ResponseWriter, r *http.Request) {

				if r.Method == http.MethodGet && r.URL.Path == "/apis/web/v1/album/1" {
					w.Write([]byte(`{"id":1,"title":"Old Album"}`))
					return
				}

				if r.Method == http.MethodGet && r.URL.Path == "/apis/web/v1/album/2" {
					w.Write([]byte(`{"id":2,"title":"New Album"}`))
					return
				}

				if r.Method == http.MethodPost {
					body, _ := io.ReadAll(r.Body)
					var req map[string]any
					Expect(json.Unmarshal(body, &req)).To(Succeed())
					Expect(req).To(HaveKey("merge_from_id"))

					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte(`{"status":"merged"}`))
					return
				}
			},
		)

		Expect(rec.Code).To(Equal(http.StatusAccepted))
	})

	It("returns bad request for invalid JSON", func() {

		mockRepo.CreateFunc = func(ctx context.Context, rule *model.Rule) error {
			Fail("Create should not be called")
			return nil
		}

		rec := run("track", "456",
			"not-json",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		)

		Expect(rec.Code).To(Equal(http.StatusBadRequest))
	})
})
