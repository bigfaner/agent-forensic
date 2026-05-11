// Shared: tab switching, active state management

function initTabs(containerSelector) {
  const container = document.querySelector(containerSelector);
  if (!container) return;
  const tabs = container.querySelectorAll('.state-tab');
  const panels = container.querySelectorAll('.state-panel');
  tabs.forEach(tab => {
    tab.addEventListener('click', () => {
      tabs.forEach(t => t.classList.remove('active'));
      panels.forEach(p => p.classList.remove('active'));
      tab.classList.add('active');
      const target = tab.dataset.target;
      const panel = container.querySelector('#' + target);
      if (panel) panel.classList.add('active');
    });
  });
}

document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('.tab-group').forEach(group => {
    initTabs('#' + group.id);
  });
});
