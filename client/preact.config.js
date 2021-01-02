module.exports = function (config) {
    if (config.devServer) {
        config.devServer.proxy = [
            {
                path: '/graphql',
                target: 'http://localhost:80/graphql',
            }
        ];
    }
};