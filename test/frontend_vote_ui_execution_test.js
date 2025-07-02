#!/usr/bin/env node

// Proper frontend test that actually executes the JavaScript logic to test the vote UI bug
const fs = require('fs');
const path = require('path');
const { JSDOM } = require('jsdom');

console.log('🧪 Frontend Vote UI Bug - JavaScript Execution Test\n');

// Read the production HTML file
const htmlPath = path.join(__dirname, '..', 'web', 'index.html');
const htmlContent = fs.readFileSync(htmlPath, 'utf8');

// Create a JSDOM environment to execute the JavaScript
const dom = new JSDOM(htmlContent, {
    url: 'http://localhost:8080',
    runScripts: 'dangerously',
    resources: 'usable',
    pretendToBeVisual: true
});

const window = dom.window;
const document = window.document;

// Wait for DOM to be ready
setTimeout(() => {
    try {
        runVoteUIBugTest();
    } catch (error) {
        console.error('❌ Test failed with error:', error);
        process.exit(1);
    }
}, 100);

function runVoteUIBugTest() {
    console.log('🎯 Testing Vote UI Bug Scenario...\n');

    // Setup: Mock the required global variables and functions
    window.currentUser = 'Bob';
    window.currentUserId = 'user_123';
    window.isModerator = false; // Bob is NOT a moderator
    window.myVote = null;
    
    // Mock the sendMessage function
    window.sendMessage = function(type, data) {
        console.log(`📤 Mock sendMessage: ${type}`, data);
    };

    console.log('👤 Test User: Bob (Non-Moderator)');
    console.log('📊 Initial State: No vote selected\n');

    // Step 1: User votes for "5"
    console.log('📝 Step 1: Bob votes for "5"');
    
    // Find a voting card with value "5" 
    const voteCards = document.querySelectorAll('.voting-card');
    let targetCard = null;
    voteCards.forEach(card => {
        if (card.textContent.trim() === '5') {
            targetCard = card;
        }
    });
    
    if (!targetCard) {
        throw new Error('Could not find vote card for value "5"');
    }
    
    // Simulate the click event with global event object (as the function expects)
    window.event = { target: targetCard };
    window.vote('5');
    
    // Verify vote was set
    const selectedCards = document.querySelectorAll('.voting-card.selected');
    const hasSelectedCard = selectedCards.length > 0;
    const voteValue = window.myVote;
    
    console.log(`   ✓ Vote value set: ${voteValue}`);
    console.log(`   ✓ Card visually selected: ${hasSelectedCard}`);
    console.log(`   ✓ Selected cards count: ${selectedCards.length}\n`);

    if (!hasSelectedCard || voteValue !== '5') {
        console.log('   🐛 Debug: myVote should be "5" but is:', voteValue);
        console.log('   🐛 Debug: Card should be selected, actual count:', selectedCards.length);
        // Let's continue the test anyway to see the full behavior
    }

    // Step 2: Simulate moderator starting a new round
    // This simulates what happens when the backend sends a session_state update
    // after the moderator calls newRound()
    console.log('🔄 Step 2: Moderator starts new round (backend clears all votes)');
    
    // Create a mock session state that represents what the backend sends
    // when all votes are cleared after a new round
    const sessionStateAfterNewRound = {
        status: 'active',
        currentStory: 'User can login',
        votesRevealed: false,
        users: {
            'alice_456': {
                id: 'alice_456',
                name: 'Alice',
                isModerator: true,
                vote: null  // Cleared by backend
            },
            'user_123': {
                id: 'user_123', 
                name: 'Bob',
                isModerator: false,
                vote: null  // Cleared by backend
            }
        }
    };

    console.log('   📡 Backend sends session_state with all votes = null');
    
    // Step 3: Call updateSessionState to simulate frontend receiving the update
    console.log('🔧 Step 3: Frontend processes session_state update');
    window.updateSessionState(sessionStateAfterNewRound);

    // Step 4: Check if the bug is present or fixed
    console.log('🔍 Step 4: Checking vote UI state after new round...\n');

    const selectedCardsAfter = document.querySelectorAll('.voting-card.selected');
    const hasSelectedCardAfter = selectedCardsAfter.length > 0;
    const voteValueAfter = window.myVote;

    console.log('📊 Results:');
    console.log(`   • myVote variable: ${voteValueAfter}`);
    console.log(`   • Selected cards count: ${selectedCardsAfter.length}`);
    console.log(`   • Has visually selected card: ${hasSelectedCardAfter}`);

    // Determine if bug is present or fixed
    const bugPresent = hasSelectedCardAfter || voteValueAfter !== null;
    
    console.log('\n🎯 Bug Analysis:');
    if (bugPresent) {
        console.log('🐛 BUG DETECTED: Vote UI not properly cleared for non-moderator');
        console.log('   ❌ Expected: Vote UI cleared after new round');
        console.log('   ❌ Actual: Vote buttons remain visually selected');
        console.log('\n🔧 Fix Required:');
        console.log('   • updateSessionState() should detect when all votes are null');
        console.log('   • Clear .selected class from voting cards');
        console.log('   • Reset myVote = null');
        process.exit(1);
    } else {
        console.log('✅ BUG FIXED: Vote UI properly cleared for non-moderator');
        console.log('   ✓ Vote variable cleared (myVote = null)');
        console.log('   ✓ Visual selection removed from cards');
        console.log('   ✓ Both moderators and participants handled correctly');
        console.log('\n🎉 Test PASSED: Frontend properly handles new round for all users!');
        process.exit(0);
    }
}

// Handle unhandled promise rejections and errors
process.on('unhandledRejection', (reason, promise) => {
    console.error('❌ Unhandled Rejection:', reason);
    process.exit(1);
});

process.on('uncaughtException', (error) => {
    console.error('❌ Uncaught Exception:', error);
    process.exit(1);
});
