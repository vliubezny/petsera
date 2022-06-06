import { fileURLToPath, URL } from "url";

import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";
import { createHtmlPlugin } from "vite-plugin-html";

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());
  const appConfig =
    mode === "production"
      ? "__APP_CONFIG__"
      : JSON.stringify({
          apiKey: env.VITE_GOOGLE_MAPS_API_KEY,
        });

  return {
    plugins: [
      vue(),
      createHtmlPlugin({
        inject: {
          data: {
            appConfig,
          },
        },
      }),
    ],
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url)),
      },
    },
    server: {
      proxy: {
        "/api": {
          target: "http://localhost:8080",
          changeOrigin: true,
        },
      },
    },
  };
});
