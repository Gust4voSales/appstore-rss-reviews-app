/**
 * Environment configuration
 * Exports environment variables for use throughout the application
 */

interface Config {
  API_URL: string;
}

const config: Config = {
  API_URL: import.meta.env.VITE_API_URL,
};

// Validate required environment variables
const validateConfig = () => {
  const requiredVars = ["API_URL"] as const;

  for (const varName of requiredVars) {
    if (!config[varName]) {
      throw new Error(`Missing required environment variable: ${varName}`);
    }
  }
};

// Validate config on import
validateConfig();

export default config;
export type { Config };
