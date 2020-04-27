module.exports = {
	publicPath: process.env.NODE_ENV === 'production' ? '/sigma-vue' : '/',
	devServer: {
		port: 20000,
	}
}