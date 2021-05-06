const path = require('path');

module.exports = {
    entry: {
        App: {
            import: [
                './src/App.tsx'
            ]
        },
    },
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, 'dist'),
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
            {
                test: /\.css$/,
                loader: 'style-loader'
            },
            {
                test: /\.css$/,
                loader: 'css-loader',
                options: {
                    modules: {
                        localIdentName: "[name]--[hash:base64:5]",
                    }
                }
            },
        ],
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js', '.css'],
    },
    devServer: {
        proxy: {
            '/api': {
                target: 'http://localhost:8000',
            },
            '/vk-callback': {
                target: 'http://localhost:8000',
            },
            '/logout': {
                target: 'http://localhost:8000',
            },
        },
        historyApiFallback: true,
        contentBase: path.join(__dirname, 'dist'),
        publicPath: '/static',
        port: 3000
    }
};
