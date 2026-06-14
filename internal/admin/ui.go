package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UIHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(adminPageHTML))
}

const adminPageHTML = `<!DOCTYPE html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
<style>

#koito-admin-content {
  --bg:      var(--color-bg);
  --bg2:     var(--color-bg-secondary);
  --bg3:     var(--color-bg-tertiary);
  --fg:      var(--color-fg);
  --fg2:     var(--color-fg-secondary);
  --fg3:     var(--color-fg-tertiary);
  --primary: var(--color-primary);
  --accent:  var(--color-accent);
  --error:   var(--color-error);
  --warning: var(--color-warning);
  --info:    var(--color-info);
  --success: var(--color-success);
  --radius:  10px;

  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  font-family: 'Jost', sans-serif;
  background: var(--bg2);
  color: var(--fg);
  font-size: 14px;
  font-weight: 400;
  line-height: 1.5;
  border: 1px solid var(--bg3);
  border-radius: var(--radius);

  width: min(96vw, 1100px);
  height: min(92vh, 800px);
  overflow: hidden;
}

#koito-admin-content * {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

#koito-admin-content .kp-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--bg3);
  flex-shrink: 0;
}

#koito-admin-content .kp-header h1 {
  font-family: 'League Spartan', sans-serif;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: .02em;
  color: var(--fg);
  margin-right: auto;
}

#koito-admin-content .kp-header h1 small {
  font-size: 12px;
  color: var(--fg3);
  font-weight: 400;
  font-family: 'Jost', sans-serif;
  margin-left: 8px;
}

#koito-admin-content .kp-close {
  width: 40px;
  height: 40px;
  min-width: 40px;
  border-radius: 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--fg3);
  line-height: 1;
  padding: 0;
  transition: color .1s ease, background .1s ease;
  font-family: 'Jost', sans-serif;
}

#koito-admin-content .kp-close:hover {
  color: var(--fg);
  background: var(--bg3);
}

#koito-admin-content .kp-toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  border-bottom: 1px solid var(--bg3);
  flex-shrink: 0;
}

#koito-admin-content .kp-toolbar .rule-count {
  color: var(--fg3);
  font-size: 13px;
  white-space: nowrap;
  flex-shrink: 0;
}

#koito-admin-content .kp-search {
  position: relative;
  flex: 1;
  max-width: 280px;
}

#koito-admin-content .kp-search input {
  width: 100%;
  padding: 8px 12px;
  background: var(--bg);
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  font-size: 14px;
  font-family: 'Jost', sans-serif;
  color: var(--fg);
  transition: border-color .1s ease;
  min-height: 40px;
}

#koito-admin-content .kp-search input:focus {
  border-color: var(--fg3);
  outline: none;
}

#koito-admin-content .kp-search input::placeholder {
  color: var(--fg3);
  opacity: .5;
}

#koito-admin-content .btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 18px;
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  cursor: pointer;
  font-size: 14px;
  font-family: 'Jost', sans-serif;
  font-weight: 500;
  background: var(--bg);
  color: var(--fg);
  transition: background .1s ease, color .1s ease, border-color .1s ease;
  -webkit-user-select: none;
  user-select: none;
  white-space: nowrap;
  min-height: 40px;
}

#koito-admin-content .btn:hover {
  background: var(--bg3);
}

#koito-admin-content .btn-primary {
  background: var(--bg);
  color: var(--primary);
  border-color: var(--primary);
}

#koito-admin-content .btn-primary:hover {
  background: var(--primary);
  color: var(--bg);
  border-color: var(--primary);
}

#koito-admin-content .btn-danger {
  background: var(--bg);
  color: var(--error);
  border-color: var(--error);
}

#koito-admin-content .btn-danger:hover {
  background: var(--error);
  color: var(--bg);
  border-color: var(--error);
}

#koito-admin-content .btn-sm {
  padding: 6px 14px;
  font-size: 13px;
  min-height: 36px;
}

#koito-admin-content .kp-body {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overscroll-behavior: contain;
  -webkit-overflow-scrolling: touch;
}

#koito-admin-content .kp-body-inner {
  padding: 0;
  min-height: 0;
}

#koito-admin-content .table-wrap {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

#koito-admin-content table {
  width: 100%;
  table-layout: fixed;
  border-collapse: collapse;
}

#koito-admin-content th,
#koito-admin-content td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid var(--bg3);
  white-space: nowrap;
  font-size: 13px;
}

#koito-admin-content th {
  background: var(--bg);
  font-weight: 600;
  font-size: 11px;
  color: var(--fg2);
  text-transform: uppercase;
  letter-spacing: .06em;
  position: sticky;
  top: 0;
  z-index: 2;
}

#koito-admin-content tr:last-child td {
  border-bottom: none;
}

#koito-admin-content tbody tr {
  transition: background .08s ease;
}

#koito-admin-content tbody tr:hover td {
  background: rgba(229, 132, 106, .04);
}

#koito-admin-content .cell {
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  display: inline-block;
  vertical-align: middle;
}

#koito-admin-content .badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  border: 1px solid transparent;
}

#koito-admin-content .badge-enabled {
  background: rgba(122, 184, 122, .15);
  color: var(--success);
  border-color: rgba(122, 184, 122, .3);
}

#koito-admin-content .badge-disabled {
  background: rgba(212, 98, 96, .15);
  color: var(--error);
  border-color: rgba(212, 98, 96, .3);
}

#koito-admin-content .badge-invalid {
  background: rgba(232, 168, 64, .15);
  color: var(--warning);
  border-color: rgba(232, 168, 64, .3);
}

#koito-admin-content .actions {
  white-space: nowrap;
}

#koito-admin-content .actions .btn + .btn {
  margin-left: 6px;
}

#koito-admin-content .empty {
  padding: 60px 20px;
  text-align: center;
  color: var(--fg3);
  font-size: 15px;
}

#koito-admin-content .mono {
  font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', monospace;
  font-size: 13px;
  color: var(--fg2);
  text-align: center;
}

#koito-admin-content .kp-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 20px;
  border-top: 1px solid var(--bg3);
  flex-shrink: 0;
  flex-wrap: wrap;
}

#koito-admin-content .kp-footer .kp-page-info {
  font-size: 12px;
  color: var(--fg3);
  white-space: nowrap;
}

#koito-admin-content .kp-footer .kp-page-btn {
  padding: 6px 12px;
  min-height: 36px;
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  background: var(--bg);
  color: var(--fg2);
  cursor: pointer;
  font-size: 12px;
  font-family: 'Jost', sans-serif;
  transition: background .1s ease, color .1s ease;
}

#koito-admin-content .kp-footer .kp-page-btn:hover:not(:disabled) {
  background: var(--bg3);
  color: var(--fg);
}

#koito-admin-content .kp-footer .kp-page-btn:disabled {
  opacity: .35;
  cursor: default;
}

#koito-admin-content .kp-footer .kp-page-btn.active {
  background: var(--primary);
  color: var(--bg);
  border-color: var(--primary);
}

#koito-admin-content .kp-footer select {
  padding: 6px 8px;
  min-height: 36px;
  background: var(--bg);
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  color: var(--fg2);
  font-size: 12px;
  font-family: 'Jost', sans-serif;
  cursor: pointer;
}

#koito-admin-content .kp-footer select:focus {
  outline: none;
  border-color: var(--fg3);
}

#koito-admin-content .kp-cards {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

#koito-admin-content .kp-card {
  background: var(--bg);
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  padding: 12px 14px;
}

#koito-admin-content .kp-card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}

#koito-admin-content .kp-card-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--fg);
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
}

#koito-admin-content .kp-card-meta {
  font-size: 12px;
  color: var(--fg3);
  margin-bottom: 8px;
  line-height: 1.5;
}

#koito-admin-content .kp-card-meta span {
  display: block;
}

#koito-admin-content .kp-card-meta .kp-label {
  color: var(--fg3);
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: .04em;
}

#koito-admin-content .kp-card-actions {
  display: flex;
  gap: 8px;
  align-items: center;
  padding-top: 10px;
  border-top: 1px solid var(--bg3);
}

#koito-admin-content .kp-card-actions .btn {
  flex: 1;
}

/* ---- Modal overlay ---- */
#koito-admin-content .kp-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, .85);
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: kpFadeIn .12s ease;
}

@keyframes kpFadeIn {
  from { opacity: 0 }
  to   { opacity: 1 }
}

#koito-admin-content .kp-modal {
  background: var(--bg2);
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  padding: 24px;
  width: min(94vw, 620px);
  max-height: min(90vh, 700px);
  overflow-y: auto;
  animation: kpScaleIn .12s ease;
}

@keyframes kpScaleIn {
  from { opacity: 0; transform: scale(.96) }
  to   { opacity: 1; transform: scale(1) }
}

#koito-admin-content .kp-modal h2 {
  font-family: 'League Spartan', sans-serif;
  font-size: 17px;
  font-weight: 700;
  margin-bottom: 18px;
  color: var(--fg);
  letter-spacing: .02em;
}

#koito-admin-content .form-row {
  display: flex;
  gap: 10px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

#koito-admin-content .form-group {
  flex: 1;
  min-width: 180px;
}

#koito-admin-content .form-group label {
  display: block;
  font-size: 11px;
  font-weight: 600;
  color: var(--fg3);
  margin-bottom: 4px;
  text-transform: uppercase;
  letter-spacing: .06em;
}

#koito-admin-content .form-group input {
  width: 100%;
  padding: 10px 12px;
  min-height: 40px;
  background: var(--bg);
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  font-size: 14px;
  font-family: 'Jost', sans-serif;
  color: var(--fg);
  transition: border-color .1s ease;
}

#koito-admin-content .form-group input:focus {
  border-color: var(--fg3);
  outline: none;
}

#koito-admin-content .form-group input::placeholder {
  color: var(--fg3);
  opacity: .5;
}

#koito-admin-content .form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 18px;
  padding-top: 16px;
  border-top: 1px solid var(--bg3);
}

#koito-admin-content .form-actions .btn {
  min-width: 90px;
}

#koito-admin-content .form-section {
  font-size: 11px;
  font-weight: 700;
  color: var(--fg3);
  margin: 16px 0 10px;
  text-transform: uppercase;
  letter-spacing: .08em;
}

#koito-admin-content .toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
  flex-shrink: 0;
}

#koito-admin-content .toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

#koito-admin-content .toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0; left: 0; right: 0; bottom: 0;
  background: var(--bg3);
  border-radius: 24px;
  transition: background .15s ease;
}

#koito-admin-content .toggle-slider::before {
  content: '';
  position: absolute;
  height: 20px;
  width: 20px;
  left: 2px;
  bottom: 2px;
  background: var(--fg);
  border-radius: 50%;
  transition: transform .15s ease;
}

#koito-admin-content .toggle-switch input:checked + .toggle-slider {
  background: var(--primary);
}

#koito-admin-content .toggle-switch input:checked + .toggle-slider::before {
  transform: translateX(20px);
}

#koito-admin-content .toggle-switch input:focus-visible + .toggle-slider {
  outline: 2px solid var(--fg3);
  outline-offset: 2px;
}

/* ---- Toast ---- */
#koito-admin-content .kp-toast {
  position: fixed;
  bottom: 24px;
  right: 24px;
  padding: 12px 22px;
  border-radius: var(--radius);
  color: var(--bg);
  font-size: 14px;
  font-weight: 500;
  z-index: 20000;
  animation: kpToastIn .15s ease;
  pointer-events: none;
}

@keyframes kpToastIn {
  from { opacity: 0; transform: translateY(10px) }
  to   { opacity: 1; transform: translateY(0) }
}

#koito-admin-content .kp-toast-success { background: var(--success) }
#koito-admin-content .kp-toast-error   { background: var(--error) }
#koito-admin-content .kp-toast-info    { background: var(--info) }

/* Mobile <= 768px */
@media (max-width: 768px) {
  #koito-admin-content {
    width: 100vw;
    height: 100vh;
    border-radius: 0;
    border: none;
  }

  #koito-admin-content .kp-header {
    padding: 10px 14px;
    gap: 8px;
  }

  #koito-admin-content .kp-header h1 {
    font-size: 16px;
  }

  #koito-admin-content .kp-header h1 small {
    display: none;
  }

  #koito-admin-content .kp-close {
    width: 44px;
    height: 44px;
    min-width: 44px;
    font-size: 28px;
  }

  #koito-admin-content .kp-toolbar {
    flex-wrap: wrap;
    padding: 10px 14px;
    gap: 8px;
  }

  #koito-admin-content .kp-toolbar .rule-count {
    font-size: 12px;
    order: 1;
  }

  #koito-admin-content .kp-search {
    max-width: none;
    order: 2;
  }

  #koito-admin-content .kp-search input {
    min-height: 44px;
    font-size: 16px;
  }

  #koito-admin-content .btn {
    min-height: 44px;
    padding: 10px 18px;
    font-size: 14px;
  }

  #koito-admin-content .btn-sm {
    min-height: 44px;
    padding: 10px 18px;
    font-size: 14px;
  }

  #koito-admin-content .kp-body {
    flex: 1;
    min-height: 0;
  }

  #koito-admin-content table {
    display: none;
  }

  #koito-admin-content .kp-cards {
    display: flex;
  }

  #koito-admin-content .kp-card-actions .btn {
    min-height: 44px;
  }

  #koito-admin-content .kp-modal {
    width: 100vw;
    max-width: 100vw;
    max-height: 100vh;
    height: 100vh;
    border-radius: 0;
    padding: 16px;
    border: none;
  }

  #koito-admin-content .kp-modal h2 {
    font-size: 16px;
    margin-bottom: 14px;
  }

  #koito-admin-content .form-row {
    flex-direction: column;
    gap: 10px;
  }

  #koito-admin-content .form-group {
    min-width: 0;
  }

  #koito-admin-content .form-group input {
    min-height: 44px;
    font-size: 16px;
  }

  #koito-admin-content .form-actions .btn {
    flex: 1;
    min-height: 44px;
  }

  #koito-admin-content .kp-footer {
    padding: 8px 14px;
    gap: 6px;
  }

  #koito-admin-content .kp-footer .kp-page-btn {
    min-height: 40px;
    padding: 8px 12px;
  }

  #koito-admin-content .kp-footer select {
    min-height: 40px;
  }

  #koito-admin-content .mono {
    font-size: 12px;
  }

  #koito-admin-content .badge {
    font-size: 11px;
    padding: 3px 8px;
  }

  #koito-admin-content .toast {
    bottom: 16px;
    right: 16px;
    left: 16px;
    text-align: center;
  }
}

/* Desktop > 768px */
@media (min-width: 769px) {
  #koito-admin-content .kp-cards {
    display: none;
  }
}

</style>

<div id="koito-admin-content">
  <div class="kp-header">
    <h1>Koito Proxy</h1>
    <button class="kp-close" id="kpCloseBtn">&times;</button>
  </div>

  <div class="kp-toolbar">
    <span class="rule-count" id="kpRuleCount"></span>
    <div class="kp-search">
      <input id="kpSearchInput" type="text" placeholder="Search rules..." oninput="kpOnSearch()">
    </div>
    <button class="btn btn-primary" onclick="kpOpenCreate()">+ New Rule</button>
  </div>

  <div class="kp-body">
    <div id="kpBodyInner" class="kp-body-inner">
      <div class="table-wrap">
        <table>
          <colgroup>
            <col style="width:12%"><col style="width:10%"><col style="width:10%">
            <col style="width:12%"><col style="width:10%"><col style="width:10%">
            <col style="width:12%"><col style="width:4%"><col style="width:6%">
            <col style="width:14%">
          </colgroup>
          <thead>
            <tr>
              <th colspan="3" style="text-align:center;border-right:1px solid var(--bg3)">Match</th>
              <th colspan="3" style="text-align:center;border-right:1px solid var(--bg3)">Replace</th>
              <th style="text-align:center" rowspan="2">MBID</th>
              <th style="text-align:center" rowspan="2">Pri</th>
              <th style="text-align:center" rowspan="2">Status</th>
              <th style="text-align:center" rowspan="2">Actions</th>
            </tr>
            <tr>
              <th style="text-align:center">Track</th>
              <th style="text-align:center">Artist</th>
              <th style="text-align:center;border-right:1px solid var(--bg3)">Release</th>
              <th style="text-align:center">Track</th>
              <th style="text-align:center">Artist</th>
              <th style="text-align:center;border-right:1px solid var(--bg3)">Release</th>
            </tr>
          </thead>
          <tbody id="kpRulesBody"></tbody>
        </table>
      </div>
      <div id="kpCardsContainer" class="kp-cards"></div>
      <div id="kpEmpty" class="empty" style="display:none">No rules defined yet.</div>
    </div>
  </div>

  <div id="kpFooter" class="kp-footer" style="display:none"></div>

  <div id="kpModal" class="kp-modal-overlay" style="display:none" onclick="if(event.target===this)kpCloseModal()">
  <div class="kp-modal">
    <h2 id="kpModalTitle">Rule</h2>
    <form id="kpRuleForm" onsubmit="return kpSaveRule(event)">
      <div class="form-section">Match Criteria</div>
      <div class="form-row">
        <div class="form-group">
          <label>Track Name</label>
          <input name="match_track_name" placeholder="e.g. Bohemian Rhapsody">
        </div>
        <div class="form-group">
          <label>Artist Name</label>
          <input name="match_artist_name" placeholder="e.g. Queen">
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label>Release Name</label>
          <input name="match_release_name" placeholder="e.g. A Night at the Opera">
        </div>
        <div class="form-group">
          <label>Artist Names (JSON)</label>
          <input name="match_artist_names" placeholder='["Queen"]'>
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label>MBID</label>
          <input name="match_mbid" placeholder="e.g. 0603a798-1842-3b30-8ba9-d1243298a8a2">
        </div>
        <div class="form-group">
          <label>Duration Bucket</label>
          <input name="match_duration_bucket" type="number" placeholder="e.g. 107 (537s/5)">
        </div>
      </div>

      <div class="form-section">Replacement Values</div>
      <div class="form-row">
        <div class="form-group">
          <label>Replace Track Name</label>
          <input name="replace_track_name" placeholder="Corrected title">
        </div>
        <div class="form-group">
          <label>Replace Artist Name</label>
          <input name="replace_artist_name" placeholder="Corrected artist">
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label>Replace Release Name</label>
          <input name="replace_release_name" placeholder="Corrected release">
        </div>
        <div class="form-group">
          <label>Replace Artist Names (JSON)</label>
          <input name="replace_artist_names" placeholder='["Corrected Artist"]'>
        </div>
      </div>
      <div class="form-row" style="align-items:center">
        <label style="display:flex;align-items:center;gap:10px;cursor:pointer;font-size:14px;color:var(--fg2);min-height:44px">
          <span class="toggle-switch">
            <input name="enabled" type="checkbox" checked>
            <span class="toggle-slider"></span>
          </span>
          Enabled
        </label>
      </div>

      <input type="hidden" name="id" value="">
      <div class="form-actions">
        <button type="button" class="btn" onclick="kpCloseModal()">Cancel</button>
        <button type="submit" class="btn btn-primary">Save</button>
      </div>
    </form>
  </div>
  </div>

  <div id="kpDeleteModal" class="kp-modal-overlay" style="display:none" onclick="if(event.target===this)kpCloseDeleteModal()">
    <div class="kp-modal" style="max-width:420px">
      <h2>Delete Rule</h2>
      <p style="margin-bottom:16px;color:var(--fg2);font-size:14px;line-height:1.5">
        Are you sure you want to delete this rule? This cannot be undone.
      </p>
      <p id="kpDeleteInfo" style="margin-bottom:18px;font-size:13px;color:var(--fg3)"></p>
      <div class="form-actions" style="border-top:none;padding-top:0;margin-top:0">
        <button class="btn" onclick="kpCloseDeleteModal()">Cancel</button>
        <button class="btn btn-danger" onclick="kpConfirmDelete()">Delete</button>
      </div>
    </div>
  </div>
</div>

<script>
var kpRules = [];
var kpDeleteId = null;
var kpEditingId = null;
var kpSearchQuery = '';
var kpPage = 1;
var kpPageSize = 10;
var kpLastMobile = null;


document.getElementById('kpCloseBtn').onclick = kpCloseAdmin;

function kpIsMobile() {
  return window.matchMedia('(max-width: 768px)').matches;
}

window.matchMedia('(max-width: 768px)').addEventListener('change', function (e) {
  kpLastMobile = e.matches;
  kpRender();
});


function kpToast(msg, type) {
  var el = document.getElementById('kpToastContainer');
  if (!el) {
    el = document.createElement('div');
    el.id = 'kpToastContainer';
    document.getElementById('koito-admin-content').appendChild(el);
  }
  var t = document.createElement('div');
  t.className = 'kp-toast kp-toast-' + type;
  t.textContent = msg;
  el.appendChild(t);
  setTimeout(function () {
    t.style.opacity = '0';
    t.style.transition = 'opacity .15s';
    setTimeout(function () { t.remove() }, 150);
  }, 3000);
}

function kpRender() {
  var tbody = document.getElementById('kpRulesBody');
  var cards = document.getElementById('kpCardsContainer');
  var empty = document.getElementById('kpEmpty');
  var footer = document.getElementById('kpFooter');
  tbody.innerHTML = '';
  cards.innerHTML = '';

  var query = kpSearchQuery.trim().toLowerCase();
  var filtered = query
    ? kpRules.filter(function (r) {
        var fields = [
          r.match_track_name, r.match_artist_name, r.match_release_name,
          r.match_mbid,
          r.replace_track_name, r.replace_artist_name, r.replace_release_name,
          String(r.priority),
          r.enabled ? 'on' : 'off',
        ];
        return fields.some(function (f) { return f && f.toLowerCase().indexOf(query) !== -1; });
      })
    : kpRules;

  var totalPages = Math.max(1, Math.ceil(filtered.length / kpPageSize));
  if (kpPage > totalPages) kpPage = totalPages;

  var start = (kpPage - 1) * kpPageSize;
  var pageItems = filtered.slice(start, start + kpPageSize);

  if (kpRules.length === 0) {
    empty.style.display = 'block';
    footer.style.display = 'none';
    document.getElementById('kpRuleCount').textContent = '';
    return;
  }

  empty.style.display = 'none';

  var countEl = document.getElementById('kpRuleCount');
  if (query) {
    countEl.textContent = filtered.length + '/' + kpRules.length + ' rule' + (kpRules.length > 1 ? 's' : '');
  } else {
    countEl.textContent = kpRules.length + ' rule' + (kpRules.length > 1 ? 's' : '');
  }

  if (kpIsMobile()) {
    pageItems.forEach(function (r) {
      var card = document.createElement('div');
      card.className = 'kp-card';

      var titleText = r.match_track_name || r.match_artist_name || r.replace_track_name || 'Untitled';
      var title = document.createElement('div');
      title.className = 'kp-card-title';
      title.textContent = titleText;

      var header = document.createElement('div');
      header.className = 'kp-card-header';
      header.appendChild(title);
      var statusBadge = r.enabled
        ? '<span class="badge badge-enabled">on</span>'
        : '<span class="badge badge-disabled">off</span>';
      if (!r.valid) statusBadge += ' <span class="badge badge-invalid">low</span>';
      var statusEl = document.createElement('div');
      statusEl.innerHTML = statusBadge;
      header.appendChild(statusEl);
      card.appendChild(header);

      var meta = document.createElement('div');
      meta.className = 'kp-card-meta';
      var parts = [];
      if (r.match_artist_name) parts.push('<span><span class="kp-label">Match Artist</span> ' + kpEsc(r.match_artist_name) + '</span>');
      if (r.match_release_name) parts.push('<span><span class="kp-label">Match Release</span> ' + kpEsc(r.match_release_name) + '</span>');
      if (r.replace_track_name) parts.push('<span><span class="kp-label">Replace Track</span> ' + kpEsc(r.replace_track_name) + '</span>');
      if (r.replace_artist_name) parts.push('<span><span class="kp-label">Replace Artist</span> ' + kpEsc(r.replace_artist_name) + '</span>');
      if (r.replace_release_name) parts.push('<span><span class="kp-label">Replace Release</span> ' + kpEsc(r.replace_release_name) + '</span>');
      if (r.match_mbid) parts.push('<span><span class="kp-label">MBID</span> ' + kpEsc(r.match_mbid) + '</span>');
      parts.push('<span><span class="kp-label">Pri</span> ' + r.priority + '</span>');
      meta.innerHTML = parts.join('');
      card.appendChild(meta);

      var acts = document.createElement('div');
      acts.className = 'kp-card-actions';
      acts.innerHTML =
        '<button class="btn btn-sm" onclick="kpEditRule(\'' + r.id + '\')">Edit</button>' +
        '<button class="btn btn-sm btn-danger" onclick="kpOpenDelete(\'' + r.id + '\')">Delete</button>';
      card.appendChild(acts);

      cards.appendChild(card);
    });
  } else {
    try {
      pageItems.forEach(function (r) {
        var tr = document.createElement('tr');
        tr.innerHTML = kpCell(r.match_track_name, ' style="text-align:center"') +
          kpCell(r.match_artist_name, ' style="text-align:center"') +
          kpCell(r.match_release_name, ' style="text-align:center"') +
          kpCell(r.replace_track_name, ' style="text-align:center"') +
          kpCell(r.replace_artist_name, ' style="text-align:center"') +
          kpCell(r.replace_release_name, ' style="text-align:center"') +
          kpCell(r.match_mbid, ' style="text-align:center"') +
          '<td class="mono">' + r.priority + '</td>' +
          '<td style="text-align:center">' +
            (r.enabled
              ? '<span class="badge badge-enabled">on</span>'
              : '<span class="badge badge-disabled">off</span>') +
            (r.valid ? '' : ' <span class="badge badge-invalid">low</span>') +
          '</td>' +
          '<td class="actions">' +
            '<button class="btn btn-sm" onclick="kpEditRule(\'' + r.id + '\')">Edit</button>' +
            '<button class="btn btn-sm btn-danger" onclick="kpOpenDelete(\'' + r.id + '\')">Del</button>' +
          '</td>';
        tbody.appendChild(tr);
      });
    } catch (e) {
      kpToast('Render error: ' + e.message, 'error');
    }
  }

  kpRenderPagination(filtered.length, totalPages);
}

function kpCell(v, s) {
  s = s || '';
  return v
    ? '<td' + s + '><span class="cell" title="' + kpEsc(v) + '">' + kpEsc(v) + '</span></td>'
    : '<td' + s + '><span class="cell" style="color:var(--fg3)">-</span></td>';
}

function kpRenderPagination(total, totalPages) {
  var el = document.getElementById('kpFooter');
  if (totalPages <= 1 && !kpSearchQuery.trim()) {
    el.style.display = 'none';
    return;
  }
  el.style.display = 'flex';

  var html = '';
  html += '<span class="kp-page-info">' + ((kpPage - 1) * kpPageSize + 1) + '-' + Math.min(kpPage * kpPageSize, total) + ' of ' + total + '</span>';

  if (!kpIsMobile()) {
    html += '<select onchange="kpPageSize=parseInt(this.value);kpPage=1;kpRender()">';
    [10, 15, 25, 50].forEach(function (n) {
      html += '<option value="' + n + '"' + (n === kpPageSize ? ' selected' : '') + '>' + n + '</option>';
    });
    html += '</select>';
  }

  html += '<button class="kp-page-btn" onclick="kpGoToPage(kpPage-1)"' + (kpPage <= 1 ? ' disabled' : '') + '>&laquo; Prev</button>';

  var maxVisible = kpIsMobile() ? 2 : 5;
  var pageStart = Math.max(1, kpPage - Math.floor(maxVisible / 2));
  var pageEnd = Math.min(totalPages, pageStart + maxVisible - 1);
  if (pageEnd - pageStart + 1 < maxVisible) {
    pageStart = Math.max(1, pageEnd - maxVisible + 1);
  }

  if (pageStart > 1) {
    html += '<button class="kp-page-btn" onclick="kpGoToPage(1)">1</button>';
    if (pageStart > 2) html += '<span class="kp-page-info">...</span>';
  }
  for (var i = pageStart; i <= pageEnd; i++) {
    html += '<button class="kp-page-btn' + (i === kpPage ? ' active' : '') + '" onclick="kpGoToPage(' + i + ')">' + i + '</button>';
  }
  if (pageEnd < totalPages) {
    if (pageEnd < totalPages - 1) html += '<span class="kp-page-info">...</span>';
    html += '<button class="kp-page-btn" onclick="kpGoToPage(' + totalPages + ')">' + totalPages + '</button>';
  }

  html += '<button class="kp-page-btn" onclick="kpGoToPage(kpPage+1)"' + (kpPage >= totalPages ? ' disabled' : '') + '>Next &raquo;</button>';

  el.innerHTML = html;
}

function kpGoToPage(p) {
  kpPage = p;
  kpRender();
}

function kpOnSearch() {
  kpSearchQuery = document.getElementById('kpSearchInput').value;
  kpPage = 1;
  kpRender();
}

function kpEsc(s) {
  var d = document.createElement('div');
  d.appendChild(document.createTextNode(s));
  return d.innerHTML;
}

function kpLoadRules() {
  fetch('/apis/admin/rules')
    .then(function (r) {
      if (r.status === 401) { throw new Error('unauthorized'); }
      return r.json();
    })
    .then(function (data) {
      kpRules = data;
      kpRender();
    })
    .catch(function (err) {
      if (err.message === 'unauthorized') {
        document.getElementById('kpRulesBody').innerHTML =
          '<tr><td colspan="10" style="text-align:center;padding:40px;color:var(--error)">' +
          'Not authenticated. Please log in to Koito first.</td></tr>';
      } else {
        kpToast('Failed to load rules: ' + err.message, 'error');
      }
    });
}

function kpOpenCreate() {
  kpEditingId = null;
  document.getElementById('kpModalTitle').textContent = 'New Rule';
  document.getElementById('kpRuleForm').reset();
  document.getElementById('kpRuleForm').querySelector('[name=enabled]').checked = true;
  document.getElementById('kpModal').style.display = 'flex';
}

function kpEditRule(id) {
  kpEditingId = id;
  var r = kpRules.find(function (x) { return x.id === id; });
  if (!r) return;

  document.getElementById('kpModalTitle').textContent = 'Edit Rule';
  document.getElementById('kpRuleForm').reset();
  kpSetField('match_track_name', r.match_track_name || '');
  kpSetField('match_artist_name', r.match_artist_name || '');
  kpSetField('match_release_name', r.match_release_name || '');
  kpSetField('match_artist_names', r.match_artist_names ? JSON.stringify(r.match_artist_names) : '');
  kpSetField('match_duration_bucket', r.match_duration_bucket || '');
  kpSetField('match_mbid', r.match_mbid || '');
  kpSetField('replace_track_name', r.replace_track_name || '');
  kpSetField('replace_artist_name', r.replace_artist_name || '');
  kpSetField('replace_release_name', r.replace_release_name || '');
  kpSetField('replace_artist_names', r.replace_artist_names ? JSON.stringify(r.replace_artist_names) : '');
  document.getElementById('kpRuleForm').querySelector('[name=enabled]').checked = r.enabled;
  document.getElementById('kpModal').style.display = 'flex';
}

function kpSetField(name, val) {
  var el = document.getElementById('kpRuleForm').querySelector('[name="' + name + '"]');
  if (el) el.value = val;
}

function kpGetField(name) {
  var el = document.getElementById('kpRuleForm').querySelector('[name="' + name + '"]');
  if (!el) return null;
  if (el.type === 'checkbox') return el.checked;
  return el.value || null;
}

function kpCloseModal() {
  document.getElementById('kpModal').style.display = 'none';
}

function kpCloseAdmin() {
  var c = document.getElementById('koito-admin-content');
  if (c) c.style.display = 'none';
  var o = document.getElementById('kp-admin-overlay');
  if (o) o.style.display = 'none';
}

function kpSaveRule(e) {
  e.preventDefault();

  var raw = {
    match_track_name: kpGetField('match_track_name'),
    match_artist_name: kpGetField('match_artist_name'),
    match_release_name: kpGetField('match_release_name'),
    match_mbid: kpGetField('match_mbid'),
    replace_track_name: kpGetField('replace_track_name'),
    replace_artist_name: kpGetField('replace_artist_name'),
    replace_release_name: kpGetField('replace_release_name'),
    enabled: kpGetField('enabled'),
  };

  var dur = kpGetField('match_duration_bucket');
  if (dur) raw.match_duration_bucket = parseInt(dur, 10);

  try {
    var an = kpGetField('match_artist_names');
    if (an) raw.match_artist_names = JSON.parse(an);
    var rn = kpGetField('replace_artist_names');
    if (rn) raw.replace_artist_names = JSON.parse(rn);
  } catch (e) {
    kpToast('Invalid JSON in artist names field', 'error');
    return;
  }

  var url = '/apis/admin/rules';
  var method = 'POST';
  if (kpEditingId) {
    url += '/' + kpEditingId;
    method = 'PUT';
  }

  fetch(url, {
    method: method,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(raw),
  })
    .then(function (r) {
      if (r.status === 401) throw new Error('unauthorized');
      if (!r.ok) throw new Error('HTTP ' + r.status);
      kpToast(kpEditingId ? 'Rule updated' : 'Rule created', 'success');
      kpCloseModal();
      kpLoadRules();
    })
    .catch(function (err) {
      if (err.message === 'unauthorized') {
        kpToast('Not authenticated', 'error');
      } else {
        kpToast('Failed to save rule: ' + err.message, 'error');
      }
    });
}

function kpOpenDelete(id) {
  kpDeleteId = id;
  var r = kpRules.find(function (x) { return x.id === id; });
  var info = r ? kpEsc(r.match_track_name || r.match_artist_name || r.id) : id;
  document.getElementById('kpDeleteInfo').textContent = 'Rule: ' + info;
  document.getElementById('kpDeleteModal').style.display = 'flex';
}

function kpCloseDeleteModal() {
  kpDeleteId = null;
  document.getElementById('kpDeleteModal').style.display = 'none';
}

function kpConfirmDelete() {
  if (!kpDeleteId) return;

  fetch('/apis/admin/rules/' + kpDeleteId, { method: 'DELETE' })
    .then(function (r) {
      if (r.status === 401) throw new Error('unauthorized');
      if (!r.ok && r.status !== 204) throw new Error('HTTP ' + r.status);
      kpToast('Rule deleted', 'success');
      kpCloseDeleteModal();
      kpLoadRules();
    })
    .catch(function (err) {
      if (err.message === 'unauthorized') {
        kpToast('Not authenticated', 'error');
      } else {
        kpToast('Failed to delete rule: ' + err.message, 'error');
      }
    });
}

kpLoadRules();
</script>`
