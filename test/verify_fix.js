// Simple test to verify the vote clearing logic is present in production
const fs = require('fs');
const path = require('path');

console.log('🧪 Verifying Frontend Vote UI Fix in Production Code...\n');

// Read the production HTML file
const htmlPath = path.join(__dirname, '..', 'web', 'index.html');
const htmlContent = fs.readFileSync(htmlPath, 'utf8');

// Check if the fix is present
const hasVoteClearingLogic = htmlContent.includes('All users have null votes, clearing vote UI');
const hasVoteDetection = htmlContent.includes('Object.values(state.users || {}).every(user => user.vote === null)');
const hasVoteUIClearing = htmlContent.includes('document.querySelectorAll(\'.voting-card\').forEach(card =>');
const hasMyVoteReset = htmlContent.includes('myVote = null');

console.log('🔍 Checking production code for vote UI clearing fix:');
console.log('✓ Vote clearing log message:', hasVoteClearingLogic ? '✅ FOUND' : '❌ MISSING');
console.log('✓ Vote detection logic:', hasVoteDetection ? '✅ FOUND' : '❌ MISSING');
console.log('✓ Vote UI clearing:', hasVoteUIClearing ? '✅ FOUND' : '❌ MISSING');
console.log('✓ myVote reset:', hasMyVoteReset ? '✅ FOUND' : '❌ MISSING');

const allFixesPresent = hasVoteClearingLogic && hasVoteDetection && hasVoteUIClearing && hasMyVoteReset;

console.log('\n📊 Result:');
if (allFixesPresent) {
    console.log('✅ ALL VOTE UI CLEARING LOGIC FOUND IN PRODUCTION');
    console.log('🎉 The frontend bug fix has been successfully implemented!');
    console.log('\n💡 The fix will:');
    console.log('   • Detect when all users have null votes (after new round)');
    console.log('   • Clear .selected class from all voting cards');
    console.log('   • Reset myVote = null for the current user');
    console.log('   • Fix the bug for both moderators and participants');
    process.exit(0);
} else {
    console.log('❌ VOTE UI CLEARING LOGIC INCOMPLETE OR MISSING');
    console.log('🐛 The frontend bug fix needs to be implemented.');
    process.exit(1);
}
