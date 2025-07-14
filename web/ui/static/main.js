// In your /static/main.js or equivalent JS file

function formatAndToggleMessage(containerElement) {
    const pre = containerElement.querySelector('pre.message-body');
    const toggleButton = containerElement.querySelector('a.toggle-msg');
    const copyButton = containerElement.querySelector('a.copy-msg'); // New copy button

    if (!pre || !toggleButton || !copyButton) {
        return;
    }

    const toggleIcon = toggleButton.querySelector('i');
    const toggleTextSpan = toggleButton.querySelector('span.toggle-text');
    const copyIcon = copyButton.querySelector('i');
    const copyTextSpan = copyButton.querySelector('span.copy-text');
    const messageTypeIndicator = containerElement.querySelector('span.message-type-indicator');

    if (!toggleIcon || !toggleTextSpan || !copyIcon || !copyTextSpan || !messageTypeIndicator) {
        return;
    }

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
        pre.style.maxHeight = '2.7em'; // Adjusted for one line + padding + border
        pre.style.overflowY = 'hidden';
        pre.style.whiteSpace = 'nowrap';
        pre.style.overflowX = 'hidden';
        pre.style.textOverflow = 'ellipsis';
        pre.classList.add('collapsed');
        pre.classList.remove('expanded');
        toggleIcon.className = 'bi bi-arrows-expand';
        if (toggleTextSpan) toggleTextSpan.textContent = 'Expand';
    }

    // Function to set styles for expanded state
    function expandMessage() {
        pre.textContent = prettyPrintedText; // Show pretty-printed/original text
        pre.style.maxHeight = '20em'; // Set a max-height for the expanded view (e.g., 20em)
        pre.style.overflowY = 'auto';
        pre.style.whiteSpace = 'pre-wrap';
        pre.style.overflowX = 'auto';
        pre.style.textOverflow = 'clip';
        pre.classList.remove('collapsed');
        pre.classList.add('expanded');
        toggleIcon.className = 'bi bi-arrows-collapse';
        if (toggleTextSpan) toggleTextSpan.textContent = 'Collapse';
    }

    // Initial state: collapsed
    collapseMessage();

    // --- Event Listeners (clone to avoid duplicates if re-initialized) ---
    const newToggleButton = toggleButton.cloneNode(true);
    toggleButton.parentNode.replaceChild(newToggleButton, toggleButton);
    
    newToggleButton.addEventListener('click', function(e) {
        e.preventDefault();
        if (pre.classList.contains('collapsed')) {
            expandMessage();
        } else {
            collapseMessage();
        }
    });

    const newCopyButton = copyButton.cloneNode(true);
    copyButton.parentNode.replaceChild(newCopyButton, copyButton);

    newCopyButton.addEventListener('click', function(e) {
        e.preventDefault();
        navigator.clipboard.writeText(prettyPrintedText).then(() => {
            const originalIconClass = newCopyButton.querySelector('i').className;
            const originalText = newCopyButton.querySelector('span.copy-text') ? newCopyButton.querySelector('span.copy-text').textContent : '';
            
            newCopyButton.querySelector('i').className = 'bi bi-check-lg text-success';
            if (newCopyButton.querySelector('span.copy-text')) newCopyButton.querySelector('span.copy-text').textContent = 'Copied!';
            
            setTimeout(() => {
                newCopyButton.querySelector('i').className = originalIconClass;
                if (newCopyButton.querySelector('span.copy-text')) newCopyButton.querySelector('span.copy-text').textContent = originalText;
            }, 2000);
        }).catch(err => {
            console.error('Failed to copy text: ', err);
            // Optionally, provide error feedback to the user
        });
    });
}

// Ensure your existing initializeMessageContentViewers and HTMX event listeners call this updated function.
// Example:
// function initializeMessageContentViewers(parentElement) {
//     if (!parentElement) return;
//     parentElement.querySelectorAll('.message-content-container').forEach(container => {
//         formatAndToggleMessage(container);
//     });
// }
// document.addEventListener('DOMContentLoaded', () => initializeMessageContentViewers(document.body));
// document.body.addEventListener('htmx:afterSwap', (event) => {
//     if (event.detail.target) initializeMessageContentViewers(event.detail.target);
// });
