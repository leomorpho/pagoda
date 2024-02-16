const config = {
  development: {
    WEBSITE_URL: "http://localhost:8000",
  },
  production: {
    WEBSITE_URL: "https://yourproductionurl.com",
  },
  // other environments...
};

// Default configuration
const defaultConfig = {
  WEBSITE_URL: "http://localhost:8000",
};

module.exports = { config, defaultConfig };
