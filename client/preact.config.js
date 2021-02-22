module.exports = function (config, env, helpers) {
    let { index } = helpers.getPluginsByName(config, 'SizePlugin')[0]
    config.plugins.splice(index, 1)
    if (config.devServer) {
        config.devServer.proxy = [
            {
                path: ['/graphql', '/wsgraphql'],
                target: 'http://localhost:80',
                ws: true
            }
        ];
    }
};