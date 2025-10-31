const cssModulesPlugin = require("esbuild-css-modules-plugin");
const esbuild = require("esbuild");

esbuild.build({
  entryPoints: ["src/index.tsx"],
  bundle: true,
  outfile: "dist/index.js",
  format: "iife",
  globalName: "app",
  jsx: "automatic",
  plugins: [
    cssModulesPlugin({
      v2: true,
      localsConvention: "camelCase",
    }),
  ],
});

