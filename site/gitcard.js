function initGitRepoCard(card) {
    const repoUrl = card.dataset.repoUrl;
    
    // Extract owner and repo from GitHub URL
    const match = repoUrl.match(/github\.com\/([^\/]+)\/([^\/]+)/);
    if (!match) {
        card.innerHTML = '<div class="repo-error">Invalid GitHub URL</div>';
        return;
    }
    
    const [, owner, repo] = match;
    const apiUrl = `https://api.github.com/repos/${owner}/${repo}`;
    
    fetch(apiUrl)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            const languageColor = getLanguageColor(data.language);
            
            // Clear card
            card.innerHTML = '';
            
            // Create avatar section
            const avatarSection = document.createElement('div');
            avatarSection.className = 'repo-avatar-section';
            
            const avatar = document.createElement('img');
            avatar.className = 'repo-avatar';
            avatar.src = data.owner.avatar_url;
            avatar.alt = data.owner.login;
            
            avatarSection.appendChild(avatar);
            
            // Create content section
            const contentSection = document.createElement('div');
            contentSection.className = 'repo-content';
            
            // Create header with repo name
            const header = document.createElement('div');
            header.className = 'repo-header';
            
            const repoLink = document.createElement('a');
            repoLink.className = 'repo-name';
            repoLink.href = data.html_url;
            repoLink.target = '_blank';
            repoLink.rel = 'noopener';
            repoLink.textContent = data.owner.login + '/' + data.name;
            
            header.appendChild(repoLink);
            contentSection.appendChild(header);
            
            // Add description if present
            if (data.description) {
                const desc = document.createElement('div');
                desc.className = 'repo-description';
                desc.textContent = data.description;
                contentSection.appendChild(desc);
            }
            
            // Add stats section
            const stats = document.createElement('div');
            stats.className = 'repo-stats';
            
            // Language stat
            if (data.language) {
                const langStat = document.createElement('div');
                langStat.className = 'repo-stat';
                
                const langColor = document.createElement('div');
                langColor.className = 'repo-language-color';
                langColor.style.backgroundColor = languageColor;
                
                const langText = document.createElement('span');
                langText.textContent = data.language;
                
                langStat.appendChild(langColor);
                langStat.appendChild(langText);
                stats.appendChild(langStat);
            }
            
            // Stars stat
            const starsStat = document.createElement('div');
            starsStat.className = 'repo-stat';
            starsStat.innerHTML = '<svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor"><path d="M8 .25a.75.75 0 01.673.418l1.882 3.815 4.21.612a.75.75 0 01.416 1.279l-3.046 2.97.719 4.192a.75.75 0 01-1.088.791L8 12.347l-3.766 1.98a.75.75 0 01-1.088-.79l.72-4.194L.818 6.374a.75.75 0 01.416-1.28l4.21-.611L7.327.668A.75.75 0 018 .25z"/></svg><span>' + data.stargazers_count + '</span>';
            stats.appendChild(starsStat);
            
            // Forks stat
            const forksStat = document.createElement('div');
            forksStat.className = 'repo-stat';
            forksStat.innerHTML = '<svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor"><path d="M5 3.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm0 2.122a2.25 2.25 0 10-1.5 0v.878A2.25 2.25 0 005.75 8.5h1.5v2.128a2.25 2.25 0 101.5 0V8.5h1.5a2.25 2.25 0 002.25-2.25v-.878a2.25 2.25 0 10-1.5 0v.878a.75.75 0 01-.75.75h-4.5A.75.75 0 015 6.25v-.878z"/></svg><span>' + data.forks_count + '</span>';
            stats.appendChild(forksStat);
            
            // License stat
            if (data.license) {
                const licenseStat = document.createElement('div');
                licenseStat.className = 'repo-stat';
                licenseStat.innerHTML = '<svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor"><path d="M8.75.75V2h.985c.304 0 .603.08.867.231l1.29.736c.038.022.08.033.124.033h2.234a.75.75 0 010 1.5h-.427l2.111 4.692a.75.75 0 01-.154.838l-.53-.53.529.531-.001.002-.002.002-.006.006-.016.015-.045.04c-.21.176-.441.327-.686.45C14.556 10.78 13.88 11 13 11a4.498 4.498 0 01-2.023-.454 3.544 3.544 0 01-.686-.45l-.045-.04-.016-.015-.006-.006-.002-.002v-.001l.529-.531-.53.53a.75.75 0 01-.154-.838L12.178 4.5h-.162c-.305 0-.604-.079-.868-.231l-1.29-.736a.245.245 0 00-.124-.033H8.75V13h2.5a.75.75 0 010 1.5h-6.5a.75.75 0 010-1.5h2.5V3.5h-.984a.245.245 0 00-.124.033l-1.289.737c-.265.15-.564.23-.869.23h-.162l2.112 4.692a.75.75 0 01-.154.838l-.53-.53.529.531-.002.002-.002.002-.006.006-.016.015-.045.04c-.21.176-.441.327-.686.45C4.556 10.78 3.88 11 3 11a4.498 4.498 0 01-2.023-.454 3.544 3.544 0 01-.686-.45l-.045-.04-.016-.015-.006-.006-.002-.002v-.001l.529-.531-.53.53a.75.75 0 01-.154-.838L2.178 4.5H1.75a.75.75 0 010-1.5h2.234a.249.249 0 00.125-.033l1.288-.737c.265-.15.564-.23.869-.23h.984V.75a.75.75 0 011.5 0z"/></svg><span>' + data.license.spdx_id + '</span>';
                stats.appendChild(licenseStat);
            }
            
            contentSection.appendChild(stats);
            
            // Add both sections to the card
            card.appendChild(avatarSection);
            card.appendChild(contentSection);
        })
        .catch(error => {
            console.error('Error fetching repository data:', error);
            card.innerHTML = `<div class="repo-error">Failed to load repository information</div>`;
        });
    
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
    
    function getLanguageColor(language) {
        const colors = {
            'JavaScript': '#f1e05a',
            'TypeScript': '#2b7489',
            'Python': '#3572A5',
            'Java': '#b07219',
            'Go': '#00ADD8',
            'Rust': '#dea584',
            'C++': '#f34b7d',
            'C': '#555555',
            'C#': '#239120',
            'PHP': '#4F5D95',
            'Ruby': '#701516',
            'Swift': '#ffac45',
            'Kotlin': '#F18E33',
            'Dart': '#00B4AB',
            'HTML': '#e34c26',
            'CSS': '#1572B6',
            'Shell': '#89e051',
            'Dockerfile': '#384d54',
            'Starlark': '#76d275'
        };
        return colors[language] || '#586069';
    }
}

// Initialize all git repo cards when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    const cards = document.querySelectorAll('.git-repo-card');
    cards.forEach(card => {
        initGitRepoCard(card);
    });
});