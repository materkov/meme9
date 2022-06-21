#!/usr/bin/env node

const cssModulesPlugin = require("esbuild-css-modules-plugin");

require("esbuild")
    .build({
        logLevel: "debug",
        entryPoints: ["src/index.tsx"],
        bundle: true,
        watch: true,
        outfile: "dist/bundle.js",
        plugins: [cssModulesPlugin()],
    })
    .catch(() => process.exit(1));
