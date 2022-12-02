package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/KiloProjects/kilonova/internal/config"
	"github.com/KiloProjects/kilonova/internal/util"
	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

// Language stuff
var serverLangs = []language.Tag{
	language.English,
	language.Romanian,
}
var langMatcher = language.NewMatcher(serverLangs)

// ValidateProblemID makes sure the problem ID is a valid uint
func (rt *Web) ValidateProblemID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		problemID, err := strconv.Atoi(chi.URLParam(r, "pbid"))
		if err != nil {
			rt.statusPage(w, r, http.StatusBadRequest, "ID invalid")
			return
		}
		problem, err1 := rt.base.Problem(r.Context(), problemID)
		if err1 != nil {
			rt.statusPage(w, r, 404, "Problema nu a fost găsită")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.ProblemKey, problem)))
	})
}

// ValidateListID makes sure the list ID is a valid uint
func (rt *Web) ValidateListID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		listID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			rt.statusPage(w, r, http.StatusBadRequest, "ID invalid")
			return
		}
		list, err1 := rt.base.ProblemList(r.Context(), listID)
		if err1 != nil {
			rt.statusPage(w, r, 404, "Lista nu a fost găsită")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.ProblemListKey, list)))
	})
}

// ValidateAttachmentID makes sure the attachment ID is a valid uint
func (rt *Web) ValidateAttachmentID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attid, err := strconv.Atoi(chi.URLParam(r, "aid"))
		if err != nil {
			http.Error(w, "ID invalid", 400)
			return
		}
		att, err1 := rt.base.Attachment(r.Context(), attid)
		if err1 != nil {
			http.Error(w, "Atașamentul nu a fost găsit", 404)
			return
		}
		if att.Private && !util.IsRProblemEditor(r) {
			rt.statusPage(w, r, 403, "Nu ai voie să accesezi acest atașament")
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.AttachmentKey, att)))
	})
}

// ValidateProblemVisible checks if the problem from context is visible from the logged in user
func (rt *Web) ValidateProblemVisible(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !util.IsRProblemVisible(r) {
			rt.statusPage(w, r, 403, "Nu ai voie să accesezi problema!")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ValidateSubmissionID puts the ID and the Submission in the router context
func (rt *Web) ValidateSubmissionID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			rt.statusPage(w, r, 400, "ID submisie invalid")
			return
		}
		sub, err1 := rt.base.Submission(r.Context(), subID, util.UserBrief(r))
		if err1 != nil {
			rt.statusPage(w, r, 400, "Submisia nu există")
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.SubKey, &sub.Submission)))
	})
}

// ValidatePasteID puts the ID and the Paste in the router context
func (rt *Web) ValidatePasteID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !config.Features.Pastes {
			rt.statusPage(w, r, 404, "Feature has been disabled by administrator.")
			return
		}
		paste, err1 := rt.base.SubmissionPaste(r.Context(), chi.URLParam(r, "id"))
		if err1 != nil {
			rt.statusPage(w, r, 400, "Paste-ul nu există")
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.PasteKey, paste)))
	})
}

func (rt *Web) mustBeAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !util.IsRAuthed(r) {
			rt.statusPage(w, r, 401, "Trebuie să fii logat")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rt *Web) mustBeProposer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !util.IsRProposer(r) {
			rt.statusPage(w, r, 401, "Trebuie să fii propunător")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rt *Web) mustBeAdmin(next http.Handler) http.Handler {
	return rt.mustBeAuthed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !util.IsRAdmin(r) {
			rt.statusPage(w, r, 401, "Trebuie să fii admin")
			return
		}
		next.ServeHTTP(w, r)
	}))
}

func (rt *Web) mustBeVisitor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if util.IsRAuthed(r) {
			rt.statusPage(w, r, 401, "Trebuie să fii delogat")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rt *Web) mustBeProblemEditor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !util.IsRProblemEditor(r) {
			rt.statusPage(w, r, 401, "Trebuie să fii autorul problemei")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getSessCookie(r *http.Request) string {
	cookie, err := r.Cookie("kn-sessionid")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (rt *Web) initSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := rt.base.GetSession(r.Context(), getSessCookie(r))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := rt.base.UserFull(r.Context(), sess)
		if user == nil {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.UserKey, user)))
	})
}

func (rt *Web) initLanguage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLang := ""
		if util.UserFull(r) != nil {
			userLang = util.UserFull(r).PreferredLanguage
		}
		// get language
		lang, _ := r.Cookie("lang")
		accept := r.Header.Get("Accept-Language")
		tag, _ := language.MatchStrings(langMatcher, lang.String(), userLang, accept)
		language, _ := tag.Base()

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.LangKey, language.String())))
	})
}
