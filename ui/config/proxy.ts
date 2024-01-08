const proxy = {
  '/api': {
    // target: 'http://starlite-first.ns-xx7juor7.svc.cluster.local:8080',
    // target: 'https://kdwqevbfhfcv.cloud.sealos.io',
    target: 'http://127.0.0.1:8080',
    
    changeOrigin: true,
    pathRewrite: { '^/api': '/api/starlite' },
  }
}

export default proxy;
