// Global date formatters
const formatDate = (dateStr) => {
    if (!dateStr || dateStr === "0001-01-01") return 'N/A';
    try {
        // Try parsing as RFC3339 first
        let date = new Date(dateStr);
        if (isNaN(date.getTime())) {
            // If that fails, try parsing as YYYY-MM-DD
            const [year, month, day] = dateStr.split('-');
            date = new Date(year, month - 1, day);
        }
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
        if (isNaN(date.getTime())) {
            // Try parsing as YYYY-MM-DD
            const [year, month, day] = dateStr.split('-');
            const parsedDate = new Date(year, month - 1, day);
            if (!isNaN(parsedDate.getTime())) {
                return `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`;
            }
            return '';
        }
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
    rechargeDate.setHours(12, 0, 0, 0); // Set to noon to avoid timezone issues
    
    // Add validation
    if (isNaN(rechargeDate.getTime())) {
        alert('Please enter a valid recharge date');
        return;
    }

    const simData = {
        name: document.getElementById('simName').value,
        number: document.getElementById('simNumber').value,
        last_recharge_date: rechargeDate.toISOString(), // This will format as RFC3339
        recharge_validity: new Date(rechargeDate.getTime() + (30 * 24 * 60 * 60 * 1000)).toISOString(),
        incoming_call_validity: new Date(rechargeDate.getTime() + (45 * 24 * 60 * 60 * 1000)).toISOString(),
        sim_expiry: new Date(rechargeDate.getTime() + (90 * 24 * 60 * 60 * 1000)).toISOString()
    };

    // Debug log
    console.log('Sending SIM data:', {
        ...simData,
        last_recharge_date_parsed: new Date(simData.last_recharge_date),
        recharge_validity_parsed: new Date(simData.recharge_validity),
        incoming_call_validity_parsed: new Date(simData.incoming_call_validity),
        sim_expiry_parsed: new Date(simData.sim_expiry)
    });

    try {
        const response = await fetch('/api/sims', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(simData)
        });
        
        // Debug log
        console.log('Response status:', response.status);
        
        if (response.ok) {
            const result = await response.json();
            console.log('Added SIM:', result);
            loadSims();
            document.getElementById('simForm').reset();
        } else {
            const error = await response.text();
            console.error('Error adding SIM:', error);
            try {
                // Try to parse error as JSON
                const errorJson = JSON.parse(error);
                alert('Failed to add SIM: ' + (errorJson.error || error));
            } catch {
                alert('Failed to add SIM: ' + error);
            }
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
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        
        // Ensure data is an array
        const sims = Array.isArray(data) ? data : [];
        console.log('Loaded SIMs:', sims); // Debug log
        
        const simList = document.getElementById('simList');
        simList.innerHTML = '';
        
        if (sims.length === 0) {
            simList.innerHTML = '<div class="alert alert-info">No SIMs found</div>';
            return;
        }
        
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
                                    onclick="openEditModal('${sim._id}', '${sim.name}', '${sim.number}', '${sim.last_recharge_date}')">
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
        console.error('Error loading SIMs:', error);
        const simList = document.getElementById('simList');
        simList.innerHTML = '<div class="alert alert-danger">Failed to load SIMs. Please try again later.</div>';
    }
}

// Open Edit Modal function
function openEditModal(id, name, number, lastRecharge) {
    try {
        // Initialize modal first
        const editModal = new bootstrap.Modal(document.getElementById('editModal'));
        
        // Set values after ensuring elements exist
        document.getElementById('editSimId').value = id;
        document.getElementById('editSimName').value = name;
        document.getElementById('editSimNumber').value = number;
        document.getElementById('editLastRechargeDate').value = formatDateForInput(lastRecharge);

        // Show the modal
        editModal.show();
    } catch (error) {
        console.error('Error in openEditModal:', error);
    }
}

// Save Edit Handler
document.getElementById('saveEdit').addEventListener('click', async () => {
    const simId = document.getElementById('editSimId').value;
    const rechargeDate = new Date(document.getElementById('editLastRechargeDate').value);
    rechargeDate.setHours(12, 0, 0, 0); // Set to noon to avoid timezone issues

    const updatedData = {
        last_recharge_date: rechargeDate.toISOString(),
        recharge_validity: new Date(rechargeDate.getTime() + (30 * 24 * 60 * 60 * 1000)).toISOString(),
        incoming_call_validity: new Date(rechargeDate.getTime() + (45 * 24 * 60 * 60 * 1000)).toISOString(),
        sim_expiry: new Date(rechargeDate.getTime() + (90 * 24 * 60 * 60 * 1000)).toISOString()
    };

    try {
        const response = await fetch(`/api/sims/${simId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(updatedData)
        });

        if (response.ok) {
            bootstrap.Modal.getInstance(document.getElementById('editModal')).hide();
            loadSims();
        } else {
            const errorText = await response.text();
            console.error('Update failed:', errorText);
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