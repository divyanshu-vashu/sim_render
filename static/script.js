// Global date formatters
const formatDate = (dateStr) => {
    if (!dateStr || dateStr === "0001-01-01") return 'N/A';
    try {
        const date = new Date(dateStr);
        if (isNaN(date.getTime())) return 'N/A';
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    } catch (error) {
        console.error('Date parsing error:', error);
        return 'N/A';
    }
};

const formatDateForInput = (dateStr) => {
    if (!dateStr) return '';
    try {
        const date = new Date(dateStr);
        if (isNaN(date.getTime())) return '';
        return date.toISOString().split('T')[0];
    } catch (error) {
        console.error('Date formatting error:', error);
        return '';
    }
};

// SIM creation handler
document.getElementById('simForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const rechargeDate = new Date(document.getElementById('rechargeDate').value);
    rechargeDate.setHours(0, 0, 0, 0);
    
    const simData = {
        name: document.getElementById('simName').value,
        number: document.getElementById('simNumber').value,
        last_recharge_date: rechargeDate.toISOString().split('T')[0],
        recharge_validity: new Date(rechargeDate.getTime() + (30 * 24 * 60 * 60 * 1000)).toISOString().split('T')[0],
        incoming_call_validity: new Date(rechargeDate.getTime() + (45 * 24 * 60 * 60 * 1000)).toISOString().split('T')[0],
        sim_expiry: new Date(rechargeDate.getTime() + (90 * 24 * 60 * 60 * 1000)).toISOString().split('T')[0]
    };

    try {
        const response = await fetch('/api/sims', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(simData)
        });
        
        if (response.ok) {
            const result = await response.json();
            console.log('Added SIM:', result);
            loadSims();
            document.getElementById('simForm').reset();
        } else {
            const error = await response.text();
            console.error('Error adding SIM:', error);
            alert('Failed to add SIM: ' + error);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to add SIM');
    }
});

// Sync and filter event listeners
document.getElementById('syncButton').addEventListener('click', () => {
    loadSims();
});

document.getElementById('filterType').addEventListener('change', () => {
    loadSims();
});

// Load SIMs function
async function loadSims() {
    try {
        const response = await fetch('/api/sims');
        const sims = await response.json();
        
        const simList = document.getElementById('simList');
        simList.innerHTML = '';
        
        sims.forEach(sim => {
            const card = document.createElement('div');
            card.className = 'card mb-3';
            card.innerHTML = `
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start">
                        <div>
                            <h5 class="card-title">${sim.name} - ${sim.number}</h5>
                            <p class="mb-2">
                                Last Recharge: ${formatDate(sim.last_recharge_date)}
                                <button class="btn btn-sm btn-outline-primary ms-2" 
                                    onclick="openEditModal(
                                        ${sim.id}, 
                                        '${sim.name}', 
                                        '${sim.number}', 
                                        '${sim.last_recharge_date}', 
                                        '${sim.recharge_validity}', 
                                        '${sim.incoming_call_validity}', 
                                        '${sim.sim_expiry}'
                                    )">
                                    Edit
                                </button>
                            </p>
                            <p>Recharge Validity: ${formatDate(sim.recharge_validity)}</p>
                            <p>Incoming Call Validity: ${formatDate(sim.incoming_call_validity)}</p>
                            <p>SIM Expiry: ${formatDate(sim.sim_expiry)}</p>
                        </div>
                    </div>
                </div>
            `;
            simList.appendChild(card);
        });
    } catch (error) {
        console.error('Error:', error);
    }
}

// Open Edit Modal function
function openEditModal(id, name, number, lastRecharge, rechargeValidity, incomingValidity, simExpiry) {
    try {
        document.getElementById('editSimId').value = id;
        document.getElementById('editSimName').value = name;
        document.getElementById('editSimNumber').value = number;
        document.getElementById('editLastRechargeDate').value = formatDateForInput(lastRecharge);
        document.getElementById('editRechargeValidity').value = formatDateForInput(rechargeValidity);
        document.getElementById('editIncomingValidity').value = formatDateForInput(incomingValidity);
        document.getElementById('editSimExpiry').value = formatDateForInput(simExpiry);

        const editModal = new bootstrap.Modal(document.getElementById('editModal'));
        editModal.show();
    } catch (error) {
        console.error('Error in openEditModal:', error);
    }
}

// Save Edit Handler
document.getElementById('saveEdit').addEventListener('click', async () => {
    const simId = document.getElementById('editSimId').value;
    const newDate = document.getElementById('editRechargeDate').value;

    try {
        const response = await fetch(`/api/sims/${simId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                lastRechargeDate: newDate
            })
        });

        if (response.ok) {
            bootstrap.Modal.getInstance(document.getElementById('editModal')).hide();
            loadSims();
        } else {
            alert('Failed to update recharge date');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to update recharge date');
    }
});

// Status functions
function getStatus(sim, currentDate) {
    const rechargeDate = new Date(sim.rechargeValidity);
    rechargeDate.setHours(0, 0, 0, 0);
    currentDate.setHours(0, 0, 0, 0);
    const daysUntilExpiry = Math.ceil((rechargeDate - currentDate) / (1000 * 60 * 60 * 24));
    
    if (daysUntilExpiry < 0) return 'Expired';
    if (daysUntilExpiry <= 3) return 'Expiring Soon';
    return 'Active';
}

function getStatusBadgeClass(sim, currentDate) {
    const rechargeDate = new Date(sim.rechargeValidity);
    const daysUntilExpiry = Math.ceil((rechargeDate - currentDate) / (1000 * 60 * 60 * 24));
    
    if (daysUntilExpiry < 0) return 'bg-danger';
    if (daysUntilExpiry <= 3) return 'bg-warning';
    return 'bg-success';
}

// Load sims when page loads
loadSims();