#!/usr/bin/env node

const cssModulesPlugin = require("esbuild-css-modules-plugin");

require("esbuild")
    .build({
        logLevel: "debug",
        entryPoints: ["src/index.tsx"],
        bundle: true,
        outfile: "dist/bundle/bundle.js",
        plugins: [cssModulesPlugin({
            localsConvention: 'camelCase',
        })],
        define: {
            'process.env.NODE_ENV': '"production"',
        }
    })
    .catch(() => process.exit(1));
