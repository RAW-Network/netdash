document.addEventListener('DOMContentLoaded', function() {
    const ctx = document.getElementById('speedChart').getContext('2d');
    
    let gradientDl = ctx.createLinearGradient(0, 0, 0, 300);
    gradientDl.addColorStop(0, 'rgba(16, 185, 129, 0.1)');
    gradientDl.addColorStop(1, 'rgba(16, 185, 129, 0)');

    let gradientUl = ctx.createLinearGradient(0, 0, 0, 300);
    gradientUl.addColorStop(0, 'rgba(59, 130, 246, 0.1)');
    gradientUl.addColorStop(1, 'rgba(59, 130, 246, 0)');

    const chart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [
                {
                    label: 'Download',
                    borderColor: '#10b981',
                    backgroundColor: gradientDl,
                    borderWidth: 2,
                    pointRadius: 0,
                    pointHoverRadius: 4,
                    pointBackgroundColor: '#10b981',
                    data: [],
                    tension: 0.4, 
                    fill: true
                },
                {
                    label: 'Upload',
                    borderColor: '#3b82f6',
                    backgroundColor: gradientUl,
                    borderWidth: 2,
                    pointRadius: 0,
                    pointHoverRadius: 4,
                    pointBackgroundColor: '#3b82f6',
                    data: [],
                    tension: 0.4,
                    fill: true
                },
                {
                    label: 'Ping',
                    borderColor: '#fbbf24',
                    borderWidth: 2,
                    borderDash: [5, 5],
                    pointRadius: 0,
                    pointHoverRadius: 4,
                    pointBackgroundColor: '#fbbf24',
                    data: [],
                    tension: 0.4
                },
                {
                    label: 'Packet Loss',
                    borderColor: '#ef4444',
                    borderWidth: 2,
                    borderDash: [3, 3],
                    pointRadius: 0,
                    pointHoverRadius: 4,
                    pointBackgroundColor: '#ef4444',
                    data: [],
                    tension: 0.4
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            interaction: {
                mode: 'index',
                intersect: false,
            },
            plugins: {
                legend: { display: false },
                tooltip: {
                    backgroundColor: '#0f172a',
                    titleColor: '#f1f5f9',
                    bodyColor: '#e2e8f0',
                    borderColor: '#1e293b',
                    borderWidth: 1,
                    padding: 12,
                    boxPadding: 4,
                    usePointStyle: true,
                    callbacks: {
                        label: function(context) {
                            let label = context.dataset.label || '';
                            if (label) {
                                label += ': ';
                            }
                            if (context.parsed.y !== null) {
                                let val = Number(context.parsed.y);
                                if (context.dataset.label === 'Packet Loss' && val === 0) {
                                    label += '0';
                                } else {
                                    label += val.toFixed(2);
                                }
                                
                                if (context.dataset.label === 'Ping') {
                                    label += ' ms';
                                } else if (context.dataset.label === 'Packet Loss') {
                                    label += '%';
                                } else {
                                    label += ' Mbps';
                                }
                            }
                            return label;
                        }
                    }
                }
            },
            scales: {
                x: {
                    ticks: { color: '#64748b', maxTicksLimit: 6, font: {size: 11} },
                    grid: { display: false }
                },
                y: {
                    beginAtZero: true,
                    ticks: { color: '#64748b', font: {size: 11} },
                    grid: { color: '#1e293b', borderDash: [4, 4] },
                    border: { display: false }
                }
            }
        }
    });

    function updateChart() {
        fetch('/api/stats')
            .then(response => response.json())
            .then(data => {
                if (!data) return;
                
                chart.data.labels = data.map(d => {
                    const date = new Date(d.created_at);
                    return date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
                });
                
                chart.data.datasets[0].data = data.map(d => d.download);
                chart.data.datasets[1].data = data.map(d => d.upload);
                chart.data.datasets[2].data = data.map(d => d.ping);
                chart.data.datasets[3].data = data.map(d => d.packet_loss);
                
                chart.update('none');
            });
    }

    function checkStatus() {
        const container = document.getElementById('status-container');
        if (!container) return;

        fetch('/partials/status')
            .then(response => response.text())
            .then(newHtml => {
                const cleanNew = newHtml.replace(/\s+/g, '');
                const cleanOld = container.innerHTML.replace(/\s+/g, '');

                if (cleanNew !== cleanOld) {
                    container.innerHTML = newHtml;
                    htmx.process(container);
                }
            })
            .catch(err => console.error("Status check failed", err));
    }

    updateChart();
    setInterval(updateChart, 15000);
    setInterval(checkStatus, 2000);
    
    document.body.addEventListener('refreshChart', function() {
        setTimeout(updateChart, 1000);
    });
});