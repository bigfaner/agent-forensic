// Shared interactions for deep-drill-analytics prototype

function showState(stateId) {
  document.querySelectorAll('.state-content').forEach(el => el.classList.remove('active'));
  document.querySelectorAll('.state-tabs .btn').forEach(el => el.classList.remove('active'));
  const target = document.getElementById(stateId);
  if (target) target.classList.add('active');
  event.target.classList.add('active');
}

function toggleOverlay(id) {
  const backdrop = document.getElementById(id + '-backdrop');
  const panel = document.getElementById(id + '-panel');
  const isActive = panel.classList.contains('active');
  if (isActive) {
    backdrop.classList.remove('active');
    panel.classList.remove('active');
  } else {
    backdrop.classList.add('active');
    panel.classList.add('active');
  }
}

function switchTab(tabGroup, tabId) {
  const group = document.querySelector(`[data-tab-group="${tabGroup}"]`);
  if (!group) return;
  group.querySelectorAll('[data-tab]').forEach(el => el.style.display = 'none');
  group.querySelectorAll('.tab-btn').forEach(el => el.classList.remove('active'));
  const target = group.querySelector(`[data-tab="${tabId}"]`);
  if (target) target.style.display = 'block';
  event.target.classList.add('active');
}
