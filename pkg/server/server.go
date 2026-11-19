package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/domehahn/skdesk/pkg/ticket"
)

type Server struct {
	mgr *ticket.Manager
}

func NewServer(mgr *ticket.Manager) *Server {
	return &Server{mgr: mgr}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tickets", s.handleTickets)
	mux.HandleFunc("/", s.handleDashboard)
	return mux
}

func (s *Server) handleTickets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tickets := s.mgr.List()
	_ = json.NewEncoder(w).Encode(tickets)
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tickets := s.mgr.List()

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>skdesk — Human-in-the-Loop Review Dashboard</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 2rem; background: #0d1117; color: #c9d1d9; }
        h1 { color: #58a6ff; }
        table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
        th, td { border: 1px solid #30363d; padding: 0.75rem; text-align: left; }
        th { background: #161b22; color: #f0f6fc; }
        .PENDING { color: #d29922; font-weight: bold; }
        .APPROVED { color: #3fb950; font-weight: bold; }
        .REJECTED { color: #f85149; font-weight: bold; }
    </style>
</head>
<body>
    <h1>skdesk — Governance Review Portal</h1>
    <p>Active security tickets requiring human-in-the-loop sign-off.</p>
    <table>
        <tr>
            <th>Ticket ID</th>
            <th>Skill Name</th>
            <th>Digest</th>
            <th>Environment</th>
            <th>Status</th>
            <th>Reviewer</th>
        </tr>`

	for _, t := range tickets {
		html += fmt.Sprintf(`<tr>
            <td><code>%s</code></td>
            <td>%s</td>
            <td><code>%s</code></td>
            <td>%s</td>
            <td class="%s">%s</td>
            <td>%s</td>
        </tr>`, t.ID, t.SkillName, t.SkillDigest, t.Environment, t.Status, t.Status, t.Reviewer)
	}

	html += `</table></body></html>`
	fmt.Fprint(w, html)
}
