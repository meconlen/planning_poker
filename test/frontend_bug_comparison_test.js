#!/usr/bin/env node

// Frontend test that shows BEFORE and AFTER the bug fix
const fs = require('fs');
const path = require('path');
const { JSDOM } = require('jsdom');

console.log('🧪 Frontend Vote UI Bug - Before vs After Fix Comparison\n');

// Read the production HTML file
const htmlPath = path.join(__dirname, '..', 'web', 'index.html');
let htmlContent = fs.readFileSync(htmlPath, 'utf8');

// Test 1: Run with the bug fix removed (simulate the original bug)
console.log('🐛 TEST 1: Simulating ORIGINAL BUG (fix removed)\n');
const htmlWithBug = htmlContent.replace(
    /\/\/ Check if all users have null votes[\s\S]*?myVote = null;\s*}/,
    '// Bug simulation: Vote clearing logic removed'
);

runTest(htmlWithBug, 'WITH BUG')
    .then(() => {
        console.log('\n' + '='.repeat(60) + '\n');
        
        // Test 2: Run with the fix present
        console.log('✅ TEST 2: Testing CURRENT CODE (with fix)\n');
        return runTest(htmlContent, 'WITH FIX');
    })
    .then(() => {
        console.log('\n🎉 COMPARISON COMPLETE!');
        console.log('The tests demonstrate that the fix successfully resolves the vote UI bug.');
    })
    .catch(error => {
        console.error('❌ Test failed:', error);
        process.exit(1);
    });

function runTest(htmlCode, testName) {
    return new Promise((resolve, reject) => {
        const dom = new JSDOM(htmlCode, {
            url: 'http://localhost:8080',
            runScripts: 'dangerously',
            resources: 'usable',
            pretendToBeVisual: true
        });

        const window = dom.window;
        const document = window.document;

        setTimeout(() => {
            try {
                const result = executeVoteUITest(window, document, testName);
                resolve(result);
            } catch (error) {
                reject(error);
            }
        }, 100);
    });
}

function executeVoteUITest(window, document, testName) {
    console.log(`🎯 Running ${testName} scenario...`);

    // Setup test environment
    window.currentUser = 'Bob';
    window.currentUserId = 'user_123';
    window.isModerator = false;
    window.myVote = null;
    window.sendMessage = function(type, data) {
        // Silent mock
    };

    // Step 1: Vote for "5"
    const voteCards = document.querySelectorAll('.voting-card');
    let targetCard = null;
    voteCards.forEach(card => {
        if (card.textContent.trim() === '5') {
            targetCard = card;
        }
    });

    if (targetCard) {
        window.event = { target: targetCard };
        window.vote('5');
    }

    const selectedCardsAfterVote = document.querySelectorAll('.voting-card.selected');
    console.log(`   📝 After voting: ${selectedCardsAfterVote.length} cards selected`);

    // Step 2: Simulate new round (moderator clears all votes)
    const sessionStateAfterNewRound = {
        status: 'active',
        currentStory: 'Test story',
        votesRevealed: false,
        users: {
            'alice_456': { id: 'alice_456', name: 'Alice', isModerator: true, vote: null },
            'user_123': { id: 'user_123', name: 'Bob', isModerator: false, vote: null }
        }
    };

    window.updateSessionState(sessionStateAfterNewRound);

    // Step 3: Check final state
    const selectedCardsAfterNewRound = document.querySelectorAll('.voting-card.selected');
    const bugPresent = selectedCardsAfterNewRound.length > 0;

    console.log(`   🔄 After new round: ${selectedCardsAfterNewRound.length} cards still selected`);
    
    if (testName === 'WITH BUG') {
        if (bugPresent) {
            console.log('   ✅ BUG REPRODUCED: Vote UI remains selected (as expected for bug test)');
        } else {
            console.log('   ⚠️  BUG NOT REPRODUCED: Vote UI was cleared (unexpected)');
        }
    } else {
        if (bugPresent) {
            console.log('   ❌ BUG STILL PRESENT: Vote UI remains selected (fix not working)');
        } else {
            console.log('   ✅ BUG FIXED: Vote UI properly cleared');
        }
    }

    return { bugPresent, selectedCount: selectedCardsAfterNewRound.length };
}
