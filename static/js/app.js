// Global variables
let currentPage = 'home';
let token = localStorage.getItem('token');
let user = JSON.parse(localStorage.getItem('user'));
let currentResumeId = null;
let resumes = [];
let skillsChart = null;
let activityChart = null;

// DOM ready
document.addEventListener('DOMContentLoaded', function() {
    // Initialize the app
    init();
    
    // Add event listeners
    setupEventListeners();
    
    // Handle routing
    handleRouting();
});

// Initialize the application
function init() {
    // Check if user is logged in
    if (token && user) {
        updateNavForLoggedInUser();
    } else {
        updateNavForLoggedOutUser();
    }
}

// Setup event listeners
function setupEventListeners() {
    // Navigation links
    document.querySelectorAll('.nav-link').forEach(link => {
        link.addEventListener('click', function(e) {
            const href = this.getAttribute('href');
            if (href.startsWith('#')) {
                e.preventDefault();
                navigateTo(href.substring(1));
            }
        });
    });
    
    // Get Started button
    document.getElementById('get-started-btn').addEventListener('click', function() {
        if (token && user) {
            navigateTo('dashboard');
        } else {
            navigateTo('register');
        }
    });
    
    // Learn More button
    document.getElementById('learn-more-btn').addEventListener('click', function() {
        window.scrollTo({
            top: document.querySelector('.row.mb-5.mt-5').offsetTop - 100,
            behavior: 'smooth'
        });
    });
    
    // Auth navigation
    document.getElementById('go-to-register').addEventListener('click', function(e) {
        e.preventDefault();
        navigateTo('register');
    });
    
    document.getElementById('go-to-login').addEventListener('click', function(e) {
        e.preventDefault();
        navigateTo('login');
    });
    
    // Back to dashboard button
    document.getElementById('back-to-dashboard').addEventListener('click', function() {
        navigateTo('dashboard');
    });
    
    // Form submissions
    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
    document.getElementById('resume-upload-form').addEventListener('submit', handleResumeUpload);
    
    // Logout
    document.getElementById('nav-logout').addEventListener('click', function(e) {
        e.preventDefault();
        handleLogout();
    });
    
    // File drop zone
    setupDropZone();
}

// Handle routing based on hash
function handleRouting() {
    window.addEventListener('hashchange', function() {
        const hash = window.location.hash.substring(1) || 'home';
        navigateTo(hash);
    });
    
    // Initial route
    const hash = window.location.hash.substring(1) || 'home';
    navigateTo(hash);
}

// Navigate to a specific page
function navigateTo(page) {
    // Check if user is logged in for protected routes
    if (['dashboard', 'upload'].includes(page) && (!token || !user)) {
        navigateTo('login');
        return;
    }
    
    // Hide all pages
    document.querySelectorAll('.page').forEach(p => {
        p.classList.remove('active');
    });
    
    // Show the selected page
    const pageElement = document.getElementById(`${page}-page`);
    if (pageElement) {
        pageElement.classList.add('active');
        currentPage = page;
        
        // Update navigation
        updateActiveNavItem(page);
        
        // Load page-specific data
        if (page === 'dashboard' && token) {
            loadDashboardData();
        }
        
        // Update URL hash
        window.location.hash = `#${page}`;
    } else {
        // Fallback to home page
        document.getElementById('home-page').classList.add('active');
        currentPage = 'home';
        updateActiveNavItem('home');
        window.location.hash = '#home';
    }
}

// Update active navigation item
function updateActiveNavItem(page) {
    document.querySelectorAll('.nav-item').forEach(item => {
        item.classList.remove('active');
    });
    
    const navItem = document.getElementById(`nav-${page}`);
    if (navItem) {
        navItem.classList.add('active');
    }
}

// Update navigation for logged in user
function updateNavForLoggedInUser() {
    document.getElementById('nav-login').classList.add('d-none');
    document.getElementById('nav-register').classList.add('d-none');
    document.getElementById('nav-logout').classList.remove('d-none');
    document.getElementById('nav-dashboard').classList.remove('d-none');
    document.getElementById('nav-upload').classList.remove('d-none');
}

// Update navigation for logged out user
function updateNavForLoggedOutUser() {
    document.getElementById('nav-login').classList.remove('d-none');
    document.getElementById('nav-register').classList.remove('d-none');
    document.getElementById('nav-logout').classList.add('d-none');
    document.getElementById('nav-dashboard').classList.add('d-none');
    document.getElementById('nav-upload').classList.add('d-none');
}

// Handle login form submission
async function handleLogin(e) {
    e.preventDefault();
    
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;
    
    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });
        
        if (!response.ok) {
            throw new Error('Login failed');
        }
        
        const data = await response.json();
        
        // Save token and user data
        localStorage.setItem('token', data.token);
        localStorage.setItem('user', JSON.stringify(data.user));
        
        // Update global variables
        token = data.token;
        user = data.user;
        
        // Update UI
        updateNavForLoggedInUser();
        
        // Navigate to dashboard
        navigateTo('dashboard');
        
        // Reset form
        document.getElementById('login-form').reset();
        
    } catch (error) {
        alert('Login failed. Please check your credentials and try again.');
        console.error('Login error:', error);
    }
}

// Handle register form submission
async function handleRegister(e) {
    e.preventDefault();
    
    const username = document.getElementById('register-username').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;
    const confirmPassword = document.getElementById('register-confirm-password').value;
    
    // Validate passwords match
    if (password !== confirmPassword) {
        alert('Passwords do not match');
        return;
    }
    
    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, email, password })
        });
        
        if (!response.ok) {
            throw new Error('Registration failed');
        }
        
        const data = await response.json();
        
        // Save token and user data
        localStorage.setItem('token', data.token);
        localStorage.setItem('user', JSON.stringify(data.user));
        
        // Update global variables
        token = data.token;
        user = data.user;
        
        // Update UI
        updateNavForLoggedInUser();
        
        // Navigate to dashboard
        navigateTo('dashboard');
        
        // Reset form
        document.getElementById('register-form').reset();
        
    } catch (error) {
        alert('Registration failed. Please try again.');
        console.error('Registration error:', error);
    }
}

// Handle logout
function handleLogout() {
    // Clear local storage
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    
    // Reset global variables
    token = null;
    user = null;
    
    // Update UI
    updateNavForLoggedOutUser();
    
    // Navigate to home
    navigateTo('home');
}

// Load dashboard data
async function loadDashboardData() {
    try {
        const response = await fetch('/api/resumes', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to load resumes');
        }
        
        const data = await response.json();
        resumes = data;
        
        // Update dashboard stats
        updateDashboardStats();
        
        // Render resumes
        renderResumes();
        
        // Update charts
        updateCharts();
        
    } catch (error) {
        console.error('Dashboard data error:', error);
        
        // Show no resumes message
        document.getElementById('no-resumes-message').style.display = 'block';
    }
}

// Update dashboard statistics
function updateDashboardStats() {
    // Resume count
    document.getElementById('resume-count').textContent = resumes.length;
    
    // Last upload date
    if (resumes.length > 0) {
        const lastUpload = new Date(resumes[0].uploaded_at);
        document.getElementById('last-upload').textContent = lastUpload.toLocaleDateString();
    } else {
        document.getElementById('last-upload').textContent = '-';
    }
    
    // Top skill
    if (resumes.length > 0) {
        const skillsMap = {};
        
        // Count occurrences of each skill
        resumes.forEach(resume => {
            resume.skills.forEach(skill => {
                skillsMap[skill] = (skillsMap[skill] || 0) + 1;
            });
        });
        
        // Find skill with highest count
        let topSkill = '-';
        let maxCount = 0;
        
        for (const skill in skillsMap) {
            if (skillsMap[skill] > maxCount) {
                maxCount = skillsMap[skill];
                topSkill = skill;
            }
        }
        
        document.getElementById('top-skill').textContent = topSkill;
    } else {
        document.getElementById('top-skill').textContent = '-';
    }
    
    // Show/hide no resumes message
    if (resumes.length === 0) {
        document.getElementById('no-resumes-message').style.display = 'block';
    } else {
        document.getElementById('no-resumes-message').style.display = 'none';
    }
}

// Render resumes in the dashboard
function renderResumes() {
    const container = document.getElementById('resumes-container');
    
    // Clear existing content except the no-resumes message
    Array.from(container.children).forEach(child => {
        if (child.id !== 'no-resumes-message') {
            container.removeChild(child);
        }
    });
    
    // Add resume cards
    resumes.forEach(resume => {
        const card = document.createElement('div');
        card.className = 'card resume-card';
        
        const uploadDate = new Date(resume.uploaded_at).toLocaleDateString();
        
        card.innerHTML = `
            <div class="card-body">
                <div class="d-flex justify-content-between align-items-center">
                    <h5 class="card-title">${resume.filename}</h5>
                    <small class="text-muted">${uploadDate}</small>
                </div>
                <div class="mt-3">
                    ${resume.skills.slice(0, 5).map(skill => `<span class="skill-badge">${skill}</span>`).join('')}
                    ${resume.skills.length > 5 ? `<span class="skill-badge">+${resume.skills.length - 5} more</span>` : ''}
                </div>
                <button class="btn btn-outline-primary mt-3 view-resume-btn" data-id="${resume.id}">View Analysis</button>
            </div>
        `;
        
        // Add event listener to view button
        card.querySelector('.view-resume-btn').addEventListener('click', function() {
            const resumeId = this.getAttribute('data-id');
            viewResumeDetails(resumeId);
        });
        
        container.appendChild(card);
    });
}

// View resume details
function viewResumeDetails(resumeId) {
    currentResumeId = resumeId;
    
    // Find the resume
    const resume = resumes.find(r => r.id === resumeId);
    
    if (!resume) {
        alert('Resume not found');
        return;
    }
    
    // Populate the details page
    document.getElementById('detail-filename').textContent = resume.filename;
    document.getElementById('detail-uploaded').textContent = new Date(resume.uploaded_at).toLocaleDateString();
    document.getElementById('detail-skills-count').textContent = resume.skills.length;
    document.getElementById('detail-content').textContent = resume.content;
    
    // Skills
    const skillsContainer = document.getElementById('detail-skills');
    skillsContainer.innerHTML = '';
    
    if (resume.skills.length > 0) {
        resume.skills.forEach(skill => {
            const badge = document.createElement('span');
            badge.className = 'skill-badge';
            badge.textContent = skill;
            skillsContainer.appendChild(badge);
        });
    } else {
        skillsContainer.innerHTML = '<p>No skills detected</p>';
    }
    
    // Education
    const educationContainer = document.getElementById('detail-education');
    educationContainer.innerHTML = '';
    
    if (resume.education.length > 0) {
        const list = document.createElement('ul');
        list.className = 'list-group';
        
        resume.education.forEach(edu => {
            const item = document.createElement('li');
            item.className = 'list-group-item';
            item.textContent = edu;
            list.appendChild(item);
        });
        
        educationContainer.appendChild(list);
    } else {
        educationContainer.innerHTML = '<p>No education details detected</p>';
    }
    
    // Experience
    const experienceContainer = document.getElementById('detail-experience');
    experienceContainer.innerHTML = '';
    
    if (resume.experience.length > 0) {
        const list = document.createElement('ul');
        list.className = 'list-group';
        
        resume.experience.forEach(exp => {
            const item = document.createElement('li');
            item.className = 'list-group-item';
            item.textContent = exp;
            list.appendChild(item);
        });
        
        experienceContainer.appendChild(list);
    } else {
        experienceContainer.innerHTML = '<p>No experience details detected</p>';
    }
    
    // Navigate to detail page
    navigateTo('resume-detail');
}

// Update charts
function updateCharts() {
    // Skills chart
    updateSkillsChart();
    
    // Activity chart
    updateActivityChart();
}

// Update skills chart
function updateSkillsChart() {
    const ctx = document.getElementById('skills-chart').getContext('2d');
    
    // Collect all skills
    const skillsMap = {};
    
    resumes.forEach(resume => {
        resume.skills.forEach(skill => {
            skillsMap[skill] = (skillsMap[skill] || 0) + 1;
        });
    });
    
    // Sort skills by count
    const sortedSkills = Object.entries(skillsMap)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 7); // Top 7 skills
    
    const labels = sortedSkills.map(item => item[0]);
    const data = sortedSkills.map(item => item[1]);
    
    // Destroy existing chart if it exists
    if (skillsChart) {
        skillsChart.destroy();
    }
    
    // Create new chart
    skillsChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Skill Frequency',
                data: data,
                backgroundColor: 'rgba(67, 97, 238, 0.7)',
                borderColor: 'rgba(67, 97, 238, 1)',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        precision: 0
                    }
                }
            }
        }
    });
}

// Update activity chart
function updateActivityChart() {
    const ctx = document.getElementById('activity-chart').getContext('2d');
    
    // Group resumes by month
    const monthlyActivity = {};
    
    resumes.forEach(resume => {
        const date = new Date(resume.uploaded_at);
        const monthYear = `${date.getMonth() + 1}/${date.getFullYear()}`;
        
        monthlyActivity[monthYear] = (monthlyActivity[monthYear] || 0) + 1;
    });
    
    // Get last 6 months
    const today = new Date();
    const labels = [];
    const data = [];
    
    for (let i = 5; i >= 0; i--) {
        const month = new Date(today.getFullYear(), today.getMonth() - i, 1);
        const monthYear = `${month.getMonth() + 1}/${month.getFullYear()}`;
        
        labels.push(monthYear);
        data.push(monthlyActivity[monthYear] || 0);
    }
    
    // Destroy existing chart if it exists
    if (activityChart) {
        activityChart.destroy();
    }
    
    // Create new chart
    activityChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Resumes Uploaded',
                data: data,
                backgroundColor: 'rgba(76, 201, 240, 0.2)',
                borderColor: 'rgba(76, 201, 240, 1)',
                borderWidth: 2,
                tension: 0.3
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        precision: 0
                    }
                }
            }
        }
    });
}

// Handle resume upload
async function handleResumeUpload(e) {
    e.preventDefault();
    
    const fileInput = document.getElementById('resume-file');
    
    if (!fileInput.files || fileInput.files.length === 0) {
        alert('Please select a PDF file to upload');
        return;
    }
    
    const file = fileInput.files[0];
    
    // Validate file type
    if (file.type !== 'application/pdf') {
        alert('Only PDF files are allowed');
        return;
    }
    
    // Show loading spinner
    document.getElementById('upload-spinner').style.display = 'inline-block';
    document.getElementById('upload-button').style.display = 'none';
    
    // Create form data
    const formData = new FormData();
    formData.append('resume', file);
    
    try {
        const response = await fetch('/api/resumes/upload', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: formData
        });
        
        if (!response.ok) {
            throw new Error('Upload failed');
        }
        
        const data = await response.json();
        
        // Hide loading spinner
        document.getElementById('upload-spinner').style.display = 'none';
        document.getElementById('upload-button').style.display = 'inline-block';
        
        // Reset form
        document.getElementById('resume-upload-form').reset();
        
        // Navigate to dashboard
        navigateTo('dashboard');
        
        // Show success message
        alert('Resume uploaded and analyzed successfully!');
        
    } catch (error) {
        // Hide loading spinner
        document.getElementById('upload-spinner').style.display = 'none';
        document.getElementById('upload-button').style.display = 'inline-block';
        
        alert('Upload failed. Please try again.');
        console.error('Upload error:', error);
    }
}

// Setup drag and drop zone
function setupDropZone() {
    const dropZone = document.getElementById('drop-zone');
    const fileInput = document.getElementById('resume-file');
    
    // Highlight drop zone when dragging over it
    ['dragenter', 'dragover'].forEach(eventName => {
        dropZone.addEventListener(eventName, function(e) {
            e.preventDefault();
            this.classList.add('active');
        });
    });
    
    // Remove highlight when dragging leaves
    ['dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, function(e) {
            e.preventDefault();
            this.classList.remove('active');
        });
    });
    
    // Handle dropped files
    dropZone.addEventListener('drop', function(e) {
        e.preventDefault();
        
        if (e.dataTransfer.files.length) {
            fileInput.files = e.dataTransfer.files;
            updateFileNameDisplay(e.dataTransfer.files[0].name);
        }
    });
    
    // Handle file input change
    fileInput.addEventListener('change', function() {
        if (this.files.length) {
            updateFileNameDisplay(this.files[0].name);
        }
    });
    
    // Click on drop zone to trigger file input
    dropZone.addEventListener('click', function() {
        fileInput.click();
    });
}

// Update file name display in drop zone
function updateFileNameDisplay(filename) {
    const dropZone = document.getElementById('drop-zone');
    const prompt = dropZone.querySelector('.drop-zone-prompt');
    
    // Check if a filename element already exists
    let filenameElement = dropZone.querySelector('.drop-zone-filename');
    
    if (!filenameElement) {
        filenameElement = document.createElement('div');
        filenameElement.className = 'drop-zone-filename';
        dropZone.appendChild(filenameElement);
    }
    
    // Update filename and hide prompt
    prompt.style.display = 'none';
    filenameElement.textContent = filename;
}
