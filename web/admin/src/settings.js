module.exports = {
  title: 'GoTribe-Admin',

  logo: require('@/assets/sidebar-logo/logo.png'),

  /**
   * @type {boolean} true | false
   * @description Whether show the settings right-panel
   */
  showSettings: true,

  /**
   * @type {boolean} true | false
   * @description Whether need tagsView
   */
  tagsView: true,

  /**
   * @type {boolean} true | false
   * @description Whether fix the header
   */
  fixedHeader: false,

  /**
   * @type {boolean} true | false
   * @description Whether show the logo in sidebar
   */
  sidebarLogo: true,

  /**
   * @type {string | array} 'production' | ['production', 'development']
   * @description Need show err logs component.
   * The default is only used in the production env
   * If you want to also use it in dev, you can pass ['production', 'development']
   */
  errorLog: 'production',

  publickey: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2+/BbwvSv2288ez8cWL9
4Vq1fgaLzcr6+wqPUfsmITnj56ougIrQZPgpdWRCcgKApyHO6S+BYbqiDYlAJxD5
+D7U0G9oZaPLvBJk/zsaU8wm6abW56L/DPrEuqw//0SWgagps4N41D8gMVLd5ThE
K4IH97/w6RyHvk/5B9djIjhVXid+56EsyZ+14ktNsI7Zsk5u0hLCBzAq2xQqKCAD
KSi0wZTIFGltzgDnzCuehWdHlL5Rdp2gJRcwkcOsXA9CRwEJtWFJAcc+2YhssZ/N
8k4eibBKIpS9dxgIR0aoOTma578lRZvRche4JKOdTxf/lfgc7oct9eoSj9bJL+bH
vwIDAQAB
-----END PUBLIC KEY-----`
}
