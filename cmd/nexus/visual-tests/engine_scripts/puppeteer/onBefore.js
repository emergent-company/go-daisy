module.exports = async (page, scenario) => {
  await page.setViewport({ width: 1440, height: 900 });
  console.log('Setting up scenario: ' + scenario.label);
};
