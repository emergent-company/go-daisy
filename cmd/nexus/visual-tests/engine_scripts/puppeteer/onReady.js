module.exports = async (page, scenario, vp) => {
  // Force consistent dimensions
  await page.setViewport({ width: 1440, height: 2000 });
  await page.evaluate(() => {
    return new Promise(resolve => {
      // Wait for fonts/icons to render
      document.fonts.ready.then(() => {
        // Clip body to viewport height so both pages match
        document.documentElement.style.overflowY = 'hidden';
        document.body.style.overflowY = 'hidden';
        document.body.style.maxHeight = '2000px';
        document.body.style.height = 'auto';
        setTimeout(resolve, 300);
      });
    });
  });
};
