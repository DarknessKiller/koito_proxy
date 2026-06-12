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
<meta name="viewport" content="width=device-width, initial-scale=1">
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
  --radius:  8px;

  display: block;
  box-sizing: border-box;
  font-family: 'Jost', sans-serif;
  background: var(--bg2);
  color: var(--fg);
  padding: 24px;
  font-size: 14px;
  font-weight: 400;
  line-height: 1.4;
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  width: 100%;
  min-width: 0;
  max-height: 85vh;
  overflow-y: auto;
  overscroll-behavior: contain;
}

#koito-admin-content * {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

#koito-admin-content h1 {
  font-family: 'League Spartan', sans-serif;
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 14px;
  display: flex;
  align-items: center;
  gap: 8px;
  letter-spacing: .02em;
  color: var(--fg);
}

#koito-admin-content h1 small {
  font-size: 12px;
  color: var(--fg3);
  font-weight: 400;
  font-family: 'Jost', sans-serif;
}

#koito-admin-content .toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

#koito-admin-content .rule-count {
  color: var(--fg3);
  font-size: 13px;
}

#koito-admin-content .btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 20px;
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  cursor: pointer;
  font-size: 13px;
  font-family: 'Jost', sans-serif;
  font-weight: 500;
  background: var(--bg);
  color: var(--fg);
  transition: background .1s ease, color .1s ease, border-color .1s ease;
  -webkit-user-select: none;
  user-select: none;
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
  padding: 4px 12px;
  font-size: 12px;
}

#koito-admin-content .kp-close {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--fg3);
  line-height: 1;
  padding: 0;
  transition: color .1s ease;
  font-family: 'Jost', sans-serif;
}

#koito-admin-content .kp-close:hover {
  color: var(--fg);
}

#koito-admin-content .table-wrap {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
}

#koito-admin-content table {
  width: 100%;
  border-collapse: collapse;
}

#koito-admin-content th,
#koito-admin-content td {
  padding: 8px 10px;
  text-align: left;
  border-bottom: 1px solid var(--bg3);
  white-space: nowrap;
  font-size: 12px;
}

#koito-admin-content th {
  background: var(--bg);
  font-weight: 600;
  font-size: 10px;
  color: var(--fg2);
  text-transform: uppercase;
  letter-spacing: .06em;
}

#koito-admin-content tr:last-child td {
  border-bottom: none;
}

#koito-admin-content tr:hover td {
  background: rgba(229, 132, 106, .04);
}

#koito-admin-content .cell {
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: inline-block;
  vertical-align: middle;
}

#koito-admin-content .badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 20px;
  font-size: 11px;
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

#koito-admin-content .badge-valid {
  background: rgba(122, 173, 212, .15);
  color: var(--info);
  border-color: rgba(122, 173, 212, .3);
}

#koito-admin-content .badge-invalid {
  background: rgba(232, 168, 64, .15);
  color: var(--warning);
  border-color: rgba(232, 168, 64, .3);
}

#koito-admin-content .actions {
  text-align: center;
  white-space: nowrap;
}

#koito-admin-content .empty {
  padding: 40px 20px;
  text-align: center;
  color: var(--fg3);
  font-size: 14px;
}

#koito-admin-content .mono {
  font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', monospace;
  font-size: 12px;
  color: var(--fg2);
  text-align: center;
}

#koito-admin-content .kp-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, .9);
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: kpFadeIn .1s ease;
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
  width: 680px;
  max-width: 95vw;
  max-height: 85vh;
  overflow-y: auto;
  animation: kpScaleIn .1s ease;
}

@keyframes kpScaleIn {
  from { opacity: 0; transform: scale(.95) }
  to   { opacity: 1; transform: scale(1) }
}

#koito-admin-content .kp-modal h2 {
  font-family: 'League Spartan', sans-serif;
  font-size: 16px;
  font-weight: 700;
  margin-bottom: 16px;
  color: var(--fg);
  letter-spacing: .02em;
}

#koito-admin-content .form-row {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

#koito-admin-content .form-group {
  flex: 1;
  min-width: 180px;
}

#koito-admin-content .form-group label {
  display: block;
  font-size: 10px;
  font-weight: 600;
  color: var(--fg3);
  margin-bottom: 3px;
  text-transform: uppercase;
  letter-spacing: .06em;
}

#koito-admin-content .form-group input {
  width: 100%;
  padding: 8px 10px;
  background: var(--bg);
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  font-size: 13px;
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
  gap: 8px;
  justify-content: flex-end;
  margin-top: 16px;
  padding-top: 14px;
  border-top: 1px solid var(--bg3);
}

#koito-admin-content .form-section {
  font-size: 10px;
  font-weight: 700;
  color: var(--fg3);
  margin: 14px 0 8px;
  text-transform: uppercase;
  letter-spacing: .08em;
}

#koito-admin-content .toggle-switch {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 22px;
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
  border-radius: 22px;
  transition: background .15s ease;
}

#koito-admin-content .toggle-slider::before {
  content: '';
  position: absolute;
  height: 18px;
  width: 18px;
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
  transform: translateX(18px);
}

#koito-admin-content .toggle-switch input:focus-visible + .toggle-slider {
  outline: 1px solid var(--fg3);
  outline-offset: 2px;
}

#koito-admin-content .kp-search {
  display: flex;
  align-items: center;
  gap: 8px;
}

#koito-admin-content .kp-search input {
  padding: 6px 10px;
  background: var(--bg);
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  font-size: 13px;
  font-family: 'Jost', sans-serif;
  color: var(--fg);
  width: 200px;
  transition: border-color .1s ease;
}

#koito-admin-content .kp-search input:focus {
  border-color: var(--fg3);
  outline: none;
}

#koito-admin-content .kp-search input::placeholder {
  color: var(--fg3);
  opacity: .5;
}

#koito-admin-content .kp-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px 0 4px;
  font-size: 12px;
  color: var(--fg2);
}

#koito-admin-content .kp-pagination .kp-page-btn {
  padding: 4px 10px;
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  background: var(--bg);
  color: var(--fg2);
  cursor: pointer;
  font-size: 12px;
  font-family: 'Jost', sans-serif;
  transition: background .1s ease, color .1s ease;
}

#koito-admin-content .kp-pagination .kp-page-btn:hover:not(:disabled) {
  background: var(--bg3);
  color: var(--fg);
}

#koito-admin-content .kp-pagination .kp-page-btn:disabled {
  opacity: .3;
  cursor: default;
}

#koito-admin-content .kp-pagination .kp-page-btn.active {
  background: var(--primary);
  color: var(--bg);
  border-color: var(--primary);
}

#koito-admin-content .kp-pagination .kp-page-info {
  margin: 0 8px;
  color: var(--fg3);
}

#koito-admin-content .kp-pagination select {
  padding: 4px 6px;
  background: var(--bg);
  border: 1px solid var(--bg3);
  border-radius: var(--radius);
  color: var(--fg2);
  font-size: 12px;
  font-family: 'Jost', sans-serif;
  cursor: pointer;
}

#koito-admin-content .kp-pagination select:focus {
  outline: none;
  border-color: var(--fg3);
}

#koito-admin-content .kp-toast {
  position: fixed;
  bottom: 24px;
  right: 24px;
  padding: 10px 20px;
  border-radius: var(--radius);
  color: var(--bg);
  font-size: 13px;
  font-weight: 500;
  z-index: 20000;
  animation: kpToastIn .15s ease;
  pointer-events: none;
}

@keyframes kpToastIn {
  from { opacity: 0; transform: translateY(8px) }
  to   { opacity: 1; transform: translateY(0) }
}

#koito-admin-content .kp-toast-success { background: var(--success) }
#koito-admin-content .kp-toast-error   { background: var(--error) }
#koito-admin-content .kp-toast-info    { background: var(--info) }

@media (max-width: 900px) {
  #koito-admin-content { padding: 12px; }
  #koito-admin-content .toolbar { flex-direction: column; gap: 10px; align-items: stretch; }
  #koito-admin-content .kp-search { flex-wrap: wrap; }
  #koito-admin-content .kp-search input { width: 100%; }
  #koito-admin-content .kp-search .rule-count { width: 100%; }
  #koito-admin-content .btn { padding: 12px 20px; font-size: 15px; justify-content: center; min-height: 44px; }
  #koito-admin-content .btn-sm { padding: 10px 16px; min-height: 40px; }
  #koito-admin-content h1 { font-size: 17px; }
  #koito-admin-content td, #koito-admin-content th { padding: 6px 5px; font-size: 11px; }
  #koito-admin-content .cell { max-width: 80px; }
  #koito-admin-content .form-row { flex-direction: column; gap: 8px; }
  #koito-admin-content .form-group { min-width: 0; }
  #koito-admin-content .kp-modal { width: 100vw; max-width: 100vw; max-height: 100vh; border-radius: 0; padding: 16px; }
  #koito-admin-content .kp-pagination { flex-wrap: wrap; justify-content: center; }
  #koito-admin-content .kp-pagination .kp-page-btn { padding: 6px 12px; min-height: 36px; }
  #koito-admin-content .kp-pagination select { min-height: 36px; }
  #koito-admin-content .mono { font-size: 11px; }
  #koito-admin-content .badge { font-size: 10px; padding: 2px 6px; }
  #koito-admin-content .actions {
    position: sticky;
    right: 0;
    background: var(--bg2);
    z-index: 2;
  }
  #koito-admin-content tr:hover td.actions {
    background: rgba(229, 132, 106, .04);
  }
  #koito-admin-content .table-wrap {
    background:
      linear-gradient(to right, transparent calc(100% - 20px), rgba(0,0,0,.15) 100%)
      0 0 / 100% 100% no-repeat;
  }
}

@media (max-width: 480px) {
  #koito-admin-content { padding: 8px; }
  #koito-admin-content table { table-layout: fixed; }
  #koito-admin-content thead tr:first-child th:nth-child(1) { width: 30% !important; }
  #koito-admin-content thead tr:first-child th:nth-child(2) { width: 30% !important; }
  #koito-admin-content thead tr:first-child th:nth-child(3) { width: 10% !important; }
  #koito-admin-content thead tr:first-child th:nth-child(4) { width: 6% !important; }
  #koito-admin-content thead tr:first-child th:nth-child(5) { width: 10% !important; }
  #koito-admin-content thead tr:first-child th:nth-child(6) { width: 14% !important; }
  #koito-admin-content td, #koito-admin-content th { padding: 3px 2px; font-size: 9px; }
  #koito-admin-content .cell { display: block; width: 100%; overflow: hidden; text-overflow: ellipsis; max-width: none; }
  #koito-admin-content .btn { padding: 8px 10px; font-size: 12px; min-height: 36px; }
  #koito-admin-content .btn-sm { padding: 6px 8px; font-size: 10px; min-height: 32px; }
  #koito-admin-content .actions .btn-sm { padding: 4px 6px; font-size: 9px; min-height: 28px; margin-right: 3px !important; }
  #koito-admin-content .kp-search input { font-size: 14px; }
  #koito-admin-content th { font-size: 8px; letter-spacing: .02em; }
  #koito-admin-content .mono { font-size: 9px; }
  #koito-admin-content .badge { font-size: 8px; padding: 1px 4px; }
  #koito-admin-content .kp-pagination { font-size: 10px; }
  #koito-admin-content .kp-pagination .kp-page-btn { padding: 4px 8px; min-height: 28px; font-size: 10px; }
  #koito-admin-content .kp-pagination select { min-height: 28px; font-size: 10px; }
}
</style>

<div id="koito-admin-content">
  <button class="kp-close" id="kpCloseBtn">&times;</button>
  <h1>Koito Proxy <small>Rule Admin</small></h1>

  <div class="toolbar">
    <div class="kp-search">
      <span class="rule-count" id="kpRuleCount"></span>
      <input id="kpSearchInput" type="text" placeholder="Search rules..." oninput="kpOnSearch()">
    </div>
    <button class="btn btn-primary" onclick="kpOpenCreate()">+ New Rule</button>
  </div>

  <div class="table-wrap">
    <table>
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

  <div id="kpPagination" class="kp-pagination" style="display:none"></div>

  <div id="kpEmpty" class="empty" style="display:none">No rules defined yet.</div>

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
        <label style="display:flex;align-items:center;gap:10px;cursor:pointer;font-size:13px;color:var(--fg2)">
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
    <div class="kp-modal" style="max-width:400px">
      <h2>Delete Rule</h2>
      <p style="margin-bottom:14px;color:var(--fg2);font-size:14px">
        Are you sure you want to delete this rule? This cannot be undone.
      </p>
      <p id="kpDeleteInfo" style="margin-bottom:16px;font-size:12px;color:var(--fg3)"></p>
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

document.getElementById('kpCloseBtn').onclick = kpCloseAdmin;

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
  var empty = document.getElementById('kpEmpty');
  var pagination = document.getElementById('kpPagination');
  tbody.innerHTML = '';

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
    pagination.style.display = 'none';
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

  function cell(v, s) {
    s = s || '';
    return v
      ? '<td' + s + '><span class="cell" title="' + kpEsc(v) + '">' + kpEsc(v) + '</span></td>'
      : '<td' + s + '><span class="cell" style="color:var(--fg3)">-</span></td>';
  }

  pageItems.forEach(function (r) {
    var tr = document.createElement('tr');
    tr.innerHTML =
      cell(r.match_track_name, ' style="text-align:center"') +
      cell(r.match_artist_name, ' style="text-align:center"') +
      cell(r.match_release_name, ' style="text-align:center"') +
      cell(r.replace_track_name, ' style="text-align:center"') +
      cell(r.replace_artist_name, ' style="text-align:center"') +
      cell(r.replace_release_name, ' style="text-align:center"') +
      cell(r.match_mbid, ' style="text-align:center"') +
      '<td class="mono">' + r.priority + '</td>' +
      '<td style="text-align:center">' +
        (r.enabled
          ? '<span class="badge badge-enabled">on</span>'
          : '<span class="badge badge-disabled">off</span>') +
        (r.valid ? '' : ' <span class="badge badge-invalid">low</span>') +
      '</td>' +
      '<td class="actions">' +
        '<button class="btn btn-sm" style="margin-right:6px" onclick="kpEditRule(\'' + r.id + '\')">Edit</button>' +
        '<button class="btn btn-sm btn-danger" onclick="kpOpenDelete(\'' + r.id + '\')">Del</button>' +
      '</td>';
    tbody.appendChild(tr);
  });

  kpRenderPagination(filtered.length, totalPages);
}

function kpRenderPagination(total, totalPages) {
  var el = document.getElementById('kpPagination');
  if (totalPages <= 1 && !kpSearchQuery.trim()) {
    el.style.display = 'none';
    return;
  }
  el.style.display = 'flex';

  var html = '';
  html += '<span class="kp-page-info">' + ((kpPage - 1) * kpPageSize + 1) + '-' + Math.min(kpPage * kpPageSize, total) + ' of ' + total + '</span>';

  html += '<select onchange="kpPageSize=parseInt(this.value);kpPage=1;kpRender()">';
  [10, 15, 25, 50].forEach(function (n) {
    html += '<option value="' + n + '"' + (n === kpPageSize ? ' selected' : '') + '>' + n + '</option>';
  });
  html += '</select>';

  html += '<button class="kp-page-btn" onclick="kpGoToPage(kpPage-1)"' + (kpPage <= 1 ? ' disabled' : '') + '>&laquo; Prev</button>';

  var maxVisible = 5;
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
