package rules_test

import (
	"context"
	"database/sql"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "modernc.org/sqlite"

	"koito_proxy/internal/model"
	"koito_proxy/internal/rules"
)

func TestRules(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rules Suite")
}

var _ = Describe("Store and Engine", func() {
	var db *sql.DB
	var store *rules.Store

	BeforeEach(func() {
		var err error
		db, err = sql.Open("sqlite", "file::memory:?mode=memory&cache=shared")
		Expect(err).NotTo(HaveOccurred())

		_, err = db.Exec(`
			CREATE TABLE overwrite_rules (
			    id INTEGER PRIMARY KEY AUTOINCREMENT,
			    match_track_name TEXT NULL,
			    match_artist_name TEXT NULL,
			    match_release_name TEXT NULL,
			    replace_track_name TEXT NULL,
			    replace_artist_name TEXT NULL,
			    replace_release_name TEXT NULL,
			    enabled INTEGER NOT NULL DEFAULT 1
			);
		`)
		Expect(err).NotTo(HaveOccurred())

		store = rules.NewStore(db)
	})

	AfterEach(func() {
		Expect(db.Close()).To(Succeed())
	})

	Describe("Store", func() {
		It("loads enabled overwrite rules", func() {
			_, err := db.Exec(`INSERT INTO overwrite_rules (match_track_name, match_artist_name, replace_track_name, enabled) VALUES (?, ?, ?, 1)`, "Track A", "Artist A", "Track B")
			Expect(err).NotTo(HaveOccurred())

			rulesList, err := store.Load(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(rulesList).To(HaveLen(1))
			Expect(rulesList[0].MatchTrackName.String).To(Equal("Track A"))
			Expect(rulesList[0].MatchArtistName.String).To(Equal("Artist A"))
			Expect(rulesList[0].ReplaceTrackName.String).To(Equal("Track B"))
		})

		It("adds new overwrite rules", func() {
			rule := model.Rule{
				MatchTrackName:   sql.NullString{String: "Old Track", Valid: true},
				ReplaceTrackName: sql.NullString{String: "New Track", Valid: true},
			}
			Expect(store.Add(context.Background(), rule)).To(Succeed())

			rulesList, err := store.Load(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(rulesList).To(HaveLen(1))
			Expect(rulesList[0].MatchTrackName.String).To(Equal("Old Track"))
			Expect(rulesList[0].ReplaceTrackName.String).To(Equal("New Track"))
		})

		It("silently ignores duplicate rules on add", func() {
			rule := model.Rule{
				MatchArtistName:   sql.NullString{String: "Artist X", Valid: true},
				ReplaceArtistName: sql.NullString{String: "Artist Y", Valid: true},
			}

			Expect(store.Add(context.Background(), rule)).To(Succeed())
			// second add should not error but also should not create a duplicate
			Expect(store.Add(context.Background(), rule)).To(Succeed())

			rulesList, err := store.Load(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(rulesList).To(HaveLen(1))
		})
	})

	Describe("Engine", func() {
		It("applies matching overwrite rules to track metadata", func() {
			rule := model.Rule{
				MatchTrackName:    sql.NullString{String: "Old Track", Valid: true},
				MatchArtistName:   sql.NullString{String: "Old Artist", Valid: true},
				ReplaceTrackName:  sql.NullString{String: "New Track", Valid: true},
				ReplaceArtistName: sql.NullString{String: "New Artist", Valid: true},
			}

			engine := rules.NewRuleEngine([]model.Rule{rule})
			metadata := &model.ListenBrainzTrackMetaData{
				TrackName:  "Old Track",
				ArtistName: "Old Artist",
			}

			engine.Apply(metadata)

			Expect(metadata.TrackName).To(Equal("New Track"))
			Expect(metadata.ArtistName).To(Equal("New Artist"))
		})

		It("applies artist-only overwrite rules", func() {
			rule := model.Rule{
				MatchArtistName:   sql.NullString{String: "Solo Artist", Valid: true},
				ReplaceArtistName: sql.NullString{String: "Renamed Artist", Valid: true},
			}

			engine := rules.NewRuleEngine([]model.Rule{rule})
			metadata := &model.ListenBrainzTrackMetaData{
				TrackName:  "Some Track",
				ArtistName: "Solo Artist",
			}

			engine.Apply(metadata)

			Expect(metadata.ArtistName).To(Equal("Renamed Artist"))
		})

		It("applies release-only overwrite rules", func() {
			rule := model.Rule{
				MatchReleaseName:   sql.NullString{String: "Old Album", Valid: true},
				ReplaceReleaseName: sql.NullString{String: "New Album", Valid: true},
			}

			engine := rules.NewRuleEngine([]model.Rule{rule})
			metadata := &model.ListenBrainzTrackMetaData{
				TrackName:   "Some Track",
				ArtistName:  "Some Artist",
				ReleaseName: "Old Album",
			}

			engine.Apply(metadata)

			Expect(metadata.ReleaseName).To(Equal("New Album"))
		})
	})
})
