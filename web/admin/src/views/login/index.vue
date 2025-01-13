<template>
  <div class="login-content" :style="{ backgroundImage: `url(${imgSrc})` }">
    <div class="container">
      <div class="leftbox">
        <img src="@/assets/images/welcome.png">
        <h3>{{ title }}</h3>
        <div class="copyright">
          © {{ currentYear }}
          由
          <a
            href="https://www.gotribe.cn"
            target="_blank"
          >GoTribe</a>
          &
          <a
            href="https://www.mactribe.cn"
            target="_blank"
          >微椒网络</a>
          强力驱动
        </div>
      </div>
      <div class="rightbox">
        <el-form ref="loginForm" :model="loginForm" :rules="loginRules" class="login-form" autocomplete="on" label-position="left">

          <el-form-item prop="username" class="inputbox">
            <img src="@/assets/images/user.png">
            <el-input
              ref="username"
              v-model="loginForm.username"
              placeholder="用户名"
              name="username"
              type="text"
              tabindex="1"
              autocomplete="on"
            />
          </el-form-item>

          <el-form-item prop="password" class="inputbox">
            <img src="@/assets/images/lock.png">
            <el-input
              :key="passwordType"
              ref="password"
              v-model="loginForm.password"
              :type="passwordType"
              placeholder="密码"
              name="password"
              tabindex="2"
              autocomplete="on"
              @keyup.native="checkCapslock"
              @blur="capsTooltip = false"
              @keyup.enter.native="handleLogin"
            />
            <span class="show-pwd" @click="showPwd">
              <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
            </span>
          </el-form-item>

        </el-form>
        <button id="my_button" @click="handleLogin">立即登录</button>
      </div>
      <div class="clearfix" />
    </div>
  </div>
</template>

<script>
import JSEncrypt from 'jsencrypt'
import defaultSettings from '@/settings'
import { mapGetters } from 'vuex'

export default {
  name: 'Login',
  computed: {
    ...mapGetters([
      'systemConfig'
    ]),
    title() {
      return this.systemConfig.title
    }
  },
  data() {
    const validatePassword = (rule, value, callback) => {
      if (value.length < 6) {
        callback(new Error('The password can not be less than 6 digits'))
      } else {
        callback()
      }
    }
    return {
      currentYear: '',
      imgSrc: require('@/assets/backgd-image/backimg.jpg'),
      loginForm: {
        username: '',
        password: ''
      },
      loginRules: {
        username: [{ required: true, trigger: 'blur' }],
        password: [
          { required: true, trigger: 'blur', validator: validatePassword }
        ]
      },
      passwordType: 'password',
      publicKey: defaultSettings.publickey,
      capsTooltip: false,
      loading: false,
      redirect: undefined,
      otherQuery: {}
    }
  },
  watch: {
    $route: {
      handler: function(route) {
        const query = route.query
        if (query) {
          this.redirect = query.redirect
          this.otherQuery = this.getOtherQuery(query)
        }
      },
      immediate: true
    }
  },
  created() {
    // window.addEventListener('storage', this.afterQRScan)
  },
  mounted() {
    if (this.loginForm.username === '') {
      this.$refs.username.focus()
    } else if (this.loginForm.password === '') {
      this.$refs.password.focus()
    }
    this.getCurYear()
  },
  destroyed() {
    // window.removeEventListener('storage', this.afterQRScan)
  },
  methods: {
    getCurYear() {
      const currentDate = new Date()
      this.currentYear = currentDate.getFullYear()
    },
    checkCapslock(e) {
      const { key } = e
      this.capsTooltip = key && key.length === 1 && key >= 'A' && key <= 'Z'
    },
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
      this.$nextTick(() => {
        this.$refs.password.focus()
      })
    },
    handleLogin() {
      this.$refs.loginForm.validate((valid) => {
        if (valid) {
          this.loading = true
          // 密码RSA加密处理
          const encryptor = new JSEncrypt()
          // 设置公钥
          encryptor.setPublicKey(this.publicKey)
          // 加密密码
          const encPassword = encryptor.encrypt(this.loginForm.password)
          const encLoginForm = {
            username: this.loginForm.username,
            password: encPassword
          }
          this.$store
            .dispatch('user/login', encLoginForm)
            .then(() => {
              this.$router.push({
                path: this.redirect || '/',
                query: this.otherQuery
              })
              this.loading = false
            })
            .catch(() => {
              this.loading = false
            })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    getOtherQuery(query) {
      return Object.keys(query).reduce((acc, cur) => {
        if (cur !== 'redirect') {
          acc[cur] = query[cur]
        }
        return acc
      }, {})
    }
  }
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg: #5c646d;
$light_gray: #fff;
$cursor: #fff;

/* reset element-ui css */
.login-content {
  .el-input {
    display: inline-block;
    width: 85%;

    input {
      background: transparent;
      border: 0px;
      -webkit-appearance: none;
    }
  }

  .el-form-item {
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 5px;
    color: #454545;
  }

  .el-form-item__content {
    display: flex;
    align-items: center;
  }
}
</style>

<style lang="scss" scoped>
.login-content {
  width: 100%;
  height: 100%;
  .container {
    position: absolute;
    top: 50%;
    left: 50%;
    margin-left: -320px;
    margin-top: -180px;
    width: 640px;
    background-color: #fff;
    z-index: 9;
    box-shadow: 0px 10px 30px rgba(8, 70, 116, 0.3);
    .leftbox {
      background: url('../../assets/images/login-left-bg.jpeg') center center;
      background-size: cover;
      float: left;
      width: 320px;
      height: 320px;
      box-sizing: border-box;
      padding: 30px;
      img {
        width: 100px;
        margin: 60px auto -10px;
        display: block;
      }
      h3 {
        color: #fff;
        font-size: 20px;
        line-height: 60px;
        text-align: center;
        margin: 0;
      }
      ::v-deep .copyright {
        font-size: 12px;
        padding-bottom: 26px;
        text-align: center;
        color: #fff;
        margin-top: 70px;
        a {
          color: #dcedff;
          text-decoration: underline;
          font-weight: bold;
          font-size: 12px;
        }
      }
    }
    .rightbox {
      float: left;
      width: 320px;
      height: 320px;
      box-sizing: border-box;
      padding: 30px;
      display: flex;
      flex-direction: column;
      justify-content: center;
      .inputbox {
        line-height: 40px;
        border-bottom: 1px solid #eee;
        margin-bottom: 20px;
        position: relative;
        img {
          height: 20px;
          vertical-align: middle;
        }
        input {
          outline: none;
          border: none;
          font-family: sans-serif;
          text-indent: 10px;
          width: 100%;
        }
        .show-pwd {
          position: absolute;
          right: 10px;
          top: 2px;
          font-size: 16px;
          color: #889aa4;
          cursor: pointer;
          user-select: none;
        }
      }
      #my_button {
        width: 100%;
        line-height: 40px;
        background-color: #34aaff;
        color: #fff;
        outline: none;
        border: none;
        margin-top: 30px;
        border-radius: 20px;
        cursor: pointer;
      }
    }
  }
}
</style>
