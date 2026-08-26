const cheerio = require('cheerio');
const { chromium } = require('playwright');

(async () => {
  const b = await chromium.launch({ headless: true });
  const p = await b.newPage();
  await p.goto('http://localhost:11001/ui/components/button', { waitUntil: 'networkidle', timeout: 15000 });
  const html = await p.content();
  await b.close();
  
  const $ = cheerio.load(html);
  const root = $('#layout-content').first();
  console.log('root length:', root.length);
  console.log('root name:', root[0]?.name);
  console.log('root class:', root.attr('class'));
  console.log('children count:', root.children().length);
  
  root.children().each((i, c) => {
    if (i < 3) {
      const el = $(c);
      console.log("  child " + i + ": <" + c.name + "> class=\"" + (el.attr('class')||'') + "\" children=" + el.children().length);
    }
  });
})();
