// Frontend Vote UI Bug Test - Node.js version
// This test uses JSDOM to simulate a browser environment and test the frontend JavaScript

const { JSDOM } = require('jsdom');

// Create a virtual DOM environment
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<head>
    <title>Frontend Vote UI Test</title>
    <style>
        .voting-card { padding: 10px; margin: 5px; border: 1px solid #ccc; }
        .voting-card.selected { background-color: #007bff; color: white; }
        .hidden { display: none; }
    </style>
</head>
<body>
    <div id="waitingRoom" class="hidden"></div>
    <div id="app">
        <input id="storyInput" />
        <div class="voting-cards">
            <div class="voting-card" data-vote="1">1</div>
            <div class="voting-card" data-vote="3">3</div>
            <div class="voting-card" data-vote="5">5</div>
        </div>
        <div id="usersGrid"></div>
        <button id="revealBtn">Reveal</button>
        <button id="newRoundBtn">New Round</button>
        <button id="setStoryBtn">Set Story</button>
        <button id="shareBtn">Share</button>
    </div>
</body>
</html>
`, { url: 'http://localhost' });

// Get the window and document objects
const window = dom.window;
const document = window.document;

// Make them global for the functions to use
global.window = window;
global.document = document;
global.console = console;
global.alert = (msg) => console.log('ALERT:', msg);

// Planning Poker JavaScript functions (extracted from web/index.html)
let currentUser = null;
let currentUserId = null;
let isModerator = false;
let myVote = null;

function sendMessage(type, data) {
    console.log('Mock sendMessage:', type, data);
}

function vote(value) {
    // Clear previous selection
    document.querySelectorAll('.voting-card').forEach(card => {
        card.classList.remove('selected');
    });
    
    // Set new selection
    myVote = value;
    // Find the card with matching vote value
    const targetCard = document.querySelector(`[data-vote="${value}"]`);
    if (targetCard) {
        targetCard.classList.add('selected');
    }
    
    sendMessage('vote', { vote: value });
}

function newRound() {
    if (!isModerator) {
        alert('Only the moderator can start a new round');
        return;
    }
    
    // Clear UI - THIS ONLY HAPPENS FOR MODERATOR
    document.querySelectorAll('.voting-card').forEach(card => {
        card.classList.remove('selected');
    });
    myVote = null;

    sendMessage('new_round');
}

function updateSessionState(state) {
    console.log('updateSessionState called with:', JSON.stringify(state, null, 2));
    
    // Check if we need to exit waiting room
    const isInWaitingRoom = !document.getElementById('waitingRoom').classList.contains('hidden');
    
    if (state.status === 'active' && isInWaitingRoom) {
        hideWaitingRoom();
    }
    
    if (isInWaitingRoom && state.status === 'waiting') {
        return;
    }
    
    if (isInWaitingRoom) {
        return;
    }
    
    // Update story
    document.getElementById('storyInput').value = state.currentStory || '';

    // Find current user and check if they're moderator
    let currentUserData = null;
    Object.values(state.users || {}).forEach(user => {
        if (user.name === currentUser) {
            currentUserData = user;
            currentUserId = user.id;
            isModerator = user.isModerator;
        }
    });

    // Show/hide moderator controls
    const revealBtn = document.getElementById('revealBtn');
    const newRoundBtn = document.getElementById('newRoundBtn');
    const setStoryBtn = document.getElementById('setStoryBtn');
    const shareBtn = document.getElementById('shareBtn');
    const storyInput = document.getElementById('storyInput');

    if (isModerator) {
        revealBtn.style.display = 'inline-block';
        newRoundBtn.style.display = 'inline-block';
        setStoryBtn.style.display = 'inline-block';
        shareBtn.style.display = 'inline-block';
        storyInput.disabled = false;
    } else {
        revealBtn.style.display = 'none';
        newRoundBtn.style.display = 'none';
        setStoryBtn.style.display = 'none';
        shareBtn.style.display = 'none';
        storyInput.disabled = true;
    }

    // BUG: This function should clear vote UI when votes are cleared
    // but it currently doesn't detect when a new round has started
    
    // The fix should go here - detect if votes were cleared and reset UI
}

function hideWaitingRoom() {
    document.getElementById('waitingRoom').classList.add('hidden');
    document.getElementById('app').classList.remove('hidden');
}

// Test Functions
function testModeratorVoteClearing() {
    console.log('\n📝 Test 1: Moderator vote clearing');
    currentUser = 'Alice';
    isModerator = true;
    
    // Moderator votes
    vote('5');
    let selectedCards = document.querySelectorAll('.voting-card.selected');
    console.log('✓ Moderator voted, selected cards:', selectedCards.length);
    console.log('✓ Moderator myVote:', myVote);
    
    // Moderator starts new round
    newRound();
    selectedCards = document.querySelectorAll('.voting-card.selected');
    console.log('✓ After new round, selected cards:', selectedCards.length);
    console.log('✓ After new round, myVote:', myVote);
    
    const passed = selectedCards.length === 0 && myVote === null;
    console.log(passed ? '✅ Test 1 PASSED: Moderator vote UI correctly cleared' : '❌ Test 1 FAILED');
    return passed;
}

function testParticipantVoteBug() {
    console.log('\n📝 Test 2: Participant vote clearing (BUG REPRODUCTION)');
    currentUser = 'Bob';
    isModerator = false;
    myVote = null; // Reset from previous test
    
    // Clear any existing selections
    document.querySelectorAll('.voting-card').forEach(card => {
        card.classList.remove('selected');
    });
    
    // Participant votes
    vote('3');
    let selectedCards = document.querySelectorAll('.voting-card.selected');
    console.log('✓ Participant voted, selected cards:', selectedCards.length);
    console.log('✓ Participant myVote:', myVote);
    
    // Simulate receiving session state update after moderator starts new round
    const sessionStateAfterNewRound = {
        status: 'active',
        currentStory: 'Test Story',
        votesRevealed: false,
        users: {
            'alice-id': {
                id: 'alice-id',
                name: 'Alice',
                isModerator: true,
                vote: null  // Backend correctly cleared vote
            },
            'bob-id': {
                id: 'bob-id', 
                name: 'Bob',
                isModerator: false,
                vote: null  // Backend correctly cleared vote
            }
        }
    };
    
    updateSessionState(sessionStateAfterNewRound);
    selectedCards = document.querySelectorAll('.voting-card.selected');
    console.log('✓ After session state update, selected cards:', selectedCards.length);
    console.log('✓ After session state update, myVote:', myVote);
    
    const bugReproduced = selectedCards.length > 0 || myVote !== null;
    if (bugReproduced) {
        console.log('❌ Test 2 FAILED: Participant vote UI NOT cleared (BUG REPRODUCED)');
        console.log('   🐛 Vote button still selected or myVote still set');
        console.log('   🔧 This confirms the bug exists!');
    } else {
        console.log('✅ Test 2 PASSED: Participant vote UI correctly cleared');
    }
    
    return !bugReproduced; // Return true if bug is NOT reproduced (i.e., fixed)
}

function runAllTests() {
    console.log('🧪 Running Frontend Vote UI Tests with JSDOM...\n');
    
    const moderatorPassed = testModeratorVoteClearing();
    const participantPassed = testParticipantVoteBug();
    
    console.log('\n📊 Test Summary:');
    console.log('Moderator vote clearing:', moderatorPassed ? '✅ WORKS' : '❌ BROKEN');
    console.log('Participant vote clearing:', participantPassed ? '✅ WORKS' : '🐛 BUG CONFIRMED');
    
    if (!participantPassed) {
        console.log('\n💡 To fix: Add vote clearing logic to updateSessionState() function');
        console.log('   The function should detect when all users have vote: null');
        console.log('   and clear the .selected class from voting cards + reset myVote');
    }
    
    // Exit with appropriate code for CI/CD
    const allPassed = moderatorPassed && participantPassed;
    process.exit(allPassed ? 0 : 1);
}

// Check if JSDOM is available
try {
    runAllTests();
} catch (error) {
    console.error('❌ Test failed to run. Install jsdom: npm install jsdom');
    console.error('Error:', error.message);
    process.exit(1);
}
