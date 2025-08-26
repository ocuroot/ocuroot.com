// Simple syntax highlighting for Starlark/Python in terminal windows
document.addEventListener('DOMContentLoaded', function() {
    const terminalCodes = document.querySelectorAll('.terminal-code code');
    
    terminalCodes.forEach(code => {
        if (code.classList.contains('language-python')) {
            highlightPython(code);
        }
    });
});

function highlightPython(element) {
    const text = element.textContent;
    const lines = text.split('\n');
    
    element.innerHTML = '';
    
    lines.forEach((line, index) => {
        const lineElement = document.createElement('div');
        lineElement.style.display = 'block';
        
        // Handle empty lines by adding a non-breaking space to preserve height
        if (line.trim() === '') {
            lineElement.innerHTML = '&nbsp;';
        } else {
            // Tokenize the line
            const tokens = tokenizePython(line);
            
            tokens.forEach(token => {
                const span = document.createElement('span');
                span.textContent = token.text;
                span.className = token.type;
                lineElement.appendChild(span);
            });
        }
        
        element.appendChild(lineElement);
    });
}

function tokenizePython(line) {
    const tokens = [];
    let i = 0;
    
    while (i < line.length) {
        const char = line[i];
        
        // Skip whitespace but preserve it
        if (/\s/.test(char)) {
            let whitespace = '';
            while (i < line.length && /\s/.test(line[i])) {
                whitespace += line[i];
                i++;
            }
            tokens.push({ text: whitespace, type: '' });
            continue;
        }
        
        // Comments
        if (char === '#') {
            tokens.push({ text: line.slice(i), type: 'py-comment' });
            break;
        }
        
        // Strings
        if (char === '"' || char === "'") {
            const quote = char;
            let string = quote;
            i++;
            while (i < line.length && line[i] !== quote) {
                if (line[i] === '\\' && i + 1 < line.length) {
                    string += line[i] + line[i + 1];
                    i += 2;
                } else {
                    string += line[i];
                    i++;
                }
            }
            if (i < line.length) {
                string += line[i];
                i++;
            }
            tokens.push({ text: string, type: 'py-string' });
            continue;
        }
        
        // Numbers
        if (/\d/.test(char)) {
            let number = '';
            while (i < line.length && /[\d.]/.test(line[i])) {
                number += line[i];
                i++;
            }
            tokens.push({ text: number, type: 'py-number' });
            continue;
        }
        
        // Identifiers and keywords
        if (/[a-zA-Z_]/.test(char)) {
            let identifier = '';
            while (i < line.length && /[a-zA-Z0-9_]/.test(line[i])) {
                identifier += line[i];
                i++;
            }
            
            const keywords = ['def', 'if', 'else', 'elif', 'for', 'while', 'return', 'import', 'from', 'as', 'class', 'try', 'except', 'finally', 'with', 'pass', 'break', 'continue', 'and', 'or', 'not', 'in', 'is', 'None', 'True', 'False'];
            
            let type = '';
            if (keywords.includes(identifier)) {
                type = 'py-keyword';
            } else if (i < line.length && line[i] === '(') {
                type = 'py-function-call';
            }
            
            tokens.push({ text: identifier, type });
            continue;
        }
        
        // Everything else (operators, punctuation)
        tokens.push({ text: char, type: '' });
        i++;
    }
    
    return tokens;
}
