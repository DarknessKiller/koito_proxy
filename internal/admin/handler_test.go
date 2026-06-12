package admin_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"

	"koito_proxy/internal/admin"
	"koito_proxy/internal/model"
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
	return nil, errors.New("mock GetByIDFunc not implemented")
}

func (m *MockRuleRepository) GetAll(ctx context.Context) ([]model.Rule, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, errors.New("mock GetAllFunc not implemented")
}

func (m *MockRuleRepository) Update(ctx context.Context, id string, rule *model.Rule) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, rule)
	}
	return errors.New("mock UpdateFunc not implemented")
}

func (m *MockRuleRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return errors.New("mock DeleteFunc not implemented")
}

var _ = Describe("Admin Handler", func() {

	var (
		mockRepo *MockRuleRepository
		engine   *rules.RuleEngine
		h        *admin.Handler
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockRepo = &MockRuleRepository{}
		engine = rules.NewRuleEngine()
		h = admin.NewHandler(mockRepo, engine)
	})

	run := func(method, path string, body any, params ...gin.Param) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)

		var req *http.Request
		if body != nil {
			b, _ := json.Marshal(body)
			req = httptest.NewRequest(method, path, bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req = httptest.NewRequest(method, path, nil)
		}
		ctx.Request = req

		if len(params) > 0 {
			ctx.Params = params
		}

		return rec
	}

	call := func(method, path string, body any, fn func(*gin.Context), params ...gin.Param) *httptest.ResponseRecorder {
		rec := run(method, path, body, params...)
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = httptest.NewRequest(method, path, nil)
		if body != nil {
			b, _ := json.Marshal(body)
			ctx.Request = httptest.NewRequest(method, path, bytes.NewReader(b))
			ctx.Request.Header.Set("Content-Type", "application/json")
		}
		if len(params) > 0 {
			ctx.Params = params
		}
		fn(ctx)
		if !ctx.Writer.Written() {
			ctx.Writer.WriteHeaderNow()
		}
		return rec
	}

	Describe("CheckAuth", func() {
		It("returns 200", func() {
			rec := call("GET", "/apis/admin/check", nil, h.CheckAuth)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var resp map[string]any
			Expect(json.Unmarshal(rec.Body.Bytes(), &resp)).To(Succeed())
			Expect(resp["ok"]).To(BeTrue())
		})
	})

	Describe("ListRules", func() {
		It("returns all rules", func() {
			mockRepo.GetAllFunc = func(ctx context.Context) ([]model.Rule, error) {
				return []model.Rule{
					{
						MatchTrackName: sql.NullString{String: "Track A", Valid: true},
						Enabled:        boolPtr(true),
					},
					{
						MatchTrackName: sql.NullString{String: "Track B", Valid: true},
						Enabled:        boolPtr(false),
					},
				}, nil
			}

			rec := call("GET", "/apis/admin/rules", nil, h.ListRules)

			Expect(rec.Code).To(Equal(http.StatusOK))
			var rules []map[string]any
			Expect(json.Unmarshal(rec.Body.Bytes(), &rules)).To(Succeed())
			Expect(rules).To(HaveLen(2))
			Expect(rules[0]["match_track_name"]).To(Equal("Track A"))
			Expect(rules[0]["enabled"]).To(BeTrue())
			Expect(rules[1]["enabled"]).To(BeFalse())
		})

		It("returns 500 when repository fails", func() {
			mockRepo.GetAllFunc = func(ctx context.Context) ([]model.Rule, error) {
				return nil, errors.New("db error")
			}

			rec := call("GET", "/apis/admin/rules", nil, h.ListRules)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Describe("CreateRule", func() {
		It("creates a rule and returns 201", func() {
			var createdID string
			mockRepo.CreateFunc = func(ctx context.Context, rule *model.Rule) error {
				createdID = rule.ID.String()
				Expect(*rule.Enabled).To(BeTrue())
				return nil
			}

			rec := call("POST", "/apis/admin/rules",
				map[string]any{
					"match_track_name":   "Old Track",
					"replace_track_name": "New Track",
					"enabled":            true,
				}, h.CreateRule)

			Expect(rec.Code).To(Equal(http.StatusCreated))
			var resp map[string]any
			Expect(json.Unmarshal(rec.Body.Bytes(), &resp)).To(Succeed())
			Expect(resp["id"]).To(Equal(createdID))
			Expect(resp["enabled"]).To(BeTrue())
		})

		It("creates a disabled rule when enabled=false", func() {
			mockRepo.CreateFunc = func(ctx context.Context, rule *model.Rule) error {
				Expect(*rule.Enabled).To(BeFalse())
				return nil
			}

			rec := call("POST", "/apis/admin/rules",
				map[string]any{
					"match_track_name":   "Old Track",
					"replace_track_name": "New Track",
					"enabled":            false,
				}, h.CreateRule)

			Expect(rec.Code).To(Equal(http.StatusCreated))
		})

		It("returns 400 for invalid JSON", func() {
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest("POST", "/apis/admin/rules",
				bytes.NewReader([]byte(`not json`)))
			ctx.Request.Header.Set("Content-Type", "application/json")
			h.CreateRule(ctx)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("GetRule", func() {
		It("returns a rule by ID", func() {
			id := ksuid.New().String()
			mockRepo.GetByIDFunc = func(ctx context.Context, rid string) (*model.Rule, error) {
				Expect(rid).To(Equal(id))
				return &model.Rule{
					MatchTrackName: sql.NullString{String: "Track", Valid: true},
					Enabled:        boolPtr(true),
				}, nil
			}

			rec := call("GET", "/apis/admin/rules/"+id, nil, h.GetRule,
				gin.Param{Key: "id", Value: id})

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("returns 404 when rule not found", func() {
			mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*model.Rule, error) {
				return nil, errors.New("not found")
			}

			rec := call("GET", "/apis/admin/rules/nonexistent", nil, h.GetRule,
				gin.Param{Key: "id", Value: "nonexistent"})

			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("UpdateRule", func() {
		It("updates a rule and returns 200", func() {
			id := ksuid.New().String()
			mockRepo.GetByIDFunc = func(ctx context.Context, rid string) (*model.Rule, error) {
				return &model.Rule{Enabled: boolPtr(true)}, nil
			}
			mockRepo.UpdateFunc = func(ctx context.Context, uid string, rule *model.Rule) error {
				Expect(uid).To(Equal(id))
				Expect(*rule.Enabled).To(BeFalse())
				return nil
			}

			rec := call("PUT", "/apis/admin/rules/"+id,
				map[string]any{
					"match_track_name": "Updated Track",
					"enabled":          false,
				}, h.UpdateRule, gin.Param{Key: "id", Value: id})

			Expect(rec.Code).To(Equal(http.StatusOK))
		})

		It("returns 404 when rule not found", func() {
			mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*model.Rule, error) {
				return nil, errors.New("not found")
			}

			rec := call("PUT", "/apis/admin/rules/nonexistent",
				map[string]any{"match_track_name": "X"},
				h.UpdateRule, gin.Param{Key: "id", Value: "nonexistent"})

			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})

		It("returns 400 for invalid JSON", func() {
			mockRepo.GetByIDFunc = func(ctx context.Context, id string) (*model.Rule, error) {
				return &model.Rule{}, nil
			}

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest("PUT", "/apis/admin/rules/123",
				bytes.NewReader([]byte(`not json`)))
			ctx.Request.Header.Set("Content-Type", "application/json")
			h.UpdateRule(ctx)

			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("DeleteRule", func() {
		It("deletes a rule and returns 204", func() {
			id := ksuid.New().String()
			mockRepo.DeleteFunc = func(ctx context.Context, did string) error {
				Expect(did).To(Equal(id))
				return nil
			}

			rec := call("DELETE", "/apis/admin/rules/"+id, nil, h.DeleteRule,
				gin.Param{Key: "id", Value: id})

			Expect(rec.Code).To(Equal(http.StatusNoContent))
		})

		It("returns 404 when rule not found", func() {
			mockRepo.DeleteFunc = func(ctx context.Context, id string) error {
				return errors.New("not found")
			}

			rec := call("DELETE", "/apis/admin/rules/nonexistent", nil, h.DeleteRule,
				gin.Param{Key: "id", Value: "nonexistent"})

			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})
	})
})

func boolPtr(v bool) *bool { return &v }
