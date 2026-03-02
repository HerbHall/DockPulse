import { defineConfig, mergeConfig } from "vitest/config";
import path from "path";
import viteConfig from "./vite.config";

export default mergeConfig(
  viteConfig,
  defineConfig({
    define: {
      __APP_VERSION__: JSON.stringify("0.0.0-test"),
    },
    resolve: {
      alias: {
        "@docker/extension-api-client": path.resolve(
          __dirname,
          "src/__mocks__/@docker/extension-api-client.ts",
        ),
      },
    },
    test: {
      environment: "jsdom",
      globals: true,
      setupFiles: "./src/test-setup.ts",
    },
  }),
);
