// In your /static/main.js or equivalent JS file

function formatAndToggleMessage(containerElement) {
    const pre = containerElement.querySelector('pre.message-body');
    const toggleButton = containerElement.querySelector('button.toggle-msg');
    const copyButton = containerElement.querySelector('button.copy-msg');

    if (!pre || !toggleButton || !copyButton) {
        console.warn('formatAndToggleMessage: Missing elements', { pre: !!pre, toggleButton: !!toggleButton, copyButton: !!copyButton });
        return;
    }

    const toggleIcon = toggleButton.querySelector('i');
    const toggleTextSpan = toggleButton.querySelector('span.toggle-text');
    const copyIcon = copyButton.querySelector('i');
    const copyTextSpan = copyButton.querySelector('span.copy-text');
    const messageTypeIndicator = containerElement.querySelector('span.message-type-indicator');

    if (!toggleIcon || !copyIcon || !messageTypeIndicator) {
        return;
    }
    // Note: toggleTextSpan and copyTextSpan are optional

    const originalRawText = pre.textContent.trim();
    let isJson = false;
    let singleLineText = originalRawText; // For collapsed view
    let prettyPrintedText = originalRawText; // For expanded view and copying

    try {
        const trimmedText = originalRawText;
        if ((trimmedText.startsWith("{") && trimmedText.endsWith("}")) || (trimmedText.startsWith("[") && trimmedText.endsWith("]"))) {
            const parsedJson = JSON.parse(trimmedText);
            singleLineText = JSON.stringify(parsedJson); // Compact JSON for single line
            prettyPrintedText = JSON.stringify(parsedJson, null, 2); // Pretty print for expansion/copy
            isJson = true;
            messageTypeIndicator.textContent = '(JSON)';
        } else {
            messageTypeIndicator.textContent = '(Text)';
        }
    } catch (e) {
        messageTypeIndicator.textContent = '(Text - Non-JSON)';
    }

    // Function to set styles for collapsed state
    function collapseMessage() {
        pre.textContent = singleLineText; // Show compact/original text
        pre.classList.add('collapsed');
        pre.classList.remove('expanded');
        // Use !important via setProperty or inline styles
        pre.style.setProperty('max-height', '3em', 'important');
        pre.style.setProperty('overflow', 'hidden', 'important');
        pre.style.setProperty('white-space', 'nowrap', 'important');
        pre.style.setProperty('text-overflow', 'ellipsis', 'important');
        if (toggleIcon) toggleIcon.className = 'bi bi-arrows-expand';
        if (toggleTextSpan) toggleTextSpan.textContent = 'Expand';
    }

    // Function to set styles for expanded state
    function expandMessage() {
        pre.textContent = prettyPrintedText; // Show pretty-printed/original text
        pre.classList.remove('collapsed');
        pre.classList.add('expanded');
        // Use !important via setProperty or inline styles
        pre.style.setProperty('max-height', '800px', 'important');
        pre.style.setProperty('overflow-y', 'auto', 'important');
        pre.style.setProperty('white-space', 'pre-wrap', 'important');
        pre.style.setProperty('text-overflow', 'clip', 'important');
        if (toggleIcon) toggleIcon.className = 'bi bi-arrows-collapse';
        if (toggleTextSpan) toggleTextSpan.textContent = 'Collapse';
    }

    // Remove any existing event listeners by cloning and replacing buttons
    const newToggleButton = toggleButton.cloneNode(true);
    toggleButton.parentNode.replaceChild(newToggleButton, toggleButton);
    
    const newCopyButton = copyButton.cloneNode(true);
    copyButton.parentNode.replaceChild(newCopyButton, copyButton);
    
    // Initial state: collapsed - ensure it starts collapsed
    collapseMessage();
    
    // Debug logging
    console.log('Message initialized:', {
        hasCollapsed: pre.classList.contains('collapsed'),
        hasExpanded: pre.classList.contains('expanded'),
        maxHeight: pre.style.maxHeight,
        toggleButtonExists: !!newToggleButton,
        copyButtonExists: !!newCopyButton
    });

    // Attach click handler to the new button
    newToggleButton.onclick = function(e) {
        e.preventDefault();
        e.stopPropagation();
        console.log('Toggle clicked, current state:', pre.classList.contains('collapsed') ? 'collapsed' : 'expanded');
        if (pre.classList.contains('collapsed')) {
            expandMessage();
        } else {
            collapseMessage();
        }
    };

    // Attach click handler to the new copy button
    newCopyButton.onclick = function(e) {
        e.preventDefault();
        e.stopPropagation();
        const textToCopy = prettyPrintedText;
        console.log('Copy clicked, text length:', textToCopy.length);
        
        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(textToCopy).then(() => {
                const originalIconClass = newCopyButton.querySelector('i').className;
                const originalText = newCopyButton.querySelector('span.copy-text') ? newCopyButton.querySelector('span.copy-text').textContent : '';
                
                newCopyButton.querySelector('i').className = 'bi bi-check-lg text-success';
                if (newCopyButton.querySelector('span.copy-text')) {
                    newCopyButton.querySelector('span.copy-text').textContent = 'Copied!';
                }
                
                setTimeout(() => {
                    newCopyButton.querySelector('i').className = originalIconClass;
                    if (newCopyButton.querySelector('span.copy-text')) {
                        newCopyButton.querySelector('span.copy-text').textContent = originalText;
                    }
                }, 2000);
            }).catch(err => {
                console.error('Failed to copy text: ', err);
                alert('Failed to copy to clipboard. Please try again.');
            });
        } else {
            // Fallback for older browsers
            const textArea = document.createElement('textarea');
            textArea.value = textToCopy;
            textArea.style.position = 'fixed';
            textArea.style.opacity = '0';
            document.body.appendChild(textArea);
            textArea.select();
            try {
                document.execCommand('copy');
                const originalIconClass = newCopyButton.querySelector('i').className;
                newCopyButton.querySelector('i').className = 'bi bi-check-lg text-success';
                setTimeout(() => {
                    newCopyButton.querySelector('i').className = originalIconClass;
                }, 2000);
            } catch (err) {
                console.error('Fallback copy failed:', err);
                alert('Failed to copy to clipboard.');
            }
            document.body.removeChild(textArea);
        }
    };
}

// Global initialization function for message viewers
function initializeMessageContentViewers(parentElement) {
    if (!parentElement) return;
    parentElement.querySelectorAll('.message-content-container').forEach(container => {
        // Only initialize if not already initialized
        if (!container.dataset.initialized) {
            formatAndToggleMessage(container);
            container.dataset.initialized = 'true';
        }
    });
}

// Note: Initialization is handled in messages.html to avoid conflicts
// This function is available for manual initialization if needed
