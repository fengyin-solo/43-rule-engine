async function load() {
  try {
    const res = await fetch('/api/stats/overview', { headers: { 'X-Api-Key': 'default-ruleengine-key' } });
    const body = await res.json();
    const stats = body.data || {};
    const container = document.getElementById('stats');
    container.innerHTML = '';
    const items = [
      { label: '规则集总数', key: 'rule_set_count' },
      { label: '已发布数', key: 'published_count' },
      { label: '总评估次数', key: 'total_evaluation' },
      { label: '平均命中率', key: 'avg_hit_ratio' },
      { label: '平均耗时(ms)', key: 'avg_duration_ms' }
    ];
    items.forEach(item => {
      const card = document.createElement('div');
      card.className = 'stat-card';
      card.innerHTML = `<div class="label">${item.label}</div><div class="value">${stats[item.key] ?? '-'}</div>`;
      container.appendChild(card);
    });
  } catch (e) {
    console.error('加载统计失败', e);
  }

  try {
    const res = await fetch('/api/rule-sets?api_key=default-ruleengine-key');
    const body = await res.json();
    const tbody = document.querySelector('#list tbody');
    tbody.innerHTML = '';
    (body.data.items || []).forEach(rs => {
      const tr = document.createElement('tr');
      tr.innerHTML = `<td>${rs.id}</td><td>${rs.name}</td><td>${rs.status}</td><td>${rs.version}</td><td>${rs.created_at}</td>`;
      tbody.appendChild(tr);
    });
  } catch (e) {
    console.error('加载规则集失败', e);
  }
}

document.getElementById('eval-form').addEventListener('submit', async function (e) {
  e.preventDefault();
  const ruleSetID = document.getElementById('eval-rule-set-id').value.trim();
  const inputRaw = document.getElementById('eval-input').value.trim();
  let input = {};
  try {
    input = JSON.parse(inputRaw);
  } catch (err) {
    document.getElementById('eval-result').textContent = 'JSON 解析失败: ' + err.message;
    return;
  }
  try {
    const res = await fetch('/api/evaluate?api_key=default-ruleengine-key', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'X-Api-Key': 'default-ruleengine-key' },
      body: JSON.stringify({ rule_set_id: ruleSetID, input })
    });
    const body = await res.json();
    document.getElementById('eval-result').textContent = JSON.stringify(body, null, 2);
  } catch (e) {
    document.getElementById('eval-result').textContent = '请求失败: ' + e.message;
  }
});

load();
