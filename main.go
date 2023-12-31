package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type RecordingData struct {
	Title       string
	Description string
	Filename    string
	Time        string
}

func GetData() map[string]RecordingData {
	data := map[string]RecordingData{}
	data["01_LifeInReview"] = RecordingData{
		Title:       "Jean Walker: Life in Review",
		Description: "Played at her memorial service, Jean gives some of the highlights of her life.",
		Filename:    "Jean_Walker_Life_in_Review.mp3",
		Time:        "4:04",
	}

	data["02_RachelIntro"] = RecordingData{
		Title:       "Rachel's Introduction",
		Description: "Humorously serious, young Rachel gives the introduction for the recordings.",
		Filename:    "Rachels_Introduction.mp3",
		Time:        "2:08",
	}

	data["03_StoryOfTheYellowDress"] = RecordingData{
		Title:       "Story of the Yellow Dresses",
		Description: "Jean wore four yellow dresses to four events: two weddings and two milestone anniversaries. Her great grandchildren wore yellow dresses and vests at her memorial service in her honor.",
		Filename:    "Story_of_the_Yellow_Dresses.mp3",
		Time:        "3:39",
	}

	data["04_StoriesAboutEdAndAlice"] = RecordingData{
		Title:       "Stories about Ed and Alice",
		Description: "Mary and Jean share about Ed and Alice Montag (their parents), the move to Sac City and life in the early days at 606 Oak Street.",
		Filename:    "Stories_about_Alice_and_Ed.mp3",
		Time:        "19:05",
	}

	data["05_MaryMontagStory"] = RecordingData{
		Title:       "Mary Montag's Story",
		Description: "Mary shares about her education, early jobs and travel.",
		Filename:    "Mary_Montags_Story.mp3",
		Time:        "24:34",
	}

	data["06_MaryMontagEpicTravels"] = RecordingData{
		Title:       "Mary Montag's Epic Travels",
		Description: "The longest of the recordings, Mary shares about her around-the-world trip and other shenanigans. It's worth the listen - if you have 58mins!",
		Filename:    "Mary_Montags_Epic_Travels.mp3",
		Time:        "58:17",
	}

	data["07_VisitingTomMontagsGrave"] = RecordingData{
		Title:       "Visiting Tom Montag's Grave",
		Description: "Capt. Tom Montag was a pilot in WWII. His plane was shot down over northern France. This recording talks about Tom’s siblings and nieces and nephews visiting his grave in France.",
		Filename:    "Visiting_Tom_Montags_Grave.mp3",
		Time:        "5:51",
	}

	data["08_AliceMontagAndCars"] = RecordingData{
		Title:       "Alice Montag and Cars",
		Description: "Alice Montag was not known for her respect of automobiles. Mary shares a few tidbits about Grandma Alice.",
		Filename:    "Alice_Montag_and_Cars.mp3",
		Time:        "5:17",
	}

	data["09_JeansStory"] = RecordingData{
		Title:       "Jean's Story, original recording",
		Description: "The unedited version of Jean’s history. You can hear her humor and charm in this longer version - along with a bit more background noise and tangents.",
		Filename:    "Jeans_Story.mp3",
		Time:        "30:26",
	}
	return data
}

//go:embed templates/*
var resources embed.FS
var t = template.Must(template.ParseFS(resources, "templates/*"))

func addDownloadHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate the server headers
		filename := strings.TrimPrefix(r.URL.Path, "/download/")
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename+"")
		w.Header().Set("Expires", "0")
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")
		fs.ServeHTTP(w, r)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Download file serving
	fsdl := WithLogging(http.FileServer(http.Dir("static/media/")))
	http.Handle("/download/", http.StripPrefix("/download/", addDownloadHeaders(fsdl)))

	// Static file serving
	fs := WithLogging(http.FileServer(http.Dir("static/")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.Handle("/brownies.html", WithLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		t.ExecuteTemplate(w, "brownies.html", nil)
	})))

	http.Handle("/", WithLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		t.ExecuteTemplate(w, "index.html", GetData())
	})))

	// Start Server
	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
