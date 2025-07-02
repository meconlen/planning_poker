// Simple test to debug the production JavaScript extraction
const fs = require('fs');
const path = require('path');

console.log('🔍 Testing JavaScript extraction...');

try {
    const htmlPath = path.join(__dirname, '..', 'web', 'index.html');
    console.log('📄 Reading HTML from:', htmlPath);
    
    const htmlContent = fs.readFileSync(htmlPath, 'utf8');
    console.log('✅ HTML file read successfully, length:', htmlContent.length);
    
    const scriptMatch = htmlContent.match(/<script[^>]*>([\s\S]*?)<\/script>/);
    console.log('🔍 Script tag found:', !!scriptMatch);
    
    if (scriptMatch) {
        console.log('📏 JavaScript code length:', scriptMatch[1].length);
        console.log('📝 First 200 characters of JS:');
        console.log(scriptMatch[1].substring(0, 200) + '...');
    } else {
        console.error('❌ No <script> tag found in HTML');
    }
} catch (error) {
    console.error('❌ Error:', error.message);
}
