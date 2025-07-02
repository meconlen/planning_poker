// Simplified Frontend Vote UI Bug Test - focuses on the specific functions we need to test
const { JSDOM } = require('jsdom');
const fs = require('fs');
const path = require('path');

console.log('🧪 Running Simplified Frontend Vote UI Test...');
console.log('📄 Using production JavaScript from web/index.html\n');

// Create a minimal JSDOM environment
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
    <div id="waitingRoom" class="hidden"></div>
    <div id="app">
        <input id="storyInput" />
        <div class="voting-cards">
            <div class="voting-card">1</div>
            <div class="voting-card">3</div>
            <div class="voting-card">5</div>
        </div>
        <button id="revealBtn">Reveal</button>
        <button id="newRoundBtn">New Round</button>
        <button id="setStoryBtn">Set Story</button>
        <button id="shareBtn">Share</button>
    </div>
    <style>
        .voting-card.selected { background-color: blue; }
        .hidden { display: none; }
    </style>
</body>
</html>
`);

// Set up global environment
global.window = dom.window;
global.document = dom.window.document;
global.console = console;

// Extract the specific functions we need from production
const htmlPath = path.join(__dirname, '..', 'web', 'index.html');
const htmlContent = fs.readFileSync(htmlPath, 'utf8');

// Extract vote function
const voteMatch = htmlContent.match(/function vote\(value\)\s*{([^}]+(?:{[^}]*}[^}]*)*?)}/);
const updateSessionStateMatch = htmlContent.match(/function updateSessionState\(state\)\s*{([^{}]*(?:{[^{}]*(?:{[^{}]*}[^{}]*)*}[^{}]*)*?)}/);
const newRoundMatch = htmlContent.match(/function newRound\(\)\s*{([^}]+(?:{[^}]*}[^}]*)*?)}/);

if (!voteMatch || !updateSessionStateMatch || !newRoundMatch) {
    console.error('❌ Could not extract required functions from production HTML');
    process.exit(1);
}

// Set up the global variables
let socket = null;
let currentSession = null;
let currentUser = null;
let currentUserId = null;
let isModerator = false;
let myVote = null;

// Mock sendMessage
function sendMessage(type, data) {
    console.log('Mock sendMessage:', type, data);
}

// Mock event
global.event = { target: null };

// Create the functions by eval
try {
    eval(`function vote(value) {${voteMatch[1]}}`);
    eval(`function updateSessionState(state) {${updateSessionStateMatch[1]}}`);
    eval(`function newRound() {${newRoundMatch[1]}}`);
    console.log('✅ Successfully extracted and loaded production functions\n');
} catch (error) {
    console.error('❌ Error loading production functions:', error.message);
    process.exit(1);
}

// Test the bug
function runTest() {
    console.log('📝 Test: Participant vote clearing after new round');
    
    // Setup as participant
    currentUser = 'Bob';
    isModerator = false;
    
    // Vote
    global.event = { target: global.document.querySelector('.voting-card') };
    vote('3');
    global.event.target.classList.add('selected'); // Simulate the UI selection
    
    let selectedCards = global.document.querySelectorAll('.voting-card.selected');
    console.log('✓ Participant voted, selected cards:', selectedCards.length);
    console.log('✓ Participant myVote:', myVote);
    
    // Simulate session state update after new round (all votes cleared)
    const sessionState = {
        status: 'active',
        currentStory: 'Test Story',
        votesRevealed: false,
        users: {
            'alice-id': { id: 'alice-id', name: 'Alice', isModerator: true, vote: null },
            'bob-id': { id: 'bob-id', name: 'Bob', isModerator: false, vote: null }
        }
    };
    
    updateSessionState(sessionState);
    
    selectedCards = global.document.querySelectorAll('.voting-card.selected');
    console.log('✓ After session state update, selected cards:', selectedCards.length);
    console.log('✓ After session state update, myVote:', myVote);
    
    const bugFixed = selectedCards.length === 0 && myVote === null;
    
    if (bugFixed) {
        console.log('✅ TEST PASSED: Vote UI correctly cleared for participant');
        return true;
    } else {
        console.log('❌ TEST FAILED: Vote UI NOT cleared for participant (BUG)');
        console.log('   🐛 This confirms the bug exists in production code');
        return false;
    }
}

const testPassed = runTest();
process.exit(testPassed ? 0 : 1);
