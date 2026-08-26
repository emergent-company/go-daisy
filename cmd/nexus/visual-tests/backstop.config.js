const path = require('path');

module.exports = {
  id: 'nexus_visual_tests',
  viewports: [
    { label: 'desktop', width: 1440, height: 900 },
  ],
  scenarios: [
    { label: 'landing', url: 'http://localhost:11001/landing', referenceUrl: 'http://localhost:11002/landing.html', misMatchThreshold: 15.0 },
    { label: 'ecommerce', url: 'http://localhost:11001/dashboards/ecommerce', referenceUrl: 'http://localhost:11002/dashboards-ecommerce.html', misMatchThreshold: 10.0 },
    { label: 'crm', url: 'http://localhost:11001/dashboards/crm', referenceUrl: 'http://localhost:11002/dashboards-crm.html', misMatchThreshold: 20.0 },
    { label: 'products', url: 'http://localhost:11001/apps/ecommerce/products', referenceUrl: 'http://localhost:11002/apps-ecommerce-products.html', misMatchThreshold: 10.0 },
    { label: 'chat', url: 'http://localhost:11001/apps/chat', referenceUrl: 'http://localhost:11002/apps-chat.html', misMatchThreshold: 10.0 },
    { label: 'file-manager', url: 'http://localhost:11001/apps/file-manager', referenceUrl: 'http://localhost:11002/apps-file-manager.html', misMatchThreshold: 10.0 },
    { label: 'settings', url: 'http://localhost:11001/pages/settings', referenceUrl: 'http://localhost:11002/pages-settings.html', misMatchThreshold: 10.0 },
    { label: 'get-help', url: 'http://localhost:11001/pages/get-help', referenceUrl: 'http://localhost:11002/pages-get-help.html', misMatchThreshold: 10.0 },
    { label: 'login', url: 'http://localhost:11001/auth/login', referenceUrl: 'http://localhost:11002/auth-login.html', misMatchThreshold: 5.0 },
    { label: 'sidebar-menu', url: 'http://localhost:11001/dashboards/ecommerce', referenceUrl: 'http://localhost:11002/dashboards-ecommerce.html', selectors: ['#sidebar-menu'], selectorExpansion: false, misMatchThreshold: 5.0 },
    { label: 'topbar', url: 'http://localhost:11001/dashboards/ecommerce', referenceUrl: 'http://localhost:11002/dashboards-ecommerce.html', selectors: ['#_layout-topbar'], selectorExpansion: false, misMatchThreshold: 5.0 },
  ],
  paths: {
    bitmaps_reference: 'backstop_data/bitmaps_reference',
    bitmaps_test: 'backstop_data/bitmaps_test',
    engine_scripts: 'backstop_data/engine_scripts',
    html_report: 'backstop_data/html_report',
    ci_report: 'backstop_data/ci_report',
  },
  report: ['browser'],
  engine: 'puppeteer',
  engineOptions: {
    args: ['--no-sandbox', '--disable-setuid-sandbox'],
  },
  asyncCaptureLimit: 3,
  debugWindow: false,
};
