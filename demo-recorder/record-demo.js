/**
 * Heywood TBS — Polished Demo Recording v3
 * BIG visible cursor, click ripple, wait-for-load, multi-role chat.
 *
 * Usage: node record-demo.js
 * Output: ./output/heywood-demo-final.mp4
 */

const { chromium } = require('playwright');
const path = require('path');
const fs = require('fs');

const SITE = 'https://heywood-tbs.nicefield-9a8db973.eastus.azurecontainerapps.io';
const OUTPUT_DIR = path.join(__dirname, 'output');
const FFMPEG = 'C:/Users/Jesse/AppData/Local/Microsoft/WinGet/Links/ffmpeg.exe';

const wait = (ms) => new Promise(r => setTimeout(r, ms));

// ─── Bezier cursor movement ─────────────────────────────────────────────────
let curX = 700, curY = 450;

function cubicBez(p0, c1, c2, p3, t) {
  const u = 1 - t;
  return {
    x: u*u*u*p0.x + 3*u*u*t*c1.x + 3*u*t*t*c2.x + t*t*t*p3.x,
    y: u*u*u*p0.y + 3*u*u*t*c1.y + 3*u*t*t*c2.y + t*t*t*p3.y,
  };
}
function ease(t) { return t < 0.5 ? 4*t*t*t : 1 - Math.pow(-2*t+2, 3)/2; }

async function glide(page, x, y, ms = 600) {
  // Animate cursor by setting position from Node.js in a loop.
  // requestAnimationFrame doesn't produce video frames in headless Playwright.
  const fromX = curX, fromY = curY;
  const steps = Math.max(15, Math.round(ms / 16));
  for (let i = 0; i <= steps; i++) {
    const t = i / steps;
    const et = t < 0.5 ? 4*t*t*t : 1 - Math.pow(-2*t+2, 3)/2;
    const cx = Math.round(fromX + (x - fromX) * et);
    const cy = Math.round(fromY + (y - fromY) * et);
    await page.evaluate(([px, py]) => {
      const c = document.getElementById('pw-c');
      if (c) { c.style.left = px + 'px'; c.style.top = py + 'px'; }
    }, [cx, cy]);
    await wait(ms / steps);
  }
  // Move actual mouse to final position for click targeting
  await page.mouse.move(Math.round(x), Math.round(y));
  curX = x; curY = y;
}

async function glideToEl(page, loc, ms = 600) {
  const b = await loc.boundingBox();
  if (!b) return;
  await glide(page, b.x + b.width*(0.35+Math.random()*0.3),
                     b.y + b.height*(0.35+Math.random()*0.3), ms);
}

async function click(page, loc, ms = 600) {
  await glideToEl(page, loc, ms);
  await wait(80);
  // Trigger click ripple
  await page.evaluate(([cx, cy]) => {
    const r = document.createElement('div');
    r.className = 'pw-ripple';
    r.style.left = cx + 'px';
    r.style.top = cy + 'px';
    document.body.appendChild(r);
    setTimeout(() => r.remove(), 600);
  }, [curX, curY]);
  await page.mouse.down();
  await wait(50);
  await page.mouse.up();
  await wait(120);
}

async function drift(page, ms = 800) {
  await glide(page, curX + (Math.random()-0.5)*30, curY + (Math.random()-0.5)*20, ms);
}

// ─── Typing ──────────────────────────────────────────────────────────────────
async function type(page, loc, text, charDelay = 42) {
  await click(page, loc, 500);  // Mouse glides to input, clicks it
  await wait(250);
  await loc.pressSequentially(text, { delay: charDelay });
}

// ─── Scroll ──────────────────────────────────────────────────────────────────
async function scrollDown(page, px = 350) {
  await page.evaluate(p => window.scrollBy({ top: p, behavior: 'smooth' }), px);
  await wait(1000);
}
async function scrollTop(page) {
  await page.evaluate(() => window.scrollTo({ top: 0, behavior: 'smooth' }));
  await wait(1000);
}
async function chatBottom(page) {
  await page.evaluate(() => {
    const el = document.querySelector('.chat-scroll');
    if (el) el.scrollTo({ top: el.scrollHeight, behavior: 'smooth' });
  });
  await wait(800);
}

// ─── Wait for content ────────────────────────────────────────────────────────
async function waitLoaded(page) {
  // Wait for skeleton animations to clear AND real content to appear
  await page.waitForFunction(() => {
    const pulses = document.querySelectorAll('[class*="animate-pulse"]');
    const text = document.body.innerText;
    return pulses.length === 0 && (
      text.includes('ACTIVE STUDENTS') || text.includes('Students') ||
      text.includes('My Record') || text.includes('Task Inbox') ||
      text.includes('At-Risk') || text.includes('Schedule') ||
      text.includes('Heywood is ready') || text.includes('Chat with Heywood')
    );
  }, { timeout: 15000 }).catch(() => {});
  await wait(600);
}

async function waitChatReady(page) {
  const inp = page.locator('input[placeholder*="Ask about"]');
  await inp.waitFor({ state: 'visible', timeout: 10000 });
  await page.waitForFunction(() => {
    const el = document.querySelector('input[placeholder*="Ask about"]');
    return el && !el.disabled;
  }, { timeout: 10000 }).catch(() => {});
  await wait(300);
}

async function waitResponse(page, maxMs = 35000) {
  // Wait for thinking dots to appear
  await page.waitForSelector('.animate-bounce, .animate-spin', { timeout: 10000 }).catch(() => {});
  await wait(1500);  // Show the "Heywood is thinking..." animation

  // Drift cursor while waiting — alive feel
  await drift(page, 600);

  // Poll until streaming finishes
  const t0 = Date.now();
  while (Date.now() - t0 < maxMs) {
    const streaming = await page.evaluate(() =>
      !!(document.querySelector('.streaming-cursor') || document.querySelector('.animate-spin'))
    );
    if (!streaming) break;
    // Subtle drift while reading
    if (Math.random() > 0.6) await drift(page, 400);
    await wait(500);
  }
  await wait(800);
}

// ─── Role switch ─────────────────────────────────────────────────────────────
async function switchRole(page, roleText) {
  const btn = page.locator('header button:has(svg.lucide-chevron-down)');
  await click(page, btn, 500);
  await wait(600);
  await click(page, page.getByText(roleText, { exact: true }), 400);
  await wait(2500);  // Cookie round-trip + re-render
}

async function clearChat(page) {
  const btn = page.locator('button:has-text("Clear")');
  try {
    if (await btn.isVisible({ timeout: 2000 })) {
      await click(page, btn, 400);
      await wait(600);
    }
  } catch {}
}

// ─── BIG BLACK CURSOR + CLICK RIPPLE (survives navigation) ─────────────────
const CURSOR_SCRIPT = `
(() => {
  function install() {
    if (document.getElementById('pw-c')) return;

    // Styles
    const style = document.createElement('style');
    style.textContent = \`
      *, *::before, *::after { cursor: none !important; }
      #pw-c {
        position: fixed; z-index: 999999; pointer-events: none;
        will-change: left, top; left: 700px; top: 450px;
        filter: drop-shadow(1px 2px 2px rgba(0,0,0,0.3));
        transition: transform 0.06s ease;
      }
      #pw-c.clicking { transform: scale(0.82); }
      .pw-ripple {
        position: fixed; z-index: 999998; pointer-events: none;
        width: 30px; height: 30px; border-radius: 50%;
        border: 3px solid #e74c3c;
        transform: translate(-50%, -50%) scale(0.3);
        opacity: 1;
        animation: pw-rip 0.5s ease-out forwards;
      }
      @keyframes pw-rip {
        to { transform: translate(-50%, -50%) scale(2.5); opacity: 0; }
      }
    \`;
    document.head.appendChild(style);

    // Cursor element — BLACK pointer, big and obvious
    const c = document.createElement('div');
    c.id = 'pw-c';
    c.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="28" height="34" viewBox="0 0 28 34">' +
      '<path d="M3 1 L3 26 L9.5 19.5 L15 31 L19 29 L13.5 17 L22 17 Z" ' +
      'fill="#111" stroke="white" stroke-width="2" stroke-linejoin="round"/></svg>';
    document.body.appendChild(c);

    // Position is driven by page.evaluate() from Playwright — no mousemove needed.
    // But keep click animation:
    document.addEventListener('mousedown', () => c.classList.add('clicking'));
    document.addEventListener('mouseup', () => c.classList.remove('clicking'));
  }
  if (document.body) install();
  else new MutationObserver((_,o) => { if(document.body){install();o.disconnect();} })
    .observe(document.documentElement, {childList:true,subtree:true});
})();
`;

// ═══════════════════════════════════════════════════════════════════════════════
async function main() {
  fs.mkdirSync(OUTPUT_DIR, { recursive: true });
  // Clean old raw recordings
  fs.readdirSync(OUTPUT_DIR)
    .filter(f => f.endsWith('.webm') && !f.includes('demo'))
    .forEach(f => fs.unlinkSync(path.join(OUTPUT_DIR, f)));

  console.log('Launching browser...');
  const browser = await chromium.launch({ headless: true });
  const ctx = await browser.newContext({
    viewport: { width: 1440, height: 900 },
    recordVideo: { dir: OUTPUT_DIR, size: { width: 1440, height: 900 } },
    colorScheme: 'light', locale: 'en-US',
  });
  await ctx.addInitScript(CURSOR_SCRIPT);
  const page = await ctx.newPage();

  const chatInput = page.locator('input[placeholder*="Ask about"]');
  const sendBtn = page.locator('form button[type="submit"]');

  try {
    // ──────────────────────────────────────────────────────────────
    // SCENE 1 — Staff Dashboard (wait for ALL data)
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 1: Staff Dashboard');
    await page.goto(SITE, { waitUntil: 'domcontentloaded', timeout: 30000 });
    await waitLoaded(page);
    await wait(1000);

    // Cursor sweeps in from off-screen
    curX = -40; curY = 300;
    await glide(page, 400, 200, 800);  // Glide to first stat card
    await wait(600);
    await glide(page, 900, 200, 700);  // Scan across stat cards
    await wait(500);
    await scrollDown(page, 350);
    await drift(page, 500);
    await scrollDown(page, 300);
    await wait(500);
    await scrollTop(page);
    await wait(500);

    // ──────────────────────────────────────────────────────────────
    // SCENE 2 — Switch to XO
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 2: XO Dashboard');
    await switchRole(page, 'Executive Officer');
    await waitLoaded(page);
    await wait(600);

    // Quick tour of XO dashboard
    await glide(page, 500, 200, 500);
    await wait(400);
    await scrollDown(page, 400);
    await drift(page, 400);
    await scrollDown(page, 350);
    await wait(400);
    await scrollTop(page);
    await wait(300);

    // ──────────────────────────────────────────────────────────────
    // SCENE 3 — XO Chat: "Okay Heywood, what do we have for today?"
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 3: XO Chat');
    await click(page, page.locator('nav a[href="/chat"]'), 600);
    await waitChatReady(page);
    await wait(400);

    await type(page, chatInput,
      'Okay Heywood, what do we have for today?', 50);
    await wait(500);
    await click(page, sendBtn, 500);  // Visible click on Send

    console.log('  Waiting for XO response...');
    await waitResponse(page, 35000);
    await chatBottom(page);
    await drift(page, 800);
    await wait(1500);

    // ──────────────────────────────────────────────────────────────
    // SCENE 4 — Tool Use: flag student + notify SPC
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 4: Tool Use');
    await waitChatReady(page);

    await type(page, chatInput,
      'Flag 2ndLt Perez for remedial land nav and notify his SPC to coordinate.', 46);
    await wait(500);
    await click(page, sendBtn, 500);

    console.log('  Waiting for tool execution...');
    await waitResponse(page, 40000);
    await chatBottom(page);
    await drift(page, 600);
    await wait(1800);

    // ──────────────────────────────────────────────────────────────
    // SCENE 5 — Task Inbox (showcase the task Heywood just created)
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 5: Task Inbox');
    await click(page, page.locator('nav a[href="/tasks"]'), 600);
    await wait(2000);
    await waitLoaded(page);
    await wait(1000);

    // Hover over the task list area to let viewer read
    await glide(page, 700, 300, 600);
    await wait(1500);

    // Click on a task card to expand details
    const taskCard = page.locator('[class*="border-l-4"]').first();
    try {
      await taskCard.waitFor({ state: 'visible', timeout: 5000 });
      await click(page, taskCard, 600);
      await wait(2500);  // Hold on task detail so viewer can read it
      await drift(page, 500);
      await wait(1500);
    } catch {}

    // ──────────────────────────────────────────────────────────────
    // SCENE 6 — Students Page
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 6: Students');
    await click(page, page.locator('nav a[href="/students"]'), 600);
    await wait(1500);
    await waitLoaded(page);
    await scrollDown(page, 300);
    await wait(500);
    await scrollTop(page);
    await wait(400);

    // ──────────────────────────────────────────────────────────────
    // SCENE 7 — At-Risk
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 7: At-Risk');
    await click(page, page.locator('nav a[href="/at-risk"]'), 600);
    await wait(1500);
    await waitLoaded(page);
    await scrollDown(page, 300);
    await wait(500);

    // ──────────────────────────────────────────────────────────────
    // SCENE 8 — Schedule
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 8: Schedule');
    await click(page, page.locator('nav a[href="/schedule"]'), 600);
    await wait(1500);
    await waitLoaded(page);
    await scrollDown(page, 300);
    await wait(600);

    // ──────────────────────────────────────────────────────────────
    // SCENE 9 — Staff Chat (peer tone)
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 9: Staff Chat');
    await switchRole(page, 'Staff Officer');
    await click(page, page.locator('nav a[href="/chat"]'), 600);
    await waitChatReady(page);
    await clearChat(page);

    await type(page, chatInput,
      'Pull up the qualification coverage gaps — any critical shortfalls for next month?', 46);
    await wait(500);
    await click(page, sendBtn, 500);

    console.log('  Waiting for Staff response...');
    await waitResponse(page, 35000);
    await chatBottom(page);
    await drift(page, 600);
    await wait(1500);

    // ──────────────────────────────────────────────────────────────
    // SCENE 10 — SPC: clicks notification bell → task inbox → chat
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 10: SPC Notifications + Chat');
    await switchRole(page, 'SPC (Alpha Co)');
    await waitLoaded(page);
    await wait(800);

    // Click the notification bell — viewer sees the badge count
    const bell = page.locator('button[title="Task inbox"]');
    try {
      await bell.waitFor({ state: 'visible', timeout: 5000 });
      await glideToEl(page, bell, 700);  // Slow approach so viewer notices badge
      await wait(1200);
      await click(page, bell, 400);  // Click the bell
      await wait(2000);
      await waitLoaded(page);

      // Now on task inbox — viewer sees the task from Heywood
      const spcTask = page.locator('[class*="border-l-4"]').first();
      try {
        await spcTask.waitFor({ state: 'visible', timeout: 5000 });
        await glide(page, 700, 300, 500);
        await wait(1500);
        await click(page, spcTask, 600);  // Click to expand
        await wait(2500);  // Hold so viewer reads the task detail
        await drift(page, 500);
        await wait(1000);
      } catch {}
    } catch {
      // Fallback: go to tasks via sidebar
      await click(page, page.locator('nav a[href="/tasks"]'), 600);
      await wait(2000);
      await waitLoaded(page);
    }

    // Now go to chat
    await click(page, page.locator('nav a[href="/chat"]'), 600);
    await waitChatReady(page);
    await clearChat(page);

    await type(page, chatInput,
      "What's the status on my Alpha Company students? Anyone falling behind?", 46);
    await wait(500);
    await click(page, sendBtn, 500);

    console.log('  Waiting for SPC response...');
    await waitResponse(page, 35000);
    await chatBottom(page);
    await drift(page, 600);
    await wait(1500);

    // ──────────────────────────────────────────────────────────────
    // SCENE 11 — Student: My Record + personal chat
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 11: Student View');
    await switchRole(page, 'Student');

    try {
      const mr = page.locator('nav a[href="/my-record"]');
      await mr.waitFor({ state: 'visible', timeout: 4000 });
      await click(page, mr, 600);
      await wait(1500);
      await waitLoaded(page);
      await scrollDown(page, 250);
      await wait(600);
      await scrollTop(page);
      await wait(400);
    } catch {}

    await click(page, page.locator('nav a[href="/chat"]'), 600);
    await waitChatReady(page);
    await clearChat(page);

    await type(page, chatInput,
      'Hey Heywood, how am I doing overall? What should I focus on to improve?', 46);
    await wait(500);
    await click(page, sendBtn, 500);

    console.log('  Waiting for Student response...');
    await waitResponse(page, 35000);
    await chatBottom(page);
    await drift(page, 600);
    await wait(1500);

    // ──────────────────────────────────────────────────────────────
    // SCENE 12 — Mobile
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 12: Mobile');
    await page.setViewportSize({ width: 390, height: 844 });
    await page.goto(SITE, { waitUntil: 'domcontentloaded', timeout: 20000 });
    await waitLoaded(page);
    curX = 195; curY = 400;
    await wait(1200);

    const menu = page.locator('button:has(svg.lucide-menu)');
    try {
      await menu.waitFor({ state: 'visible', timeout: 5000 });
      await click(page, menu, 500);
      await wait(1000);
      await click(page, page.locator('nav a[href="/students"]'), 500);
      await wait(1500);
      await waitLoaded(page);
    } catch {}
    await scrollDown(page, 200);
    await wait(800);

    // ──────────────────────────────────────────────────────────────
    // SCENE 13 — Final: XO Dashboard hero shot
    // ──────────────────────────────────────────────────────────────
    console.log('Scene 13: Final Hero');
    await page.setViewportSize({ width: 1440, height: 900 });
    await page.goto(SITE, { waitUntil: 'domcontentloaded', timeout: 20000 });
    curX = 700; curY = 450;
    await switchRole(page, 'Executive Officer');
    await waitLoaded(page);
    await glide(page, 700, 350, 800);
    await wait(2500);

    console.log('All scenes recorded!');

  } catch (err) {
    console.error('ERROR:', err.message);
    console.error(err.stack?.split('\n').slice(0,4).join('\n'));
  }

  await ctx.close();
  await browser.close();

  // ─── Rename + Convert ──────────────────────────────────────────
  const files = fs.readdirSync(OUTPUT_DIR).filter(f => f.endsWith('.webm') && !f.includes('demo'));
  if (!files.length) { console.log('No video found!'); return; }

  const src = path.join(OUTPUT_DIR, files.sort().pop());
  const webm = path.join(OUTPUT_DIR, 'heywood-demo-final.webm');
  const mp4 = path.join(OUTPUT_DIR, 'heywood-demo-final.mp4');
  try { fs.unlinkSync(webm); } catch {}
  try { fs.unlinkSync(mp4); } catch {}
  fs.renameSync(src, webm);
  console.log(`WebM: ${(fs.statSync(webm).size/1048576).toFixed(1)} MB`);

  console.log('Converting to MP4...');
  try {
    require('child_process').execSync(
      `"${FFMPEG}" -i "${webm}" -c:v libx264 -crf 18 -preset slow -pix_fmt yuv420p -movflags +faststart "${mp4}" -y`,
      { stdio: 'pipe' }
    );
    console.log(`MP4: ${mp4} (${(fs.statSync(mp4).size/1048576).toFixed(1)} MB)`);
  } catch { console.log('MP4 failed, webm available'); }
}

main().catch(console.error);
