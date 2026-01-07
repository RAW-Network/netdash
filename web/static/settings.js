document.addEventListener('DOMContentLoaded', function() {
    const els = {
        schedule: document.getElementById('scheduleSelect'),
        cron: document.getElementById('cronInput'),
        history: document.getElementById('historyInput'),
        server: document.getElementById('serverInput'),
        bar: document.getElementById('actionBar'),
        reset: document.getElementById('resetBtn'),
        modal: document.getElementById('deleteModal'),
        modalContent: document.getElementById('modalContent')
    };

    const form = document.getElementById('settingsForm');
    const initial = {
        cron: form.dataset.initialCron,
        server: form.dataset.initialServer,
        history: form.dataset.initialHistory
    };

    function init() {
        const isCustom = !Array.from(els.schedule.options).some(o => o.value === initial.cron);
        els.schedule.value = isCustom ? 'custom' : initial.cron;
        els.cron.value = initial.cron;
        
        toggleCron();
        checkChanges();
    }

    function toggleCron() {
        if (els.schedule.value !== 'custom') {
            els.cron.value = els.schedule.value;
            els.cron.classList.add('hidden');
        } else {
            els.cron.classList.remove('hidden');
        }
    }

    function checkChanges() {
        const current = {
            cron: els.schedule.value === 'custom' ? els.cron.value.trim() : els.schedule.value,
            server: els.server.value.trim(),
            history: els.history.value
        };

        const hasChanged = JSON.stringify(current) !== JSON.stringify(initial);
        
        if (hasChanged) {
            els.bar.classList.remove('max-h-0', 'opacity-0');
            els.bar.classList.add('max-h-24', 'opacity-100');
        } else {
            els.bar.classList.add('max-h-0', 'opacity-0');
            els.bar.classList.remove('max-h-24', 'opacity-100');
        }
    }

    els.schedule.addEventListener('change', () => { toggleCron(); checkChanges(); });
    els.cron.addEventListener('input', checkChanges);
    els.history.addEventListener('input', checkChanges);
    els.server.addEventListener('input', checkChanges);
    
    els.reset.addEventListener('click', () => {
        els.server.value = initial.server;
        els.history.value = initial.history;
        init();
    });

    window.openModal = function() {
        els.modal.classList.remove('opacity-0', 'pointer-events-none');
        setTimeout(() => els.modalContent.classList.remove('scale-95'), 10);
    }

    window.closeModal = function() {
        els.modalContent.classList.add('scale-95');
        els.modal.classList.add('opacity-0', 'pointer-events-none');
    }

    window.confirmClear = function() {
        document.getElementById('clearForm').submit();
    }

    init();
});