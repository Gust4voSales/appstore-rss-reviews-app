/**
 * Environment configuration
 * Exports environment variables for use throughout the application
 */

interface Config {
  API_URL: string;
}

const config: Config = {
  // Remember to use the IP of the machine running the server
  // If running on a different machine or emulator, use the IP of the machine running the server
  API_URL: process.env.EXPO_PUBLIC_API_URL || 'http://127.0.0.1:8080',
};

// Validate required environment variables
const validateConfig = () => {
  const requiredVars = ['API_URL'] as const;

  for (const varName of requiredVars) {
    if (!config[varName]) {
      throw new Error(`Missing required environment variable: ${varName}`);
    }
  }
};

// Validate config on import
validateConfig();

export default config;
export { Config };
