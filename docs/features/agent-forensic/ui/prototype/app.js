// Shared behaviors for agent-forensic TUI prototype

// Tab key simulation — cycle panel focus
document.addEventListener('keydown', function(e) {
  if (e.key === 'Tab') {
    e.preventDefault();
    const panels = document.querySelectorAll('.panel');
    const focused = document.querySelector('.panel.focused');
    if (!focused) { panels[0]?.classList.add('focused'); return; }
    focused.classList.remove('focused');
    const idx = Array.from(panels).indexOf(focused);
    panels[(idx + 1) % panels.length]?.classList.add('focused');
  }
});

// Toggle modal
function openModal(id) {
  document.getElementById(id).style.display = 'flex';
  document.querySelector('.status-bar').innerHTML =
    '<span class="key">j/k</span>:select  <span class="key">Enter</span>:jump  <span class="key">Esc</span>:close';
}
function closeModal(id) {
  document.getElementById(id).style.display = 'none';
  restoreMainStatusBar();
}
function restoreMainStatusBar() {
  const sb = document.querySelector('.status-bar');
  if (!sb) return;
  sb.innerHTML = '<span class="key">1</span>:sess <span class="key">2</span>:call <span class="key">j/k</span>:nav <span class="key">Enter</span>:expand <span class="key">Tab</span>:detail <span class="key">/</span>:search <span class="key">n/p</span>:replay <span class="key">d</span>:diag <span class="key">s</span>:stats <span class="key">m</span>:mon <span class="monitor-on">监听:开</span> <span class="key">q</span>:quit';
}

// Escape to close modals
document.addEventListener('keydown', function(e) {
  if (e.key === 'Escape') {
    const modals = document.querySelectorAll('.modal-overlay');
    modals.forEach(m => { if (m.style.display === 'flex') m.style.display = 'none'; });
    const pickers = document.querySelectorAll('.session-picker-overlay');
    pickers.forEach(p => p.style.display = 'none');
    restoreMainStatusBar();
  }
});

// Session row selection
function selectSession(el) {
  document.querySelectorAll('.session-row').forEach(r => r.classList.remove('selected'));
  el.classList.add('selected');
}

// Tree node selection + expand
function toggleNode(el) {
  document.querySelectorAll('.tree-node').forEach(n => n.classList.remove('selected'));
  el.classList.toggle('selected');
  const children = el.nextElementSibling;
  if (children && children.classList.contains('tree-children')) {
    children.style.display = children.style.display === 'none' ? 'block' : 'none';
  }
}
function selectNode(el) {
  document.querySelectorAll('.tree-node').forEach(n => n.classList.remove('selected'));
  el.classList.add('selected');
}

// Toggle monitoring
let monitorOn = true;
function toggleMonitor() {
  monitorOn = !monitorOn;
  const el = document.getElementById('monitor-status');
  if (!el) return;
  el.textContent = monitorOn ? '监听:开' : '监听:关';
  el.className = monitorOn ? 'monitor-on' : 'monitor-off';
}
