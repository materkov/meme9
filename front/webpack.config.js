const path = require('path');

module.exports = {
    entry: {
        React: ['react', 'react-dom'],
        Global: {
            import: [
                './src/entrypoints/Global.tsx',
                './src/DataFetcher.ts',
                './src/JsFetcher.ts',
                './src/RouteResolver.ts',
            ],
            dependOn: ['React'],
        },
        LoginPage: {
            import: './src/entrypoints/LoginPage.tsx',
            dependOn: ['React', 'Global'],
        },
        PostPage: {
            import: './src/entrypoints/PostPage.tsx',
            dependOn: ['React', 'Global'],
        },
        UserPage: {
            import: './src/entrypoints/UserPage.tsx',
            dependOn: ['React', 'Global'],
        },
        Composer: {
            import: './src/entrypoints/Composer.tsx',
            dependOn: ['React', 'Global'],
        },
        Feed: {
            import: './src/entrypoints/Feed.tsx',
            dependOn: ['React', 'Global'],
        },
        Index: {
            import: './src/entrypoints/Index.tsx',
            dependOn: ['React', 'Global'],
        },
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
        ],
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
    },
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, 'dist'),
    },
    devServer: {
        proxy: {
            '/': {
                target: 'http://localhost:8000',
                bypass: function (req, res, proxyOptions) {
                    if (req.originalUrl.startsWith('/static')) {
                        return req.originalUrl.substring(8);
                    }
                    if (req.originalUrl.startsWith('/test.html')) {
                        return req.originalUrl;
                    }
                    return null;
                }
            },
        },
        publicPath: '/static',
        contentBase: path.join(__dirname, 'dist'),
        port: 3000
    }
};
